# Техническое задание на разработку универсального инструмента реверс-инжиниринга Go-приложений

**Объединённый инструмент** (кодовое название — **GoREveal**) должен агрегировать возможности GoReSym, GoRE и redress в модульную платформу с интеграциями в IDA Pro и Ghidra, конвейером деобфускации garble и структурированным хранилищем результатов анализа. Проект закрывает ключевую рыночную нишу: на сегодня **не существует единого инструмента**, способного выполнить полный цикл от извлечения символов до псевдо-Go декомпиляции с деобфускацией — аналитики вынуждены вручную оркестрировать 5–7 разрозненных утилит. Ниже — полная спецификация, покрывающая архитектуру, интерфейсы модулей, фазы разработки, оценку трудозатрат и технические риски.

---

## 1. Анализ предметной области и конкурентная карта

Экосистема реверс-инжиниринга Go-бинарников фрагментирована. GoReSym (Mandiant, MIT) — самый быстрый парсер pclntab (1–5 секунд), работающий на форке исходников рантайма Go и выдающий JSON. GoRE (goretk, AGPL-3.0) — единственная библиотека с полным программным API для анализа Go-бинарников: пакеты, типы, методы, классификация main/stdlib/vendor. Redress (goretk, AGPL-3.0) — CLI поверх GoRE с уникальной возможностью «проекции исходного дерева» (source tree projection), восстанавливающей структуру каталогов и файлов исходного проекта по метаданным бинарника.

**Покрытие возможностей ключевых инструментов:**

| Возможность | GoReSym | GoRE | redress | GoResolver | GoStringUngarbler |
|---|:---:|:---:|:---:|:---:|:---:|
| Парсинг pclntab | ✅ | ✅ | ✅ | — | — |
| Восстановление типов | ✅ | ✅ | ✅ | — | — |
| Извлечение строк | ✅ (v3.2+) | — | — | — | ✅ (garble) |
| Классификация пакетов | Частичная | ✅ | ✅ | ✅ | — |
| Проекция исходного дерева | — | Данные | ✅ | — | — |
| JSON-выход | ✅ | N/A (библиотека) | — | ✅ | ✅ |
| CFG-деобфускация garble | — | — | — | ✅ | — |
| Строковая деобфускация | — | — | — | — | ✅ |
| API для программного доступа | Ограниченный | ✅ | — | — | — |

**Критические рыночные пробелы**, закрываемые проектом: отсутствие единого фреймворка, отсутствие queryable-хранилища результатов анализа (SQL), отсутствие интегрированного конвейера деобфускации, разрыв между инструментами извлечения и инструментами аннотации в IDA/Ghidra.

Коммерческие инструменты: IDA Pro 9.x имеет встроенный golang-плагин с поддержкой `__golang` calling convention (ABI0 и ABIInternal) и парсингом pclntab/typelinks, но не выполняет деобфускацию и не создаёт Go-специфичных высокоуровневых представлений. Ghidra с 10.3 включает `GolangSymbolAnalyzer` и `GoRttiMapper`, однако поддержка версий Go уже — до 1.23 против go1.6–1.25 у GolangAnalyzerExtension. JEB Decompiler имеет Python-скрипты для Go 1.13+, но без деобфускации. Binary Ninja поддерживает `.gopclntab` через community-плагины.

---

## 2. Архитектура системы

### 2.1. Общая модульная архитектура

Архитектура строится по принципу **«Core Library + CLI + тонкие плагины»** — паттерн, доказавший эффективность в GoReSym (Go CLI → JSON → IDAPython import) и валидированный в проектах rizin, angr и BinExport.

