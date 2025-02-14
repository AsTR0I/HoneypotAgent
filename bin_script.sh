#!/bin/bash

# Директории
SRC_DIR="./honeypot-client/src"
BIN_DIR="./honeypot-client/bin"
BUILD_DIR="./honeypot-client/builds"
ARCHITECTURES=("amd64" "386")  # Архитектуры для компиляции

# Платформы для компиляции
PLATFORMS=("linux" "windows" "freebsd")

# Создаем директории, если их нет
mkdir -p $BIN_DIR
mkdir -p $BUILD_DIR

# Компиляция и архивирование
for PLATFORM in "${PLATFORMS[@]}"; do
  for ARCH in "${ARCHITECTURES[@]}"; do
    # Компиляция
    BIN_NAME="HoneypotAgent${PLATFORM^}${ARCH}"
    OUTPUT_PATH="${BIN_DIR}/${BIN_NAME}"

    echo "Компиляция для ${PLATFORM}/${ARCH}..."
    GOOS=$PLATFORM GOARCH=$ARCH go build -o $OUTPUT_PATH $SRC_DIR

    if [ $? -eq 0 ]; then
      echo "Компиляция завершена: ${OUTPUT_PATH}"

      # Архивирование в build
      BUILD_ARCHIVE="${BUILD_DIR}/${BIN_NAME}.tar.gz"
      echo "Архивирование в ${BUILD_ARCHIVE}..."
      tar -czvf $BUILD_ARCHIVE -C $BIN_DIR $BIN_NAME
      echo "Архив ${BUILD_ARCHIVE} создан."
    else
      echo "Ошибка компиляции для ${PLATFORM}/${ARCH}"
    fi
  done
done

echo "Процесс завершен."
