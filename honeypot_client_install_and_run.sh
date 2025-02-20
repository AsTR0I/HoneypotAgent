#!/bin/bash

# Вывод приветствия
echo "🚀 Добро пожаловать, я honeypot-client, начинаем установку..."

# Определение ОС
OS=$(uname -s)
ARCH=$(uname -m)

echo "🔍 Определена ОС: $OS"
echo "🔍 Определена архитектура: $ARCH"

# Определение ссылки на бинарник в зависимости от ОС и архитектуры
if [ "$OS" == "Linux" ]; then
    case "$ARCH" in
        "x86_64")
            BIN_URL="https://github.com/AsTR0I/HoneypotAgent/blob/main/honeypot-client/bin/linux/amd64/HoneypotAgent?raw=true"
            ;;
        "i386")
            BIN_URL="https://github.com/AsTR0I/HoneypotAgent/blob/main/honeypot-client/bin/linux/386/HoneypotAgent?raw=true"
            ;;
        *)
            echo "❌ Неизвестная архитектура для Linux: $ARCH. Скрипт завершает работу."
            exit 1
            ;;
    esac
elif [ "$OS" == "FreeBSD" ]; then
    case "$ARCH" in
        "x86_64")
            BIN_URL="https://github.com/AsTR0I/HoneypotAgent/blob/main/honeypot-client/bin/freebsd/amd64/HoneypotAgent?raw=true"
            ;;
        "i386")
            BIN_URL="https://github.com/AsTR0I/HoneypotAgent/blob/main/honeypot-client/bin/freebsd/386/HoneypotAgent?raw=true"
            ;;
        *)
            echo "❌ Неизвестная архитектура для FreeBSD: $ARCH. Скрипт завершает работу."
            exit 1
            ;;
    esac
elif [ "$OS" == "Darwin" ]; then
    echo "❌ macOS не поддерживается. Скрипт завершает работу."
    exit 1
elif [ "$OS" == "Windows" ]; then
    case "$ARCH" in
        "x86_64")
            BIN_URL="https://github.com/AsTR0I/HoneypotAgent/blob/main/honeypot-client/bin/windows/amd64/HoneypotAgent?raw=true"
            ;;
        *)
            echo "❌ Неизвестная архитектура для Windows: $ARCH. Скрипт завершает работу."
            exit 1
            ;;
    esac
else
    echo "❌ Неизвестная операционная система: $OS. Скрипт завершает работу."
    exit 1
fi

# Печать ссылки на бинарник
echo "🔗 Ссылка для загрузки: $BIN_URL"

# Пример скачивания файла
echo "📥 Скачиваем файл..."
curl -L -o HoneypotAgent "$BIN_URL"

# Проверка скачивания
if [ $? -eq 0 ]; then
    echo "✅ Бинарник успешно скачан."
else
    echo "❌ Ошибка при скачивании бинарника."
    exit 1
fi