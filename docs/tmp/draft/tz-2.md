# GoReveal — Техническое задание на разработку универсального инструмента реверс-инжиниринга Go-приложений

**Версия документа**: 2.0
**Дата**: 19 марта 2026
**Статус**: Для разработки

---

## 1. Введение и цель проекта

### 1.1. Назначение

GoReveal — модульная платформа для автоматизированного реверс-инжиниринга приложений, скомпилированных на языке Go. Инструмент объединяет функциональность трёх open-source проектов: GoReSym (Mandiant), GoRE library (goretk) и redress (goretk) в единый продукт с интеграциями IDA Pro (C++ plugin), Ghidra (Java extension), конвейером деобфускации garble и структурированным хранилищем результатов в PostgreSQL.

### 1.2. Проблема

На сегодня **не существует единого инструмента**, способного выполнить полный цикл от извлечения символов до псевдо-Go декомпиляции с деобфускацией. Аналитики вынуждены вручную оркестрировать 5–7 разрозненных утилит (GoReSym для символов, GoRE для типов, redress для source tree projection, GoResolver для деобфускации CFG, GoStringUngarbler для строк, отдельные скрипты импорта для IDA и Ghidra), теряя время и контекст между инструментами.

### 1.3. Целевые пользователи

- Malware-аналитики (incident response, threat intelligence)
- Security-исследователи и пентестеры
- Разработчики, анализирующие legacy Go-код без исходников

---

## 2. Конкурентная карта и рыночные пробелы

### 2.1. Покрытие возможностей существующих инструментов

| Возможность | GoReSym | GoRE | redress | GoResolver | GoStringUngarbler | **GoReveal** |
|---|:---:|:---:|:---:|:---:|:---:|:---:|
| Парсинг pclntab (Go 1.2–1.26) | ✅ | ✅ | ✅ | — | — | ✅ |
| Восстановление типов/структур | ✅ | ✅ | ✅ | — | — | ✅ |
| Извлечение строк | ✅ (v3.2+) | — | — | — | ✅ (garble) | ✅ |
| Классификация пакетов (main/std/vendor) | Частичная | ✅ | ✅ | ✅ | — | ✅ |
| Проекция исходного дерева | — | Данные | ✅ | — | — | ✅ |
| CFG-деобфускация garble | — | — | — | ✅ | — | ✅ |
| Строковая деобфускация garble | — | — | — | — | ✅ | ✅ |
| C++ IDA-плагин (нативный SDK) | — | — | — | Только Python | — | ✅ |
| Ghidra extension | — | — | — | Python script | — | ✅ |
| Queryable storage (PostgreSQL) | — | — | — | — | — | ✅ |
| Генерация отчётов | — | — | — | — | — | ✅ |
| Единый CLI | — | — | ✅ (ограниченный) | — | — | ✅ |

### 2.2. Коммерческие инструменты

IDA Pro 9.x имеет встроенный golang-плагин с поддержкой `__golang` calling convention (ABI0 и ABIInternal) и парсингом pclntab/typelinks, но не выполняет деобфускацию, не создаёт высокоуровневых представлений и не экспортирует в queryable storage. Ghidra 11.3+ включает `GolangSymbolAnalyzer` и `GoRttiMapper` — покрытие до Go 1.23, ограниченная поддержка типов. JEB Decompiler имеет Python-скрипты для Go 1.13+, без деобфускации. Binary Ninja — community-плагины.

### 2.3. Уникальное ценностное предложение GoReveal

1. **Единый конвейер**: от pclntab до аннотированного бинарника в IDA/Ghidra за одну команду
2. **Нативный C++ IDA-плагин**: in-process интеграция через IDA SDK, включая Hex-Rays decompiler API — в отличие от всех существующих Go RE-инструментов, использующих IDAPython
3. **PostgreSQL-хранилище**: SQL-запросы по результатам анализа, корреляция между бинарниками, pg_trgm для фаззи-матчинга обфусцированных имён, pgvector для семантического поиска
4. **Интегрированная деобфускация garble**: строки (GoStringUngarbler) + имена (GoResolver) + объединение результатов

---

## 3. Архитектура системы

### 3.1. Общая модульная архитектура

```
┌─────────────────────────────────────────────────────────────────┐
│                      Клиентский уровень                         │
│  ┌──────────┐  ┌────────────────┐  ┌────────────────────────┐  │
│  │  CLI/TUI  │  │ IDA C++ Plugin │  │ Ghidra Java Extension  │  │
│  │  (Go)     │  │ (IDA SDK)      │  │ (AbstractAnalyzer)     │  │
│  └─────┬─────┘  └───────┬────────┘  └───────────┬────────────┘  │
│        │                │                        │               │
│        │ JSON/Protobuf через subprocess / c-shared              │
└────────┼────────────────┼────────────────────────┼───────────────┘
         │                │                        │
┌────────▼────────────────▼────────────────────────▼───────────────┐
│                   Core Go Binary (goreveal)                       │
│  ┌────────────┐ ┌──────────┐ ┌──────────┐ ┌───────────────────┐ │
│  │ pclntab    │ │moduledata│ │  types   │ │   strings         │ │
│  │ parser     │ │ parser   │ │ recovery │ │   extractor       │ │
│  ├────────────┤ ├──────────┤ ├──────────┤ ├───────────────────┤ │
│  │ function   │ │ package  │ │  CFG     │ │  deobfuscation    │ │
│  │ recovery   │ │ classify │ │ builder  │ │  pipeline         │ │
│  ├────────────┤ ├──────────┤ ├──────────┤ ├───────────────────┤ │
│  │ source tree│ │ pseudo-Go│ │ xref     │ │  report           │ │
│  │ projection │ │ renderer │ │ engine   │ │  generator        │ │
│  └────────────┘ └──────────┘ └──────────┘ └───────────────────┘ │
├──────────────────────────────────────────────────────────────────┤
│            Binary Format Layer (PE / ELF / Mach-O)               │
├──────────────────────────────────────────────────────────────────┤
│   Storage: PostgreSQL (primary) │ SQLite (portable) │ JSON/Proto │
└──────────────────────────────────────────────────────────────────┘
```

### 3.2. Ключевые архитектурные решения

**AD-1: Subprocess-граница между Core и плагинами (MVP)**

IDA C++ Plugin и Ghidra Extension общаются с Core через вызов subprocess (`goreveal` CLI → stdout JSON или запись в файл). Это решает проблему лицензирования: Core под AGPL-3.0 (из-за GoRE), а плагины — Apache 2.0 как отдельные произведения, вызывающие внешний процесс.

**AD-2: Путь миграции к in-process интеграции**

Core компилируется как C shared library через `-buildmode=c-shared`, экспортируя C API. IDA C++ Plugin может `dlopen`/`LoadLibrary` эту библиотеку для in-process вызовов, устраняя subprocess overhead. Переход осуществляется после стабилизации Core API (Фаза 4+).

**AD-3: PostgreSQL — primary storage**

PostgreSQL выбран как основное хранилище по причинам: pg_trgm для фаззи-матчинга garble-обфусцированных имён, JSONB для расширяемых метаданных, pgvector для семантического поиска (Фаза 4+), concurrent access для командной работы, materialized views для предвычисленных аналитик, LISTEN/NOTIFY для уведомления плагинов о новых результатах. SQLite — portable fallback для автономной работы.

