# GoReveal Design Guide: архитектура Go-платформы для реверс-инжиниринга

**GoReveal — это платформа реверс-инжиниринга Go-бинарников, которая объединяет подходы gore, GoReSym и redress в единую MIT-лицензированную библиотеку с CLI, gRPC API и плагинами для IDA/Ghidra.** Этот guide определяет выбор библиотек, структуру проекта, паттерны и CI/CD-пайплайн для разработки на Go 1.26 (2026). Все рекомендации основаны на текущем состоянии экосистемы Go, benchmark-данных и практиках ведущих проектов. Фокус — на минимальных зависимостях, cross-compilation без CGo и максимальной производительности при анализе бинарников размером 100MB+.

---

## 1. Стек библиотек: обоснованный выбор

### PostgreSQL: pgx/v5 + sqlc + Atlas

**jackc/pgx v5.8+** — единственный рекомендуемый PostgreSQL-драйвер для новых проектов. Нативный интерфейс pgx работает **~3× быстрее** пути через `database/sql`, поддерживает ~70 типов PostgreSQL, LISTEN/NOTIFY, COPY-протокол, логическую репликацию. lib/pq официально в maintenance mode — GitLab, Temporal, OpenTofu мигрировали на pgx.

**sqlc v1.30** генерирует нативный pgx-код при указании `sql_package: "pgx/v5"` в `sqlc.yaml`. Генерируемый `Queries` struct принимает и `*pgx.Conn`, и `*pgxpool.Pool` через интерфейс `DBTX` — никакой абстракции `database/sql` не требуется. Типы pgtype (`pgtype.Text` вместо `sql.NullString`) используются автоматически. **pgxpool** — встроенный connection pool с настраиваемыми min/max соединениями, health-check и idle-timeout.

Для миграций — **Atlas (ariga.io)**: декларативный подход «Terraform для баз данных» с auto-planned миграциями, автоматическим откатом при ошибках, lint и dry-run в CI. Atlas импортирует существующие миграции goose/golang-migrate. Для data-миграций (Go-код) можно комбинировать Atlas с goose.

```yaml
# sqlc.yaml
version: "2"
sql:
  - engine: "postgresql"
    queries: "queries/"
    schema: "schema/"
    gen:
      go:
        package: "db"
        out: "internal/db"
        sql_package: "pgx/v5"
        emit_json_tags: true
```

### SQLite: modernc.org/sqlite — без CGo

**modernc.org/sqlite v1.45+** — чистый Go-транспиляция SQLite через ccgo. Для CLI-инструмента это критично: `GOOS=linux GOARCH=arm64 go build` работает без C-компилятора. Benchmark-данные (август 2025): **для запросов modernc сравним или быстрее mattn** (760ms vs 1018ms на простых query), при конкурентном доступе — **на 24% быстрее**. Bulk-вставки медленнее в 2-3.5×, но для RE-инструмента с embedded кэшем это не bottleneck. Альтернатива: **ncruces/go-sqlite3** (WASM через wazero) — показывает **168,757 reads/sec** vs 53,598 у modernc при конкурентных чтениях.

**Решение для GoReveal:** modernc.org/sqlite для локального кэша анализа бинарников, cross-compilation для всех платформ одной командой.

### CLI: cobra для сложного RE-инструмента

**spf13/cobra (~43,300 stars)** остаётся доминантным для сложных CLI (kubectl, Docker CLI, Hugo, GitHub CLI). Для GoReveal с вложенными командами (`goreveal analyze`, `goreveal serve`, `goreveal export --format=ida`), автогенерацией shell-completions (bash/zsh/fish/powershell) и man-страниц — cobra оптимален. **kong** (~3,000 stars) привлекателен struct-based подходом, но экосистема cobra (code generation, Viper-интеграция) критична для RE-инструмента с десятками флагов.

**urfave/cli v3.6.1** — сильный вариант для проще структурированных CLI, с нативной поддержкой `context.Context`. Но для GoReveal масштаб команд требует cobra.

