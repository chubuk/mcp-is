#!/usr/bin/env bash
# deploy.sh — comprehensive deployment script for MCP-IS
#
# Usage:
#   ./scripts/deploy.sh <command> [options]
#
# Commands:
#   install   Full installation on clean server
#   update    Update existing installation
#   check     Pre-flight checks before deploy
#   rollback  Rollback to previous version
#   status    Check current status
#   backup    Manual backup
#
# Options:
#   --quiet, -q      Suppress non-essential output
#   --verbose, -v    Enable verbose/debug output
#   --no-color       Disable colored output
#   --yes, -y        Answer yes to all prompts
#   --dry-run        Show what would be done without executing
#
# Examples:
#   ./scripts/deploy.sh install --yes
#   ./scripts/deploy.sh update --verbose
#   ./scripts/deploy.sh check
#   ./scripts/deploy.sh backup --dry-run

set -euo pipefail

# ============================================================================
# Конфигурация по умолчанию
# ============================================================================

# Имя проекта и основные пути
readonly PROJECT_NAME="mcp-is"
readonly PROJECT_DIR="${MCP_IS_DIR:-/opt/${PROJECT_NAME}}"
readonly COMPOSE_FILE="docker-compose.yaml"
readonly BACKUP_DIR="${BACKUP_DIR:-/opt/${PROJECT_NAME}/backups}"
readonly LOG_DIR="${LOG_DIR:-/var/log/${PROJECT_NAME}}"
readonly ENV_FILE=".env"
readonly ENV_EXAMPLE=".env.example"
readonly DOCKER_IMAGES=(
    "ghcr.io/chubuk/mcp-is:latest"
    "ghcr.io/chubuk/mcp-is-nginx:latest"
    "postgres:15-alpine"
)

# Порты по умолчанию
readonly DEFAULT_HTTP_PORT=80
readonly DEFAULT_HTTPS_PORT=443
readonly DEFAULT_APP_PORT=8080

# Пороговые значения
readonly MIN_DISK_SPACE_MB=1024      # 1 GB минимум
readonly MIN_MEMORY_MB=512           # 512 MB минимум
readonly MAX_LOAD_THRESHOLD=4.0      # Максимальная загрузка перед бэкапом
readonly BACKUP_RETENTION_DAYS=7     # Хранение бэкапов

# Версия скрипта
readonly SCRIPT_VERSION="1.0.0"

# ============================================================================
# Флаги и опции
# ============================================================================

COMMAND=""
QUIET=false
VERBOSE=false
NO_COLOR=false
AUTO_YES=false
DRY_RUN=false
FORCE=false

# Цвета (только для TTY)
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Проверка TTY и --no-color
if [[ ! -t 1 ]] || [[ "${NO_COLOR}" == "true" ]]; then
    RED=""
    GREEN=""
    YELLOW=""
    BLUE=""
    NC=""
fi

# ============================================================================
# Утилиты вывода
# ============================================================================

# Форматированный вывод
print_status() {
    local level="$1"
    local message="$2"
    local timestamp
    timestamp=$(date '+%Y-%m-%d %H:%M:%S')

    case "$level" in
        OK)
            [[ "$QUIET" == "true" ]] && return
            echo -e "${GREEN}[OK]${NC}    ${timestamp} ${message}"
            ;;
        WARN)
            echo -e "${YELLOW}[WARN]${NC}  ${timestamp} ${message}"
            ;;
        FAIL)
            echo -e "${RED}[FAIL]${NC}  ${timestamp} ${message}"
            ;;
        INFO)
            [[ "$QUIET" == "true" ]] && return
            echo -e "${BLUE}[INFO]${NC}  ${timestamp} ${message}"
            ;;
        DEBUG)
            [[ "$VERBOSE" == "true" ]] && echo -e "[DEBUG] ${timestamp} ${message}"
            ;;
        STEP)
            [[ "$QUIET" == "true" ]] && return
            echo ""
            echo -e "${BLUE}==>${NC} ${message}"
            ;;
    esac
}

# Вывод ошибки и выход
die() {
    local message="$1"
    local exit_code="${2:-1}"
    print_status "FAIL" "$message"
    exit "$exit_code"
}

# Запрос подтверждения
confirm() {
    local message="${1:-Continue?}"
    if [[ "$AUTO_YES" == "true" ]]; then
        return 0
    fi
    if [[ "$DRY_RUN" == "true" ]]; then
        print_status "INFO" "Dry run: skipping confirmation"
        return 1
    fi
    read -rp "$(echo -e "${YELLOW}${message} [y/N]: ${NC}")" response
    case "$response" in
        [yY][eE][sS]|[yY]) return 0 ;;
        *) return 1 ;;
    esac
}