**AD-4: Протокол обмена данными**

Protocol Buffers (primary, по модели BinExport Google) + JSON (secondary, human-readable). Protobuf обеспечивает типобезопасность, обратную совместимость и компактность. JSON — для ad-hoc анализа, отладки и интеграции с инструментами не поддерживающими Protobuf.

---

## 4. Спецификация модулей

### 4.1. Модуль 1: Core Library (`goreveal-core`)

**Язык**: Go 1.26
**Лицензия**: AGPL-3.0 (из-за зависимости от GoRE)
**Ключевые зависимости**: `goretk/gore` (AGPL-3.0), `spf13/cobra`, `protocolbuffers/protobuf-go`, `lib/pq`, `mattn/go-sqlite3`

#### 4.1.1. Публичный Go API

```go
package goreveal

import "context"

// Analyzer — центральный объект анализа
type Analyzer struct { /* private fields */ }

// Open открывает бинарник для анализа
func Open(path string, opts ...Option) (*Analyzer, error)
func (a *Analyzer) Close() error

// Версия Go и метаданные
func (a *Analyzer) GoVersion() string
func (a *Analyzer) BuildInfo() (*BuildInfo, error)
func (a *Analyzer) ModuleData() (*ModuleData, error)

// Извлечение данных
func (a *Analyzer) Functions() ([]*Function, error)
func (a *Analyzer) Types() ([]*GoType, error)
func (a *Analyzer) Packages() ([]*Package, error)
func (a *Analyzer) Strings() ([]*GoString, error)
func (a *Analyzer) Interfaces() ([]*Interface, error)
func (a *Analyzer) SourceTree() (*SourceProjection, error)
func (a *Analyzer) CrossReferences() ([]*CrossReference, error)

// Деобфускация
func (a *Analyzer) DetectObfuscation() (*ObfuscationInfo, error)
func (a *Analyzer) Deobfuscate(ctx context.Context, opts DeobfuscationOptions) (*DeobfuscationResult, error)

// Экспорт
func (a *Analyzer) ExportJSON(w io.Writer, opts ExportOptions) error
func (a *Analyzer) ExportProtobuf(w io.Writer) error
func (a *Analyzer) ExportPostgreSQL(connStr string) error
func (a *Analyzer) ExportSQLite(dbPath string) error

// Генерация отчётов
func (a *Analyzer) GenerateReport(w io.Writer, format ReportFormat) error
```

#### 4.1.2. Ключевые структуры данных

```go
type Function struct {
    Name        string `json:"name"`
    StartAddr   uint64 `json:"start_addr"`
    EndAddr     uint64 `json:"end_addr"`
    PackageName string `json:"package_name"`
    SourceFile  string `json:"source_file"`
    StartLine   int    `json:"start_line"`
    EndLine     int    `json:"end_line"`
    Receiver    string `json:"receiver,omitempty"`
    IsMethod    bool   `json:"is_method"`
    PackageType string `json:"package_type"` // main|stdlib|vendor|unknown
    ArgSize     int    `json:"arg_size"`
    // C-совместимая сигнатура для IDA
    CSignature  string `json:"c_signature,omitempty"`
}

type GoType struct {
    Name          string       `json:"name"`
    Addr          uint64       `json:"addr"`
    Kind          string       `json:"kind"` // struct|interface|func|map|slice|chan|ptr|...
    KindID        uint8        `json:"kind_id"`
    PackagePath   string       `json:"package_path"`
    Size          int          `json:"size"`
    Alignment     int          `json:"alignment"`
    Fields        []TypeField  `json:"fields,omitempty"`
    Methods       []TypeMethod `json:"methods,omitempty"`
    // C-совместимая декларация для IDA ida_typeinf
    CDecl         string       `json:"c_decl"`
    // IDA TIL-совместимое определение
    IDATypeInfo   string       `json:"ida_type_info,omitempty"`
}

type TypeField struct {
    Name       string `json:"name"`
    TypeName   string `json:"type_name"`
    Offset     int    `json:"offset"`
    Size       int    `json:"size"`
    Tag        string `json:"tag,omitempty"`
    Anonymous  bool   `json:"anonymous"`
}

type Package struct {
    Name           string `json:"name"`
    Path           string `json:"path"`
    Classification string `json:"classification"` // main|stdlib|vendor|unknown
    Version        string `json:"version,omitempty"`
    FunctionCount  int    `json:"function_count"`
    TypeCount      int    `json:"type_count"`
}

type GoString struct {
    Value           string `json:"value"`
    Addr            uint64 `json:"addr"`
    Length          int    `json:"length"`
    ReferencingFunc string `json:"referencing_func,omitempty"`
    IsDeobfuscated  bool   `json:"is_deobfuscated"`
}

type CrossReference struct {
    FromAddr uint64 `json:"from_addr"`
    ToAddr   uint64 `json:"to_addr"`
    Type     string `json:"type"` // call|data|string|type
    FromFunc string `json:"from_func,omitempty"`
    ToFunc   string `json:"to_func,omitempty"`
}

type BuildInfo struct {
    GoVersion    string            `json:"go_version"`
    ModulePath   string            `json:"module_path"`
    Dependencies []Dependency      `json:"dependencies"`
    Settings     map[string]string `json:"settings"`
    BuildID      string            `json:"build_id"`
    GoRoot       string            `json:"go_root"`
    MainRoot     string            `json:"main_root"`
}

type ObfuscationInfo struct {
    IsObfuscated     bool   `json:"is_obfuscated"`
    ObfuscatorName   string `json:"obfuscator_name"` // garble|gobfuscate|unknown
    ObfuscatorVersion string `json:"obfuscator_version,omitempty"`
    HasLiterals      bool   `json:"has_obfuscated_literals"`
    HasControlFlow   bool   `json:"has_control_flow_flattening"`
    Confidence       float64 `json:"confidence"` // 0.0–1.0
}

type DeobfuscationOptions struct {
    StringUngarble  bool   `json:"string_ungarble"`
    CFGResolve      bool   `json:"cfg_resolve"`
    Timeout         time.Duration
    GoResolverPath  string // путь к goresolver binary
    UnicornEmulate  bool   // использовать Unicorn для эмуляции
}

type SourceProjection struct {
    MainPackage string              `json:"main_package"`
    Directories []SourceDirectory   `json:"directories"`
}

type SourceDirectory struct {
    Path  string       `json:"path"`
    Files []SourceFile `json:"files"`
}

type SourceFile struct {
    Path      string   `json:"path"`
    Functions []string `json:"functions"`
    LineRange [2]int   `json:"line_range"` // [start, end]
}

type ReportFormat int
const (
    ReportMarkdown ReportFormat = iota
    ReportHTML
    ReportJSON
)
```

#### 4.1.3. C Shared Library API (для in-process IDA интеграции, Фаза 4+)

