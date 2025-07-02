package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"mig/address"
	"mig/brand"
	"mig/user"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки файла .env: %v", err)
	}

	if len(os.Args) < 2 {
		log.Fatalf("Не указан аргумент. Используйте:  'user', 'password', 'brand', 'address' или 'validate'")
	}

	arg := os.Args[1]

	switch arg {
	case "user":
		log.Println("Запуск миграции пользователей")
		user.Migrate()
	case "password":
		log.Println("Запуск миграции паролей пользователей")
		user.MigratePasswords()
	case "address":
		log.Println("Запуск миграции адресов")
		address.Migrate()
	case "validate":
		log.Println("Запуск валидации адресов")
		address.Validate()
	case "brand":
		log.Println("Запуск миграции brand")
		brand.Migrate()
	default:
		log.Fatalf("Неизвестный аргумент: %s. Используйте: 'user', 'password', 'brand', 'address' или 'validate'", arg)
	}

	log.Println("Миграция завершена успешно")
}
