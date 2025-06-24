package user

/*
устанавливаются дефолтные значения имен для юзеров без имен
если телефон > 20 устанавливает в nil
*/
import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
	"unicode/utf8"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

const (
	batchSize        = 10000
	defaultFirstName = ""
	defaultLastName  = ""
)

type OldCustomer struct {
	Id       int
	Email    string
	FullName *string
	Company  *string
	Phone    *string
}

func (o *OldCustomer) FirstName() string {
	if o.FullName == nil || *o.FullName == "" {
		return defaultFirstName
	}
	parts := strings.Fields(*o.FullName)
	if len(parts) > 0 {
		return parts[0] // Первый элемент — это имя
	}
	return defaultFirstName
}

func (o *OldCustomer) LastName() string {
	if o.FullName == nil || *o.FullName == "" {
		return defaultLastName
	}
	parts := strings.Fields(*o.FullName)
	if len(parts) > 1 {
		return strings.Join(parts[1:], " ") // Все, кроме первого элемента — это фамилия
	}
	return defaultLastName
}

type NewCustomer struct {
	Id        int
	Email     string
	FirstName string
	LastName  string
	Phone     *string
	Company   *string
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

	offset := 150000
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
		time.Sleep(1 * time.Second)
	}
}

// fetchCustomers выбирает пакет покупателей из MySQL
func fetchCustomers(mysqlDB *sql.DB, offset, limit int) ([]NewCustomer, error) {
	query := `
		SELECT 
		    u.ID AS id,
			u.user_email AS email,
			(SELECT meta_value FROM wp_usermeta WHERE user_id = u.ID AND meta_key = 'full_name' LIMIT 1) AS fullname,
			(SELECT meta_value FROM wp_usermeta WHERE user_id = u.ID AND meta_key = 'registered_user_phone' LIMIT 1) AS phone,
			(SELECT meta_value FROM wp_usermeta WHERE user_id = u.ID AND meta_key = 'company_name' LIMIT 1) AS company
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

		err := rows.Scan(&oldCustomer.Id, &oldCustomer.Email, &oldCustomer.FullName, &oldCustomer.Phone, &oldCustomer.Company)
		if err != nil {
			return nil, fmt.Errorf("ошибка чтения строки: %w", err)
		}

		if oldCustomer.Phone != nil && utf8.RuneCountInString(*oldCustomer.Phone) > 20 {
			oldCustomer.Phone = nil
		}
		if oldCustomer.Company != nil && utf8.RuneCountInString(*oldCustomer.Company) > 100 {
			oldCustomer.Company = nil
		}
		if utf8.RuneCountInString(oldCustomer.Email) > 100 {
			continue
		}

		newCustomer := NewCustomer{
			Id:        oldCustomer.Id,
			Email:     oldCustomer.Email,
			FirstName: trimString(oldCustomer.FirstName(), 100),
			LastName:  trimString(oldCustomer.LastName(), 100),
			Phone:     oldCustomer.Phone,
			Company:   oldCustomer.Company,
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
		INSERT INTO customers (id, email, firstname, lastname, phone, company)
		VALUES %s
		ON CONFLICT (email) DO NOTHING;
	`

	// Подготавливаем placeholders для батчевой вставки
	placeholders := make([]string, len(customers))
	fieldsCnt := 6
	values := make([]interface{}, 0, len(customers)*fieldsCnt)

	for i, customer := range customers {
		placeholders[i] = fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d)",
			i*fieldsCnt+1, i*fieldsCnt+2, i*fieldsCnt+3, i*fieldsCnt+4, i*fieldsCnt+5, i*fieldsCnt+6)
		values = append(values, customer.Id, customer.Email, customer.FirstName, customer.LastName, customer.Phone, customer.Company)
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

func trimString(input string, limit int) string {
	if utf8.RuneCountInString(input) > limit {
		// Преобразуем строку в руны для корректной обрезки
		runes := []rune(input)
		return string(runes[:limit])
	}
	return input
}
