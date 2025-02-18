package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/viper"
)

const (
	Version = "2025.02.14"
	RN      = "\r\n"
	R       = "\r"
	N       = "\n"
	LogFile = "honeypot_agent.log"
	LogDir  = "./logs"
)

type Config struct {
	ListenPorts struct {
		TCP []int `mapstructure:"tcp"`
		UDP []int `mapstructure:"udp"`
	} `mapstructure:"listenports"`
	Trigger struct {
		URL string `json:"url"`
	} `json:"trigger"`
}

// Загрузка конфигурации
func loadConfig() (*Config, error) {
	viper.SetConfigName("HoneypotAgent") // Имя конфигурационного файла без расширения
	viper.SetConfigType("json")          // Указываем расширения
	viper.AddConfigPath(".")             // Ищем в текущей директории

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("Ошибка при чтении конфигурации: %v", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("Ошибка при распаковке конфигурации: %v", err)
	}

	return &config, nil
}

// Функция для отправки триггера
func sendTrigger(config *Config, address, source string) {
	hostname, err := os.Hostname()
	if err != nil {
		logToFile(fmt.Sprintf("Ошибка при получении имени хоста: %v", err))
		return
	}

	// только IP (без порта)
	addressParts := strings.Split(address, ":")
	ipAddress := addressParts[0]

	if config.Trigger.URL == "" {
		logToFile("Ошибка: URL в конфигурации пустой")
		return
	}
	url := config.Trigger.URL

	payload := map[string]string{
		"address": ipAddress,
		"source":  "agent:%"+hostname+"%"
	}

	// Сериализация payload в JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		logToFile(fmt.Sprintf("Ошибка при сериализации JSON: %v", err))
		return
	}
	logToFile(fmt.Sprintf("JSON Payload:", string(jsonPayload))) // Логируем JSON перед отправкой
	// Логируем отправку запроса
	logToFile(fmt.Sprintf("Отправка HTTP-запроса на адрес %s с address=%s и hostname=%s", url, ipAddress, hostname))

	// Отправка HTTP POST запроса
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewReader(jsonPayload))
	if err != nil {
		logToFile(fmt.Sprintf("Ошибка при создании HTTP запроса: %v", err))
		return
	}
	// Устанавливаем заголовок Content-Type как application/json
	req.Header.Set("Content-Type", "application/json")

	// Отправляем запрос и получаем ответ
	resp, err := client.Do(req)
	if err != nil {
		logToFile(fmt.Sprintf("Ошибка при отправке HTTP запроса: %v", err))
		return
	}
	defer resp.Body.Close()

	// Логируем результат
	if resp.StatusCode == http.StatusOK {
		logToFile(fmt.Sprintf("Успешный ответ от сервера: %v", resp.Status))
	} else {
		logToFile(fmt.Sprintf("Ошибка при ответе от сервера: %v", resp.Status))
	}
}

// Слушаем TCP порт
func listenTCP(port int, config *Config) {
	address := fmt.Sprintf(":%d", port)
	ln, err := net.Listen("tcp", address)

	if err != nil {
		logToFile(fmt.Sprintf("Не удаётся слушать TCP порт %d: %v", port, err))
		return
	}

	defer ln.Close()

	logToFile(fmt.Sprintf("Слушаем TCP порт: %d", port))
	for {
		conn, err := ln.Accept()
		if err != nil {
			logToFile(fmt.Sprintf("Ошибка при принятии соединения на порту %d: %v", port, err))
			continue
		}

		go safeHandleRequest(conn, "tcp", port, config)
	}
}

// Слушаем UDP порт
func listenUDP(port int, config *Config) {
	address := fmt.Sprintf(":%d", port)
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		log.Printf("Ошибка при разрешении UDP адреса: %v\n", err)
		return
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Printf("Ошибка при прослушивании UDP порта %d: %v\n", port, err)
		return
	}
	defer conn.Close()

	logToFile(fmt.Sprintf("Слушаем UDP порт: %d", port))

	for {
		buffer := make([]byte, 1024)
		_, remoteAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Printf("Ошибка при чтении с порта %d: %v\n", port, err)
			continue
		}

		go safeHandleRequest(nil, "udp", port, config, remoteAddr.String())
	}
}

// Обработка запроса
func handleRequest(conn net.Conn, protocol string, port int, config *Config, addr ...string) {
	source := fmt.Sprintf("%s:%d", protocol, port)
	if len(addr) > 0 {
		source = addr[0]
	}

	if conn == nil {
		logToFile("Ошибка: соединение отсутствует")
		return
	}

	logToFile(fmt.Sprintf("Обработан запрос с адреса %s", conn.RemoteAddr()))

	// Отправка триггера с конфигом
	sendTrigger(config, conn.RemoteAddr().String(), source)

	if err := conn.Close(); err != nil {
		logToFile(fmt.Sprintf("Ошибка при закрытии соединения: %v", err))
	}
}

// Обработчик с безопасной горутиной (программа не завершает работу)
func safeHandleRequest(conn net.Conn, protocol string, port int, config *Config, addr ...string) {
	defer func() {
		if r := recover(); r != nil {
			logToFile(fmt.Sprintf("Паника в горутине: %v", r))
		}
	}()
	handleRequest(conn, protocol, port, config, addr...)
}

// Запись в файл логов
func logToFile(message string) {
	// Получаем текущую дату
	currentDate := time.Now()

	// Создаем лог-файл, если он не существует
	if _, err := os.Stat(LogDir); os.IsNotExist(err) {
		err := os.Mkdir(LogDir, 0755)
		if err != nil {
			fmt.Printf("Ошибка при создании директории логов: %v\n", err)
			return
		}
	}

	logFilePath := filepath.Join(LogDir, LogFile)
	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Ошибка при открытии лог-файла: %v\n", err)
		return
	}
	defer file.Close()

	// Удаляем старые логи
	removeOldLogs()

	// Создаем лог-запись
	logger := log.New(file, "", 0)
	logger.SetFlags(0)
	logger.Println(currentDate.Format("2006-01-02 15:04:05") + " - " + message)
}

// Удаление старых логов, старше 7 дней
func removeOldLogs() {
	files, err := os.ReadDir(LogDir)
	if err != nil {
		fmt.Printf("Ошибка при чтении директории логов: %v\n", err)
		return
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		// Получаем время последней модификации файла
		fileInfo, err := file.Info()
		if err != nil {
			continue
		}

		// Проверяем, если файл старше 7 дней, то удаляем его
		if time.Since(fileInfo.ModTime()).Hours() > 7*24 {
			err := os.Remove(filepath.Join(LogDir, file.Name()))
			if err != nil {
				fmt.Printf("Ошибка при удалении старого лог-файла %v: %v\n", file.Name(), err)
			} else {
				fmt.Printf("Удален старый лог-файл: %v\n", file.Name())
			}
		}
	}
}

func main() {
	// Выводим информацию о версии и платформе
	fmt.Printf("Honeypot Agent Version: %s\n", Version)
	fmt.Printf("Platform: %s/%s\n", runtime.GOOS, runtime.GOARCH)

	// Загружаем конфигурацию
	config, err := loadConfig()
	if err != nil {
		logToFile(fmt.Sprintf("Ошибка загрузки конфигурации: %v", err))
		return
	}

	// Слушаем TCP порты
	for _, port := range config.ListenPorts.TCP {
		go listenTCP(port, config)
	}

	// Слушаем UDP порты
	for _, port := range config.ListenPorts.UDP {
		go listenUDP(port, config)
	}

	select {} // Ожидаем завершения горутин
}