```
┌─────────────────────────────────────────────────────────────┐
│                     Клиентский уровень                       │
│  ┌──────────┐  ┌──────────────┐  ┌────────────────────────┐│
│  │  CLI/TUI  │  │ IDA Plugin   │  │ Ghidra Extension       ││
│  │  (Go)     │  │ (IDAPython)  │  │ (Java Analyzer)        ││
│  └─────┬─────┘  └──────┬───────┘  └───────────┬────────────┘│
│        │               │                       │             │
│        │     JSON/Protobuf через subprocess    │             │
└────────┼───────────────┼───────────────────────┼─────────────┘
         │               │                       │
┌────────▼───────────────▼───────────────────────▼─────────────┐
│                  Core Go Binary (goreveal)                     │
│  ┌────────────┐ ┌──────────┐ ┌──────────┐ ┌───────────────┐ │
│  │ pclntab    │ │moduledata│ │  types   │ │   strings     │ │
│  │ parser     │ │ parser   │ │ recovery │ │   extractor   │ │
│  ├────────────┤ ├──────────┤ ├──────────┤ ├───────────────┤ │
│  │ function   │ │ package  │ │  CFG     │ │  deobfuscation│ │
│  │ recovery   │ │ classify │ │ builder  │ │  pipeline     │ │
│  └────────────┘ └──────────┘ └──────────┘ └───────────────┘ │
├──────────────────────────────────────────────────────────────┤
│           Binary Format Layer (PE / ELF / Mach-O)            │
├──────────────────────────────────────────────────────────────┤
│     Storage Layer: SQLite + JSON + Protocol Buffers          │
└──────────────────────────────────────────────────────────────┘
```

**Ключевое архитектурное решение**: IDA-плагин и Ghidra-расширение общаются с Core через **вызов subprocess** (Go CLI → stdout JSON или запись в файл), а не через library linking. Это решает проблему лицензирования: если ядро под AGPL-3.0, тонкие плагины остаются MIT/Apache 2.0 как отдельные произведения, вызывающие внешний процесс. Именно так работает GoReSym сегодня.

### 2.2. Модули и их интерфейсы

#### Модуль 1: Core Library (`goreveal-core`)

**Язык**: Go 1.22+ (с прицелом на Go 1.26 при стабилизации).
**Зависимости**: GoRE library (AGPL-3.0), элементы GoReSym-подхода (форк runtime-парсеров, MIT-совместимая часть).

**Публичный Go API**:
```go
package goreveal

// Открытие и анализ бинарника
type Analyzer struct { ... }
func Open(path string, opts ...Option) (*Analyzer, error)
func (a *Analyzer) Close() error

// Извлечение данных
func (a *Analyzer) Functions() ([]*Function, error)
func (a *Analyzer) Types() ([]*GoType, error)
func (a *Analyzer) Packages() ([]*Package, error)
func (a *Analyzer) Strings() ([]*GoString, error)
func (a *Analyzer) Interfaces() ([]*Interface, error)
func (a *Analyzer) BuildInfo() (*BuildInfo, error)
func (a *Analyzer) SourceTree() (*SourceProjection, error)
func (a *Analyzer) ModuleData() (*ModuleData, error)

// Экспорт результатов
func (a *Analyzer) ExportJSON(w io.Writer) error
func (a *Analyzer) ExportProtobuf(w io.Writer) error
func (a *Analyzer) ExportSQLite(dbPath string) error

// Деобфускация (опциональный pipeline)
func (a *Analyzer) DetectObfuscation() (*ObfuscationInfo, error)
func (a *Analyzer) Deobfuscate(opts DeobfuscationOptions) (*DeobfuscationResult, error)
```

**Ключевые структуры данных**:
```go
type Function struct {
    Name        string   `json:"name"`
    StartAddr   uint64   `json:"start_addr"`
    EndAddr     uint64   `json:"end_addr"`
    PackageName string   `json:"package_name"`
    SourceFile  string   `json:"source_file"`
    StartLine   int      `json:"start_line"`
    EndLine     int      `json:"end_line"`
    Receiver    string   `json:"receiver,omitempty"`
    IsMethod    bool     `json:"is_method"`
    PackageType string   `json:"package_type"` // main|stdlib|vendor|unknown
}

type GoType struct {
    Name        string       `json:"name"`
    Addr        uint64       `json:"addr"`
    Kind        string       `json:"kind"` // struct|interface|func|map|slice|chan|...
    PackagePath string       `json:"package_path"`
    Size        int          `json:"size"`
    Fields      []TypeField  `json:"fields,omitempty"`
    Methods     []TypeMethod `json:"methods,omitempty"`
    CDecl       string       `json:"c_decl"` // C-совместимое объявление для IDA
}

type Package struct {
    Name           string `json:"name"`
    Path           string `json:"path"`
    Classification string `json:"classification"` // main|stdlib|vendor|unknown
    Version        string `json:"version,omitempty"`
    FunctionCount  int    `json:"function_count"`
}
```

