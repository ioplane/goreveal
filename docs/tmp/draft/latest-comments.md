# GoReveal Strategic Review — 2026-03-31

> Контекст: это результат внешнего стратегического ревью текущего состояния GoReveal, проведённого на основе полной документации репозитория (README.md + все docs/architecture/ + все docs/plans/), юридического анализа RE-бизнеса из Испании, технического design guide по библиотекам Go 1.26, SIMD-оптимизациям, и обзора headless RE-стека (Rizin/angr/Z3/Frida/capa).
>
> Задача: добавить этот документ в `docs/plans/2026-03-31-goreveal-strategic-review.md` и обновить ссылки в связанных документах.

---

## Часть 1. Текущая оценка

### Что подтверждено

Проект на 69% завершённости (platform 81%, accuracy engine 46%). Sprint 12 — правильный primary lane. Sprint 7 — правильно в maintenance mode. Sprint 11 — completed checkpoint. Все архитектурные решения в документации согласованы и непротиворечивы.

Bounded bridge chain на canonical ELF fixture (19 slices) — достаточный proof-of-concept. `sprint12-pcln-checkpoint.md` корректно диагностирует точку остановки.

### Что требует действий

Пять категорий решений, от блокирующих до стратегических.

---

## Часть 2. Немедленные действия (1–2 недели)

### 2.1 Deferred continuation tasks

Из `deferred-continuation.md` — пять задач остаются активными:

1. Compact runtime trust/evidence summary enum на `schema.RuntimeMetadata`
   - Предложенные значения: `symbol_backed | go_module_fallback | section_heuristic | absent`
   - Это должен быть один поле, не набор bool-флагов
2. Expose через `goreveal inspect runtime`
3. Red-green тесты для rich и stripped fixtures
4. Решение по зеркалированию в `export ida` / `export ghidra`
5. Reassess: runtime trust UX vs capability transfer

Рекомендация: начать с задачи 1, остальные идут последовательно.

### 2.2 Решение по лицензии — блокер публичного релиза

Текущее состояние: лицензия не зафиксирована в репозитории.

Рекомендация: **MIT** для всего репозитория (core + plugins + CLI).

Обоснование:
- Clean-room boundary строго соблюдается — MIT не создаёт конфликтов
- Hex-Rays (Бельгия) работает по аналогичной модели — коммерческий продукт из ЕС
- MIT позволяет closed-source IDA/Ghidra плагины и коммерческие интеграции
- MIT совместим с будущим dual-license для server mode (MIT core + commercial server)
- AGPL baseline tools (gore, redress) не копируются — только поведенческие сравнения

Действие: добавить `LICENSE` файл (MIT) и `license: MIT` в README badges.

Альтернатива если решишь иначе: Apache 2.0 (patent grant, но длиннее). Не рекомендуется AGPL — закроет closed-source интеграции.

---

## Часть 3. Стратегические действия (1–3 месяца)

### 3.1 Третий fixture — Windows PE Go binary

Текущее состояние: два fixture (canonical ELF + stripped ELF). Оба — одна платформа, один формат.

Рекомендация: добавить `corpus/fixtures/go-pe-buildinfo-windows-amd64/fixture.exe`.

Обоснование:
- Валидирует кросс-форматную архитектуру core/
- Открывает рынок Windows Go malware (большинство Go malware = PE)
- `debug/pe` из stdlib уже работает — низкий технический риск
- PE Go бинарники сохраняют `.gopclntab`, `.typelink`, build info — те же структуры

Создание fixture: скомпилировать тот же Go-код что canonical ELF, но `GOOS=windows GOARCH=amd64 go build`. Можно в текущем Podman dev container.

Scope первого PE slice:
- format detection = PE ✓
- build info через `debug/buildinfo` ✓
- runtime metadata: `.gopclntab`, `.typelink` addresses через PE section headers
- НЕ полный moduledata parse — только bounded section evidence аналогично первому ELF slice

### 3.2 Go code peeling — killer feature №1

Из `market-killer-features-brainstorm.md` — уже определён как «extremely high PM value».

Рекомендация: начать MVP после PE fixture.

Архитектура MVP:
- Новый модуль `engine/peeling/` (не в core — это enrichment, не recovery)
- Классификация функций: `user | stdlib | runtime | third_party`
- Источники данных (уже в schema):
  - `import_path` + `module_local` из function metadata
  - `build_info.path` для определения module boundary
  - `source_file` из pclntab line table
- Дополнительный слой (новый): Go stdlib/runtime function fingerprints
  - Сгенерировать fingerprints из эталонных Go builds (Go 1.20–1.26)
  - Хранить как встроенную sigdb-подобную таблицу
  - Матчинг по `xxhash` первых N байт функции + размер + имя паттерн