```go
package main

import "C"

//export GoReveal_Open
func GoReveal_Open(path *C.char) C.int { /* returns handle */ }

//export GoReveal_Analyze
func GoReveal_Analyze(handle C.int) *C.char { /* returns JSON */ }

//export GoReveal_GetFunctions
func GoReveal_GetFunctions(handle C.int) *C.char { /* returns JSON array */ }

//export GoReveal_GetTypes
func GoReveal_GetTypes(handle C.int) *C.char { /* returns JSON array */ }

//export GoReveal_Free
func GoReveal_Free(ptr *C.char) { /* frees C string */ }

//export GoReveal_Close
func GoReveal_Close(handle C.int) { /* closes analyzer */ }
```

Сборка: `CGO_ENABLED=1 go build -buildmode=c-shared -o libgoreveal.so ./cmd/libgoreveal`

### 4.2. Модуль 2: CLI (`goreveal`)

**Язык**: Go 1.26
**Фреймворк**: `spf13/cobra`

#### 4.2.1. Команды

```bash
# === Основной анализ ===

# Полный анализ с JSON-выходом
goreveal analyze /path/to/binary
goreveal analyze --format=json --output=result.json /path/to/binary
goreveal analyze --format=proto --output=result.pb /path/to/binary

# Экспорт в PostgreSQL (primary storage)
goreveal analyze --export-pg="postgres://user:pass@localhost:5432/goreveal" /path/to/binary

# Экспорт в SQLite (portable fallback)
goreveal analyze --export-sqlite=analysis.db /path/to/binary

# === Секционные команды ===

# Функции с фильтрацией
goreveal functions /path/to/binary
goreveal functions --package-type=main /path/to/binary
goreveal functions --package=crypto/tls /path/to/binary

# Типы с C-декларациями для IDA
goreveal types /path/to/binary
goreveal types --c-decl --ida-compat /path/to/binary

# Строки
goreveal strings /path/to/binary

# Проекция исходного дерева
goreveal source-tree /path/to/binary

# Информация о сборке
goreveal info /path/to/binary

# === Деобфускация ===

# Полный конвейер
goreveal deobfuscate /path/to/binary
goreveal deobfuscate --string-ungarble --cfg-resolve /path/to/binary

# Только строки
goreveal deobfuscate --string-ungarble /path/to/binary

# Только CFG similarity
goreveal deobfuscate --cfg-resolve /path/to/binary

# === Отчёты ===

goreveal report --format=markdown --output=report.md /path/to/binary
goreveal report --format=html --output=report.html /path/to/binary

# === Сервер для live-интеграции ===

goreveal serve --socket=/tmp/goreveal.sock
goreveal serve --port=50051  # gRPC
```

#### 4.2.2. Глобальные флаги

```
--verbose, -v         Подробный вывод
--quiet, -q           Только ошибки
--format              Формат вывода: json|proto|text|table (default: json)
--output, -o          Файл вывода (default: stdout)
--export-pg           PostgreSQL connection string
--export-sqlite       Путь к SQLite файлу
--timeout             Таймаут анализа (default: 5m)
--workers             Количество параллельных worker'ов (default: runtime.NumCPU())
```

### 4.3. Модуль 3: IDA Pro C++ Plugin (`goreveal-ida`)

**Язык**: C++17
**Целевая версия IDA**: 9.0+
**IDA SDK**: ida-sdk-90 (C++ headers + libraries)
**Лицензия**: Apache 2.0 (отдельное произведение, subprocess-граница)

#### 4.3.1. Обоснование выбора C++ вместо IDAPython

| Критерий | C++ IDA SDK | IDAPython |
|---|---|---|
| **Производительность** | Нативная, без GIL | Python GIL, overhead интерпретатора |
| **Hex-Rays API** | Полный доступ (microcode, ctree) | Частичный (без microcode transforms) |
| **Типовая система** | Прямая работа с `tinfo_t` через C++ API | Обёртка, менее гибкая |
| **Calling convention** | Полный контроль через `set_cc()` | Ограниченный |
| **Deployment** | Один .dll/.so файл | Директория Python-файлов |
| **Отладка** | C++ debugger, valgrind | Python debugger |
| **Сообщество Go RE** | 0% существующих инструментов | 100% (AlphaGolang, go_parser, etc.) |
| **Время разработки** | 2–3× дольше | Быстрее |

Выбор C++ обоснован двумя killer-features: полный доступ к **Hex-Rays decompiler API** (microcode transforms для улучшения декомпиляции Go-кода) и **нативная производительность** при обработке бинарников с 50 000+ функциями.

#### 4.3.2. Архитектура плагина

```
goreveal-ida/
├── src/
│   ├── plugin_main.cpp          // ida_plugin_entry, PLUGIN struct
│   ├── goreveal_plugin.h        // GoRevealPlugin : public plugmod_t
│   ├── goreveal_plugin.cpp      // Реализация основного класса
│   ├── core_bridge.h            // Интерфейс к Go Core (subprocess/dlopen)
│   ├── core_bridge_subprocess.cpp  // MVP: вызов goreveal CLI
│   ├── core_bridge_inprocess.cpp   // Фаза 4+: dlopen libgoreveal.so
│   ├── function_annotator.cpp   // Создание функций, переименование
│   ├── type_importer.cpp        // ida_typeinf: структуры, интерфейсы
│   ├── string_annotator.cpp     // Аннотация Go-строк
│   ├── calling_convention.cpp   // Go ABI0/ABIInternal CC
│   ├── folder_organizer.cpp     // ida_dirtree: пакеты → папки
│   ├── hexrays_optimizer.cpp    // Hex-Rays microcode transforms (опционально)
│   └── json_parser.cpp          // rapidjson/nlohmann парсинг
├── include/
│   └── nlohmann/json.hpp        // Header-only JSON парсер
├── CMakeLists.txt               // Сборка под Windows/Linux/macOS
├── ida-plugin.json              // IDA 9.0+ метаданные плагина
└── README.md
```

#### 4.3.3. Спецификация основного класса плагина