#### Модуль 2: CLI (`goreveal`)

**Язык**: Go. **Фреймворк**: `cobra` (аналогично redress).

```bash
# Полный анализ с JSON-выходом
goreveal analyze --format=json --output=result.json /path/to/binary

# Только функции
goreveal functions --package-type=main /path/to/binary

# Типы с C-декларациями (для импорта в IDA)
goreveal types --c-decl /path/to/binary

# Проекция исходного дерева (аналог redress source)
goreveal source-tree /path/to/binary

# Деобфускация garble
goreveal deobfuscate --string-ungarble --cfg-resolve /path/to/binary

# Экспорт в SQLite
goreveal analyze --export-sqlite=analysis.db /path/to/binary

# Экспорт в PostgreSQL
goreveal analyze --export-pg="postgres://user:pass@host/db" /path/to/binary

# Режим сервера для IDA/Ghidra интеграции
goreveal serve --socket=/tmp/goreveal.sock
```

#### Модуль 3: IDA Pro Plugin (`goreveal-ida`)

**Язык**: Python 3.13+ (IDAPython). **Целевая версия IDA**: 9.0+.
**Архитектура**: Plugin (`plugin_t` subclass), не loader и не processor module. Все существующие Go RE инструменты для IDA — плагины (AlphaGolang, go_parser, IDAGolangHelper, GoReSym import script, GoFastAnalyzer).

**Обоснование выбора IDAPython**: 100% community Go RE инструментов используют IDAPython. Hex-Rays подтверждает: «IDAPython plugins are faster and easier to develop... and offer almost the same powerful capabilities as native C++ plugins». C++ SDK нужен только для processor/loader модулей, которые не требуются.

**Ключевые IDAPython модули**:

- **`ida_typeinf`** — центральный модуль в IDA 9.0+ (заменил удалённые `ida_struct` и `ida_enum`). `tinfo_t` — основной объект типа. `idc_parse_types()` для парсинга C-деклараций. `tinfo_t.add_udm()` для добавления полей структур
- **`ida_funcs`** — `add_func(start, end)` для создания функций, обнаруженных через pclntab
- **`ida_name`** — `set_name(ea, name, SN_NOWARN|SN_NOCHECK|SN_FORCE)` — паттерн GoReSym
- **`ida_bytes`** — `del_items()` для очистки неверных интерпретаций, `get_byte/dword/qword()` для чтения
- **`ida_dirtree`** — (IDA 7.6+) организация функций в директории по Go-пакетам (паттерн AlphaGolang)

**Критические изменения IDA 9.0**:
- Модули `ida_struct` и `ida_enum` **полностью удалены**; вся работа со структурами через `ida_typeinf`
- `get_inf_structure()` удалён → индивидуальные `inf_get_*()` / `inf_set_*()`
- Новая возможность **idalib** — запуск IDA как библиотеки в headless-режиме для пакетного анализа
- Рекомендуется `ida-plugin.json` для метаданных плагина

**Алгоритм работы плагина**:
1. При загрузке бинарника (`ev_newfile` hook) определить, является ли он Go-бинарником
2. Запустить `goreveal analyze --format=json` как subprocess
3. Распарсить JSON-результат
4. Определить Go-специфичные типы через `idc_parse_types()`:
   - `BUILTIN_STRING{char *ptr; size_t len;}`
   - `BUILTIN_INTERFACE{void *tab; void *data;}`
   - `go_slice{void *array; size_t len; size_t cap;}`
   - `complex64_t{float real; float imag;}`
5. Создать функции через `ida_funcs.add_func()` для адресов из pclntab, не обнаруженных IDA
6. Переименовать функции через `ida_name.set_name()` с флагами `SN_FORCE`
7. Импортировать типы через `idc_parse_types(typ['CDecl'], HTI_PAKDEF|HTI_DCL)`
8. Организовать функции в папки по пакетам через `ida_dirtree`
9. Применить Go calling convention: для Go 1.17+ используется register-based ABI (RAX, RBX, RCX, RDI, RSI, R8, R9, R10, R11). IDA 8.1+ поддерживает `__golang` CC нативно; для edge cases — `__usercall` с явным указанием регистров

