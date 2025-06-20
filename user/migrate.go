package user

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

const (
	batchSize = 1000
)

type OldCustomer struct {
	Email    string
	FullName *string
	Phone    *string
}

func (o *OldCustomer) FirstName() string {
	if o.FullName == nil || *o.FullName == "" {
		return ""
	}
	parts := strings.Fields(*o.FullName)
	if len(parts) > 0 {
		return parts[0] // Первый элемент — это имя
	}
	return ""
}

func (o *OldCustomer) LastName() string {
	if o.FullName == nil || *o.FullName == "" {
		return ""
	}
	parts := strings.Fields(*o.FullName)
	if len(parts) > 1 {
		return strings.Join(parts[1:], " ") // Все, кроме первого элемента — это фамилия
	}
	return ""
}

type NewCustomer struct {
	Email     string
	FirstName string
	LastName  string
	Phone     *string
}

func Migrate() {
	// Подключаемся к MySQL
	mysqlDB, err := sql.Open("mysql", os.Getenv("LEGACY_DL_DSN"))
	if err != nil {
		log.Fatalf("Ошибка подключения к MySQL: %v", err)
	}
	defer mysqlDB.Close()

	// Подключаемся к PostgreSQL
	postgresDB, err := sql.Open("postgres", os.Getenv("NEW_DL_DSN"))
	if err != nil {
		log.Fatalf("Ошибка подключения к PostgreSQL: %v", err)
	}
	defer postgresDB.Close()

	offset := 0
	for {
		// Получение батча данных из MySQL
		customers, err := fetchCustomers(mysqlDB, offset, batchSize)
		if err != nil {
			log.Fatalf("Ошибка получения данных из MySQL: %v", err)
		}

		// Если больше записей нет, выходим
		if len(customers) == 0 {
			log.Println("Все данные успешно перенесены!")
			break
		}

		// Вставляем батч данных в PostgreSQL
		err = insertCustomers(postgresDB, customers)
		if err != nil {
			log.Fatalf("Ошибка вставки данных в PostgreSQL: %v", err)
		}

		offset += batchSize
		log.Printf("Перенесено %d записей...", offset)
	}
}

// fetchCustomers выбирает пакет покупателей из MySQL
func fetchCustomers(mysqlDB *sql.DB, offset, limit int) ([]NewCustomer, error) {
	query := `
		SELECT 
			u.user_email AS email,
			(SELECT meta_value FROM wp_usermeta WHERE user_id = u.ID AND meta_key = 'full_name' LIMIT 1) AS fullname,
			(SELECT meta_value FROM wp_usermeta WHERE user_id = u.ID AND meta_key = 'registered_user_phone' LIMIT 1) AS phone
		FROM wp_users u
		ORDER BY u.ID
		LIMIT ? OFFSET ?;`

	// Выполняем запрос с лимитом и смещением
	rows, err := mysqlDB.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	defer rows.Close()

	// Читаем данные пакета
	customers := make([]NewCustomer, 0)
	for rows.Next() {
		var oldCustomer OldCustomer

		err := rows.Scan(&oldCustomer.Email, &oldCustomer.FullName, &oldCustomer.Phone)
		if err != nil {
			return nil, fmt.Errorf("ошибка чтения строки: %w", err)
		}
		newCustomer := NewCustomer{
			Email:     oldCustomer.Email,
			FirstName: oldCustomer.FirstName(),
			LastName:  oldCustomer.LastName(),
			Phone:     oldCustomer.Phone,
		}
		customers = append(customers, newCustomer)
	}

	return customers, rows.Err()
}

// insertCustomers выполняет батчевую вставку записей в PostgreSQL
func insertCustomers(postgresDB *sql.DB, customers []NewCustomer) error {
	if len(customers) == 0 {
		return nil
	}

	// SQL-запрос для вставки
	query := `
		INSERT INTO customers (email, firstname, lastname, phone)
		VALUES %s
		ON CONFLICT (email) DO NOTHING;
	`

	// Подготавливаем placeholders для батчевой вставки
	placeholders := make([]string, len(customers))
	values := make([]interface{}, 0, len(customers)*4)

	for i, customer := range customers {
		placeholders[i] = fmt.Sprintf("($%d, $%d, $%d, $%d)",
			i*4+1, i*4+2, i*4+3, i*4+4)
		values = append(values, customer.Email, customer.FirstName, customer.LastName, customer.Phone)
	}

	// Формируем финальный SQL-запрос
	finalQuery := fmt.Sprintf(query, strings.Join(placeholders, ", "))

	// Выполняем запрос
	_, err := postgresDB.Exec(finalQuery, values...)
	if err != nil {
		return fmt.Errorf("ошибка при выполнении вставки: %w", err)
	}

	return nil
}