```cpp
// goreveal_plugin.h
#pragma once
#include <ida.hpp>
#include <idp.hpp>
#include <loader.hpp>
#include <kernwin.hpp>
#include <funcs.hpp>
#include <name.hpp>
#include <typeinf.hpp>
#include <bytes.hpp>
#include <dirtree.hpp>
#include <auto.hpp>

class GoRevealPlugin : public plugmod_t {
public:
    // Точка входа: вызывается при активации плагина
    bool idaapi run(size_t arg) override;

    // Проверка: является ли бинарник Go-бинарником
    static bool is_go_binary();

private:
    // Мост к Go Core
    std::string invoke_core(const std::string& binary_path, 
                           const std::vector<std::string>& args);

    // Шаг 1: Создание Go built-in типов
    void create_builtin_types();
    //   - go_string_t { char *ptr; size_t len; }
    //   - go_slice_t  { void *array; size_t len; size_t cap; }
    //   - go_iface_t  { void *tab; void *data; }
    //   - go_eface_t  { void *_type; void *data; }
    //   - go_complex64_t  { float real; float imag; }
    //   - go_complex128_t { double real; double imag; }

    // Шаг 2: Импорт функций из pclntab
    int import_functions(const nlohmann::json& funcs);
    //   - ida_funcs::add_func(start, end) для не-обнаруженных IDA функций
    //   - ida_name::set_name(ea, name, SN_FORCE|SN_NOWARN|SN_NOCHECK)
    //     Паттерн GoReSym: замена '.' на '_' в именах (IDA ограничение)
    //   - Установка функций-границ для runtime-функций

    // Шаг 3: Применение Go calling convention
    void apply_calling_conventions(const nlohmann::json& funcs);
    //   Go 1.17+ ABIInternal (amd64):
    //     Целые аргументы: RAX, RBX, RCX, RDI, RSI, R8, R9, R10, R11
    //     Float: XMM0–XMM14
    //     R14 = *g (текущая горутина), RDX = closure context
    //   Go <1.17 ABI0:
    //     Все аргументы через стек
    //   IDA 9.0: использовать __golang CC если доступна,
    //     иначе __usercall с явным указанием регистров

    // Шаг 4: Импорт типов через ida_typeinf
    int import_types(const nlohmann::json& types);
    //   - Парсинг CDecl через idc_parse_types(cdecl, HTI_PAKDEF|HTI_DCL)
    //   - Для структур: tinfo_t + add_udm() для каждого поля
    //   - Для интерфейсов: itab-структура с fun[] массивом
    //   - Категоризация в type library: "/golang/runtime/", "/golang/user/"

    // Шаг 5: Аннотация строк
    int annotate_strings(const nlohmann::json& strings);
    //   - Создание строковых литералов через ida_bytes
    //   - Комментарии с декодированными строками на адресах использования
    //   - Для deobfuscated строк: комментарий "garble deobfuscated: <value>"

    // Шаг 6: Организация в папки по пакетам
    void organize_folders(const nlohmann::json& packages);
    //   - ida_dirtree: создание дерева main/ stdlib/ vendor/ unknown/
    //   - Перемещение функций в соответствующие папки
    //   Паттерн AlphaGolang Step 3

    // Шаг 7: Cross-references
    void import_xrefs(const nlohmann::json& xrefs);

    // Опционально: Hex-Rays microcode optimization
    void optimize_hexrays_output();
    //   - Замена runtime.newobject + type_ptr на "new(TypeName)"
    //   - Свёртка Go string concat паттернов
    //   - Улучшение slice/map/channel операций в декомпиляции
};

// IDA plugin entry point
plugmod_t* idaapi init();
```

#### 4.3.4. Сборка (CMakeLists.txt)

```cmake
cmake_minimum_required(VERSION 3.20)
project(goreveal-ida LANGUAGES CXX)

set(CMAKE_CXX_STANDARD 17)

# IDA SDK
set(IDA_SDK_DIR "" CACHE PATH "Path to IDA SDK")
set(IDA_INSTALL_DIR "" CACHE PATH "Path to IDA installation")

# Платформо-зависимые настройки
if(WIN32)
    set(IDA_LIB_SUFFIX ".lib")
    set(PLUGIN_EXT ".dll")
elseif(APPLE)
    set(IDA_LIB_SUFFIX ".dylib")
    set(PLUGIN_EXT ".dylib")
else()
    set(IDA_LIB_SUFFIX ".so")
    set(PLUGIN_EXT ".so")
endif()

# Целевая архитектура: 64-bit
add_definitions(-D__X64__ -D__EA64__)

include_directories(
    ${IDA_SDK_DIR}/include
    ${CMAKE_SOURCE_DIR}/include
)

link_directories(${IDA_SDK_DIR}/lib/x64_linux_gcc_64)  # Платформо-зависимо

add_library(goreveal_ida MODULE
    src/plugin_main.cpp
    src/goreveal_plugin.cpp
    src/core_bridge_subprocess.cpp
    src/function_annotator.cpp
    src/type_importer.cpp
    src/string_annotator.cpp
    src/calling_convention.cpp
    src/folder_organizer.cpp
    src/json_parser.cpp
)

target_link_libraries(goreveal_ida ida pro)

set_target_properties(goreveal_ida PROPERTIES
    PREFIX ""
    SUFFIX "${PLUGIN_EXT}"
    OUTPUT_NAME "goreveal"
)

# Установка
install(TARGETS goreveal_ida DESTINATION ${IDA_INSTALL_DIR}/plugins)
```

#### 4.3.5. ida-plugin.json (IDA 9.0+)

```json
{
    "name": "GoReveal",
    "description": "Universal Go binary reverse engineering platform",
    "version": "1.0.0",
    "author": "GoReveal Team",
    "url": "https://github.com/...",
    "type": "plugin",
    "min_ida_version": "9.0",
    "platforms": ["win64", "linux64", "macos_arm64"]
}
```

### 4.4. Модуль 4: Ghidra Extension (`goreveal-ghidra`)

**Язык**: Java 21+
**Целевая версия Ghidra**: 11.3+
**Лицензия**: Apache 2.0

#### 4.4.1. Архитектура расширения

```
goreveal-ghidra/
├── src/main/java/goreveal/
│   ├── GoRevealAnalyzer.java         // extends AbstractAnalyzer
│   ├── GoRevealImportScript.java     // extends GhidraScript (headless)
│   ├── bridge/
│   │   └── CoreBridge.java           // Subprocess вызов goreveal CLI
│   ├── types/
│   │   ├── GoTypeFactory.java        // Создание Go-типов через DataTypeManager
│   │   ├── GoBuiltinTypes.java       // string, slice, interface, channel
│   │   └── GoStructBuilder.java      // Построение struct из JSON
│   ├── annotators/
│   │   ├── FunctionAnnotator.java    // Создание/переименование функций
│   │   ├── StringAnnotator.java      // Go-строки в листинге
│   │   └── CallingConventionApplier.java
│   └── util/
│       └── JsonParser.java           // com.google.gson парсинг
├── data/
│   └── go_builtin_types.gdt          // Предкомпилированные Go runtime типы
├── extension.properties
├── Module.manifest
└── build.gradle
```

#### 4.4.2. GoRevealAnalyzer

```java
public class GoRevealAnalyzer extends AbstractAnalyzer {
    
    public GoRevealAnalyzer() {
        super("GoReveal", "Go binary analysis via GoReveal platform",
              AnalyzerType.BYTE_ANALYZER);
        // Приоритет: после встроенного GolangSymbolAnalyzer
        setPriority(AnalysisPriority.DATA_TYPE_PROPOGATION.after());
        setSupportsOneTimeAnalysis();
    }

    @Override
    public boolean canAnalyze(Program program) {
        // Проверка Go build ID magic: "\xff Go buildinf:" 
        // или pclntab magic (0xf0ffffff, 0xf1ffffff, etc.)
        return GoDetector.isGoBinary(program);
    }

    @Override
    public boolean added(Program program, AddressSetView set,
                         TaskMonitor monitor, MessageLog log) {
        // 1. Запуск goreveal CLI как subprocess
        String jsonResult = CoreBridge.analyze(
            program.getExecutablePath(), monitor);
        
        // 2. Парсинг JSON
        GoAnalysis analysis = JsonParser.parse(jsonResult);
        
        // 3. Создание Go built-in типов
        DataTypeManager dtm = program.getDataTypeManager();
        int txId = dtm.startTransaction("GoReveal Import");
        try {
            GoTypeFactory.createBuiltinTypes(dtm);
            
            // 4. Импорт пользовательских типов
            for (GoType type : analysis.getTypes()) {
                GoStructBuilder.createType(dtm, type);
            }
            
            // 5. Импорт функций
            FunctionAnnotator.importFunctions(
                program, analysis.getFunctions(), monitor);
            
            // 6. Аннотация строк
            StringAnnotator.annotateStrings(
                program, analysis.getStrings(), monitor);
            
            dtm.endTransaction(txId, true);
        } catch (Exception e) {
            dtm.endTransaction(txId, false);
            log.appendMsg("GoReveal", "Error: " + e.getMessage());
            return false;
        }
        return true;
    }
}
```

