
# HoneypotAgent

## Описание

**HoneypotAgent** – это агент для мониторинга сетевой активности, который прослушивает заданные TCP и UDP порты и отправляет уведомления о входящих соединениях на API.

## Функциональность

- Читает конфигурацию из JSON-файла `HoneypotAgent.json`.
- Прослушивает указанные в конфигурации TCP и UDP порты.
- Фиксирует входящие соединения и отправляет HTTP-запрос на сервер триггеров.
- Ведёт логирование событий в файл `logs/honeypot_agent.log`.
- Автоматически удаляет старые логи (старше 7 дней).

## Формат конфигурации (JSON)

```json
{
  "listenports": {
    "tcp": [80, 443],
    "udp": [53, 161]
  },
  "trigger": {
    "url": "http://example.com/trigger"
  }
}
```

- **listenports.tcp** – список TCP-портов для прослушивания.
- **listenports.udp** – список UDP-портов для прослушивания.
- **trigger.url** – API, куда отправляются уведомления о входящих соединениях.

## Логирование

- Логи записываются в `logs/honeypot_agent.log`.
- Формат записи: `YYYY-MM-DD HH:MM:SS - Сообщение`.
- Если лог-файл отсутствует, он создаётся автоматически.
- Логи старше 7 дней удаляются автоматически.

## Запуск

1. Перемещаемся в нужную директорию:

   ```bash
   mkdir /root/honeypotagent/ && cd /root/honeypotagent/
   ```

2. Скачиваем программу:

   - Для Linux:
     - 64-бит: `wget https://github.com/AsTR0I/HoneypotAgent/blob/main/honeypot-client/builds/HoneypotAgentLinuxamd64.tar.gz`
     - 32-бит: `wget https://github.com/AsTR0I/HoneypotAgent/blob/main/honeypot-client/builds/HoneypotAgentLinux386.tar.gz`
   - Для FreeBSD:
     - 64-бит: `wget https://github.com/AsTR0I/HoneypotAgent/blob/main/honeypot-client/builds/HoneypotAgentLinuxamd64.tar.gz`
     - 32-бит: `wget https://github.com/AsTR0I/HoneypotAgent/blob/main/honeypot-client/builds/HoneypotAgentFreebsd386.tar.gz`
   - Для Windows:
     - 64-бит: `wget https://github.com/AsTR0I/HoneypotAgent/blob/main/honeypot-client/builds/HoneypotAgentWindowsamd64.tar.gz`
     - 32-бит: `wget https://github.com/AsTR0I/HoneypotAgent/blob/main/honeypot-client/builds/HoneypotAgentWindows386.tar.gz`

3. Рапаковываем архив:

   ```bash
   tar -xzf HoneypotAgentLinuxamd64.tar.gz
   ```

4. Создаём и редактируем конфигурационный файл в формате JSON:

   Пример конфигурации:

   ```json
   {
     "listenports": {
       "tcp": [21, 22, 54321],
       "udp": [5, 6, 7]
     },
     "trigger": {
       "url": "http://honey.cocobri.ru:8088/add-host?token=95a62fbd-76e3-46f2-b454-12d22679916f"
     }
   }
   ```

5. Запускаем агент:

   ```bash
   ./HoneypotAgentLinuxamd64
   ```