# Выполнение команды с dry-run поддержкой
run_cmd() {
    if [[ "$DRY_RUN" == "true" ]]; then
        print_status "INFO" "Dry run: $*"
        return 0
    fi
    if [[ "$VERBOSE" == "true" ]]; then
        print_status "DEBUG" "Running: $*"
        "$@"
    else
        "$@" 2>/dev/null
    fi
}

# ============================================================================
# Проверки системы
# ============================================================================

# Проверка совместимости ОС и архитектуры
check_os_compatibility() {
    print_status "STEP" "Проверка совместимости ОС..."

    # Определение ОС
    local os_name=""
    local os_version=""
    local os_family=""

    if [[ -f /etc/os-release ]]; then
        # shellcheck disable=SC1091
        source /etc/os-release
        os_name="${ID:-unknown}"
        os_version="${VERSION_ID:-unknown}"
        os_family="${ID_LIKE:-$ID}"
    elif [[ -f /etc/redhat-release ]]; then
        os_name="centos"
        os_version=$(grep -oE '[0-9]+' /etc/redhat-release | head -1)
        os_family="rhel"
    else
        die "Не удалось определить операционную систему"
    fi

    print_status "DEBUG" "ОС: ${os_name} ${os_version} (семейство: ${os_family})"

    # Проверка поддерживаемых ОС
    case "$os_family" in
        debian|ubuntu)
            print_status "OK" "ОС поддерживается: ${os_name} ${os_version}"
            ;;
        rhel|fedora|centos)
            print_status "OK" "ОС поддерживается: ${os_name} ${os_version}"
            ;;
        *)
            die "ОС '${os_name}' (${os_family}) не поддерживается. Поддерживаются: Ubuntu/Debian/CentOS/Fedora"
            ;;
    esac

    # Проверка архитектуры
    local arch
    arch=$(uname -m)
    case "$arch" in
        x86_64|amd64)
            print_status "OK" "Архитектура: amd64"
            ;;
        aarch64|arm64)
            print_status "OK" "Архитектура: arm64"
            ;;
        armv7l|armhf)
            print_status "OK" "Архитектура: armv7"
            ;;
        *)
            die "Архитектура '${arch}' не поддерживается. Поддерживаются: amd64, arm64, armv7"
            ;;
    esac

    # Проверка доступных ресурсов
    local total_memory_mb
    total_memory_mb=$(awk '/MemTotal/ {printf "%d", $2/1024}' /proc/meminfo)

    if [[ "$total_memory_mb" -lt "$MIN_MEMORY_MB" ]]; then
        die "Недостаточно памяти: ${total_memory_mb} MB (минимум ${MIN_MEMORY_MB} MB)"
    fi
    print_status "OK" "Память: ${total_memory_mb} MB (минимум ${MIN_MEMORY_MB} MB)"

    local cpu_cores
    cpu_cores=$(nproc 2>/dev/null || echo 1)
    if [[ "$cpu_cores" -lt 1 ]]; then
        die "Недостаточно CPU: ${cpu_cores} ядер (минимум 1)"
    fi
    print_status "OK" "CPU: ${cpu_cores} ядер"

    return 0
}

# Проверка Docker и Docker Compose
check_docker() {
    print_status "STEP" "Проверка Docker..."

    # Проверка Docker
    if ! command -v docker &>/dev/null; then
        die "Docker не установлен. Установите: https://docs.docker.com/engine/install/"
    fi

    local docker_version
    docker_version=$(docker version --format '{{.Server.Version}}' 2>/dev/null || echo "unknown")
    print_status "OK" "Docker версия: ${docker_version}"

    # Проверка Docker daemon
    if ! docker info &>/dev/null; then
        die "Docker daemon не запущен или нет прав доступа. Попробуйте: sudo systemctl start docker"
    fi
    print_status "OK" "Docker daemon запущен"

    # Проверка Docker Compose
    local compose_cmd=""
    if docker compose version &>/dev/null; then
        compose_cmd="docker compose"
    elif command -v docker-compose &>/dev/null; then
        compose_cmd="docker-compose"
    else
        die "Docker Compose не установлен. Установите: https://docs.docker.com/compose/install/"
    fi

    local compose_version
    compose_version=$($compose_cmd version --short 2>/dev/null || echo "unknown")
    print_status "OK" "Docker Compose: ${compose_version}"

    # Экспорт для использования в других функциях
    export COMPOSE_CMD="$compose_cmd"

    return 0
}

# Проверка доступности портов
check_ports() {
    print_status "STEP" "Проверка портов..."

    local ports_to_check=("${DEFAULT_HTTP_PORT}" "${DEFAULT_HTTPS_PORT}")

    for port in "${ports_to_check[@]}"; do
        if ss -tlnp 2>/dev/null | grep -q ":${port} "; then
            local process
            process=$(ss -tlnp 2>/dev/null | grep ":${port} " | awk '{print $6}' | head -1)
            print_status "WARN" "Порт ${port} уже занят: ${process}"
        else
            print_status "OK" "Порт ${port} свободен"
        fi
    done

    return 0
}

