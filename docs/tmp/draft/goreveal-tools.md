Прямого аналога Javassist для Go не существует — и это фундаментальное отличие платформ. Javassist работает с JVM-байткодом (высокоуровневый, стековый, стандартизированный), а Go компилируется в нативный машинный код. Но есть инструменты на разных уровнях, покрывающие аналогичные задачи.

### Уровень 1: Манипуляция исходным кодом (AST) — ближайший аналог

Это самый близкий к Javassist подход в Go, потому что AST-трансформации до компиляции = Javassist-трансформации до загрузки класса.

**go/ast + go/parser + go/printer** (stdlib) — чтение, модификация и запись Go-исходников. Позволяет добавлять/удалять функции, менять сигнатуры, инжектировать код. Ограничение: теряет комментарии и форматирование при roundtrip.

**dave/dst** (Decorated Syntax Tree) — решает главную проблему go/ast: **сохраняет комментарии, пустые строки и форматирование** при модификации. Для инструментации Go-кода это практически обязательная замена go/ast.

**jennifer** (`dave/jennifer`) — программная генерация Go-кода через fluent API, аналогично тому как Javassist позволяет строить классы программно:

```go
f := jen.NewFile("main")
f.Func().Id("main").Params().Block(
    jen.Qual("fmt", "Println").Call(jen.Lit("hello")),
)
// Генерирует: func main() { fmt.Println("hello") }
```

### Уровень 2: Compile-time инструментация — аналог -javaagent

**`-toolexec`** — механизм Go compiler, позволяющий подменять инструменты сборки. Именно так работает **garble** (обфускатор): перехватывает вызовы `compile` и `link`, модифицирует AST между этапами. Это closest equivalent к Javassist `ClassFileTransformer` — перехват на этапе компиляции.

**golang.org/x/tools/go/analysis** — фреймворк для статического анализа и трансформации кода (базис для `staticcheck`, `golangci-lint`). Может генерировать `SuggestedFix` — автоматические правки кода.

**golang.org/x/tools/go/ssa** — SSA (Static Single Assignment) представление Go-программ. Аналог Soot Jimple для Java — промежуточное представление для анализа потока данных, call graph construction, pointer analysis.

### Уровень 3: Runtime monkey-patching — аналог runtime instrumentation

**agiledragon/gomonkey v2** (~2,100 stars) — runtime-патчинг функций через подмену первых байт machine code на JMP к замене. Работает через `syscall.Mprotect` для снятия write-protection с code pages:

```go
patches := gomonkey.ApplyFunc(targetFunc, replacementFunc)
defer patches.Reset()
```

Ограничения: только amd64/arm64, ломается при инлайнинге (`-gcflags="-l"` обязателен), не thread-safe, только для тестов.

**bou.ke/monkey** — пионер подхода, но не поддерживается с 2022. gomonkey — его духовный наследник.

### Уровень 4: Бинарная инструментация — для GoReveal

Для твоего проекта реверс-инжиниринга ключевые инструменты:

**go-delve/delve** internals — ptrace-based инструментация, breakpoints через `INT3` injection, goroutine-aware. Delve API можно использовать программно для динамической инструментации Go-бинарников.

**Frida** — как мы обсуждали, работает с Go через JVMTI-аналогичный подход, но требует asm-трамплинов из-за нестандартного ABI.

**eBPF uprobes** — наименее инвазивный метод: kernel-level перехват вызовов функций без модификации бинарника.

### Сводка: что использовать в GoReveal

| Javassist-задача | Go-эквивалент | Где в GoReveal |
|---|---|---|
| Чтение .class → модификация → запись | `dave/dst` (AST roundtrip) | Генерация псевдо-Go из бинарника |
| Программная генерация классов | `dave/jennifer` (code gen) | Генерация C-деклараций для IDA |
| ClassFileTransformer (compile-time) | `-toolexec` + AST transforms | Понимание как работает garble |
| Runtime instrumentation | `gomonkey` / Frida / eBPF | Динамический анализ |
| SSA/IR анализ (как Soot) | `golang.org/x/tools/go/ssa` | Статический анализ потока данных |

Для GoReveal конкретно: **dave/dst** и **jennifer** полезны для модуля генерации псевдо-Go кода (аналог redress source-tree projection), а **golang.org/x/tools/go/ssa** — для будущего модуля статического анализа декомпилированного кода.