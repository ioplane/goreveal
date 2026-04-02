# Комплексная защита коммерческого ПО на Go: от бинарника до лицензии

Go-приложения в модели OSS + Enterprise Edition требуют трёхуровневой стратегии защиты: обфускация бинарников, архитектурное разделение кодовой базы и криптографическая система лицензирования. Ни один из уровней не является самодостаточным — только их сочетание создаёт экономически оправданный барьер для нелицензионного использования. Это исследование охватывает инструменты, паттерны и конкретные решения для Go 1.26 с фокусом на практическую реализацию.

---

## Часть 1. Защита Go-бинарников от реверс-инжиниринга

### Go-бинарники — самые «прозрачные» среди компилируемых языков

Go-бинарники **значительно более уязвимы к реверс-инжинирингу**, чем аналоги на C/C++ и Rust. Причина — рантайм Go требует метаданных, которые невозможно удалить без нарушения работоспособности программы.

**Ключевые структуры, облегчающие RE:**

- **PCLNTAB (`.gopclntab`)** — таблица соответствия программных счётчиков именам функций и строкам исходного кода. Присутствует с Go 1.2 и **не может быть удалена** — без неё сборщик мусора и рантайм перестают работать. Даже в stripped-бинарнике «Hello World» содержится свыше 1 000 именованных функций.
- **Moduledata** — внутренняя таблица рантайма с раскладкой исполняемого файла, ссылками на типы и pclntab.
- **Type Information (RTTI)** — полные дескрипторы всех типов (структуры, интерфейсы, каналы, срезы). Необходимы для рефлексии, interface dispatch и GC.
- **Buildinfo (Go 1.18+)** — версия компилятора, зависимости, флаги сборки, GOOS/GOARCH, иногда Git-информация.
- **Строковая таблица** — Go не использует null-terminated строки; линкер размещает их последовательно с указателем + длиной, образуя легко извлекаемые блобы. По умолчанию встраиваются полные пути к исходным файлам и домашний каталог пользователя сборки.

Для сравнения: в C/C++ с флагом `-s` удаляются практически все символы, и рантайм в них не нуждается. Rust при stripping ведёт себя аналогично C. В Go же даже с `-ldflags "-s -w"` сохраняются **pclntab, moduledata, type info и имена функций**.

### Инструменты RE для Go-бинарников: ландшафт угроз

Понимание инструментов атакующего — основа грамотной защиты.

| Инструмент | Назначение | Особенности |
|---|---|---|
| **GoReSym** (Mandiant/FLARE) | Извлечение метаданных, имён функций, типов | Работает со stripped и запакованными бинарниками; интегрирован в VirusTotal; поддерживает все архитектуры |
| **GoResolver** (Volexity, 2025) | Деобфускация garble через CFG-similarity | Сравнивает графы потока управления обфусцированного бинарника с чистыми шаблонами; плагины для IDA Pro и Ghidra |
| **GoStringUngarbler** (Mandiant, 2025) | Автоматическая расшифровка garble `-literals` | Обрабатывает все три типа трансформаций garble (stack, seed, split) |
| **Redress** (GoRETK) | Восстановление пакетов, типов, проекции исходного кода | CLI + плагин для radare2; представлен на DEF CON |
| **AlphaGolang** (SentinelOne) | Набор IDAPython-скриптов для пошагового анализа | Восстановление pclntab, именование функций, очистка RTYPE |
| **IDA Pro 8.3+** | Нативный анализ Go-бинарников | Реконструкция PCLN-таблицы, Lumina для идентификации |
| **Ghidra** | Бесплатная альтернатива IDA | Находит ~63% функций в stripped Go-бинарниках; кастомные скрипты для pclntab |

### Garble — основной инструмент обфускации

**Garble** (`mvdan.cc/garble`) — единственный активно поддерживаемый обфускатор для Go (2024–2026). Gobfuscate устарел и не работает с Go modules.

**Возможности garble:**

| Функция | Флаг/директива | Статус |
|---|---|---|
| Обфускация имён (функций, пакетов, типов) | По умолчанию | Стабильно |
| Шифрование строковых литералов | `-literals` | Стабильно |
| Удаление имён файлов, строк паники | `-tiny` | Стабильно |
| Детерминированный/случайный seed | `-seed=random` | Стабильно |
| Обфускация потока управления | `GARBLE_EXPERIMENTAL_CONTROLFLOW=1` + `//garble:controlflow` | Экспериментально |