**Передача данных CLI → Plugin**: JSON промежуточный формат (через subprocess + файл или stdout). Это battle-tested паттерн GoReSym. Для продвинутого использования — gRPC/Connect RPC через Unix domain socket (паттерн `ida-headless-mcp`).

#### Модуль 4: Ghidra Extension (`goreveal-ghidra`)

**Язык**: Java (Analyzer extension). **Целевая версия Ghidra**: 11.3+.

**Обоснование**: Java Analyzer — оптимальный подход, так как он интегрируется в auto-analysis pipeline Ghidra, имеет полный доступ к DataTypeManager API и переживает обновления Ghidra. GhidraScript (Java/Python) — для потребления JSON из внешнего CLI.

**Архитектура расширения**:
```
goreveal-ghidra/
  src/main/java/goreveal/
    GoRevealAnalyzer.java       // extends AbstractAnalyzer
    GoRevealImportScript.java   // extends GhidraScript
    types/GoTypeFactory.java    // создание Go-типов через DataTypeManager
    util/JsonParser.java        // парсинг JSON от CLI
  data/
    known_types.gdt             // предкомпилированные Go runtime типы
  extension.properties
  build.gradle
```

**GoRevealAnalyzer** (AnalyzerType: `FUNCTION_ANALYZER`, priority после встроенного `GolangSymbolAnalyzer`):
- `canAnalyze(Program)` — проверка Go build ID magic (`\xff Go buildinf:`) или pclntab magic (`0xf0ffffff` для Go 1.20+)
- `added()` — запуск `goreveal` CLI как subprocess, парсинг JSON, применение результатов
- Создание типов через `DataTypeManager` с `CategoryPath("/golang")`:
  - `StructureDataType` для GoString (ptr + len), GoSlice (array + len + cap), GoIface (tab + data)
  - `PointerDataType`, `ArrayDataType` для составных типов
  - Всё в транзакции: `dtm.openTransaction("GoReveal Import")`

**Headless-режим** для CI/CD:
```bash
analyzeHeadless /tmp/projects GoProject \
  -import /path/to/go_binary \
  -postScript GoRevealImport.java /path/to/analysis.json
```

**Совместимость с существующими инструментами**: Расширение дополняет GolangAnalyzerExtension (go1.6–1.25) и ghostrings (P-Code-based строковый анализ), не дублируя их. Приоритет анализа — после встроенного Go-анализатора Ghidra.

#### Модуль 5: Storage Layer

**Промежуточный формат**: Protocol Buffers (primary, по примеру BinExport Google) + JSON (secondary, human-readable).

**Protobuf-схема** (`goreveal.proto`):
```protobuf
syntax = "proto3";
package goreveal.v1;

message GoAnalysis {
  BinaryInfo binary_info = 1;
  repeated Function functions = 2;
  repeated GoType types = 3;
  repeated Package packages = 4;
  repeated GoString strings = 5;
  repeated Interface interfaces = 6;
  repeated CrossReference xrefs = 7;
  BuildInfo build_info = 8;
  SourceProjection source_tree = 9;
  ObfuscationInfo obfuscation = 10;
}
// (Полные определения сообщений — в приложении)
```

**SQLite-схема** для queryable-хранилища:

```sql
CREATE TABLE binaries (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  filename TEXT NOT NULL,
  sha256 TEXT UNIQUE NOT NULL,
  os TEXT, arch TEXT, endianness TEXT,
  go_version TEXT, format TEXT,
  build_id TEXT, go_root TEXT, main_root TEXT,
  analyzed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE functions (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  binary_id INTEGER REFERENCES binaries(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  start_addr INTEGER NOT NULL,
  end_addr INTEGER,
  package_id INTEGER REFERENCES packages(id),
  source_file TEXT, start_line INTEGER, end_line INTEGER,
  receiver TEXT, is_method BOOLEAN DEFAULT FALSE,
  UNIQUE(binary_id, start_addr)
);

CREATE TABLE types (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  binary_id INTEGER REFERENCES binaries(id) ON DELETE CASCADE,
  name TEXT NOT NULL, addr INTEGER,
  kind TEXT, package_path TEXT, size INTEGER,
  c_decl TEXT
);

CREATE TABLE type_fields (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  type_id INTEGER REFERENCES types(id) ON DELETE CASCADE,
  name TEXT, field_type TEXT,
  offset INTEGER, size INTEGER, tag TEXT
);

CREATE TABLE packages (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  binary_id INTEGER REFERENCES binaries(id) ON DELETE CASCADE,
  name TEXT NOT NULL, path TEXT,
  classification TEXT CHECK(classification IN ('main','stdlib','vendor','unknown')),
  version TEXT,
  UNIQUE(binary_id, path)
);

CREATE TABLE strings (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  binary_id INTEGER REFERENCES binaries(id) ON DELETE CASCADE,
  value TEXT, addr INTEGER, length INTEGER,
  referencing_func_id INTEGER REFERENCES functions(id)
);

CREATE TABLE cross_references (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  binary_id INTEGER REFERENCES binaries(id) ON DELETE CASCADE,
  from_addr INTEGER, to_addr INTEGER,
  xref_type TEXT CHECK(xref_type IN ('call','data','string','type'))
);
```

Индексы: `(binary_id, start_addr)` для функций, `(binary_id)` для типов, `(from_addr)` и `(to_addr)` для xref. **PostgreSQL**: аналогичная схема с `JSONB`-полями для расширяемых метаданных (`build_settings JSONB` в `binaries`).

---

## 3. Конвейер деобфускации garble

Garble — доминирующий обфускатор Go (5300+ звёзд на GitHub, BSD-3-Clause). Он применяет **деструктивные** трансформации (маглинг имён, удаление debug/DWARF, стирание BuildInfo) и **обратимые** (`-literals`: XOR/ADD/SUB шифрование строк с 5 вариантами; экспериментальное сплющивание control flow).

**Критические слабости garble**, эксплуатируемые деобфускаторами: pclntab и ModuleData **невозможно удалить** (Go runtime требует их для работы); обфусцированные имена **консистентны внутри пакета** (идентификация одной функции раскрывает весь пакет); `-literals` обратимы через эмуляцию, так как расшифрованная строка обязана существовать в рантайме.

### Интегрированный pipeline деобфускации

```
Вход: garble-обфусцированный Go-бинарник
        │
   ┌────▼────────────────────────┐
   │ СТАДИЯ 1: GoReSym-движок   │  ~1-5 сек
   │ Извлечение pclntab, типов, │
   │ обфусцированных имён,      │
   │ определение версии Go       │
   └────────────┬────────────────┘
                │
   ┌────────────▼────────────────┐
   │ СТАДИЯ 2: Строковая деобф.  │  ~1-10 мин
   │ Паттерн GoStringUngarbler:  │
   │ regex→prologue→Unicorn      │
   │ emulation→binary patching   │
   │ Выход: патченый бинарник    │
   │ + дамп расшифрованных строк │
   └────────────┬────────────────┘
                │
   ┌────────────▼────────────────┐
   │ СТАДИЯ 3: CFG-деобфускация  │  ~10-60 мин
   │ Паттерн GoResolver:         │
   │ fingerprint Go version →    │
   │ bootstrap reference →       │
   │ CFG similarity → name match │
   │ + package propagation       │
   └────────────┬────────────────┘
                │
   ┌────────────▼────────────────┐
   │ СТАДИЯ 4: Объединение       │
   │ Merge результатов стадий    │
   │ 1+2+3 в единый GoAnalysis   │
   │ JSON/Protobuf → import      │
   │ в IDA/Ghidra                │
   └─────────────────────────────┘
```

**Стратегия интеграции**: оркестрационный слой, а не форк. GoStringUngarbler (Apache 2.0) и GoResolver (Python, pip-installable) вызываются программно из Python-обёртки. GoReSym вызывается как subprocess. Все выходные данные — JSON, легко объединяемы. Это позволяет получать обновления из upstream без maintenance burden.