#### 4.4.3. Headless-режим для CI/CD

```bash
analyzeHeadless /tmp/projects GoProject \
  -import /path/to/go_binary \
  -postScript GoRevealImport.java /path/to/analysis.json \
  -scriptLog /tmp/goreveal.log
```

### 4.5. Модуль 5: PostgreSQL Storage Layer

#### 4.5.1. Полная схема базы данных

```sql
-- === Расширения ===
CREATE EXTENSION IF NOT EXISTS pg_trgm;       -- Фаззи-поиск по именам
CREATE EXTENSION IF NOT EXISTS pgcrypto;      -- SHA-256 хеширование
-- CREATE EXTENSION IF NOT EXISTS vector;     -- Фаза 4+: семантический поиск

-- === Основные таблицы ===

CREATE TABLE binaries (
    id              SERIAL PRIMARY KEY,
    filename        TEXT NOT NULL,
    sha256          TEXT UNIQUE NOT NULL,
    file_size       BIGINT,
    os              TEXT,          -- linux|windows|darwin
    arch            TEXT,          -- amd64|arm64|386
    endianness      TEXT,          -- little|big
    go_version      TEXT,
    format          TEXT,          -- elf|pe|macho
    build_id        TEXT,
    go_root         TEXT,
    main_root       TEXT,
    module_path     TEXT,
    is_stripped      BOOLEAN DEFAULT FALSE,
    is_obfuscated   BOOLEAN DEFAULT FALSE,
    obfuscator_name TEXT,
    build_settings  JSONB,        -- Расширяемые build settings
    raw_metadata    JSONB,        -- Полный GoReSym-подобный вывод
    analyzed_at     TIMESTAMPTZ DEFAULT NOW(),
    analysis_version TEXT         -- Версия goreveal
);

CREATE TABLE packages (
    id              SERIAL PRIMARY KEY,
    binary_id       INTEGER NOT NULL REFERENCES binaries(id) ON DELETE CASCADE,
    name            TEXT NOT NULL,
    path            TEXT,
    classification  TEXT CHECK(classification IN ('main','stdlib','vendor','unknown')),
    version         TEXT,
    function_count  INTEGER DEFAULT 0,
    type_count      INTEGER DEFAULT 0,
    UNIQUE(binary_id, path)
);

CREATE TABLE functions (
    id              SERIAL PRIMARY KEY,
    binary_id       INTEGER NOT NULL REFERENCES binaries(id) ON DELETE CASCADE,
    package_id      INTEGER REFERENCES packages(id) ON DELETE SET NULL,
    name            TEXT NOT NULL,
    full_name       TEXT NOT NULL,
    start_addr      BIGINT NOT NULL,
    end_addr        BIGINT,
    source_file     TEXT,
    start_line      INTEGER,
    end_line         INTEGER,
    receiver        TEXT,
    is_method       BOOLEAN DEFAULT FALSE,
    arg_size        INTEGER,
    c_signature     TEXT,
    -- Для garble-деобфускации
    original_name   TEXT,         -- Обфусцированное имя (до деобфускации)
    deobfuscated    BOOLEAN DEFAULT FALSE,
    deobfuscation_confidence REAL,
    UNIQUE(binary_id, start_addr)
);

CREATE TABLE types (
    id              SERIAL PRIMARY KEY,
    binary_id       INTEGER NOT NULL REFERENCES binaries(id) ON DELETE CASCADE,
    name            TEXT NOT NULL,
    addr            BIGINT,
    kind            TEXT,         -- struct|interface|func|map|slice|chan|ptr|...
    kind_id         SMALLINT,
    package_path    TEXT,
    size            INTEGER,
    alignment       INTEGER,
    c_decl          TEXT,         -- C-совместимая декларация
    reconstructed   TEXT,         -- Восстановленный Go-код
    metadata        JSONB         -- Дополнительные метаданные типа
);

CREATE TABLE type_fields (
    id              SERIAL PRIMARY KEY,
    type_id         INTEGER NOT NULL REFERENCES types(id) ON DELETE CASCADE,
    name            TEXT,
    field_type_name TEXT,
    field_type_id   INTEGER REFERENCES types(id),
    offset          INTEGER,
    size            INTEGER,
    tag             TEXT,
    anonymous       BOOLEAN DEFAULT FALSE,
    field_order     INTEGER       -- Порядок поля в структуре
);

CREATE TABLE type_methods (
    id              SERIAL PRIMARY KEY,
    type_id         INTEGER NOT NULL REFERENCES types(id) ON DELETE CASCADE,
    name            TEXT NOT NULL,
    func_id         INTEGER REFERENCES functions(id),
    signature       TEXT,
    is_exported     BOOLEAN DEFAULT FALSE
);

CREATE TABLE strings (
    id              SERIAL PRIMARY KEY,
    binary_id       INTEGER NOT NULL REFERENCES binaries(id) ON DELETE CASCADE,
    value           TEXT,
    addr            BIGINT,
    length          INTEGER,
    referencing_func_id INTEGER REFERENCES functions(id),
    is_deobfuscated BOOLEAN DEFAULT FALSE,
    original_encrypted BYTEA       -- Зашифрованное значение (для garble)
);

CREATE TABLE cross_references (
    id              SERIAL PRIMARY KEY,
    binary_id       INTEGER NOT NULL REFERENCES binaries(id) ON DELETE CASCADE,
    from_addr       BIGINT NOT NULL,
    to_addr         BIGINT NOT NULL,
    xref_type       TEXT CHECK(xref_type IN ('call','data','string','type')),
    from_func_id    INTEGER REFERENCES functions(id),
    to_func_id      INTEGER REFERENCES functions(id)
);

CREATE TABLE dependencies (
    id              SERIAL PRIMARY KEY,
    binary_id       INTEGER NOT NULL REFERENCES binaries(id) ON DELETE CASCADE,
    path            TEXT NOT NULL,
    version         TEXT,
    sum             TEXT,         -- Go module checksum
    is_replaced     BOOLEAN DEFAULT FALSE,
    replace_path    TEXT,
    replace_version TEXT,
    UNIQUE(binary_id, path)
);

CREATE TABLE source_files (
    id              SERIAL PRIMARY KEY,
    binary_id       INTEGER NOT NULL REFERENCES binaries(id) ON DELETE CASCADE,
    filepath        TEXT NOT NULL,
    package_id      INTEGER REFERENCES packages(id),
    function_count  INTEGER DEFAULT 0
);

-- === Индексы ===

-- Функции: основные запросы
CREATE INDEX idx_functions_binary ON functions(binary_id);
CREATE INDEX idx_functions_addr ON functions(binary_id, start_addr);
CREATE INDEX idx_functions_package ON functions(package_id);
CREATE INDEX idx_functions_name_trgm ON functions USING gin(full_name gin_trgm_ops);

-- Типы
CREATE INDEX idx_types_binary ON types(binary_id);
CREATE INDEX idx_types_kind ON types(binary_id, kind);
CREATE INDEX idx_types_name_trgm ON types USING gin(name gin_trgm_ops);

-- Строки: полнотекстовый + trigramный поиск
CREATE INDEX idx_strings_binary ON strings(binary_id);
CREATE INDEX idx_strings_value_trgm ON strings USING gin(value gin_trgm_ops);

-- Cross-references
CREATE INDEX idx_xref_from ON cross_references(binary_id, from_addr);
CREATE INDEX idx_xref_to ON cross_references(binary_id, to_addr);

-- Зависимости
CREATE INDEX idx_deps_path ON dependencies(path);

-- === Materialized views ===

CREATE MATERIALIZED VIEW mv_package_summary AS
SELECT 
    p.binary_id,
    p.name AS package_name,
    p.path AS package_path,
    p.classification,
    COUNT(DISTINCT f.id) AS func_count,
    COUNT(DISTINCT t.id) AS type_count,
    MIN(f.start_addr) AS min_addr,
    MAX(f.end_addr) AS max_addr
FROM packages p
LEFT JOIN functions f ON f.package_id = p.id
LEFT JOIN types t ON t.package_path = p.path AND t.binary_id = p.binary_id
GROUP BY p.binary_id, p.name, p.path, p.classification;

CREATE MATERIALIZED VIEW mv_binary_overview AS
SELECT 
    b.id,
    b.filename,
    b.go_version,
    b.os || '/' || b.arch AS platform,
    b.is_obfuscated,
    b.obfuscator_name,
    COUNT(DISTINCT f.id) FILTER (WHERE f.package_id IN 
        (SELECT id FROM packages WHERE classification = 'main')) AS user_func_count,
    COUNT(DISTINCT f.id) AS total_func_count,
    COUNT(DISTINCT t.id) AS type_count,
    COUNT(DISTINCT d.id) AS dep_count
FROM binaries b
LEFT JOIN functions f ON f.binary_id = b.id
LEFT JOIN types t ON t.binary_id = b.id
LEFT JOIN dependencies d ON d.binary_id = b.id
GROUP BY b.id;

-- Обновление materialized views
-- REFRESH MATERIALIZED VIEW CONCURRENTLY mv_package_summary;
-- REFRESH MATERIALIZED VIEW CONCURRENTLY mv_binary_overview;

-- === Фаза 4+: pgvector для семантического поиска ===
-- ALTER TABLE functions ADD COLUMN embedding vector(768);
-- CREATE INDEX ON functions USING hnsw (embedding vector_cosine_ops);
```