# Проверка DNS резолвинга
check_dns() {
    local domain="${1:-}"
    local server_ip
    server_ip=$(curl -s ifconfig.me 2>/dev/null || curl -s ipinfo.io/ip 2>/dev/null || echo "")

    if [[ -z "$domain" ]]; then
        print_status "INFO" "Домен не указан, пропуск проверки DNS"
        return 0
    fi

    print_status "STEP" "Проверка DNS для ${domain}..."

    if [[ -z "$server_ip" ]]; then
        print_status "WARN" "Не удалось определить внешний IP сервера"
        return 0
    fi

    local domain_ip
    domain_ip=$(dig +short "$domain" 2>/dev/null | head -1 || nslookup "$domain" 2>/dev/null | grep 'Address:' | tail -1 | awk '{print $2}' || echo "")

    if [[ -z "$domain_ip" ]]; then
        print_status "WARN" "Не удалось разрезолвить домен ${domain}"
        return 0
    fi

    if [[ "$domain_ip" == "$server_ip" ]]; then
        print_status "OK" "DNS: ${domain} -> ${domain_ip} (совпадает с сервером)"
    else
        print_status "WARN" "DNS: ${domain} -> ${domain_ip} (сервер: ${server_ip})"
    fi

    return 0
}

# Проверка дискового пространства
check_disk() {
    print_status "STEP" "Проверка дискового пространства..."

    local check_path="${1:-$PROJECT_DIR}"
    local available_mb
    available_mb=$(df -m "$check_path" 2>/dev/null | awk 'NR==2 {print $4}' || echo "0")

    if [[ "$available_mb" -lt "$MIN_DISK_SPACE_MB" ]]; then
        die "Недостаточно места на диске: ${available_mb} MB (минимум ${MIN_DISK_SPACE_MB} MB)"
    fi
    print_status "OK" "Дисковое пространство: ${available_mb} MB доступно"

    # Проверка inodes
    local available_inodes
    available_inodes=$(df -i "$check_path" 2>/dev/null | awk 'NR==2 {print $4}' || echo "0")
    if [[ "$available_inodes" -lt 1000 ]]; then
        print_status "WARN" "Мало свободных inodes: ${available_inodes}"
    fi

    return 0
}

# Проверка состояния контейнеров
check_containers() {
    print_status "STEP" "Проверка контейнеров..."

    local all_healthy=true

    for service in db app nginx; do
        local container_name="mcp_is_${service}"
        local status
        status=$(docker inspect --format='{{.State.Status}}' "$container_name" 2>/dev/null || echo "not_found")

        case "$status" in
            running)
                local health
                health=$(docker inspect --format='{{if .State.Health}}{{.State.Health.Status}}{{else}}no_healthcheck{{end}}' "$container_name" 2>/dev/null || echo "unknown")
                case "$health" in
                    healthy)
                        print_status "OK" "Контейнер ${service}: запущен (здоров)"
                        ;;
                    unhealthy)
                        print_status "WARN" "Контейнер ${service}: запущен (нездоров)"
                        all_healthy=false
                        ;;
                    no_healthcheck)
                        print_status "OK" "Контейнер ${service}: запущен (без healthcheck)"
                        ;;
                    *)
                        print_status "WARN" "Контейнер ${service}: запущен (статус: ${health})"
                        all_healthy=false
                        ;;
                esac
                ;;
            exited|dead)
                print_status "FAIL" "Контейнер ${service}: остановлен (${status})"
                all_healthy=false
                ;;
            not_found)
                print_status "INFO" "Контейнер ${service}: не найден"
                ;;
            *)
                print_status "WARN" "Контейнер ${service}: неизвестный статус (${status})"
                all_healthy=false
                ;;
        esac
    done

    if [[ "$all_healthy" == "false" ]]; then
        return 1
    fi
    return 0
}