**Порядок стадий обоснован**: GoReSym первым (самый быстрый, даёт базовые метаданные); GoStringUngarbler вторым (патчит бинарник до загрузки в SRE-инструмент, чтобы строки были видны в декомпиляции); GoResolver третьим (самый тяжёлый, бенефитирует от частично деобфусцированного бинарника). Импорт в IDA/Ghidra — последним.

---

## 4. Лицензионная стратегия

### Анализ лицензий компонентов

| Компонент | Лицензия | Влияние |
|---|---|---|
| GoRE (goretk/gore) | **AGPL-3.0** | Вирусная — при `import` в Go static linking делает весь бинарник AGPL |
| redress (goretk) | **AGPL-3.0** | Аналогично GoRE |
| GoReSym (Mandiant) | **MIT** | Максимально пермиссивная |
| GoStringUngarbler (Mandiant) | **Apache 2.0** | Пермиссивная, patent grant |
| GoResolver (Volexity) | Требует уточнения | Вероятно пермиссивная |
| IDA SDK | **MIT** (недавно открыт) | Плагины распространяются свободно |
| Ghidra | **Apache 2.0** | Расширения без ограничений |

### Рекомендуемая стратегия: модульное лицензирование (Вариант B)

**Ключевой прецедент**: Mandiant при создании GoReSym **намеренно отказался от зависимости на GoRE** именно из-за AGPL, вместо этого форкнув BSD-лицензированные исходники рантайма Go. Это проверенный путь к максимальной лицензионной гибкости.

**Рекомендуемая архитектура лицензий**:

- **Core Go Binary** (`goreveal`): **AGPL-3.0** (при использовании GoRE как зависимости). Распространяется как отдельный исполняемый файл с полным исходным кодом
- **IDA Plugin** (`goreveal-ida`): **Apache 2.0** — отдельное произведение, вызывающее Core через subprocess. Не содержит AGPL-кода, только парсинг JSON и IDAPython API calls
- **Ghidra Extension** (`goreveal-ghidra`): **Apache 2.0** — аналогично, отдельное произведение
- **Protobuf-определения и JSON Schema**: **Apache 2.0** — по модели Zitadel, чтобы клиентский код не наследовал AGPL

**Альтернативный вариант** (если AGPL неприемлем для организации): построить Core на основе GoReSym (MIT) + собственная реализация недостающей функциональности GoRE (классификация пакетов, source tree projection, методы). Это **значительно увеличивает трудозатраты** (оценка: +3–4 месяца), но позволяет лицензировать весь проект под MIT/Apache 2.0.

**Для корпоративных пользователей**: AGPL **не ограничивает внутреннее использование**. Команда безопасности может свободно использовать AGPL-инструмент для внутреннего анализа без обязательств по раскрытию. Проблема возникает только при распространении или предоставлении как сервиса. Google запрещает AGPL, но большинство security-команд более лояльны к AGPL, так как RE-инструменты — end-user приложения для внутреннего использования.

---

## 5. Технический стек и инфраструктура

### Языки и зависимости

| Компонент | Язык | Ключевые зависимости |
|---|---|---|
| Core Library + CLI | Go 1.22+ | `goretk/gore`, `spf13/cobra`, `protocolbuffers/protobuf-go`, `mattn/go-sqlite3`, `lib/pq` |
| IDA Plugin | Python 3.13+ | IDAPython (ida_funcs, ida_name, ida_typeinf, ida_bytes, ida_dirtree) |
| Ghidra Extension | Java 21+ | Ghidra API (AbstractAnalyzer, DataTypeManager, StructureDataType) |
| Deobfuscation wrapper | Python 3.12+ | GoStringUngarbler, GoResolver (pip), unicorn-engine |

### Build system и CI/CD

- **Go**: стандартный `go build`, `go test`, `goreleaser` для cross-platform релизов (PE/ELF/Mach-O × amd64/arm64)
- **IDA Plugin**: `pip install` с `setup.py`/`pyproject.toml`, установка в `<IDADIR>/plugins/`
- **Ghidra Extension**: Gradle build (`gradle -PGHIDRA_INSTALL_DIR=...`), выход — ZIP для `File → Install Extensions`
- **CI/CD**: GitHub Actions с матрицей: Go 1.22/1.23/1.24, IDA 9.0/9.1+ (mock tests), Ghidra 11.3+ (headless tests)
- **Тестирование**: набор Go-бинарников разных версий (go1.16–1.26), архитектур (amd64, arm64), обфускаторов (garble v0.11–0.13), форматов (PE, ELF, Mach-O) — **~50 тестовых сэмплов**

