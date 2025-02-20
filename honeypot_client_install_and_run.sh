#!/bin/bash

# Вывод приветствия
echo "🚀 Добро пожаловать, я honeypot-client, начинаем установку..."

# Определение ОС
OS=$(uname -s)
ARCH=$(uname -m)

echo "🔍 Определена ОС: $OS"
echo "🔍 Определена архитектура: $ARCH"

# Определение ссылки на ZIP-файл в зависимости от ОС и архитектуры
if [ "$OS" == "Linux" ]; then
    case "$ARCH" in
        "x86_64")
            ZIP_URL="https://github.com/AsTR0I/HoneypotAgent/blob/main/honeypot-client/builds/HoneypotAgentLinux386.tar.gz"
            ;;
        "aarch64")
            ZIP_URL="https://github.com/AsTR0I/HoneypotAgent/blob/main/honeypot-client/builds/HoneypotAgentLinuxamd64.tar.gz"
        ;;
    *)
        echo "❌ Неизвестная архитектура: $ARCH. Скрипт завершает работу."
            exit 1
            ;;
    esac