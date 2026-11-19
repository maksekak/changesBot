# Руководство по установке ChangesBot через systemd (консоль)

Этот гайд описывает, как развернуть ChangesBot как системную службу на Linux с помощью `systemd`, используя только терминал.

## 1. Требования

- Linux с `systemd` (Ubuntu, Debian, CentOS, RHEL и т.д.)
- Скомпилированный бинарный файл бота для Linux
- Права `sudo` (root)

## 2. Подготовка на сервере

Выполняйте все команды ниже на целевом Linux-сервере.

### 2.1. Создание пользователя для бота

```bash
sudo useradd -r -m -s /bin/bash changesbot
```

### 2.2. Создание директории и копирование файлов

```bash
sudo mkdir -p /opt/changesBot
sudo cp changesBotBuild /opt/changesBot/
sudo cp .env /opt/changesBot/
sudo cp groups.json /opt/changesBot/   # если используется
sudo cp week_state.json /opt/changesBot/ # если используется
sudo cp *.xlsx /opt/changesBot/        # файлы расписания
sudo chown -R changesbot:changesbot /opt/changesBot
sudo chmod +x /opt/changesBot/changesBotBuild
```

> **Важно:** бинарный файл должен быть скомпилирован под Linux.  
> Пример сборки на Windows для Linux:
> ```bash
> GOOS=linux GOARCH=amd64 go build -o changesBotBuild .
> ```

### 2.3. Проверка .env

Убедитесь, что файл `/opt/changesBot/.env` существует и содержит токен:

```bash
cat /opt/changesBot/.env
```

Ожидаемый формат:
```
token=ВАШ_ТОКЕН_БОТА
```

## 3. Установка systemd-сервиса

### 3.1. Копирование unit-файла

```bash
sudo cp changesbot.service /etc/systemd/system/changesbot.service
```

### 3.2. Перечитать конфигурацию systemd

```bash
sudo systemctl daemon-reload
```

### 3.3. Включить автозапуск и запустить сервис

```bash
sudo systemctl enable changesbot
sudo systemctl start changesbot
```

### 3.4. Проверить статус

```bash
sudo systemctl status changesbot
```

Пример успешного вывода:
```
● changesbot.service - ChangesBot — Telegram бот расписания колледжа
     Loaded: loaded (/etc/systemd/system/changesbot.service; enabled; vendor preset: enabled)
     Active: active (running) since ...
```

## 4. Полезные команды управления

| Команда | Описание |
|---------|----------|
| `sudo systemctl start changesbot` | Запустить бота |
| `sudo systemctl stop changesbot` | Остановить бота |
| `sudo systemctl restart changesbot` | Перезапустить бота |
| `sudo systemctl enable changesbot` | Включить автозапуск при старте системы |
| `sudo systemctl disable changesbot` | Выключить автозапуск |
| `sudo systemctl status changesbot` | Показать текущий статус |
| `journalctl -u changesbot -f` | Посмотреть логи в реальном времени |
| `journalctl -u changesbot --since "1 hour ago"` | Логи за последний час |

## 5. Удаление сервиса

```bash
sudo systemctl stop changesbot
sudo systemctl disable changesbot
sudo rm /etc/systemd/system/changesbot.service
sudo systemctl daemon-reload
sudo systemctl reset-failed
```

## 6. Решение проблем

### Бот не запускается
1. Проверьте логи: `journalctl -u changesbot -n 50 --no-pager`
2. Убедитесь, что `.env` читается и содержит верный токен
3. Проверьте права на бинарный файл: `ls -l /opt/changesBot/changesBotBuild`

### Ошибка прав доступа к файлам
```bash
sudo chown -R changesbot:changesbot /opt/changesBot
```