# Проверка SSL сертификата
check_ssl() {
    local domain="${1:-}"

    if [[ -z "$domain" ]] || [[ "$domain" == "localhost" ]]; then
        print_status "INFO" "SSL проверка пропущена (домен не задан)"
        return 0
    fi

    print_status "STEP" "Проверка SSL сертификата..."

    local cert_path="/etc/letsencrypt/live/${domain}/fullchain.pem"

    if [[ ! -f "$cert_path" ]]; then
        print_status "WARN" "SSL сертификат не найден для ${domain}"
        return 0
    fi

    local expiry_date
    expiry_date=$(openssl x509 -in "$cert_path" -noout -enddate 2>/dev/null | cut -d= -f2 || echo "")

    if [[ -z "$expiry_date" ]]; then
        print_status "WARN" "Не удалось прочитать SSL сертификат"
        return 0
    fi

    local expiry_epoch
    expiry_epoch=$(date -d "$expiry_date" +%s 2>/dev/null || echo "0")
    local current_epoch
    current_epoch=$(date +%s)
    local days_left=$(( (expiry_epoch - current_epoch) / 86400 ))

    if [[ "$days_left" -lt 0 ]]; then
        print_status "FAIL" "SSL сертификат истёк ${days_left#-} дней назад"
        return 1
    elif [[ "$days_left" -lt 30 ]]; then
        print_status "WARN" "SSL сертификат истекает через ${days_left} дней"
    else
        print_status "OK" "SSL сертификат действителен (${days_left} дней)"
    fi

    return 0
}

# Проверка последнего бэкапа
check_backup() {
    print_status "STEP" "Проверка бэкапов..."

    if [[ ! -d "$BACKUP_DIR" ]]; then
        print_status "WARN" "Директория бэкапов не найдена: ${BACKUP_DIR}"
        return 0
    fi

    local backup_count
    backup_count=$(find "$BACKUP_DIR" -name "db_*.sql.gz" 2>/dev/null | wc -l)

    if [[ "$backup_count" -eq 0 ]]; then
        print_status "WARN" "Бэкапы не найдены"
        return 0
    fi

    local last_backup
    last_backup=$(ls -t "$BACKUP_DIR"/db_*.sql.gz 2>/dev/null | head -1 || echo "")

    if [[ -n "$last_backup" ]]; then
        local last_backup_name
        last_backup_name=$(basename "$last_backup")
        local last_backup_date
        last_backup_date=$(stat -c %Y "$last_backup" 2>/dev/null || echo "0")
        local current_date
        current_date=$(date +%s)
        local hours_since=$(( (current_date - last_backup_date) / 3600 ))

        print_status "OK" "Последний бэкап: ${last_backup_name} (${hours_since}ч назад)"
        print_status "INFO" "Всего бэкапов: ${backup_count}"
    fi

    return 0
}

# ============================================================================
# Генерация конфигурации
# ============================================================================

# Генерация случайного пароля
generate_password() {
    local length="${1:-32}"
    openssl rand -base64 "$length" | tr -d '/+=' | head -c "$length"
}

# Генерация .env файла
generate_env() {
    print_status "STEP" "Генерация .env файла..."

    local env_path="${PROJECT_DIR}/${ENV_FILE}"

    # Проверка существования .env.example
    local example_path="${PROJECT_DIR}/${ENV_EXAMPLE}"
    if [[ ! -f "$example_path" ]]; then
        die "${ENV_EXAMPLE} не найден в ${PROJECT_DIR}"
    fi

    # Если .env уже существует, делаем резервную копию
    if [[ -f "$env_path" ]]; then
        local backup_name="${env_path}.backup.$(date +%Y%m%d_%H%M%S)"
        run_cmd cp "$env_path" "$backup_name"
        print_status "INFO" "Существующий .env сохранён как $(basename "$backup_name")"
    fi

    # Генерация случайных паролей
    local db_password
    db_password=$(generate_password 24)
    local jwt_secret
    jwt_secret=$(generate_password 32)
    local encryption_key
    encryption_key=$(generate_password 32)

    # Создание .env из шаблона
    if [[ "$DRY_RUN" == "true" ]]; then
        print_status "INFO" "Dry run: будет создан .env с:"
        print_status "INFO" "  DB_PASSWORD=${db_password}"
        print_status "INFO" "  JWT_SECRET=${jwt_secret}"
        print_status "INFO" "  ENCRYPTION_KEY=${encryption_key}"
        return 0
    fi

    # Копируем шаблон
    run_cmd cp "$example_path" "$env_path"

    # Замена значений
    run_cmd sed -i "s|^DB_PASSWORD=.*|DB_PASSWORD=${db_password}|" "$env_path"
    run_cmd sed -i "s|^JWT_SECRET=.*|JWT_SECRET=${jwt_secret}|" "$env_path"
    run_cmd sed -i "s|^ENCRYPTION_KEY=.*|ENCRYPTION_KEY=${encryption_key}|" "$env_path"

    # Установка прав доступа
    run_cmd chmod 600 "$env_path"

    print_status "OK" "Файл .env создан с безопасными паролями"
    print_status "INFO" "DB_PASSWORD: ${db_password}"
    print_status "INFO" "JWT_SECRET: ${jwt_secret}"
    print_status "INFO" "ENCRYPTION_KEY: ${encryption_key}"

    return 0
}

# ============================================================================
# SSL и Let's Encrypt
# ============================================================================

