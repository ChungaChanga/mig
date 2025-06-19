package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	address "mig/adresses"
	"mig/users"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки файла .env: %v", err)
	}

	if len(os.Args) < 2 {
		log.Fatalf("Не указан аргумент. Используйте: 'users' или 'address'")
	}

	arg := os.Args[1]

	switch arg {
	case "users":
		log.Println("Запуск миграции пользователей")
		users.Migrate()
	case "address":
		log.Println("Запуск миграции адресов")
		address.Migrate()
	default:
		log.Fatalf("Неизвестный аргумент: %s. Используйте: 'users' или 'address'", arg)
	}

	log.Println("Миграция завершена успешно")
}