### Logging: slog + tint (dev) + JSONHandler (prod)

**log/slog** (stdlib Go 1.21+) — стандарт логирования в 2026. Benchmarks показывают **~40 B/op, 1 alloc/op** — конкурентоспособно с zerolog и zap. Ключевой паттерн: slog как единый API, handler-ы подменяются без изменения кода:

- **Разработка:** `lmittmann/tint` — colorized console output без зависимостей
- **Production:** `slog.NewJSONHandler` или `zapslog.Handler` (zap как backend) для максимальной производительности
- **Контекст:** `veqryn/slog-context` для OpenTelemetry TraceID/SpanID injection

zerolog и zap остаются relевантны как slog-backends: zerolog (~40 B/op, 0 allocs) — самый быстрый, zap (~168 B/op) — самый кастомизируемый. Но **slog API избавляет от lock-in**.

```go
// Инициализация slog для GoReveal
func initLogger(dev bool) *slog.Logger {
    if dev {
        return slog.New(tint.NewHandler(os.Stderr, &tint.Options{
            Level: slog.LevelDebug,
            TimeFormat: time.Kitchen,
        }))
    }
    return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
        Level: slog.LevelInfo,
    }))
}
```

### Error handling: cockroachdb/errors для стэк-трейсов

Стандартный `errors` + `fmt.Errorf` с `%w` покрывает базовые потребности (wrapping, `errors.Is`, `errors.As`, `errors.Join`). **Но stdlib не даёт stack traces** — критично для RE-инструмента с глубокими call-стеками парсинга бинарников.

**cockroachdb/errors** — drop-in replacement для `errors` и `fmt.Errorf`, добавляющий: автоматические stack traces через `%+v`, Sentry-интеграцию (`report.BuildSentryReport()`), PII-safe redacted messages, protobuf encode/decode для сетевого transport. Импортируется **8,364 пакетами**.

```go
import "github.com/cockroachdb/errors"

func parsePclntab(data []byte) (*Pclntab, error) {
    if len(data) < 16 {
        return nil, errors.New("pclntab too short") // автоматический stack trace
    }
    // ...
    return nil, errors.Wrap(err, "parsing function table")
}
```

`hashicorp/go-multierror` заменён `errors.Join` (Go 1.20+). **Не используйте его в новом коде.**

### Configuration: koanf v2 для multi-source, caarlos0/env для env-only

**koanf v2 (~3,900 stars)** заменяет viper с фундаментальными преимуществами: корректная case-sensitivity (viper lowercase всё), модульные зависимости (**бинарник ~3× меньше** viper-эквивалента), чистый Provider/Parser интерфейс, `Get()` возвращает копии а не ссылки. viper (~30,100 stars) страдает от форсированной case-insensitivity и огромного dependency tree.

Для GoReveal CLI: **caarlos0/env v11** для environment variables (zero-dependency, struct tags) + YAML/TOML через koanf для файловой конфигурации сервера.

### Serialization: stdlib JSON → encoding/json/v2

**encoding/json/v2** принят в Go 1.25 как experimental (`GOEXPERIMENT=jsonv2`). Benchmark: **до 2.1× быстрее Sonic**, ~1.8× быстрее v1 на реальных данных. True streaming, case-sensitive matching, `omitzero` tag. Рабочая группа (с ноября 2025) готовит стабилизацию — вероятно, в Go 1.26 или 1.27.

**Сейчас используем encoding/json v1** (стабильный, zero-dependency), переходим на v2 при стабилизации. json-iterator **не поддерживается**, goccy/go-json **pre-v1 с паниками**, bytedance/sonic **ограничен платформой** (JIT только amd64). Для RE CLI-инструмента JSON не является bottleneck.

**Protocol Buffers:** только `google.golang.org/protobuf` — gogo/protobuf **официально deprecated**. Кодогенерация через **buf** (`buf generate`).

### Binary parsing: stdlib + saferwall/pe