#### 4.5.2. Ключевые SQL-запросы для AI-анализа

```sql
-- Фаззи-поиск по garble-обфусцированным именам (pg_trgm)
SELECT full_name, similarity(full_name, 'crypto') AS sim
FROM functions WHERE binary_id = 1
AND full_name % 'crypto'
ORDER BY sim DESC LIMIT 20;

-- Корреляция зависимостей между бинарниками
SELECT d1.path, b1.filename AS binary1, b2.filename AS binary2
FROM dependencies d1
JOIN dependencies d2 ON d1.path = d2.path AND d1.binary_id != d2.binary_id
JOIN binaries b1 ON b1.id = d1.binary_id
JOIN binaries b2 ON b2.id = d2.binary_id;

-- Функции использующие crypto-пакеты (потенциальный C2/ransomware)
SELECT f.full_name, p.path, f.start_addr
FROM functions f
JOIN packages p ON f.package_id = p.id
WHERE p.path LIKE 'crypto/%' AND p.classification != 'stdlib'
AND f.binary_id = 1;

-- Обзор бинарника
SELECT * FROM mv_binary_overview WHERE id = 1;
```

---

## 5. Конвейер деобфускации garble

### 5.1. Четырёхстадийный pipeline

```
Вход: garble-обфусцированный Go-бинарник
        │
   ┌────▼────────────────────────┐
   │ СТАДИЯ 1: Извлечение мета   │  ~1-5 сек
   │ GoReveal Core:               │
   │ - pclntab, типы, buildinfo   │
   │ - Определение версии Go      │
   │ - Детекция обфускации        │
   │   (garble fingerprint)       │
   └────────────┬────────────────┘
                │
   ┌────────────▼────────────────┐
   │ СТАДИЯ 2: Строковая деобф.  │  ~1-10 мин
   │ GoStringUngarbler (Apache):  │
   │ - Regex → пролог функций    │
   │ - Unicorn эмуляция          │
   │ - Выход: расшифрованные     │
   │   строки + адреса           │
   └────────────┬────────────────┘
                │
   ┌────────────▼────────────────┐
   │ СТАДИЯ 3: CFG-деобфускация  │  ~10-60 мин
   │ GoResolver (pip):            │
   │ - Fingerprint Go version     │
   │ - Bootstrap reference        │
   │ - CFG similarity match       │
   │ - Package propagation        │
   └────────────┬────────────────┘
                │
   ┌────────────▼────────────────┐
   │ СТАДИЯ 4: Объединение       │
   │ - Merge стадий 1+2+3        │
   │ - Запись в PostgreSQL        │
   │ - JSON/Protobuf для IDA     │
   │ - Генерация отчёта          │
   └─────────────────────────────┘
```

### 5.2. Стратегия интеграции

GoStringUngarbler (Apache 2.0) и GoResolver вызываются как **subprocess** из Go-оркестратора. Это позволяет: получать обновления из upstream без maintenance burden, не зависеть от Python-зависимостей в Go Core (только при деобфускации), параллелизировать стадии 2 и 3 (они независимы).

### 5.3. Детекция обфускации

```go
func (a *Analyzer) DetectObfuscation() (*ObfuscationInfo, error) {
    info := &ObfuscationInfo{}
    
    // Признаки garble:
    // 1. Отсутствие BuildInfo (garble -tiny)
    // 2. Имена функций — хеши фиксированной длины
    //    (garble использует content-addressable hashing)
    // 3. Консистентный префикс на все функции одного пакета
    // 4. Отсутствие debug info и source file paths
    // 5. Паттерны -literals: XOR/ADD/SUB пролог перед slicebytetostring
    
    // Эвристика: если >80% user-функций имеют имена из hex-символов
    // длиной 8–16 символов → garble с высокой вероятностью
    
    return info, nil
}
```

---

## 6. Генерация отчётов

### 6.1. Markdown-отчёт (пример структуры)

```markdown
# GoReveal Analysis Report

## Binary Overview
- **File**: suspicious_binary
- **SHA-256**: abc123...
- **Go Version**: go1.21.5
- **Platform**: linux/amd64
- **Obfuscated**: Yes (garble, confidence: 0.95)

## Build Dependencies
| Module | Version | Notes |
|--------|---------|-------|
| github.com/gorilla/websocket | v1.5.0 | WebSocket library |
| golang.org/x/crypto | v0.14.0 | Extended crypto |

## Package Summary
- **User packages**: 5 (23 functions)
- **Vendor packages**: 12 (187 functions)  
- **Stdlib packages**: 142 (4,561 functions)

## Suspicious Indicators
- Uses `crypto/aes` + `crypto/cipher` (encryption capability)
- Uses `net/http` + `gorilla/websocket` (network communication)
- Contains hardcoded IP addresses: 192.168.1.100:8443
- 3 functions with `exec.Command` calls (command execution)

## Source Tree Projection
[tree structure]

## Deobfuscation Results
- Recovered 187/210 function names (89%)
- Decrypted 45 string literals
```