# Настройка Let's Encrypt
setup_letsencrypt() {
    local domain="${1:-}"
    local email="${2:-}"

    if [[ -z "$domain" ]] || [[ "$domain" == "localhost" ]]; then
        print_status "INFO" "Let's Encrypt пропущен (домен не задан)"
        return 0
    fi

    if [[ -z "$email" ]]; then
        print_status "WARN" "Email для Let's Encrypt не задан"
        return 0
    fi

    print_status "STEP" "Настройка Let's Encrypt..."

    # Проверка DNS перед выдачей сертификата
    if ! check_dns "$domain"; then
        print_status "WARN" "DNS не настроен, пропуск Let's Encrypt"
        return 0
    fi

    # Запуск certbot через Docker
    if [[ "$DRY_RUN" == "true" ]]; then
        print_status "INFO" "Dry run: будет запущен certbot для ${domain}"
        return 0
    fi

    # Получение сертификата через certbot
    docker run --rm \
        -v "${PROJECT_DIR}/certbot/conf:/etc/letsencrypt" \
        -v "${PROJECT_DIR}/certbot/www:/var/www/certbot" \
        certbot/certbot certonly \
        --webroot \
        --webroot-path=/var/www/certbot \
        --email "$email" \
        --agree-tos \
        --no-eff-email \
        -d "$domain" \
        --force-renewal 2>/dev/null || {
            print_status "WARN" "Не удалось получить сертификат. Проверьте DNS настройки."
            return 0
        }

    print_status "OK" "SSL сертификат получен для ${domain}"

    # Настройка автоматического обновления
    setup_certbot_cron "$domain"

    return 0
}

# Настройка cron для обновления сертификатов
setup_certbot_cron() {
    local domain="${1:-}"

    local cron_cmd="0 3 * * 1 docker run --rm -v ${PROJECT_DIR}/certbot/conf:/etc/letsencrypt -v ${PROJECT_DIR}/certbot/www:/var/www/certbot certbot/certbot renew --quiet --post-hook 'docker exec mcp_is_nginx nginx -s reload' 2>&1 | logger -t certbot"

    # Проверяем, нет ли уже такой задачи
    if crontab -l 2>/dev/null | grep -qF "certbot"; then
        print_status "INFO" "Cron для certbot уже настроен"
    else
        if confirm "Добавить cron для обновления SSL?"; then
            (crontab -l 2>/dev/null; echo "$cron_cmd") | crontab -
            print_status "OK" "Cron для обновления SSL добавлен"
        fi
    fi
}

# ============================================================================
# Управление бэкапами
# ============================================================================

# Настройка cron для бэкапов
setup_cron() {
    print_status "STEP" "Настройка cron для бэкапов..."

    local cron_cmd="0 2 * * * ${PROJECT_DIR}/scripts/backup.sh ${BACKUP_DIR} 2>&1 | logger -t mcp-is-backup"

    # Проверяем, нет ли уже такой задачи
    if crontab -l 2>/dev/null | grep -qF "mcp-is-backup"; then
        print_status "INFO" "Cron для бэкапов уже настроен"
        return 0
    fi

    if confirm "Добавить cron для ежедневных бэкапов (02:00)?"; then
        if [[ "$DRY_RUN" == "true" ]]; then
            print_status "INFO" "Dry run: будет добавлен cron: ${cron_cmd}"
            return 0
        fi

        (crontab -l 2>/dev/null; echo "$cron_cmd") | crontab -
        print_status "OK" "Cron для бэкапов добавлен"
    fi

    return 0
}

# Загрузка конфигурации из .env
load_backup_config() {
    local env_file="${PROJECT_DIR}/${ENV_FILE}"

    if [[ -f "$env_file" ]]; then
        # shellcheck disable=SC1090
        set -a
        # shellcheck disable=SC1091
        source "$env_file"
        set +a
        print_status "DEBUG" "Конфигурация загружена из ${env_file}"
    else
        print_status "WARN" "Файл .env не найден, используются значения по умолчанию"
    fi

    # Установка значений по умолчанию
    export DB_PASSWORD="${DB_PASSWORD:-password}"
    export SERVER_PORT="${SERVER_PORT:-8080}"
    export DOMAIN="${DOMAIN:-localhost}"
    export LETSENCRYPT_EMAIL="${LETSENCRYPT_EMAIL:-}"

    return 0
}

# ============================================================================
# Команды
# ============================================================================

# Предварительная проверка
cmd_check() {
    print_status "STEP" "Предварительная проверка системы"

    local has_errors=false

    check_os_compatibility || has_errors=true
    check_docker || has_errors=true
    check_ports || has_errors=true
    check_disk "${PROJECT_DIR}" || has_errors=true
    check_containers || has_errors=true

    # Загрузка конфигурации для проверки DNS и SSL
    load_backup_config
    check_dns "${DOMAIN}" || has_errors=true
    check_ssl "${DOMAIN}" || has_errors=true
    check_backup || has_errors=true

    echo ""
    if [[ "$has_errors" == "true" ]]; then
        print_status "WARN" "Обнаружены проблемы. Проверьте вывод выше."
        return 1
    else
        print_status "OK" "Все проверки пройдены успешно"
        return 0
    fi
}