Go stdlib покрывает базу: `debug/elf`, `debug/pe`, `debug/macho`, `debug/dwarf`, `debug/gosym` (PC→line, function lookup через `.gopclntab`). **saferwall/pe** добавляет: Rich Header, ImpHash, authentihash, entropy per section, полный data directory parsing, anomaly detection — критично для анализа обфусцированных/упакованных бинарников. `encoding/binary` — основной инструмент для чтения структур с контролем byte order.

### HTTP/gRPC: ConnectRPC + Huma

**ConnectRPC (connectrpc.com/connect)** — лучший выбор для GoReveal serve mode. Одним handler-ом даёт **три протокола**: Connect (curl-friendly JSON через HTTP/1.1+), gRPC, gRPC-Web. Не нужен Envoy proxy для gRPC-Web. Работает через стандартный `net/http` — совместим с любым middleware.

```bash
# curl-friendly вызов RE API через Connect
curl --header "Content-Type: application/json" \
  --data '{"path": "/bin/malware"}' \
  https://goreveal.local/api.v1.AnalyzerService/Analyze
```

**Huma v2** (Go 1.25+) — опционально для REST с OpenAPI 3.1 документацией. Адаптеры для BunRouter (`humabunrouter`), chi, http.ServeMux. Для GoReveal: ConnectRPC для основного API (protobuf), Huma для дополнительных REST endpoints.

### Concurrency: errgroup + semaphore

**errgroup.Group** с `SetLimit(n)` — это bounded worker pool из stdlib. Покрывает **90% use cases** параллельного анализа секций бинарника. **semaphore.NewWeighted** для rate-limiting тяжёлых операций. **singleflight.Group** для дедупликации конкурентных запросов к одному бинарнику.

```go
g, ctx := errgroup.WithContext(ctx)
g.SetLimit(runtime.NumCPU())
for _, section := range binary.Sections() {
    g.Go(func() error {
        return analyzeSection(ctx, section)
    })
}
return g.Wait()
```

**uber-go/goleak** — обязателен в тестах: `defer goleak.VerifyNone(t)`.

### Build/Release: GoReleaser + cosign + syft

**GoReleaser v2.5+** — индустриальный стандарт. Multi-platform builds (`CGO_ENABLED=0`), Docker multi-arch images, Homebrew/Scoop/Winget publishing, SBOM через **syft**, подпись через **cosign** (Sigstore keyless). Для CGo (если понадобится): `CC="zig cc"` — документированная стратегия с примером `goreleaser-example-zig-cgo`.

### Testing: testify + go-cmp + stdlib fuzzing

**stretchr/testify v1.10** — де-факто стандарт (18,512 importers для `require`). `assert`/`require` для утверждений, `mock` для моков. **google/go-cmp** — для сложных struct-сравнений с `cmp.Diff()`. Stdlib **fuzzing** (`testing.F`) — **критичен для парсеров бинарных форматов** GoReveal. Go 1.24 добавил `t.Context()` и `testing.B.Loop()`.

---

## 2. Структура проекта GoReveal

### Монорепо с Go workspaces

Для GoReveal (core + CLI + IDA plugin + Ghidra extension + proto) — **монорепо** с `go.work`. Атомарные изменения через все компоненты, единый CI/CD, простой рефакторинг shared-кода.