### 6.2. HTML-отчёт

Интерактивный HTML с навигацией по пакетам, поиском по функциям, визуализацией source tree. Генерируется через Go `html/template`.

---

## 7. Лицензионная стратегия

### 7.1. Архитектура лицензий

| Модуль | Лицензия | Обоснование |
|---|---|---|
| Core Library + CLI | **AGPL-3.0** | GoRE зависимость (AGPL-3.0) при Go static linking делает весь бинарник AGPL |
| IDA C++ Plugin | **Apache 2.0** | Отдельное произведение, subprocess-граница, не содержит AGPL-код |
| Ghidra Extension | **Apache 2.0** | Аналогично IDA Plugin |
| Protobuf/JSON Schema | **Apache 2.0** | Чтобы клиентский код не наследовал AGPL |
| Документация | **CC BY 4.0** | Стандарт для технической документации |

### 7.2. Ключевые правила

1. IDA Plugin и Ghidra Extension **никогда** не импортируют Go-код напрямую — только потребляют JSON/Protobuf через файл или stdout
2. При переходе на in-process (c-shared), IDA C++ Plugin загружает `libgoreveal.so` через `dlopen` — динамическая линковка с AGPL-библиотекой не требует AGPL для вызывающего кода при соблюдении subprocess-аналогичной семантики (спорно, но практика Zitadel/MongoDB подтверждает)
3. Всем пользователям Core доступен полный исходный код (требование AGPL)

### 7.3. Альтернативный путь (MIT/Apache 2.0 для всего)

Если AGPL неприемлем: реализовать Core на основе GoReSym (MIT) + собственная реализация функциональности GoRE. Дополнительные трудозатраты: **+12–16 недель**. Ключевые компоненты для переписывания: парсинг moduledata, классификация пакетов, source tree projection, метод-ресивер привязка.

---

## 8. План разработки по фазам

### Фаза 1: Foundation (8 недель)

| # | Задача | Оценка | Результат |
|---|---|---|---|
| 1.1 | Инициализация проекта, CI/CD, набор тестовых бинарников (50+) | 1 нед | Go module, GitHub Actions, тестовая матрица Go 1.16–1.26 |
| 1.2 | Core Library: интеграция GoRE, парсинг pclntab | 2 нед | `Open()`, `Functions()`, `Packages()` |
| 1.3 | Core Library: типы, строки, source tree, buildinfo | 2 нед | `Types()`, `Strings()`, `SourceTree()`, `BuildInfo()` |
| 1.4 | Core Library: C-декларации типов для IDA | 1 нед | `CDecl` и `CSignature` поля в структурах |
| 1.5 | CLI: все базовые команды | 1 нед | `goreveal analyze/functions/types/strings/source-tree/info` |
| 1.6 | Protobuf + JSON Schema определения | 1 нед | `goreveal.proto`, JSON Schema v1 |

**Milestone 1**: CLI-инструмент, анализирующий Go-бинарники Go 1.16–1.26, выдающий JSON/Protobuf с полными метаданными.

### Фаза 2: PostgreSQL + IDA C++ Plugin (8 недель)

| # | Задача | Оценка | Результат |
|---|---|---|---|
| 2.1 | PostgreSQL: полная схема, миграции, экспорт из Core | 2 нед | `goreveal analyze --export-pg`, все таблицы + индексы + mv |
| 2.2 | SQLite fallback | 0.5 нед | `goreveal analyze --export-sqlite` |
| 2.3 | IDA C++ Plugin: каркас, сборка CMake, core_bridge subprocess | 1.5 нед | Плагин загружается, вызывает goreveal CLI |
| 2.4 | IDA Plugin: function import + naming | 1 нед | Функции создаются и переименовываются |
| 2.5 | IDA Plugin: type import (ida_typeinf, tinfo_t) | 1.5 нед | Все Go-типы как C-структуры в IDA |
| 2.6 | IDA Plugin: calling convention (ABI0 + ABIInternal) | 0.5 нед | Правильная декомпиляция аргументов |
| 2.7 | IDA Plugin: string annotation + folder organization | 0.5 нед | Строки аннотированы, функции в папках |
| 2.8 | IDA Plugin: тестирование на 20+ реальных бинарниках | 0.5 нед | Regression test suite |

**Milestone 2**: Полная цепочка CLI → PostgreSQL + CLI → IDA Pro 9.0+ аннотация (C++ native plugin).

### Фаза 3: Ghidra + Деобфускация + Отчёты (8 недель)

| # | Задача | Оценка | Результат |
|---|---|---|---|
| 3.1 | Ghidra Analyzer extension (Java) | 2.5 нед | Auto-analysis при обнаружении Go-бинарника |
| 3.2 | Ghidra: Go DataTypes + string annotation | 1 нед | Полная аннотация в Ghidra 11.3+ |
| 3.3 | Детекция обфускации (garble fingerprint) | 0.5 нед | `goreveal analyze` сообщает об обфускации |
| 3.4 | Интеграция GoStringUngarbler (строковая деобфускация) | 1.5 нед | `goreveal deobfuscate --string-ungarble` |
| 3.5 | Интеграция GoResolver (CFG similarity) | 1.5 нед | `goreveal deobfuscate --cfg-resolve` |
| 3.6 | Генерация отчётов (Markdown + HTML) | 1 нед | `goreveal report --format=markdown/html` |

**Milestone 3**: Полнофункциональный инструмент с IDA + Ghidra + деобфускация garble + отчёты.

### Фаза 4: Advanced Features + In-Process Integration (6 недель)

| # | Задача | Оценка | Результат |
|---|---|---|---|
| 4.1 | C Shared Library (libgoreveal.so) | 1.5 нед | `-buildmode=c-shared`, C API |
| 4.2 | IDA Plugin: переход на in-process через dlopen | 1 нед | Устранение subprocess overhead |
| 4.3 | IDA Plugin: Hex-Rays microcode optimizer | 1.5 нед | Улучшенная декомпиляция Go-конструкций |
| 4.4 | Cross-reference engine | 1 нед | Таблица xref в PostgreSQL |
| 4.5 | Документация, примеры, расширенный test suite | 1 нед | API docs, tutorial, 50+ бинарников |

### Суммарная оценка

| Фаза | Длительность | Трудозатраты |
|---|---|---|
| 1. Foundation | 8 нед | 8 чел-нед |
| 2. PostgreSQL + IDA C++ | 8 нед | 8 чел-нед |
| 3. Ghidra + Deobfuscation + Reports | 8 нед | 8 чел-нед |
| 4. Advanced + In-Process | 6 нед | 6 чел-нед |
| **Итого** | **~30 недель** | **~30 чел-нед** |

При одном full-time разработчике — **~7.5 месяцев**. При двух — **~4.5 месяца** с параллелизацией (Фаза 2: один на PostgreSQL, другой на IDA; Фаза 3: один на Ghidra, другой на деобфускацию).