# Полная установка
cmd_install() {
    print_status "STEP" "Полная установка ${PROJECT_NAME}"

    # Создание директорий
    print_status "INFO" "Создание директорий..."
    run_cmd mkdir -p "${PROJECT_DIR}"
    run_cmd mkdir -p "${BACKUP_DIR}"
    run_cmd mkdir -p "${LOG_DIR}"

    # Проверки
    cmd_check || die "Предварительные проверки не пройдены"

    # Клонирование/обновление репозитория
    if [[ -d "${PROJECT_DIR}/.git" ]]; then
        print_status "INFO" "Репозиторий уже существует, пропуск клонирования"
    else
        if [[ "$DRY_RUN" == "true" ]]; then
            print_status "INFO" "Dry run: будет клонирован репозиторий"
        else
            print_status "INFO" "Клонирование репозитория..."
            git clone https://github.com/chubuk/mcp-is.git "${PROJECT_DIR}"
        fi
    fi

    # Генерация .env
    generate_env

    # Загрузка переменных
    load_backup_config

    # Сборка и запуск контейнеров
    print_status "STEP" "Сборка и запуск контейнеров..."

    if [[ "$DRY_RUN" == "true" ]]; then
        print_status "INFO" "Dry run: будет выполнен 'docker compose up -d --build'"
    else
        cd "${PROJECT_DIR}"
        $COMPOSE_CMD up -d --build
    fi

    # Ожидание готовности
    wait_for_containers

    # Настройка Let's Encrypt
    setup_letsencrypt "${DOMAIN}" "${LETSENCRYPT_EMAIL}"

    # Настройка cron для бэкапов
    setup_cron

    # Проверка состояния
    check_containers

    echo ""
    print_status "OK" "Установка завершена успешно!"
    print_status "INFO" "Приложение доступно на: http://localhost:${SERVER_PORT}"
    if [[ "${DOMAIN:-}" != "localhost" ]]; then
        print_status "INFO" "Домен: https://${DOMAIN}"
    fi
}

# Обновление
cmd_update() {
    print_status "STEP" "Обновление ${PROJECT_NAME}"

    # Проверка наличия установки
    if [[ ! -d "${PROJECT_DIR}/.git" ]]; then
        die "Проект не найден в ${PROJECT_DIR}. Выполните сначала: install"
    fi

    # Бэкап перед обновлением
    if confirm "Создать бэкап перед обновлением?"; then
        cmd_backup || print_status "WARN" "Бэкап не удался, продолжаем обновление"
    fi

    # Сохранение текущей версии
    local old_version
    old_version=$(git -C "${PROJECT_DIR}" rev-parse HEAD 2>/dev/null || echo "unknown")
    print_status "INFO" "Текущая версия: ${old_version:0:8}"

    # Обновление кода
    print_status "INFO" "Получение обновлений..."
    run_cmd git -C "${PROJECT_DIR}" fetch origin
    run_cmd git -C "${PROJECT_DIR}" reset --hard origin/main

    local new_version
    new_version=$(git -C "${PROJECT_DIR}" rev-parse HEAD 2>/dev/null || echo "unknown")
    print_status "INFO" "Новая версия: ${new_version:0:8}"

    # Проверка изменений в docker-compose.yaml
    if git -C "${PROJECT_DIR}" diff "$old_version" "$new_version" -- docker-compose.yaml &>/dev/null; then
        print_status "INFO" "docker-compose.yaml не изменился"
    else
        print_status "INFO" "docker-compose.yaml обновлён"
    fi

    # Пересборка и запуск
    print_status "STEP" "Пересборка и запуск контейнеров..."

    cd "${PROJECT_DIR}"
    run_cmd $COMPOSE_CMD up -d --build --force-recreate

    # Ожидание готовности
    wait_for_containers

    # Проверка состояния
    check_containers

    echo ""
    print_status "OK" "Обновление завершено успешно!"
}