```
goreveal/
├── go.work                    # Go workspace
├── core/                      # go.mod — основная библиотека анализа
│   ├── analyzer.go            # Точка входа: AnalyzeFile(), NewAnalyzer()
│   ├── binary.go              # BinaryFile interface + ELF/PE/Mach-O
│   ├── pclntab.go             # Парсинг PC line table
│   ├── moduledata.go          # Moduledata extraction
│   ├── types.go               # Type information recovery
│   ├── strings.go             # Go string extraction
│   ├── functions.go           # Function/method recovery
│   └── internal/
│       ├── gosym/             # Форк debug/gosym с расширениями (из GoReSym)
│       ├── elfutil/           # ELF helpers: section scanning, mmap
│       ├── peutil/            # PE helpers + saferwall/pe integration
│       └── machoutil/         # Mach-O helpers
├── cmd/
│   └── goreveal/              # go.mod — CLI-приложение
│       ├── main.go
│       └── internal/
│           ├── analyze.go     # cobra: analyze command
│           ├── serve.go       # cobra: serve command (ConnectRPC)
│           └── export.go      # cobra: export command (IDA/Ghidra JSON)
├── proto/                     # go.mod — protobuf definitions
│   └── api/v1/
│       ├── analyzer.proto
│       └── buf.gen.yaml
├── plugins/
│   ├── ida/                   # Python: IDA Pro integration
│   │   └── goreveal_ida.py
│   └── ghidra/                # Java: Ghidra extension
│       └── GoRevealScript.java
├── internal/                  # Shared internal utilities (monorepo-level)
│   ├── version/               # Build version injection
│   └── testutil/              # Shared test helpers
├── .golangci.yml
├── .goreleaser.yaml
├── Makefile
└── buf.yaml
```

### Почему не pkg/, почему internal/