- CLI: `goreveal inspect functions --user-only` или `goreveal peel <binary>`
- Export: добавить `classification` field в export ida/ghidra payloads

Зависимости от design guide:
- `cespare/xxhash/v2` для fingerprinting (10–13 ГБ/с, zero-alloc)
- SoA layout для таблицы функций если профилирование покажет необходимость

### 3.3 Серверный стек — зафиксировать решения

Из `runtime-modes-and-storage-ideas.md` — решения уже правильные. Фиксация:

| Компонент | Решение | Обоснование |
|-----------|---------|-------------|
| Server DB | PostgreSQL 18 | pg_trgm, JSONB, pgvector ready |
| Go driver | pgx/v5 + sqlc | Zero-CGo, type-safe queries |
| Migrations | goose (не Atlas) | Embeddable Go library, embed.FS |
| API | ConnectRPC | Три протокола одним handler-ом |
| Queue | River | PostgreSQL-backed, транзакционный |
| Object storage | S3-compatible (Garage) | Для артефактов и blob-ов |
| Local DB | SQLite (modernc.org/sqlite) | Zero-CGo, текущий storage уже работает |
| Config | koanf v2 | Multi-source config для server mode |
| Errors | cockroachdb/errors | Stack traces для parser debugging |

Действие: создать `docs/architecture/2026-03-31-goreveal-server-stack-decision.md` с этой таблицей.

### 3.4 Дополнения к golangci-lint

Текущее состояние: golangci-lint импортирован из gobfd, `make lint` green.

Рекомендация: добавить два новых linter-а из golangci-lint v2:
- `modernize` — автоматические Go 1.26 идиомы
- `sloglint` — проверка корректности slog-вызовов (проект использует slog)

---

## Часть 4. Deferred backlog обновления

### 4.1 Thin Rizin adapter

Из анализа чата про headless RE-стек: Rizin + rz-pipe подтверждён как валидный headless движок.

Рекомендация: добавить в deferred backlog наравне с JEB и Binary Ninja.

Формат: `goreveal export rizin <binary>` → JSON payload.

Преимущество: `rzpipe-go` — нативный Go биндинг, GoReveal может общаться с Rizin in-process без Python/Java моста. Технически проще чем IDA/Ghidra адаптеры.

Приоритет: ниже JEB и Binary Ninja (меньше рынок). Добавить как строку в `capability-transfer-plan.md` и `feature-map.md`.

### 4.2 Sprint 13 деобфускация — конкретизированный путь

Текущее состояние: scaffold + два прохода (function-name refinement + string-segment extraction).

Для garble string encryption нужен constraint solver. Два этапа:

**Этап A (ближайший)**: External orchestration.
- GoReveal вызывает angr через subprocess для конкретных задач
- Паттерн: аналогично текущим `scripts/baseline/` wrappers
- Добавляет Python-зависимость, но быстро в разработке

**Этап B (долгосрочный)**: Native Go Z3.
- `aclements/go-z3` (Austin Clements, Go runtime team в Google)
- Более типобезопасный чем `mitchellh/go-z3`
- GoReveal решает garble constraints in-process
- CGo-зависимость от Z3 — единственное разумное исключение из zero-CGo стратегии

Рекомендация: не начинать ни один этап до PE fixture + code peeling. Добавить как заметку в Sprint 13 секцию `scrum-implementation-plan.md`.

### 4.3 MCP interop с IDA/Ghidra MCP servers

Из `agent-mcp-and-artifact-transfer-ideas.md` — GoReveal MCP surfaces уже спроектированы.

Дополнение: GoReveal MCP tools должны **дополнять** IDA/Ghidra MCP, не конкурировать.

Двухступенчатый workflow для AI-агентов:
1. `goreveal mcp → analyze_binary` → canonical schema с Go-specific metadata
2. Агент передаёт результат через `ida-pro-mcp → import_goreveal_analysis`
3. Или через `ghidra-mcp → apply_annotations`

GoReveal = Go-native knowledge source. IDA/Ghidra = analyst workspace.

Добавить как секцию в `agent-mcp-and-artifact-transfer-ideas.md`.

### 4.4 Go version tracking — function-level diff

Из `market-killer-features-brainstorm.md` — killer feature №2.

Конкретизация: GoReveal уже имеет SQLite persistence + `goreveal diff sqlite`. Текущий diff — schema-level. Для version tracking нужен **function-level diff с similarity scoring**.

Подход: Diaphora-подобный (export SQLite → match functions by hash/signature → score similarity). GoReveal schema уже содержит функции с адресами, пакетами, source metadata — достаточно для первого function-level matching без декомпилятора.

Не начинать до code peeling MVP. Добавить как заметку в `market-killer-features-brainstorm.md`.

---

