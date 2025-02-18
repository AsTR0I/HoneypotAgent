HoneypotAgent

Описание

HoneypotAgent – это агент для мониторинга сетевой активности, который прослушивает заданные TCP и UDP порты и отправляет уведомления о входящих соединениях на api.

Функциональность

Читает конфигурацию из JSON-файла HoneypotAgent.json.
Прослушивает указанные в конфигурации TCP и UDP порты.
Фиксирует входящие соединения и отправляет HTTP-запрос на сервер триггеров.
Ведёт логирование событий в файл logs/honeypot_agent.log.
Автоматически удаляет старые логи (старше 7 дней).

Формат конфигурации (json)

{
  "listenports": {
    "tcp": [80, 443],
    "udp": [53, 161]
  },
  "trigger": {
    "url": "http://example.com/trigger"
  }
}

listenports.tcp – список TCP-портов для прослушивания.
listenports.udp – список UDP-портов для прослушивания.
trigger.url – api, куда отправляются уведомления о входящих соединениях.

Логирование

Логи записываются в logs/honeypot_agent.log.
Формат записи: YYYY-MM-DD HH:MM:SS - Сообщение
Если лог-файл отсутствует, он создаётся автоматически.
Логи старше 7 дней удаляются автоматически

Запуск

Перемещаемся в нужную директорию: mkdir /root/honeypotagent/ && cd /root/honeypotagent/

Скачиваем программу:
    linux:
        64: wget https://github.com/AsTR0I/HoneypotAgent/blob/main/honeypot-client/builds/HoneypotAgentLinuxamd64.tar.gz
        32: wget https://github.com/AsTR0I/HoneypotAgent/blob/main/honeypot-client/builds/HoneypotAgentLinux386.tar.gz
    freeBSD:
        64: wget https://github.com/AsTR0I/HoneypotAgent/blob/main/honeypot-client/builds/HoneypotAgentLinuxamd64.tar.gz
        32: wget https://github.com/AsTR0I/HoneypotAgent/blob/main/honeypot-client/builds/HoneypotAgentFreebsd386.tar.gz
    win:
        64: https://github.com/AsTR0I/HoneypotAgent/blob/main/honeypot-client/builds/HoneypotAgentWindowsamd64.tar.gz
        32: https://github.com/AsTR0I/HoneypotAgent/blob/main/honeypot-client/builds/HoneypotAgentWindows386.tar.gz

Рапаковываем: tar -xzf HoneypotAgentLinuxamd64.tar.gz

Создаём и редактируем config в json формате:
    Пример конфига:
        {
            "listenports": {
                "tcp": [21, 22, 54321],
                "udp": [5, 6, 7]
            },
            "trigger": {
                "url": "http://honey.cocobri.ru:8088/add-host?token=95a62fbd-76e3-46f2-b454-12d22679916f"
            }
        }
        
запускаем ./HoneypotAgentLinuxamd64