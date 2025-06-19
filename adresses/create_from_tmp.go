package address

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

// Mock структура для демонстрации AddressManagement
type AddressManagement struct{}

func (m *AddressManagement) AddShippingAddress(ctx context.Context, customerId int, address *ShippingAddress) error {
	// Имитация вызова метода для работы с адресами доставки
	fmt.Printf("Shipping address added for customer %d: %+v\n", customerId, address)
	return nil // Вернуть nil для успешного выполнения
}

func (m *AddressManagement) AddBillingAddress(ctx context.Context, customerId int, address *BillingAddress) error {
	// Имитация вызова метода для работы с биллинговыми адресами
	fmt.Printf("Billing address added for customer %d: %+v\n", customerId, address)
	return nil // Вернуть nil для успешного выполнения
}

// Структуры для адресов
type ShippingAddress struct {
	PostalCode      *string
	CountryCode     string
	SubdivisionCode *string
	SubdivisionName *string
	CityName        *string
	AddressLine1    string
	AddressLine2    *string
	Fullname        string
}

type BillingAddress struct {
	PostalCode      *string
	CountryCode     string
	SubdivisionCode *string
	SubdivisionName *string
	CityName        *string
	AddressLine1    string
	AddressLine2    *string
	Fullname        string
}

func main() {
	// Получить строку подключения из переменной окружения
	dsn := os.Getenv("NEW_DL_DSN")
	if dsn == "" {
		log.Fatal("Переменная окружения NEW_DL_DSN не задана")
	}

	// Подключение к базе данных
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	defer db.Close()

	ctx := context.Background()
	addressManagement := &AddressManagement{} // Ваш экземпляр AddressManagement

	// Запрос для получения адресов со статусом 'wait'
	rows, err := db.QueryContext(ctx, `
		SELECT id, customer_email, postal_code, country_code, subdivision_code, 
			   subdivision_name, city_name, address_line1, address_line2, firstname, lastname, type 
		FROM customers_addresses_tmp
		WHERE status = 'wait'
	`)
	if err != nil {
		log.Fatalf("Ошибка запроса к базе данных: %v", err)
	}
	defer rows.Close()

	// Обработка каждой строки
	for rows.Next() {
		var (
			id              int
			customerEmail   string
			countryCode     string
			addressLine1    string
			addressType     string
			postalCode      sql.NullString
			subdivisionCode sql.NullString
			subdivisionName sql.NullString
			cityName        sql.NullString
			addressLine2    sql.NullString
			firstname       sql.NullString
			lastname        sql.NullString
		)

		// Чтение одной строки
		if err := rows.Scan(&id, &customerEmail, &postalCode, &countryCode, &subdivisionCode,
			&subdivisionName, &cityName, &addressLine1, &addressLine2, &firstname, &lastname, &addressType); err != nil {
			log.Printf("Ошибка чтения строки: %v", err)
			continue
		}

		// Получение customer_id из таблицы customers
		var customerId int
		err := db.QueryRowContext(ctx, `
			SELECT id 
			FROM customers 
			WHERE email = $1`, customerEmail).Scan(&customerId)
		if errors.Is(err, sql.ErrNoRows) {
			// Обновить статус и комментарий об ошибке
			updateStatus(ctx, db, id, "error", "Клиент не найден")
			continue
		} else if err != nil {
			log.Printf("Ошибка получения customer_id: %v", err)
			updateStatus(ctx, db, id, "error", err.Error())
			continue
		}

		fullname := strings.TrimSpace(fmt.Sprintf("%s %s", firstname.String, lastname.String))
		var addressError error

		// Добавление адреса в зависимости от типа
		if addressType == "shipping" {
			shippingAddress := &ShippingAddress{
				PostalCode:      nullStringPointer(postalCode),
				CountryCode:     countryCode,
				SubdivisionCode: nullStringPointer(subdivisionCode),
				SubdivisionName: nullStringPointer(subdivisionName),
				CityName:        nullStringPointer(cityName),
				AddressLine1:    addressLine1,
				AddressLine2:    nullStringPointer(addressLine2),
				Fullname:        fullname,
			}
			addressError = addressManagement.AddShippingAddress(ctx, customerId, shippingAddress)
		} else if addressType == "billing" {
			billingAddress := &BillingAddress{
				PostalCode:      nullStringPointer(postalCode),
				CountryCode:     countryCode,
				SubdivisionCode: nullStringPointer(subdivisionCode),
				SubdivisionName: nullStringPointer(subdivisionName),
				CityName:        nullStringPointer(cityName),
				AddressLine1:    addressLine1,
				AddressLine2:    nullStringPointer(addressLine2),
				Fullname:        fullname,
			}
			addressError = addressManagement.AddBillingAddress(ctx, customerId, billingAddress)
		} else {
			addressError = fmt.Errorf("неизвестный тип адреса: %s", addressType)
		}

		// Обработка результата добавления
		if addressError != nil {
			updateStatus(ctx, db, id, "error", addressError.Error())
			continue
		}

		// Успешно добавлено
		updateStatus(ctx, db, id, "done", "")
	}

	// Проверка ошибок завершения обработчика строк
	if err = rows.Err(); err != nil {
		log.Fatalf("Ошибка обработки строк: %v", err)
	}
}

// Функция для обновления статуса записи
func updateStatus(ctx context.Context, db *sql.DB, id int, status, comment string) {
	_, err := db.ExecContext(ctx, `
		UPDATE customers_addresses_tmp
		SET status = $1, status_comment = $2
		WHERE id = $3
	`, status, comment, id)
	if err != nil {
		log.Printf("Ошибка обновления статуса для id %d: %v", id, err)
	}
}

// Помощник для работы с sql.NullString
func nullStringPointer(ns sql.NullString) *string {
	if ns.Valid {
		return &ns.String
	}
	return nil
}