**Russ Cox (Go Team Lead) напрямую критиковал** `golang-standards/project-layout` (issue #117): «Это не стандарт Go». Официальное руководство `go.dev/doc/modules/layout` рекомендует `internal/` (compiler-enforced import restriction) и `cmd/`. **pkg/ устарел** — новые проекты его не используют. Публичный API GoReveal — пакет `core/` на верхнем уровне модуля; всё остальное — в `internal/`.

### Dependency injection: manual (конструкторы)

**Manual DI — консенсус Go-сообщества в 2026.** google/wire в beta (v0.3.0, feature-frozen), uber-go/fx — runtime reflection с overhead. Для GoReveal масштаб не требует DI-фреймворка:

```go
func NewAnalyzer(file BinaryFile, logger *slog.Logger, opts ...Option) *Analyzer {
    a := &Analyzer{file: file, logger: logger, workers: runtime.NumCPU()}
    for _, opt := range opts {
        opt(a)
    }
    return a
}
```

---

## 3. Паттерны Go 1.26 для GoReveal

### Iterators (range-over-func) для streaming анализа

Go 1.23 стабилизировал range-over-func — **идеально для streaming результатов RE-анализа** без аллокации промежуточных слайсов. GoReSym v3.2 уже использует каналы для streaming pclntab candidates (решение OOM). Итераторы — лучший API:

```go
// iter.Seq2 для streaming функций из бинарника
func (a *Analyzer) Functions() iter.Seq2[int, Function] {
    return func(yield func(int, Function) bool) {
        for i, entry := range a.funcTab {
            fn := a.parseFunction(entry)
            if !yield(i, fn) {
                return
            }
        }
    }
}

// Использование: lazy evaluation, zero intermediate allocations
for i, fn := range analyzer.Functions() {
    if fn.Package == "main" {
        fmt.Printf("%d: %s at 0x%x\n", i, fn.Name, fn.EntryPC)
    }
}
```

Stdlib интеграция: `slices.All`, `slices.Sorted`, `maps.Keys`. Утилиты: **samber/lo** `it/` sub-package с 100+ lazy-функциями для Go 1.23+ sequences.

### Generics: где реально полезно

Полезные generic-паттерны для GoReveal:

- **Generic cache:** `type Cache[K comparable, V any] struct { ... }` для кэширования результатов парсинга
- **Утилиты:** `samber/lo` (~300+ type-safe helpers) — `lo.Map`, `lo.Filter`, `lo.GroupBy` для обработки коллекций типов/функций
- **Result type:** `samber/mo` — `Result[T]`, `Option[T]` для nullable результатов парсинга
- **Generic repository:** `type Repository[T any, ID comparable] interface` для storage layer

**Не используйте generics для:** сложных type hierarchies, generic methods на типах (ограничение Go), всего что проще написать конкретно.

### Error handling: sentinel + typed + wrapping

Три уровня ошибок GoReveal:

```go
// Sentinel errors — стабильные, проверяемые условия
var (
    ErrUnsupportedFormat = errors.New("unsupported binary format")
    ErrPclntabNotFound   = errors.New("pclntab not found")
    ErrCorruptedBinary   = errors.New("corrupted binary")
)

// Typed errors — структурированные данные
type ParseError struct {
    Section string
    Offset  uint64
    Cause   error
}
func (e *ParseError) Error() string { return fmt.Sprintf("parse %s at 0x%x: %v", e.Section, e.Offset, e.Cause) }
func (e *ParseError) Unwrap() error { return e.Cause }

// Wrapping — контекст на boundaries
return errors.Wrapf(err, "analyzing %s section at offset 0x%x", section.Name, section.Offset)
```

### Functional options для конфигурации анализатора

**Functional options** (Dave Cheney / Rob Pike) — доминантный паттерн Go для optional config. Uber Go Style Guide рекомендует его:

```go
type Option func(*Analyzer)

func WithWorkers(n int) Option     { return func(a *Analyzer) { a.workers = n } }
func WithLogger(l *slog.Logger) Option { return func(a *Analyzer) { a.logger = l } }
func WithTypes(b bool) Option      { return func(a *Analyzer) { a.extractTypes = b } }
func WithStrings(b bool) Option    { return func(a *Analyzer) { a.extractStrings = b } }
func WithMmap(b bool) Option       { return func(a *Analyzer) { a.useMmap = b } }

result, err := goreveal.AnalyzeFile("binary.exe",
    goreveal.WithWorkers(8),
    goreveal.WithTypes(true),
    goreveal.WithStrings(true),
    goreveal.WithMmap(true),
)
```

### Interface design: маленькие интерфейсы на стороне consumer

```go
// Определяем интерфейсы на стороне потребителя, не провайдера
type BinaryFile interface {
    ReadAt(offset uint64, size uint64) ([]byte, error)
    ReadVirtual(addr uint64, size uint64) ([]byte, error)
    Sections() []Section
    Format() FileFormat
    ByteOrder() binary.ByteOrder
    Close() error
}

// Compile-time проверка
var _ BinaryFile = (*ELFFile)(nil)
var _ BinaryFile = (*PEFile)(nil)
var _ BinaryFile = (*MachOFile)(nil)
```

### Repository/Service + транзакции для storage

```go
type AnalysisRepository interface {
    Save(ctx context.Context, result *AnalysisResult) error
    GetByHash(ctx context.Context, hash string) (*AnalysisResult, error)
    List(ctx context.Context, filter Filter) ([]*AnalysisResult, error)
}

// Atomic callback pattern для транзакций
type Store interface {
    Atomic(ctx context.Context, fn func(Store) error) error
    Analyses() AnalysisRepository
    Binaries() BinaryRepository
}
```

---

## 4. Архитектура RE-движка: уроки gore и GoReSym

### Два подхода к парсингу Go runtime структур

**GoRE (goretk/gore)** использует code-generated structs для каждой версии Go moduledata (`moduledata_1_7_32`, `moduledata_1_16_64` и т.д.), читая их через `encoding/binary.Read()`. Подход проще в поддержке, но ограничен: **только little-endian**, нет resilience к обфускации, загружает целые секции в память.

**GoReSym (mandiant/GoReSym)** форкает код Go runtime, переименовывает `internal/` пути и экспортирует структуры. Даёт **идентичный парсинг runtime** (точность = Go compiler), поддерживает все endianness, ARM64/x86, byte scanning для модифицированного pclntab magic. Минус: **ручной merge upstream** при каждом релизе Go.

**Рекомендация для GoReveal: гибридный подход.** Используем GoReSym-style forked runtime для pclntab и type parsing (где точность критична), gore-style generated structs для moduledata (меняется чаще всего, проще обновлять).

### Ключевые Go runtime структуры

**pclntab** — PC line table, карта адресов→функций→файлов→строк. Magic values определяют версию формата: `0xFFFFFFFB` (Go 1.2–1.15), `0xFFFFFFFA` (Go 1.16–1.17), `0xFFFFFFF1` (Go 1.18–1.19), `0xFFFFFFF0` (Go 1.20+). В stripped бинарниках данные остаются в `.gopclntab` (ELF) / `__gopclntab` (Mach-O).

**moduledata** — layout исполняемого образа: text, data, bss, types, typelinks, functab, itablinks. В stripped ELF — в `.noptrdata`, в PE — требует brute-force поиск по адресу pclntab в data-секциях. **Меняется с каждой версией Go** — основная причина разной поддержки версий.

**buildinfo** (`.go.buildinfo`) — 32-байтная структура с magic `\xff Go buildinf:` (Go 1.18+), содержит module info и build settings.

### Memory-mapped file access для больших бинарников

**edsrzf/mmap-go** — cross-platform mmap (Linux/macOS/Windows), Read/Write/CoW/Exec modes. Benchmark: **25-38× быстрее** стандартного I/O при random access, **1.8 GB/s** throughput. Для GoReveal это критично: RE-анализ требует random access к разным секциям 100MB+ бинарника.

```go
f, _ := os.Open(path)
data, _ := mmap.Map(f, mmap.RDONLY, 0)
defer mmap.Unmap(data)
// Zero-copy доступ к любому offset бинарника
pclntab := data[section.Offset : section.Offset+section.Size]
```

GoReSym v3.2 решил OOM-проблемы через **streaming pclntab candidates** по каналам вместо хранения всех в памяти. GoReveal должен использовать тот же подход: mmap + streaming + bounded concurrency.

### Параллелизация анализа: pipeline pattern

```
Open Binary → Detect Version → Parse Pclntab → Find Moduledata
                                     ↓
                    ┌────────────────┼────────────────┐
                    ↓                ↓                ↓
            Extract Functions  Extract Types   Extract Strings
                    ↓                ↓                ↓
                    └────────────────┼────────────────┘
                                     ↓
                              Output Results
```

Секции анализируются параллельно через `errgroup.SetLimit(runtime.NumCPU())`. Каждый worker возвращает результаты через channel. Type parsing рекурсивен (типы ссылаются на другие типы), но independent typelink offsets можно распараллелить.

---

## 5. Унификация gore + GoReSym + redress

### Текущий ландшафт и лицензионные ограничения

| Проект | Роль | Лицензия | Stars | Ограничения |
|--------|------|----------|-------|-------------|
| **gore** | Библиотека | AGPL-3.0 | ~530 | Только little-endian, нет mmap |
| **redress** | CLI (обёртка gore) | AGPL-3.0 | ~1,140 | Зависит от gore |
| **GoReSym** | CLI + парсер | MIT | ~746 | Не библиотека, fork maintenance |

**AGPL-3.0 gore/redress несовместим** с closed-source IDA/Ghidra плагинами. GoReveal core **должен быть MIT** — это позволяет интеграцию с любыми consumers. Код GoReSym (MIT) может быть использован напрямую.

### Единый API core-библиотеки

```go
package goreveal

// Главный интерфейс для consumers (IDA, Ghidra, CLI)
type AnalysisResult struct {
    Metadata  Metadata    // Arch, OS, GoVersion, BuildID
    Functions []Function  // Name, EntryPC, EndPC, Package, SourceFile
    Types     []GoType    // Kind, Name, Fields, Methods
    Strings   []GoString  // Value, Address, Length
    Packages  []Package   // Name, Functions, Types, FilePaths
    Errors    []error     // Non-fatal parsing errors
}

// Простой API
result, err := goreveal.AnalyzeFile("binary", goreveal.DefaultOptions())

// Streaming API для больших бинарников
analyzer, err := goreveal.NewAnalyzer("binary", goreveal.WithMmap(true))
for fn := range analyzer.Functions() { /* ... */ }
for typ := range analyzer.Types() { /* ... */ }

// Consumer plugin interface
type Consumer interface {
    OnMetadata(Metadata) error
    OnFunction(Function) error
    OnType(GoType) error
    OnString(GoString) error
    Finalize() error
}
analyzer.AddConsumer(&IDAJSONConsumer{w: file})
analyzer.AddConsumer(&GhidraConsumer{conn: socket})
analyzer.Run(ctx)
```

### Стратегия рефакторинга

- **Файловый слой:** Единый `BinaryFile` interface с ELF/PE/Mach-O реализациями, backed by mmap (edsrzf/mmap-go). Lazy loading секций.
- **Parser слой:** GoReSym-based forked runtime code для pclntab/gosym/types (MIT-совместимо). Go-generated moduledata structs для быстрого обновления при новых версиях Go.
- **Analysis слой:** Concurrent pipeline с errgroup. Streaming results через iter.Seq2.
- **Output слой:** JSON (encoding/json), Protobuf (ConnectRPC), IDA Python script, Ghidra Java bridge.
- **Shared internal/:** Единый `internal/` пакет для byte-level utilities, version detection, error types, logging setup.
- **Единый подход:** cockroachdb/errors для error handling, slog для logging, functional options для конфигурации — во всех модулях монорепо.

---

## 6. Linting и CI/CD pipeline

### golangci-lint v2: конфигурация для GoReveal

golangci-lint **v2.11.3** (март 2026) — мажорная переработка. v1 end-of-life. Миграция: `golangci-lint migrate`.

```yaml
# .golangci.yml
version: "2"

formatters:
  enable:
    - goimports
    - golines
  settings:
    goimports:
      local-prefixes:
        - github.com/goreveal
    golines:
      max-len: 120

linters:
  default: standard
  enable:
    # Корректность
    - govet
    - staticcheck
    - errcheck
    - bodyclose
    - contextcheck
    - errorlint
    - nilerr
    - sqlclosecheck
    # Безопасность
    - gosec
    # Качество
    - revive
    - gocritic
    - unconvert
    - unparam
    - exhaustive
    - nolintlint
    # Modern Go (новые в v2)
    - modernize          # Предлагает современные Go-идиомы
    - exptostd           # x/exp → stdlib замены
    - intrange           # Integer range for-loops
    - fatcontext         # Nested contexts в циклах
    # Performance
    - perfsprint
    # slog
    - sloglint           # Консистентность log/slog
  settings:
    errcheck:
      check-type-assertions: true
    govet:
      enable-all: true
    staticcheck:
      checks: ["all"]
    exhaustive:
      default-signifies-exhaustive: true
  exclusions:
    presets:
      - comments
      - std-error-handling
      - common-false-positives
    rules:
      - path: _test\.go
        linters: [funlen, gocyclo, errcheck]
```

**Ключевые новые линтеры v2:** `modernize` (предлагает modern Go idioms), `exptostd` (заменяет x/exp на stdlib), `fatcontext` (nested contexts), `sloglint` (slog consistency). **NilAway (uber)** — опционально через module plugin; ещё produces false positives но ценен для nil-safety.

### Formatting: gofumpt через gopls

**gofumpt** — строгий superset gofmt, де-факто стандарт 2026. Vitess и многие проекты перешли на него. Интеграция через gopls:

```json
// VS Code settings.json
{ "gopls": { "formatting.gofumpt": true } }
```

### GitHub Actions: полный CI/CD pipeline

```yaml
name: CI
on:
  push: { branches: [main] }
  pull_request: { branches: [main] }

env:
  GO_VERSION: "1.26"
  GOLANGCI_LINT_VERSION: v2.11

jobs:
  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v5
      - uses: actions/setup-go@v6
        with: { go-version: "${{ env.GO_VERSION }}" }
      - uses: golangci/golangci-lint-action@v9
        with: { version: "${{ env.GOLANGCI_LINT_VERSION }}" }

  test:
    strategy:
      matrix:
        go: ["1.25.x", "1.26.x"]
        os: [ubuntu-latest, macos-latest, windows-latest]
    runs-on: ${{ matrix.os }}
    steps:
      - uses: actions/checkout@v5
      - uses: actions/setup-go@v6
        with: { go-version: "${{ matrix.go }}" }
      - run: go test -v -race -coverprofile=coverage.out ./...
      - uses: codecov/codecov-action@v4
        if: matrix.os == 'ubuntu-latest' && matrix.go == '1.26.x'

  security:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v5
      - uses: actions/setup-go@v6
        with: { go-version: "${{ env.GO_VERSION }}" }
      - run: |
          go install golang.org/x/vuln/cmd/govulncheck@latest
          govulncheck ./...
      - uses: aquasecurity/trivy-action@master
        with: { scan-type: fs, severity: "CRITICAL,HIGH" }

  release:
    if: startsWith(github.ref, 'refs/tags/v')
    needs: [lint, test, security]
    runs-on: ubuntu-latest
    permissions: { contents: write }
    steps:
      - uses: actions/checkout@v5
        with: { fetch-depth: 0 }
      - uses: actions/setup-go@v6
        with: { go-version: "${{ env.GO_VERSION }}" }
      - uses: goreleaser/goreleaser-action@v7
        with: { version: latest, args: "release --clean" }
        env: { GITHUB_TOKEN: "${{ secrets.GITHUB_TOKEN }}" }
```

**actions/setup-go v6** автоматически кэширует GOCACHE и GOMODCACHE по go.sum — отдельная настройка кэша не нужна. **govulncheck** использует static analysis — репортит только уязвимости в функциях, которые реально вызывает код (низкий noise). **Trivy** дополняет SBOM-генерацией (CycloneDX/SPDX).

### Pre-commit hooks и Makefile

```yaml
# .pre-commit-config.yaml
repos:
  - repo: https://github.com/golangci/golangci-lint
    rev: v2.11.3
    hooks:
      - id: golangci-lint
  - repo: https://github.com/tekwizely/pre-commit-golang
    rev: v0.8.3
    hooks:
      - id: go-fumpt
      - id: go-mod-tidy
```

```makefile
.PHONY: generate fmt lint test build

generate:
	buf generate
	sqlc generate
	go generate ./...

fmt:
	gofumpt -w .

lint:
	golangci-lint run --timeout=5m

test:
	go test -v -race -count=1 ./...

test-fuzz:
	go test -fuzz=FuzzParsePclntab -fuzztime=30s ./core/...

build:
	CGO_ENABLED=0 go build -o bin/goreveal ./cmd/goreveal

check: generate fmt lint test
	@echo "✅ All checks passed"
```

---

## Заключение: ключевые решения и что делать иначе

GoReveal строится на трёх фундаментальных решениях. **Первое — MIT-лицензированный core**, написанный с нуля на базе GoReSym-подхода (forked runtime для парсеров) и gore API-дизайна (чистый Go-пакет для consumers). AGPL gore не подходит для closed-source IDA/Ghidra плагинов. **Второе — zero-CGo стек**: modernc/sqlite, `CGO_ENABLED=0` builds, GoReleaser для всех платформ одной командой. Это радикально упрощает cross-compilation и deployment. **Третье — streaming-first архитектура**: mmap для файлового доступа, iter.Seq2 для API, channels для pipeline — решает OOM-проблемы, которые GoReSym исправлял в v3.2.

Минимальный набор внешних зависимостей за пределами stdlib: **pgx/v5** (PostgreSQL), **cobra** (CLI), **cockroachdb/errors** (stack traces), **connectrpc** (API), **testify** (тесты), **edsrzf/mmap-go** (memory-mapped I/O). Всё остальное — stdlib Go 1.26: slog, encoding/json, errors.Join, errgroup, debug/elf|pe|macho, iter. Этот стек обеспечивает баланс между минимальными зависимостями и production-ready функциональностью для серьёзного RE-инструмента.