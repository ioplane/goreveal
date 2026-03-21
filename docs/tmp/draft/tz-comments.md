### Go ↔ C/C++ и IDA SDK

Go через CGo действительно умеет вызывать C напрямую (`import "C"`), но с C++ ситуация сложнее — нужны `extern "C"` обёртки. Для IDA SDK есть три реалистичных варианта:

**Вариант A: Go → C shared library → IDA C++ plugin.** Компилируешь Core как `libgoreveal.so` через `-buildmode=c-shared`, экспортируя C API (`GoReveal_Open`, `GoReveal_Analyze`, `GoReveal_GetFunctionsJSON`). Пишешь тонкий IDA C++ plugin, который `dlopen`/`LoadLibrary` загружает эту библиотеку и вызывает её API, а результаты применяет через IDA SDK напрямую. Преимущество — in-process, без overhead subprocess'а, прямой доступ к IDA SDK C++ API (ida_typeinf, ida_funcs, Hex-Rays decompiler API). Недостаток — CGo усложняет cross-compilation, нужно собирать .so/.dll/.dylib под каждую платформу.

**Вариант B: Нативный C++ IDA plugin + Go CLI subprocess** (текущий вариант в ТЗ, но plugin на C++ вместо Python). C++ plugin вызывает Go CLI, парсит JSON, применяет результаты через SDK. Быстрее IDAPython при работе с типами и аннотациями — `tinfo_t` и `idc_parse_types` в C++ работают без Python-overhead. Проще в deployment — plugin.dll + goreveal бинарник.

**Вариант C: Гибрид — C++ plugin для performance-critical частей + IDAPython для rapid prototyping.** C++ plugin создаёт типы, применяет calling conventions, работает с Hex-Rays API. IDAPython-скрипты для интерактивных задач (пользовательские запросы, визуализация, ad-hoc анализ). Так работают серьёзные коммерческие плагины (например, HexRaysPyTools).

**Рекомендация**: Вариант A для MVP будет overengineering. Начни с **Вариант B** (C++ plugin + subprocess), затем мигрируй на **Вариант A** (in-process через c-shared) когда subprocess станет bottleneck. C++ plugin даёт полный доступ к Hex-Rays decompiler API (`hexrays_sdk`), чего IDAPython не может — например, трансформацию microcode для улучшения декомпиляции Go-кода.

### PostgreSQL как primary storage

Полностью согласен — для твоего стека это правильный выбор. Вот что меняется в ТЗ:

PostgreSQL становится **основным хранилищем**, SQLite остаётся как **portable fallback** (для аналитиков без PG). Конкретные преимущества PG для этого проекта: `pg_trgm` для фаззи-поиска по garble-обфусцированным именам, `JSONB` для расширяемых метаданных без миграций схемы (build settings, custom annotations), `pgvector` для семантического поиска в Фазе 4+, `LISTEN/NOTIFY` для real-time уведомлений IDA/Ghidra плагинов о новых результатах анализа, полноценный concurrent access (несколько аналитиков на одной базе), materialized views для предвычисленных сводок.

Обновлённая команда CLI:

```bash
# Primary: PostgreSQL
goreveal analyze --export-pg="postgres://user:pass@localhost/goreveal" /path/to/binary

# Fallback: SQLite
goreveal analyze --export-sqlite=analysis.db /path/to/binary
```
