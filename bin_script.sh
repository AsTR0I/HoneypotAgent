#!/bin/bash

# Директории
SRC_DIR="./honeypot-client/src"
BIN_DIR="./honeypot-client/bin"
BUILD_DIR="./honeypot-client/builds"

# Архитектуры для компиляции (FreeBSD 386 убираем)
ARCHITECTURES=("amd64" "386")

# Платформы для компиляции (Windows, Linux, FreeBSD)
PLATFORMS=("linux" "windows" "freebsd")

# Исключаем FreeBSD 386, так как Go не поддерживает его
EXCLUDE_FREEBSD_386=true

# Компиляция и архивирование
for PLATFORM in "${PLATFORMS[@]}"; do
  for ARCH in "${ARCHITECTURES[@]}"; do
    # Пропускаем FreeBSD 386, если нужно
    if [[ "$PLATFORM" == "freebsd" && "$ARCH" == "386" && "$EXCLUDE_FREEBSD_386" == true ]]; then
      echo "⚠️ Пропускаем компиляцию FreeBSD 386 (не поддерживается)"
      continue
    fi

    # Создаём директории, если их нет
    BIN_PATH="${BIN_DIR}/${PLATFORM}/${ARCH}"
    BUILD_PATH="${BUILD_DIR}/${PLATFORM}/${ARCH}"
    mkdir -p "$BIN_PATH" "$BUILD_PATH"

    # Определяем имя бинарника
    BIN_NAME="HoneypotAgent"
    [[ "$PLATFORM" == "windows" ]] && BIN_NAME="HoneypotAgent.exe"

    OUTPUT_PATH="${BIN_PATH}/${BIN_NAME}"

    echo "🚀 Компиляция для ${PLATFORM}/${ARCH}..."
    
    # Исправление флага ldflags, теперь это должно быть частью параметров go build
    GOOS=$PLATFORM GOARCH=$ARCH go build -ldflags "-s -w" -o "$OUTPUT_PATH" "$SRC_DIR"

    if [ $? -eq 0 ]; then
      echo "✅ Компиляция завершена: ${OUTPUT_PATH}"

      # Архивирование
      if [[ "$PLATFORM" == "windows" ]]; then
        BUILD_ARCHIVE="${BUILD_PATH}/HoneypotAgent.zip"
        echo "📦 Архивирование в ${BUILD_ARCHIVE}..."
        zip -j "$BUILD_ARCHIVE" "$OUTPUT_PATH"
      else
        BUILD_ARCHIVE="${BUILD_PATH}/HoneypotAgent.tar.gz"
        echo "📦 Архивирование в ${BUILD_ARCHIVE}..."
        tar -czvf "$BUILD_ARCHIVE" -C "$BIN_PATH" "$BIN_NAME"
      fi

      echo "🎉 Архив ${BUILD_ARCHIVE} создан."
    else
      echo "❌ Ошибка компиляции для ${PLATFORM}/${ARCH}"
    fi
  done
done

echo "✅ Процесс завершен."
