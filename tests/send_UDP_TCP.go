package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

// Функция для отправки данных через UDP
func sendUDP(address string, port int, message string) {
	// Формируем адрес для подключения
	addr := fmt.Sprintf("%s:%d", address, port)

	// Разбираем адрес в формат UDP
	remoteAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		log.Fatalf("Ошибка при разрешении адреса: %v\n", err)
		os.Exit(1)
	}

	// Получаем UDP-соединение
	conn, err := net.DialUDP("udp", nil, remoteAddr)
	if err != nil {
		log.Fatalf("Ошибка при подключении к UDP серверу: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	// Отправляем сообщение
	_, err = conn.Write([]byte(message))
	if err != nil {
		log.Fatalf("Ошибка при отправке данных: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Данные успешно отправлены на сервер через UDP!")
}

// Функция для отправки данных через TCP
func sendTCP(address string, port int, message string) {
	// Формируем адрес для подключения
	addr := fmt.Sprintf("%s:%d", address, port)

	// Разбираем адрес для TCP
	remoteAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		log.Fatalf("Ошибка при разрешении адреса: %v\n", err)
		os.Exit(1)
	}

	// Получаем TCP-соединение
	conn, err := net.DialTCP("tcp", nil, remoteAddr)
	if err != nil {
		log.Fatalf("Ошибка при подключении к TCP серверу: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	// Отправляем сообщение
	_, err = conn.Write([]byte(message))
	if err != nil {
		log.Fatalf("Ошибка при отправке данных: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Данные успешно отправлены на сервер через TCP!")
}

func main() {
	// Указываем адрес и порт для отправки данных
	address := "127.0.0.1"
	udpPort := 12345
	tcpPort := 54321
	message := "Hello, Server!"

	// Отправляем данные через UDP
	sendUDP(address, udpPort, message)

	// Отправляем данные через TCP
	sendTCP(address, tcpPort, message)
}