---

## 9. Технические риски и митигация

### 9.1. Высокий приоритет

| Риск | Вероятность | Влияние | Митигация |
|---|---|---|---|
| Go runtime-структуры меняются между версиями (pclntab, _type) | Высокая | Высокое | Форк runtime-парсеров по версиям (подход GoReSym). CI-мониторинг Go releases. Автотесты на матрице Go 1.16–1.26 |
| AGPL-заражение плагинов при нарушении subprocess-границы | Средняя | Критическое | Строгая subprocess-граница. Code review с AGPL-checklist. Плагины **никогда** не импортируют Go-код. Отдельные git-репозитории |
| IDA 9.x API breaking changes | Средняя | Высокое | Таргетирование только IDA 9.0+. Abstraction layer для ida_typeinf. Тестирование на каждом IDA update |
| C++ IDA plugin: сложность отладки + cross-platform | Высокая | Среднее | CMake кроссплатформенная сборка. CI матрица: win64, linux64, macOS_arm64. Address Sanitizer в debug builds |

### 9.2. Средний приоритет

| Риск | Вероятность | Влияние | Митигация |
|---|---|---|---|
| GoResolver: неизвестная лицензия, вычислительная тяжесть | Средняя | Среднее | Проверить LICENSE перед интеграцией. Timeout + параллелизация. Режимы «быстрый» (stdlib only) и «полный» |
| garble эволюционирует (control flow flattening) | Высокая | Среднее | Модульная деобфускация — добавление новых стадий. Мониторинг garble releases |
| CGo cross-compilation для libgoreveal.so | Средняя | Среднее | Только в Фазе 4. Docker-based cross-compilation. Fallback на subprocess |
| PostgreSQL: overhead развёртывания для solo-аналитиков | Низкая | Низкое | SQLite fallback всегда доступен. Docker Compose для быстрого запуска PG |

### 9.3. Низкий приоритет

| Риск | Вероятность | Влияние | Митигация |
|---|---|---|---|
| Ghidra API backward compatibility | Низкая | Среднее | Ghidra 11.3+. Dual mode: analyzer + fallback GhidraScript |
| Объём тестового набора | Средняя | Низкое | GoStrap (Volexity) для генерации reference-бинарников. CI pipeline |
| Hex-Rays microcode API complexity | Высокая | Низкое (опционально) | Фаза 4, отдельный модуль. Начать с простых pattern-based трансформаций |

---

## 10. Требования к окружению разработки

### 10.1. Обязательные инструменты

| Инструмент | Версия | Назначение |
|---|---|---|
| Go | 1.26+ | Core Library, CLI |
| GCC/Clang | 12+ / 15+ | CGo, IDA C++ Plugin |
| CMake | 3.20+ | Сборка IDA Plugin |
| IDA Pro SDK | 9.0+ | C++ headers и libraries |
| IDA Pro | 9.0+ | Тестирование плагина |
| Java JDK | 21+ | Ghidra Extension |
| Gradle | 8.0+ | Сборка Ghidra Extension |
| Ghidra | 11.3+ | Тестирование расширения |
| Python | 3.13+ | GoStringUngarbler, GoResolver |
| PostgreSQL | 16+ | Primary storage |
| Docker | 24+ | CI/CD, cross-compilation |
| protoc | 3.21+ | Protobuf compilation |

### 10.2. Рекомендуемая структура репозиториев

```
github.com/<org>/goreveal           # Core + CLI (AGPL-3.0)
github.com/<org>/goreveal-ida       # IDA C++ Plugin (Apache 2.0)
github.com/<org>/goreveal-ghidra    # Ghidra Extension (Apache 2.0)
github.com/<org>/goreveal-proto     # Protobuf + JSON Schema (Apache 2.0)
github.com/<org>/goreveal-testdata  # Тестовые бинарники (MIT)
```

Раздельные репозитории обеспечивают чистое лицензионное разделение.

---

## 11. Критерии приёмки MVP (Milestone 3)

### 11.1. Функциональные критерии

- [ ] CLI анализирует Go-бинарники версий Go 1.16–1.26 (ELF, PE, Mach-O)
- [ ] Восстановление ≥99% функций из pclntab (паритет с GoReSym)
- [ ] Восстановление ≥95% типов включая поля структур и методы
- [ ] Извлечение строк с корректной длиной (не null-terminated)
- [ ] Классификация пакетов: main, stdlib, vendor, unknown
- [ ] Проекция исходного дерева (source tree)
- [ ] Экспорт в PostgreSQL с полной схемой и индексами
- [ ] Экспорт в SQLite (portable fallback)
- [ ] IDA C++ Plugin: автоматическая аннотация бинарника за <30 секунд (типичный 15 МБ бинарник)
- [ ] IDA Plugin: корректное применение Go calling convention (ABI0 + ABIInternal)
- [ ] IDA Plugin: все Go-типы импортированы как C-структуры
- [ ] Ghidra Extension: автоматический анализ при открытии Go-бинарника
- [ ] Детекция garble обфускации с confidence ≥0.9
- [ ] Строковая деобфускация garble `-literals` (5 режимов)
- [ ] CFG-деобфускация через GoResolver
- [ ] Генерация Markdown/HTML отчётов

### 11.2. Нефункциональные критерии

- [ ] Время анализа типичного 15 МБ бинарника: <10 секунд (без деобфускации)
- [ ] Потребление памяти: <500 МБ для бинарников до 50 МБ
- [ ] Тестовое покрытие: ≥70% для Core Library
- [ ] Тестовый набор: 50+ бинарников (Go 1.16–1.26 × amd64/arm64 × clean/garble)
- [ ] Документация: README, API reference, tutorial для каждого модуля
- [ ] CI/CD: автоматическая сборка и тестирование на каждый PR

---

## Приложения

### A. Protobuf Schema (goreveal.proto) — краткая версия

```protobuf
syntax = "proto3";
package goreveal.v1;

message GoAnalysis {
  BinaryInfo binary = 1;
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

// Полные определения — в отдельном файле goreveal.proto
```

### B. Docker Compose для разработки

```yaml
version: "3.9"
services:
  postgres:
    image: pgvector/pgvector:pg16
    environment:
      POSTGRES_DB: goreveal
      POSTGRES_USER: goreveal
      POSTGRES_PASSWORD: devpass
    ports:
      - "5432:5432"
    volumes:
      - pgdata:/var/lib/postgresql/data
      - ./migrations:/docker-entrypoint-initdb.d

volumes:
  pgdata:
```

### C. Ссылки на upstream-проекты

| Проект | URL | Лицензия |
|---|---|---|
| GoReSym | github.com/mandiant/GoReSym | MIT |
| GoRE (gore) | github.com/goretk/gore | AGPL-3.0 |
| redress | github.com/goretk/redress | AGPL-3.0 |
| GoResolver | github.com/volexity/GoResolver | Уточнить |
| GoStringUngarbler | github.com/mandiant/gostringungarbler | Apache 2.0 |
| IDA SDK | hex-rays.com/ida-sdk | MIT (открыт) |
| Ghidra | github.com/NationalSecurityAgency/ghidra | Apache 2.0 |