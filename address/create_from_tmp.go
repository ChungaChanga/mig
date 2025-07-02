package address

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	"google.golang.org/grpc"
	"mig/adapters"
	"mig/address/model"
	addressProto "mig/api/ausweis/proto/address"
)

func Validate() {
	ausweisSocket := os.Getenv("AUSWEIS_SOCKET")
	if ausweisSocket == "" {
		log.Fatal("Не указана переменная окружения AUSWEIS_SOCKET")
	}
	conn, err := grpc.Dial(ausweisSocket, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// Создаем реальный клиент, сгенерированный из .proto-файла
	client := addressProto.NewAddressServiceClient(conn)
	addressService := adapters.NewAddressService(client)
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

	// Запрос для получения адресов со статусом 'wait'
	rows, err := db.QueryContext(ctx, `
		SELECT id, customer_id, postal_code, country_code, subdivision_code, 
			   subdivision_name, city_name, address_line1, address_line2, firstname, lastname, residential, liftgate, type 
		FROM customers_addresses_tmp_diff1
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
			customerId      int
			countryCode     string
			addressLine1    string
			addressType     string
			postalCode      sql.NullString
			subdivisionCode sql.NullString
			subdivisionName sql.NullString
			cityName        sql.NullString
			addressLine2    sql.NullString
			firstname       string
			lastname        string
			fullname        string
			residential     bool
			liftgate        bool
		)

		// Чтение одной строки
		if err := rows.Scan(&id, &customerId, &postalCode, &countryCode, &subdivisionCode,
			&subdivisionName, &cityName, &addressLine1, &addressLine2,
			&firstname, &lastname, &residential, &liftgate, &addressType); err != nil {
			log.Printf("Ошибка чтения строки: %v", err)
			continue
		}
		fullname = firstname + " " + lastname
		// Получение customer_id из таблицы customers

		var addrId int
		var addressError error

		// Добавление адреса в зависимости от типа
		if addressType == "shipping" || addressType == "default" {
			shippingAddress := &model.ShippingAddress{
				Address: model.Address{
					Id:              0,
					CustomerId:      customerId,
					PostalCode:      nullStringPointer(postalCode),
					CountryCode:     countryCode,
					SubdivisionCode: nullStringPointer(subdivisionCode),
					SubdivisionName: nullStringPointer(subdivisionName),
					CityName:        *nullStringPointer(cityName),
					AddressLine1:    addressLine1,
					AddressLine2:    nullStringPointer(addressLine2),
				},
				Fullname:        fullname,
				IsResidential:   residential,
				RequestLiftgate: liftgate,
			}
			addrId, addressError = addressService.CreateShippingAddress(ctx, customerId, shippingAddress)
		} else if addressType == "billing" {
			billingAddress := &model.BillingAddress{
				Address: model.Address{
					Id:              0,
					CustomerId:      customerId,
					PostalCode:      nullStringPointer(postalCode),
					CountryCode:     countryCode,
					SubdivisionCode: nullStringPointer(subdivisionCode),
					SubdivisionName: nullStringPointer(subdivisionName),
					CityName:        *nullStringPointer(cityName),
					AddressLine1:    addressLine1,
					AddressLine2:    nullStringPointer(addressLine2),
				},
				Fullname: fullname,
			}
			addrId, addressError = addressService.CreateBillingAddress(ctx, customerId, billingAddress)
		} else {
			addressError = fmt.Errorf("неизвестный тип адреса: %s", addressType)
		}

		// Обработка результата добавления
		if addressError != nil {
			updateStatus(ctx, db, id, "error", addressError.Error())
			continue
		}

		// Успешно добавлено
		updateStatus(ctx, db, id, "done", strconv.Itoa(addrId))
	}

	// Проверка ошибок завершения обработчика строк
	if err = rows.Err(); err != nil {
		log.Fatalf("Ошибка обработки строк: %v", err)
	}
}

// Функция для обновления статуса записи
func updateStatus(ctx context.Context, db *sql.DB, id int, status, comment string) {
	_, err := db.ExecContext(ctx, `
		UPDATE customers_addresses_tmp_diff1
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