**Рекомендуемая команда сборки:**

```bash
GARBLE_EXPERIMENTAL_CONTROLFLOW=1 \
garble -literals -tiny -seed=random \
  build -trimpath -buildmode=pie \
  -ldflags="-s -w -buildid=" \
  -o myapp ./cmd/myapp
```

**Контроль потока управления** — самая мощная функция garble. Активируется per-функция через директиву:

```go
//garble:controlflow block_splits=max junk_jumps=max flatten_passes=max
func verifyLicense(key string) bool {
    // критическая логика
}
```

Параметры `block_splits` разбивают SSA-блоки, `junk_jumps` вставляют мёртвый код, `flatten_passes` уплощают CFG, `flatten_hardening` (xor/delegate_table) защищает диспетчер от статического анализа.

**Влияние на производительность:** сборка замедляется **примерно в 2 раза** (два прохода: оригинальный + обфусцированный). `-tiny` уменьшает бинарник на **2–15%**. Кэш Go build полностью поддерживается.

**Критические ограничения garble:**
- **Экспортируемые методы никогда не обфусцируются** (могут требоваться интерфейсами)
- **Go plugins не поддерживаются** (issue #87)
- Константы (`const`) не обфусцируются (вычисляются при компиляции)
- Пакеты с `//go:linkname` в рантайм — не обфусцируются

**Важнейшее ограничение**: инструменты деобфускации **уже существуют**. GoStringUngarbler от Mandiant автоматически расшифровывает `-literals`, GoResolver от Volexity восстанавливает имена функций через CFG-similarity. Garble повышает стоимость RE, но не делает его невозможным.

### Stripping символов и удаление отладочной информации

| Флаг | Что удаляет | Команда |
|---|---|---|
| `-ldflags="-w"` | DWARF-таблицу (отладочная информация) | `go build -ldflags="-w"` |
| `-ldflags="-s"` | Таблицу символов + отладочную информацию | `go build -ldflags="-s"` |
| `-trimpath` | Полные пути к файлам → относительные пути модулей | `go build -trimpath` |
| `-ldflags="-buildid="` | Уникальный Build ID | `go build -ldflags="-buildid="` |

**Максимальный stripping:**

```bash
go build -trimpath -ldflags="-w -s -buildid=" -o myapp
```

Уменьшает размер бинарника на **25–37%** (на реальных проектах масштаба Istio pilot-agent: с ~100 МБ до ~75 МБ). Однако **pclntab, type info и имена функций сохраняются**. Использование GNU `strip` на Go-бинарниках **категорически не рекомендуется** — может вызвать случайные паники (документировано командой Docker/Moby).

### Анти-дебаггинг и шифрование строк

**Библиотека `biter777/antideb`** — готовое решение для Linux, включающее обнаружение ptrace, программных точек останова (INT3), отладочных переменных окружения, родительского процесса (gdb, strace), LD_PRELOAD-хуков и отключённого ASLR:

```go
import "github.com/biter777/antideb"

func main() {
    go antideb.Detect(true) // panic при обнаружении дебаггера
    // основная логика
}
```

Дополнительно: проверка `/proc/self/status` на `TracerPid != 0`, timing-based детекция (замедление при single-stepping). Но все userspace-техники обходимы — они лишь **повышают стоимость** RE.

**Шифрование строк** через garble `-literals` использует три типа трансформаций: stack (XOR с ключом на стеке), seed (цепная зависимость байтов), split (фрагментация по switch-блокам). Для дополнительной защиты критических секретов используйте кастомное шифрование через `go:generate`:

```go
//go:generate go run ./cmd/encrypt-strings
var encryptedAPIKey = []byte{0x4a, 0x2b, ...} // предварительно зашифровано
```

### Runtime integrity checks и использование CGO

**Самопроверка бинарника** выполняется через вычисление SHA-256 собственного исполняемого файла при старте и сравнение с хешем, встроенным через `-ldflags`. Go 1.24+ предлагает нативную самопроверку целостности в рамках FIPS 140-3 (`GODEBUG=fips140=on`), логику которой можно изучить для кастомных реализаций.

**CGO для критических секций** — компромиссный подход. C-код не несёт Go-метаданных и может использовать LLVM-обфускаторы (Hikari, O-LLVM). Однако цена высока: CGO-вызов стоит **~100 нс против ~2 нс** для Go-вызова, кросс-компиляция усложняется, бинарник становится динамически слинкованным, теряется воспроизводимость сборки. Рекомендация: используйте CGO **только** для самых критичных фрагментов (проверка лицензии, деривация ключей), минимизируя границу взаимодействия.

### PIE и другие build-time опции

Сборка с `-buildmode=pie` включает ASLR (рандомизация адресов), усложняя статический анализ. Влияние на производительность — менее 1%. Полная «защищённая» сборка:

```bash
CGO_ENABLED=0 \
garble -literals -tiny -seed=random \
  build -trimpath -buildmode=pie \
  -ldflags="-s -w -buildid= -X main.version=1.0.0" \
  -o myapp ./cmd/myapp
```

**Многослойная стратегия защиты бинарника (по приоритету):**

1. garble с `-literals -tiny` — наибольший эффект при минимальных затратах
2. Build flags: `-trimpath -ldflags="-s -w -buildid="` — без влияния на производительность, всегда используйте
3. PIE: `-buildmode=pie` — минимальное влияние
4. Контроль потока: `GARBLE_EXPERIMENTAL_CONTROLFLOW=1` — экспериментально, но эффективно
5. Антидебаггинг: `antideb` или кастомные проверки
6. Runtime integrity: самопроверка хеша
7. CGO для критичных секций — крайняя мера при максимальных требованиях

---

## Часть 2. Архитектурные подходы для модели OSS + Enterprise Edition

### Моно-репозиторий с директориальным разделением — индустриальный стандарт

Подавляющее большинство успешных Go-проектов с моделью OSS+EE используют **единый моно-репозиторий** с разделением кода по директориям. CockroachDB, Grafana, Teleport и (исторически) Sourcegraph следуют этому паттерну. HashiCorp — исключение, исторически хранивший enterprise-код в приватном репозитории.

**Рекомендуемая структура:**

```
project/
├── cmd/
│   ├── server/          # Основная точка входа
│   └── server-oss/      # OSS-only точка входа
├── pkg/
│   ├── core/            # Apache 2.0 / OSS
│   ├── auth/            # OSS-аутентификация
│   └── ee/              # Enterprise (проприетарная лицензия)
│       ├── sso/         # Enterprise SSO
│       ├── audit/       # Enterprise аудит
│       └── rbac/        # Enterprise RBAC
├── internal/
├── licenses/
│   ├── APACHE2.txt
│   └── ENTERPRISE.txt
└── Makefile             # make build-oss, make build-enterprise
```

CockroachDB реализует это элегантно: все enterprise-пакеты живут в `pkg/ccl/` и используют суффикс `ccl` (например, `changefeedccl`, `backupccl`). Правило аудита простое — ни один Apache 2.0 пакет не должен импортировать `ccl`-пакет.

Преимущества моно-репо: атомарные кросс-модульные изменения, единый CI/CD, простой рефакторинг, один граф зависимостей. Отдельные репозитории создают сложность управления зависимостями, но обеспечивают более жёсткую изоляцию кода.

### Build tags для условной компиляции — идиоматический Go-подход

Go build tags (`//go:build`) — основной механизм условной компиляции enterprise-функциональности. Паттерн «заглушка + реализация»:

**OSS-заглушка** (`sso_oss.go`):
```go
//go:build !enterprise

package auth

func NewSSOProvider() SSOProvider {
    return &noopSSO{}
}

type noopSSO struct{}
func (n *noopSSO) Authenticate(ctx context.Context, token string) (*User, error) {
    return nil, fmt.Errorf("SSO доступен в Enterprise Edition")
}
```

**Enterprise-реализация** (`sso_enterprise.go`):
```go
//go:build enterprise

package auth

func NewSSOProvider() SSOProvider {
    return &samlSSO{} // полная SAML-реализация
}
```

**Сборка:**
```bash
go build -o myapp ./cmd/server                    # OSS
go build -tags enterprise -o myapp-ee ./cmd/server # Enterprise
```

**Продвинутый паттерн — Feature Registry через init():**

```go
// registry.go (OSS)
package features

var DefaultRegistry = &FeatureRegistry{features: make(map[string]Feature)}

func Register(name string, f Feature) {
    DefaultRegistry.features[name] = f
}

// ee_audit.go (Enterprise)
//go:build enterprise

package features

func init() {
    Register("audit-log", &EnterpriseAuditLog{})
}
```

Enterprise-фичи регистрируются при инициализации пакета. В рантайме проверяется одновременно: (a) скомпилирована ли фича и (b) разрешает ли лицензия.

### Как это делают лидеры рынка

| Проект | Разделение кода | Лицензия | Гейтинг фич |
|---|---|---|---|
| **CockroachDB** | Моно-репо, `pkg/ccl/` с суффиксом `ccl` | BSL 1.1 (core) + CCL (enterprise); BSL → Apache 2.0 через 4 года | Рантайм-проверка enterprise-лицензии; `make buildoss` — чистая Apache 2.0 сборка |
| **Teleport** | Моно-репо; `/api` — Apache 2.0, остальное — AGPL v3 | AGPL v3 (исходники); коммерческая лицензия для pre-built CE (v16+: ограничение <100 сотрудников, <$10M ARR) | Enterprise-фичи (SAML/OIDC, FedRAMP) за лицензионным ключом |
| **Grafana** | Моно-репо; enterprise в отдельной директории | AGPL v3 (core); коммерческая (enterprise) | JWT-файл `license.jwt`; автообновление токена каждые 24 часа |
| **HashiCorp** (Vault, Consul) | Публичный + приватный репозиторий | BSL 1.1 (с августа 2023; ранее MPL 2.0) | `.hclic` файлы; проверка при старте; нельзя запустить версию, вышедшую после истечения лицензии |
| **MinIO** | Моно-репо | AGPL v3; коммерческая лицензия для тех, кто не может соблюдать AGPL | Коммерческие фичи через отдельные предложения |
| **Sourcegraph** | Было моно-репо с `enterprise/` | Apache 2.0 → Sourcegraph Enterprise → полностью приватный (2024) | Было директориальное разделение; OSS-сборка исключала `enterprise/` |

### Source-available лицензии: сравнительный анализ

Для модели OSS+EE с closed-source EE-бинарником выбор лицензии критичен для OSS-части.

| Параметр | BSL 1.1 | SSPL v1 | ELv2 | FSL 1.1 | AGPL v3 |
|---|---|---|---|---|---|
| Одобрена OSI | Нет | Нет | Нет | Нет | **Да** |
| Конвертация в OSS | Да (GPLv2+, до 4 лет) | Нет | Нет | Да (Apache/MIT, 2 года) | Уже OSS |
| Защита от облачных провайдеров | Через Additional Use Grant | Очень сильная (весь стек) | Умеренная (запрет managed service) | Умеренная (запрет конкуренции) | Умеренная (copyleft) |
| Сложность для compliance | Высокая (каждая BSL — уникальная) | Средняя | **Низкая** (очень простая) | **Низкая** (стандартизирована) | Средняя |
| Кто использует | HashiCorp, CockroachDB, MariaDB | MongoDB, Redis, Graylog | Elastic | Sentry, Codecov | Grafana, MinIO |

**BSL 1.1** — самая распространённая среди Go-проектов. Ключевая особенность: «Additional Use Grant» — настраиваемое поле, определяющее, что разрешено. У HashiCorp оно запрещает «конкурирующие предложения», у CockroachDB — «database as a service». Каждая BSL фактически уникальна, что усложняет compliance-ревью.

**FSL 1.1** (Functional Source License, Sentry, 2023) — эволюция BSL с двумя важными преимуществами: стандартизирована (нет переменных полей), конвертация через **2 года** (вместо 4), причём в пермиссивную лицензию (Apache/MIT), а не copyleft.

**AGPL v3** — выбор Grafana. Остаётся единственной OSI-одобренной опцией с некоторой защитой. Многие корпоративные юридические отделы **опасаются AGPL**, что парадоксально работает как мягкий коммерческий рычаг — компании, нервничающие из-за copyleft-обязательств, покупают коммерческую лицензию.

**Рекомендация для описанного случая**: **AGPL v3** для OSS-части (если важна community-экосистема) или **FSL 1.1** (если приоритет — бизнес-защита). Для EE-части — проприетарная лицензия.

### Plugin-архитектура: когда и зачем

**Нативный Go plugin (`plugin` package)** — непригоден для продакшена. Только Linux/macOS/FreeBSD, нет Windows, обязательно `CGO_ENABLED=1`, хост и плагин должны быть собраны одной версией Go с идентичными зависимостями, нет выгрузки плагинов.

**HashiCorp go-plugin** (`hashicorp/go-plugin`) — зрелое решение для межпроцессных плагинов через gRPC. Изоляция процессов, кросс-язычность, кросс-платформенность, двунаправленная коммуникация, встроенное логирование. Используется в Terraform, Vault, Consul, Nomad — десятки миллионов пользователей. Однако для модели OSS+EE **build tags превосходят плагины** — нет смысла добавлять runtime-сложность для того, что является compile-time решением. go-plugin оправдан для сторонней расширяемости (как Terraform-провайдеры), а не для внутренних feature tiers.

### Watermarking бинарников для трекинга утечек

Наиболее практичный подход — **встраивание уникальных идентификаторов при сборке** через `-ldflags`:

```bash
go build -ldflags "-X main.customerID=ACME-Corp -X main.buildID=$(uuidgen)" -o myapp-ee
```

Дополнительные техники: встраивание данных через `embed`, аппендинг подписанных метаданных после маркера конца Go-бинарника (`[Binary][Separator][JSON metadata][Signature]`), NOP-sled вариации для уникального отпечатка в коде. Watermarking — механизм **детекции** (отслеживание утечек постфактум), а не **предотвращения**.

---

## Часть 3. Системы лицензирования для Go-приложений

### Обзор платформ лицензирования

| Платформа | Go SDK | Offline-поддержка | Feature flags | Цена (начальная) | Особенности |
|---|---|---|---|---|---|
| **keygen.sh** | Официальный `keygen-go/v2` | Да (signed license files) | Да (entitlements) | $49/мес (1K ALU); CE бесплатно | Fair Source; self-hosted вариант; Ed25519; используют Spotify, Sennheiser |
| **Cryptlex** | `lexactivator-go` (CGO!) | Да | Да (Entitlement Sets) | ~$35/мес | Enterprise-grade; требует CGO |
| **LemonSqueezy** | Нет (REST API) | Ограничено | Нет нативного | Транзакционная комиссия | MoR (Merchant of Record); куплен Shopify; простые сценарии |
| **LicenseSpring** | Нет (REST API) | Отличная | Да | $19/мес | Anti-clock-tampering встроен; нет Go SDK |
| **Собственное решение** | — | Полный контроль | Полный контроль | Время разработки | `hyperboloide/lk` + JWT; Keygen CE |

**Рекомендация**: для широкого рынка — **keygen.sh** (облако или self-hosted CE) как платформа с лучшим Go SDK и полным набором функций. Для максимального контроля — собственное решение на базе Ed25519 + JWT.

### Генерация и валидация лицензионных ключей

**Ed25519 — оптимальный выбор для Go**: компактные подписи (64 байта против 256+ у RSA-2048), быстрая верификация, защита от padding-атак, встроен в стандартную библиотеку Go.

**JWT-подход (используется Grafana Enterprise):**

```go
type LicenseClaims struct {
    jwt.RegisteredClaims
    CustomerID string            `json:"customer_id"`
    Plan       string            `json:"plan"`
    Features   map[string]bool   `json:"features"`
    Limits     map[string]int    `json:"limits"`
    MachineID  string            `json:"machine_id,omitempty"`
}

// Верификация на стороне клиента
func VerifyLicense(tokenString string, publicKey ed25519.PublicKey) (*LicenseClaims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &LicenseClaims{},
        func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodEd25519); !ok {
                return nil, fmt.Errorf("неожиданный метод подписи")
            }
            return publicKey, nil
        })
    if err != nil {
        return nil, err
    }
    claims := token.Claims.(*LicenseClaims)
    if !token.Valid {
        return nil, errors.New("невалидная лицензия")
    }
    return claims, nil
}
```

**Hardware fingerprinting** — библиотека `denisbrodbeck/machineid`: кросс-платформенная, без CGO, использует OS-level machine UUID (стабильнее MAC/BIOS/CPU, особенно в VM). `ProtectedID("com.mycompany.myapp")` возвращает HMAC-SHA256 от machine ID, безопасный для передачи.

```go
import "github.com/denisbrodbeck/machineid"

fingerprint, err := machineid.ProtectedID("com.mycompany.myapp")
```

### Offline vs Online валидация: гибридный подход побеждает

**Offline** (подписанные файлы лицензий): работает без интернета, нет задержек, приватность. Но нельзя отозвать лицензию в реальном времени, уязвимость к манипуляциям с часами.

**Online** (phone-home): мгновенный отзыв, точный учёт, динамические entitlements. Но зависимость от сервера, проблемы приватности, простой сервера = простой клиента.

**Рекомендуемый гибридный паттерн:**

1. При запуске — попытка онлайн-валидации
2. При успехе — сохранение подписанного снимка лицензии на диск
3. При недоступности сервера — проверка подписи кэшированной лицензии
4. Grace period **7–30 дней** оффлайн-работы до обязательного check-in
5. Трёхступенчатая деградация: Silent Fail → Warning → Forced Validation

### Feature flags в лицензии: два подхода

**JSON claims (рекомендуется)** — гибко и расширяемо:
```go
type LicenseFeatures struct {
    Plan     string          `json:"plan"`
    Features map[string]bool `json:"features"`
    Limits   map[string]int  `json:"limits"`
    Modules  []string        `json:"modules"`
}
```

**Bitmask (компактно, для встраиваемых систем):**
```go
const (
    FeatureSSO       uint64 = 1 << 0
    FeatureAudit     uint64 = 1 << 1
    FeatureRBAC      uint64 = 1 << 2
    FeatureAnalytics uint64 = 1 << 3
)

func HasFeature(license, feature uint64) bool {
    return license & feature != 0
}
```

JSON claims предпочтительнее для широкого рынка: проще добавлять новые фичи без изменения формата, человекочитаемо, естественно встраивается в JWT.

### Защита от clock tampering

Комбинация нескольких техник:

1. **Файловая метка** — запись timestamp при каждой проверке лицензии. Если `modTime` файла в будущем — часы откручены назад.
2. **NTP-верификация** — при наличии сети запрос к pool.ntp.org, флаг при расхождении >24 часов.
3. **Монотонные часы** — `time.Now()` в Go содержит монотонный компонент; если wall clock двигается назад, а монотонный продолжает — детекция.
4. **Проверка created_at** — системное время всегда должно быть ≥ `created_at` лицензии (встроено в подпись, нельзя изменить без инвалидации).

Для оффлайн-устройств **абсолютная защита от clock tampering невозможна**. Лучшая стратегия — комбинация файловых меток, проверки created_at и периодического онлайн check-in.

### Graceful degradation при истечении лицензии

Паттерн, используемый HashiCorp, Grafana и Teleport:

```go
type LicenseState int
const (
    Valid           LicenseState = iota // Полный функционал
    GracePeriod                         // Истекла, но в grace period (7-14 дней)
    ExpiredReadOnly                     // Только чтение
    ExpiredDisabled                     // Enterprise-фичи отключены
    Invalid                             // Полная блокировка
)
```

**Критические принципы:**
- **Никогда не уничтожать данные пользователя** при истечении лицензии
- Экспорт/бэкап **всегда доступен** независимо от статуса
- HashiCorp Vault: нельзя запустить/перезапустить/unseal версии, вышедшие после истечения, но запущенные инстансы продолжают работать
- Grafana: short-lived JWT токены обновляются каждые 24 часа; оффлайн-инстансы требуют связи с командой Grafana
- Teleport: Auth Service проверяет лицензию каждый час; кластер продолжает работать при недоступности api.teleport.sh

### Go-библиотеки для лицензирования

**`hyperboloide/lk`** — минималистичная библиотека подписи лицензий на ECDSA (P-384). Включает CLI `lkgen`. Чисто криптографическая — нет feature flags, онлайн-валидации, fingerprinting. Подходит как фундамент для кастомной системы:

```go
privateKey, _ := lk.NewPrivateKey()
license, _ := lk.NewLicense(privateKey, jsonBytes)
encoded, _ := license.ToB64String() // дистрибутировать клиенту

// Верификация
publicKey := privateKey.GetPublicKey()
ok, _ := license.Verify(publicKey)
```

**`keygen-sh/keygen-go/v2`** — полноценный SDK: валидация лицензий, активация машин, проверка entitlements, подпись ответов Ed25519, оффлайн license files. Чистый Go, без CGO:

```go
keygen.Account = "your-account-id"
keygen.Product = "your-product-id"
keygen.LicenseKey = "customer-license-key"
keygen.PublicKey = "your-ed25519-public-key"

license, err := keygen.Validate(fingerprint)
switch {
case err == keygen.ErrLicenseNotActivated:
    // активация
case err == keygen.ErrLicenseExpired:
    // деградация
}
```

**Стандартная библиотека Go** — `crypto/ed25519`, `crypto/ecdsa`, `crypto/rsa` + `golang-jwt/jwt/v5` достаточны для построения полноценной кастомной системы лицензирования.

### Телеметрия и phone-home: правовые рамки

**GDPR (ЕС):** телеметрические данные (IP, MAC, device ID) — персональные данные. Правовые основания: **легитимный интерес** (ст. 6(1)(f)) для проверки лицензии — как правило, законно, но требует balancing test; **согласие** (ст. 6(1)(a)) — более защищённо, но требует opt-in. Microsoft получила претензии от голландского регулятора за расплывчатые цели телеметрии.

**Практические рекомендации:** минимизация данных (только ключ, хеш fingerprint, версия, timestamp), анонимизация идентификаторов, раскрытие в EULA и privacy policy, opt-in для аналитики использования (лицензионная валидация может опираться на легитимный интерес).

### Юридическая защита: DMCA, EULA, ToS

**DMCA § 1201** запрещает обход «технических средств защиты, эффективно контролирующих доступ». Лицензионные ключи, шифрование, подпись кода — всё это квалифицируется как TPM. Создание и распространение инструментов для взлома лицензии **уголовно наказуемо**. Исключения: обратная разработка для интероперабельности (§1201(f)), исследование безопасности.

**Чеклист для EULA:**
- Запрет реверс-инжиниринга (enforceable как контракт)
- Точное определение лицензируемых фич, количества мест, срока действия
- Раскрытие телеметрии
- Условия прекращения лицензии
- Ограничение ответственности
- Применимая юрисдикция
- Явный запрет обхода лицензионной валидации, модификации лицензионных файлов, распространения keygen'ов

---

## Рекомендуемая стратегия для описанного сценария

Для Go 1.26 приложения с моделью OSS (GitHub) + Enterprise (closed-source бинарник), feature-flag лицензированием и широким рынком:

**Архитектура кода:** моно-репозиторий с `pkg/ee/` для enterprise-кода. Build tags (`//go:build enterprise`) с интерфейсным паттерном Feature Registry. CI собирает и тестирует оба варианта.

**Лицензия OSS-части:** AGPL v3 (если важна экосистема и community) или FSL 1.1 (если приоритет — бизнес-защита с конвертацией в Apache 2.0 через 2 года). CLA для всех контрибьюторов. Регистрация товарного знака.

**Сборка EE-бинарника:**
```bash
GARBLE_EXPERIMENTAL_CONTROLFLOW=1 \
garble -literals -tiny -seed=random \
  build -tags enterprise -trimpath -buildmode=pie \
  -ldflags="-s -w -buildid= -X main.customerID=${CUSTOMER} -X main.buildID=$(uuidgen)" \
  -o myapp-ee ./cmd/server
```

**Система лицензирования:** Ed25519-подписанные JWT-файлы с JSON claims для feature flags. Гибридная валидация (online-first + offline cache с grace period 14 дней). Платформа — keygen.sh (cloud или self-hosted CE) или собственная на `crypto/ed25519` + `golang-jwt/jwt/v5` + `denisbrodbeck/machineid`.

**Деградация:** Valid → Grace Period (14 дней, предупреждения) → Read-Only (enterprise-фичи отключены, экспорт доступен) → Invalid (только базовый функционал).

**Юридическая защита:** EULA с anti-circumvention clause, ссылка на DMCA, disclosure телеметрии, CLA, товарный знак.

Ни одна техническая мера не обеспечивает абсолютной защиты. Цель — поднять стоимость обхода выше стоимости лицензии, создав **экономически нецелесообразный** барьер для нелицензионного использования.