### CGo-интеграция

CGo потребуется для:
- Интеграции с Unicorn Engine (C-библиотека) для встроенной эмуляции строковой деобфускации (альтернатива: вызов Python GoStringUngarbler через subprocess)
- Потенциально — для экспорта Core Library как C shared library (`-buildmode=c-shared`) для прямой загрузки в non-Go плагины

**Риск**: CGo усложняет cross-compilation и увеличивает время сборки. Рекомендация — использовать CGo только при доказанной необходимости, по умолчанию предпочитая subprocess-архитектуру.

---

## 6. План разработки по фазам

### Фаза 1: Foundation (8 недель)

| Задача | Оценка | Результат |
|---|---|---|
| Инициализация проекта, CI/CD, тестовый набор бинарников | 1 нед | Go module, GitHub Actions, 20+ тестовых PE/ELF/Mach-O |
| Core Library: интеграция GoRE, парсинг pclntab, функции, пакеты | 3 нед | `Open()`, `Functions()`, `Packages()`, `Types()` |
| Core Library: интеграция GoReSym-подхода для строк и типов | 2 нед | `Strings()`, `Types()` с C-декларациями |
| CLI: базовые команды (analyze, functions, types, source-tree) | 1 нед | `goreveal analyze --format=json` |
| JSON Schema + Protobuf определения | 1 нед | `goreveal.proto`, JSON Schema v1 |

**Milestone 1**: CLI-инструмент, анализирующий Go-бинарники go1.16–1.26, выдающий JSON с функциями, типами, пакетами, строками, проекцией исходного дерева.

### Фаза 2: Storage + IDA Integration (6 недель)

| Задача | Оценка | Результат |
|---|---|---|
| SQLite экспорт с полной схемой | 1.5 нед | `goreveal analyze --export-sqlite` |
| PostgreSQL экспорт | 1 нед | `goreveal analyze --export-pg` |
| IDA Plugin: базовая версия (subprocess → JSON → import) | 2 нед | Импорт функций, имён, базовых типов |
| IDA Plugin: Go-типы (ida_typeinf), calling convention, ida_dirtree | 1.5 нед | Полная аннотация Go-бинарника в IDA 9.0+ |

**Milestone 2**: Полная цепочка CLI → SQLite/PostgreSQL + CLI → IDA Pro аннотация.

### Фаза 3: Ghidra + Deobfuscation (8 недель)

| Задача | Оценка | Результат |
|---|---|---|
| Ghidra Analyzer extension (Java) | 3 нед | Auto-analysis при обнаружении Go-бинарника |
| Ghidra import script + Go DataTypes | 1.5 нед | Полная аннотация в Ghidra 11.3+ |
| Конвейер деобфускации: интеграция GoStringUngarbler | 2 нед | `goreveal deobfuscate --string-ungarble` |
| Конвейер деобфускации: интеграция GoResolver (CFG) | 1.5 нед | `goreveal deobfuscate --cfg-resolve` |

**Milestone 3**: Полнофункциональный инструмент с IDA + Ghidra + деобфускация garble.

### Фаза 4: Polish + Advanced Features (4 недели)

| Задача | Оценка | Результат |
|---|---|---|
| gRPC/Unix socket сервер для live-интеграции | 1.5 нед | `goreveal serve` |
| Cross-reference engine | 1 нед | Таблица xref в SQLite/PG |
| Документация, примеры, тесты на 50+ бинарников | 1.5 нед | README, API docs, test suite |

### Суммарная оценка трудозатрат

| Фаза | Длительность | Трудозатраты (человеко-недели) |
|---|---|---|
| 1. Foundation | 8 нед | 8 |
| 2. Storage + IDA | 6 нед | 6 |
| 3. Ghidra + Deobfuscation | 8 нед | 8 |
| 4. Polish | 4 нед | 4 |
| **Итого** | **~26 недель** | **~26 чел-нед** |

При одном full-time разработчике — **~6.5 месяцев** до полной функциональности. При двух — **~4 месяца** с параллелизацией Фаз 2 и 3.