## Часть 5. Бизнес-решения

### 5.1 SL в Испании — до публичного релиза

Из юридического анализа RE-бизнеса:
- Продажа GoReveal лицензий из Испании полностью легальна
- Hex-Rays (Бельгия, BCE: 0873.473.914) — прямой прецедент RE-инструмента в ЕС
- GoReveal = analytical tool, не intrusion software. Классификация EAR99
- SL (не autónomo) — для ограничения ответственности при RE-консалтинге

Действия:
1. Зарегистрировать SL до публичного релиза
2. PI Insurance €1–2M (Hiscox или Berkley España)
3. Geo-blocking для sanctioned countries в лицензии и платёжном процессе
4. Дифференцированный compliance workflow для US-клиентов (DMCA circuit split Bowers vs Vault)

### 5.2 Санкционные ограничения

Бинарное правило: Россия и Иран **полностью закрыты** для любых услуг (EU Regulation 833/2014 с поправками).

Китай: усиленный due diligence, License Exception ACE не применяется к государственным конечным пользователям.

Действие: добавить sanctions compliance секцию в будущий EULA/Terms of Service.

### 5.3 NIS2/CRA — window of opportunity

~160,000 организаций ЕС обязаны проводить security testing по NIS2. CRA (полностью применим с декабря 2027) обязывает производителей иметь процессы управления уязвимостями.

GoReveal + Go code peeling + RE-консалтинг = прямое попадание в NIS2 compliance workflow.

Рекомендация: таргетировать EU NIS2-обязанные организации как первых клиентов после публичного бета.

---

## Часть 6. Чего НЕ делать

Подтверждение позиций из документации + дополнения:

1. **Не расширять bounded bridges на том же fixture** — 19 slices достаточно, нужен PE
2. **Не начинать SIMD до профилирования** — правило проекта корректное
3. **Не строить rich TUI/web UI** — пока schema не стабилизирована
4. **Не клонировать GoResolver CFG engine** — оставить как external orchestration
5. **Не начинать server mode до лицензии + PE fixture** — tech debt на архитектуре
6. **Не публиковать GoReveal до регистрации SL** — ответственность autónomo неограничена
7. **Не менять package/type heuristics** — пока runtime-semantic truth fixture-local
8. **Не строить Go metadata knowledge network** — пока нет code peeling + version tracking базы
9. **Не добавлять CGo зависимости** — кроме Z3 в Sprint 13 (если дойдёт)

---

## Часть 7. Рекомендуемый порядок выполнения

```
1. Runtime trust summary enum          [Sprint 12, 2-3 дня]
2. Решение по лицензии (MIT)           [Блокер, 1 день]
3. PE fixture + bounded evidence        [Sprint 12, 1-2 недели]
4. Code peeling MVP                     [Новый epic, 2-3 недели]
5. Server stack scaffold                [Sprint 10, 2-3 недели]
6. Публичный бета → EU NIS2 targeting  [Бизнес, параллельно]
7. SL регистрация + PI insurance        [Бизнес, параллельно]
```

---

## Часть 8. Файлы для обновления

При добавлении этого документа, обновить ссылки в:

- `docs/plans/2026-03-20-goreveal-deferred-continuation.md` — добавить ссылку на этот review
- `docs/plans/2026-03-20-goreveal-functional-assessment.md` — добавить ссылку
- `docs/plans/2026-03-19-goreveal-scrum-implementation-plan.md` — добавить в «Quantified roadmap checkpoint»
- `docs/plans/2026-03-20-goreveal-market-killer-features-brainstorm.md` — добавить заметки по code peeling fingerprinting и version tracking function-level diff
- `docs/plans/2026-03-21-goreveal-agent-mcp-and-artifact-transfer-ideas.md` — добавить секцию MCP interop
- `docs/plans/2026-03-19-goreveal-capability-transfer-plan.md` — добавить Rizin в deferred backlog
- `docs/plans/2026-03-19-goreveal-feature-map.md` — добавить Rizin adapter, code peeling, PE fixture в mindmap
- `README.md` — добавить ссылку на этот документ в Documentation table

---

## Источники

Этот review основан на:
- Полная документация репозитория GoReveal (22 документа)
- Юридический анализ RE-бизнеса из Испании (Директива 2009/24/EC, LPI, DMCA, экспортный контроль)
- Design guide по библиотекам и паттернам Go 1.26 для GoReveal
- SIMD и CPU utilization guide для RE binary parsing
- Анализ headless RE-стека (Rizin/rz-ghidra/sigdb, angr/cle/claripy, Z3 Go bindings, capa, Frida, Diaphora)
- C/C++ RE complete guide (ELF/PE/Mach-O, IDA 9.3, Ghidra 12.0, Binary Ninja 5.2)