# Откат
cmd_rollback() {
    print_status "STEP" "Откат к предыдущей версии"

    # Проверка наличия установки
    if [[ ! -d "${PROJECT_DIR}/.git" ]]; then
        die "Проект не найден в ${PROJECT_DIR}"
    fi

    # Показ доступных версий
    print_status "INFO" "Последние 10 коммитов:"
    git -C "${PROJECT_DIR}" log --oneline -10

    # Запрос версии для отката
    local target_commit="${1:-}"
    if [[ -z "$target_commit" ]]; then
        target_commit=$(git -C "${PROJECT_DIR}" rev-parse HEAD~1 2>/dev/null || echo "")
        if [[ -z "$target_commit" ]]; then
            die "Не удалось определить предыдущую версию"
        fi
        print_status "INFO" "Откат к: ${target_commit:0:8}"
    fi

    # Подтверждение
    if ! confirm "Откатить к ${target_commit:0:8}?"; then
        print_status "INFO" "Откат отменён"
        return 0
    fi

    # Бэкап перед откатом
    cmd_backup || true

    # Выполнение отката
    run_cmd git -C "${PROJECT_DIR}" reset --hard "$target_commit"

    # Пересборка и запуск
    print_status "STEP" "Пересборка после отката..."
    cd "${PROJECT_DIR}"
    run_cmd $COMPOSE_CMD up -d --build --force-recreate

    # Ожидание готовности
    wait_for_containers

    # Проверка состояния
    check_containers

    echo ""
    print_status "OK" "Откат завершён успешно!"
}

# Проверка статуса
cmd_status() {
    print_status "STEP" "Статус ${PROJECT_NAME}"

    # Информация о версии
    if [[ -d "${PROJECT_DIR}/.git" ]]; then
        local version
        version=$(git -C "${PROJECT_DIR}" describe --tags --always 2>/dev/null || echo "unknown")
        local commit
        commit=$(git -C "${PROJECT_DIR}" rev-parse HEAD 2>/dev/null || echo "unknown")
        print_status "INFO" "Версия: ${version} (${commit:0:8})"
    else
        print_status "WARN" "Репозиторий не найден"
    fi

    # Проверка контейнеров
    check_containers

    # Проверка использования ресурсов
    echo ""
    print_status "STEP" "Использование ресурсов"
    docker stats --no-stream --format "table {{.Name}}\t{{.CPUPerc}}\t{{.MemUsage}}\t{{.NetIO}}" 2>/dev/null || true

    # Проверка бэкапов
    check_backup

    # Проверка SSL
    load_backup_config
    check_ssl "${DOMAIN}"

    return 0
}

# Ручной бэкап
cmd_backup() {
    print_status "STEP" "Создание бэкапа"

    # Проверка нагрузки перед бэкапом
    local load_1min
    load_1min=$(awk '{print $1}' /proc/loadavg)
    local load_threshold="${MAX_LOAD_THRESHOLD}"

    if awk "BEGIN {exit !($load_1min > $load_threshold)}" 2>/dev/null; then
        print_status "WARN" "Высокая загрузка системы: ${load_1min} (порог: ${load_threshold})"
        if ! confirm "Продолжить бэкап при высокой нагрузке?"; then
            return 1
        fi
    else
        print_status "OK" "Загрузка системы: ${load_1min}"
    fi

    # Проверка контейнера БД
    if ! docker ps --format '{{.Names}}' | grep -q "mcp_is_db"; then
        die "Контейнер БД не запущен"
    fi

    # Создание директории для бэкапов
    mkdir -p "${BACKUP_DIR}"

    local date
    date=$(date +%Y%m%d_%H%M%S)
    local backup_file="${BACKUP_DIR}/db_${date}.sql.gz"

    if [[ "$DRY_RUN" == "true" ]]; then
        print_status "INFO" "Dry run: будет создан бэкап ${backup_file}"
        return 0
    fi

    print_status "INFO" "Создание дампа базы данных..."
    docker exec mcp_is_db pg_dump -U admin mcp_is | gzip > "$backup_file"

    local size
    size=$(du -h "$backup_file" | cut -f1)
    print_status "OK" "Бэкап создан: ${backup_file} (${size})"

    # Очистка старых бэкапов
    print_status "INFO" "Очистка бэкапов старше ${BACKUP_RETENTION_DAYS} дней..."
    local deleted
    deleted=$(find "$BACKUP_DIR" -name "db_*.sql.gz" -mtime +${BACKUP_RETENTION_DAYS} -delete -print 2>/dev/null | wc -l)
    print_status "INFO" "Удалено ${deleted} старых бэкапов"

    local remaining
    remaining=$(find "$BACKUP_DIR" -name "db_*.sql.gz" | wc -l)
    print_status "INFO" "Всего бэкапов: ${remaining}"

    return 0
}