---

## 7. Технические риски и митигация

### Высокий приоритет

**Риск 1: Нестабильность Go runtime-структур между версиями**. pclntab менял формат в go1.2, 1.16, 1.18, 1.20. `runtime._type` переехал в `abi` пакет в go1.21. Каждый major-релиз Go потенциально ломает парсеры.
*Митигация*: Копировать подход GoReSym — форкать исходники runtime для каждой поддерживаемой версии. Автотесты на матрице Go-версий. CI-мониторинг новых релизов Go.

**Риск 2: AGPL-заражение при неправильных границах модулей**. Если IDA/Ghidra плагины случайно импортируют AGPL-код (через shared library или copy-paste), они становятся AGPL.
*Митигация*: Строгая subprocess-граница. Плагины **никогда** не импортируют Go-код напрямую — только потребляют JSON через файл или stdout. Code review с AGPL-checklist.

**Риск 3: IDA 9.x API breaking changes**. Удаление `ida_struct`, `ida_enum` и 20+ функций ломает совместимость с плагинами для IDA 7/8.
*Митигация*: Таргетирование исключительно IDA 9.0+. Использование `ida_typeinf` для всех типовых операций. Не поддерживать IDA 7/8 — нет смысла при быстрой миграции сообщества.

### Средний приоритет

**Риск 4: GoResolver — неизвестная лицензия и вычислительная тяжесть**. CFG similarity на больших бинарниках может занимать часы. Лицензия репозитория требует уточнения.
*Митигация*: Предварительно проверить LICENSE в репозитории. Реализовать timeout и параллелизацию CFG-сравнения. Предложить режим «быстрый» (только stdlib-matching) и «полный».

**Риск 5: garble эволюционирует**. Экспериментальный control flow flattening (`GARBLE_EXPERIMENTAL_CONTROLFLOW=1`) не обрабатывается ни одним существующим инструментом.
*Митигация*: Мониторинг garble releases. Модульная архитектура деобфускации позволяет добавлять новые стадии. В Фазе 4+ — исследование deflattening на основе symbolic execution (SENinja/angr подходы).

**Риск 6: CGo cross-compilation**. Если потребуется CGo для Unicorn Engine, cross-compilation для macOS/Windows/Linux × amd64/arm64 усложняется значительно.
*Митигация*: По умолчанию — subprocess-вызов Python GoStringUngarbler. CGo только для оптимизированной встроенной эмуляции в поздних фазах.

### Низкий приоритет

**Риск 7: Ghidra API backward compatibility**. Ghidra не гарантирует стабильность Java API между major-версиями.
*Митигация*: Таргетирование Ghidra 11.3+. Ghidra scripts (не compiled analyzers) более устойчивы к обновлениям — по опыту Trellix. При необходимости — dual mode: analyzer + fallback script.

**Риск 8: Объём тестового набора**. 50+ бинарников для тестирования — компиляция под разные Go-версии, архитектуры, обфускаторы.
*Митигация*: Использовать GoStrap (Volexity) для автоматической генерации reference-бинарников. CI pipeline с матрицей Go-версий.

---

## Заключение: уникальное ценностное предложение и архитектурные принципы

GoReveal становится **первым** инструментом, объединяющим полный цикл анализа Go-бинарника в единую платформу: от pclntab-парсинга до деобфускации garble и импорта аннотаций в IDA/Ghidra. Три архитектурных принципа обеспечивают успех проекта. Первый — **subprocess-граница** между AGPL-ядром и пермиссивно-лицензированными плагинами, решающая лицензионную дилемму без компромиссов в функциональности. Второй — **оркестрация, а не реимплементация** деобфускационных инструментов (GoStringUngarbler, GoResolver), позволяющая получать upstream-обновления. Третий — **queryable storage** (SQLite/PostgreSQL), отсутствующий у всех существующих инструментов и открывающий возможности для аналитики: корреляция пакетов между сэмплами, поиск по типам и строкам, кластеризация малвари по сигнатурам Go-зависимостей. Именно этот третий компонент — структурированная база анализа — трансформирует разрозненные CLI-утилиты в платформу для масштабного исследования Go-угроз.