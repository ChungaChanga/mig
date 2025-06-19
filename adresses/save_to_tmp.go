package address

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/elliotchance/phpserialize"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

// Batch size for migration
const batchSize = 100

type Address struct {
	Address1        string `php:"address_1"`
	Address2        string `php:"address_2"`
	City            string `php:"city"`
	Country         string `php:"country"`
	FullName        string `php:"full_name"`
	PostCode        string `php:"postcode"`
	State           string `php:"state"`
	AdrHash         string `php:"adr_hash"`
	IsResidential   int    `php:"is_residential"`
	Phone           string `php:"phone"`
	DefaultBilling  int    `php:"default_billing"`
	DefaultShipping int    `php:"default_shipping"`
}

func convertMapToSlice(raw map[interface{}]interface{}) ([]Address, error) {
	var addresses []Address

	// Цикл по ключам карты
	for key, value := range raw {
		// Пропустим непонятный элемент, например, "full_name"
		if key == "full_name" {
			continue
		}

		// Проверяем, что значение является картой (map)
		entry, ok := value.(map[interface{}]interface{})
		if !ok {
			log.Printf("Пропускаем несоответствующий элемент: ключ=%s, значение=%v", key, value)
			continue
		}
		log.Print(entry)
		// Преобразуем карту в структуру Address
		address := Address{
			Address1:        getString(entry["address_1"]),
			Address2:        getString(entry["address_2"]),
			City:            getString(entry["city"]),
			Country:         getString(entry["country"]),
			FullName:        getString(entry["full_name"]),
			PostCode:        getString(entry["postcode"]),
			State:           getString(entry["state"]),
			AdrHash:         getString(entry["adr_hash"]),
			IsResidential:   getInt(entry["is_residential"]),
			Phone:           getString(entry["phone"]),
			DefaultBilling:  getInt(entry["default_billing"]),
			DefaultShipping: getInt(entry["default_shipping"]),
		}

		// Добавляем адрес в срез
		addresses = append(addresses, address)
	}
	return addresses, nil
}

// Помощник для извлечения строки
func getString(val interface{}) string {
	if str, ok := val.(string); ok {
		return str
	}
	return ""
}

func getInt(val interface{}) int {
	switch v := val.(type) {
	case int:
		return v
	case int64:
		return int(v) // Явное преобразование int64 -> int
	default:
		return 0
	}
}

// Entry point
func Migrate() {
	// Получение DSN из переменных окружения
	legacyDSN := os.Getenv("LEGACY_DL_DSN")
	newDSN := os.Getenv("NEW_DL_DSN")

	if legacyDSN == "" || newDSN == "" {
		log.Fatal("Не указаны переменные окружения LEGACY_DL_DSN или NEW_DL_DSN")
	}

	// Подключение к базам данных
	legacyDB, err := sql.Open("mysql", legacyDSN)
	if err != nil {
		log.Fatalf("Ошибка подключения к LEGACY DL: %v", err)
	}
	defer legacyDB.Close()

	newDB, err := sql.Open("postgres", newDSN)
	if err != nil {
		log.Fatalf("Ошибка подключения к NEW DL: %v", err)
	}
	defer newDB.Close()

	// Миграция данных
	err = migrateAddresses(legacyDB, newDB)
	if err != nil {
		log.Fatalf("Ошибка миграции данных: %v", err)
	}
}

// migrateAddresses выполняет миграцию данных из MySQL в PostgreSQL
func migrateAddresses(legacyDB, newDB *sql.DB) error {
	offset := 0
	for {
		// Получение пакета данных из LEGACY DL
		log.Printf("Загрузка пакета данных. Offset: %d", offset)
		rows, err := legacyDB.Query(`
            SELECT wp_users.user_email, wp_usermeta.meta_value
            FROM wp_users
            JOIN wp_usermeta ON wp_users.ID = wp_usermeta.user_id
            WHERE wp_usermeta.meta_key = '_dliquid_address_book'
            AND wp_users.ID IN (select user_id from wp_usermeta where meta_key='reseller_certs')
            ORDER BY wp_users.ID
            LIMIT ? OFFSET ?`, batchSize, offset)
		if err != nil {
			return fmt.Errorf("Ошибка выполнения запроса к LEGACY DL: %v", err)
		}

		var batch []interface{}
		count := 0
		for rows.Next() {
			var email string
			var metaValue string
			err = rows.Scan(&email, &metaValue)
			if err != nil {
				return fmt.Errorf("Ошибка чтения строки: %v", err)
			}

			var rawData map[interface{}]interface{}

			// Попытка десериализации
			err := phpserialize.Unmarshal([]byte(metaValue), &rawData)
			if err != nil {
				log.Fatalf("Ошибка десериализации: %v", err)
			}

			// Обработать десериализованные данные
			addresses, err := convertMapToSlice(rawData)
			if err != nil {
				log.Fatalf("Ошибка обработки адресов: %v", err)
			}

			// Преобразование адресов в пакет для вставки
			for _, addr := range addresses {
				parts := strings.Split(addr.FullName, " ")
				firstName := ""
				lastName := ""
				if len(parts) > 0 {
					firstName = parts[0]
				}
				if len(parts) > 1 {
					lastName = strings.Join(parts[1:], " ")
				}

				// Определение типа (type): "billing", "shipping" или "default"
				addrType := "default"
				if addr.DefaultBilling == 1 {
					addrType = "billing"
				}
				if addr.DefaultShipping == 1 {
					addrType = "shipping"
				}
				meta := map[string]string{
					"fullname": firstName + " " + lastName,
				}

				// Конвертируем в JSON
				metaJson, err := json.Marshal(meta)
				if err != nil {
					log.Fatalf("Ошибка при создании JSON: %v", err)
				}

				batch = append(batch, email, addr.PostCode, addr.Country, addr.State, "", addr.City, addr.Address1, addr.Address2, metaJson, addrType)
				count++
			}
		}
		_ = rows.Close()

		// Если нет данных, выходим
		if count == 0 {
			break
		}

		// Вставка данных в NEW DL
		log.Printf("Вставка %d записей в NEW DL", count)
		query := buildInsertQuery(count)
		_, err = newDB.Exec(query, batch...)
		if err != nil {
			return fmt.Errorf("Ошибка вставки в NEW DL: %v", err)
		}

		// Увеличение смещения для следующего пакета
		offset += batchSize
	}

	log.Println("Миграция завершена успешно")
	return nil
}

// buildInsertQuery создаёт SQL-запрос для вставки данных в PostgreSQL
func buildInsertQuery(count int) string {
	values := []string{}
	for i := 0; i < count; i++ {
		j := i * 10
		values = append(values, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)", j+1, j+2, j+3, j+4, j+5, j+6, j+7, j+8, j+9, j+10))
	}
	return fmt.Sprintf(`
        INSERT INTO customers_addresses_tmp (
            customer_email, postal_code, country_code,
            subdivision_code, subdivision_name, city_name,
            address_line1, address_line2, meta, type
        ) VALUES %s
        ON CONFLICT (customer_email, postal_code, country_code, address_line1, address_line2)
        DO NOTHING`, strings.Join(values, ", "))
}