# Ожидание готовности контейнеров
wait_for_containers() {
    if [[ "$DRY_RUN" == "true" ]]; then
        print_status "INFO" "Dry run: ожидание контейнеров пропущено"
        return 0
    fi

    print_status "INFO" "Ожидание готовности контейнеров..."
    local timeout=120
    local interval=5
    local elapsed=0

    while [[ $elapsed -lt $timeout ]]; do
        local all_healthy=true

        for service in db app nginx; do
            local container="mcp_is_${service}"
            local status
            status=$(docker inspect --format='{{.State.Status}}' "$container" 2>/dev/null || echo "not_found")

            if [[ "$status" != "running" ]]; then
                all_healthy=false
                break
            fi

            # Проверка healthcheck (если есть)
            local health
            health=$(docker inspect --format='{{if .State.Health}}{{.State.Health.Status}}{{else}}no_healthcheck{{end}}' "$container" 2>/dev/null || echo "unknown")

            if [[ "$health" == "unhealthy" ]]; then
                all_healthy=false
                break
            fi
        done

        if [[ "$all_healthy" == "true" ]]; then
            print_status "OK" "Все контейнеры готовы"
            return 0
        fi

        sleep "$interval"
        elapsed=$((elapsed + interval))
        print_status "DEBUG" "Ожидание... (${elapsed}s/${timeout}s)"
    done

    print_status "WARN" "Таймаут ожидания контейнеров"
    return 0
}

# ============================================================================
# Парсинг аргументов
# ============================================================================

usage() {
    cat <<EOF
${BLUE}${PROJECT_NAME} deploy script v${SCRIPT_VERSION}${NC}

Usage:
  $(basename "$0") <command> [options]

Commands:
  install   Full installation on clean server
  update    Update existing installation
  check     Pre-flight checks before deploy
  rollback  Rollback to previous version (optional: commit hash)
  status    Check current status
  backup    Manual backup

Options:
  --quiet, -q      Suppress non-essential output
  --verbose, -v    Enable verbose/debug output
  --no-color       Disable colored output
  --yes, -y        Answer yes to all prompts
  --dry-run        Show what would be done without executing
  --force, -f      Force operation without confirmations

Examples:
  $(basename "$0") install --yes
  $(basename "$0") update --verbose
  $(basename "$0") check
  $(basename "$0") rollback abc123
  $(basename "$0") backup --dry-run

Environment variables:
  MCP_IS_DIR       Project directory (default: /opt/mcp-is)
  BACKUP_DIR       Backup directory (default: /opt/mcp-is/backups)
  LOG_DIR          Log directory (default: /var/log/mcp-is)
EOF
}

parse_args() {
    while [[ $# -gt 0 ]]; do
        case "$1" in
            install|update|check|rollback|status|backup)
                COMMAND="$1"
                shift
                ;;
            --quiet|-q)
                QUIET=true
                shift
                ;;
            --verbose|-v)
                VERBOSE=true
                shift
                ;;
            --no-color)
                NO_COLOR=true
                RED=""
                GREEN=""
                YELLOW=""
                BLUE=""
                NC=""
                shift
                ;;
            --yes|-y)
                AUTO_YES=true
                shift
                ;;
            --dry-run)
                DRY_RUN=true
                shift
                ;;
            --force|-f)
                FORCE=true
                shift
                ;;
            --help|-h)
                usage
                exit 0
                ;;
            --version|-V)
                echo "${PROJECT_NAME} deploy script v${SCRIPT_VERSION}"
                exit 0
                ;;
            *)
                echo "Неизвестная опция: $1"
                usage
                exit 1
                ;;
        esac
    done

    if [[ -z "$COMMAND" ]]; then
        echo "Не указана команда"
        echo ""
        usage
        exit 1
    fi
}

# ============================================================================
# Основная функция
# ============================================================================

main() {
    parse_args "$@"

    # Проверка root для некоторых операций
    if [[ $EUID -ne 0 ]]; then
        case "$COMMAND" in
            install|rollback)
                if [[ "$DRY_RUN" != "true" ]]; then
                    print_status "WARN" "Для установки/отката рекомендуется root привилегии"
                fi
                ;;
        esac
    fi

    # Загрузка конфигурации
    if [[ -f "${PROJECT_DIR}/${ENV_FILE}" ]]; then
        load_backup_config
    fi

    # Выполнение команды
    case "$COMMAND" in
        install)
            cmd_install
            ;;
        update)
            cmd_update
            ;;
        check)
            cmd_check
            ;;
        rollback)
            # Передаём аргумент после флагов как коммит для отката
            # (parse_args уже снял флаги, теперь ищем оставшийся аргумент)
            local commit_hash=""
            # Ищем коммит в оставшихся позиционных аргументах
            for arg in "$@"; do
                if [[ "$arg" != "--"* ]] && [[ "$arg" != "-"* ]]; then
                    commit_hash="$arg"
                    break
                fi
            done
            if [[ -n "$commit_hash" ]]; then
                cmd_rollback "$commit_hash"
            else
                cmd_rollback
            fi
            ;;
        status)
            cmd_status
            ;;
        backup)
            cmd_backup
            ;;
        *)
            echo "Неизвестная команда: ${COMMAND}"
            usage
            exit 1
            ;;
    esac
}

# ============================================================================
# Точка входа
# ============================================================================

main "$@"
