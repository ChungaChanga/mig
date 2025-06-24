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

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

type Customer struct {
	Id       int
	Password string
}

func MigratePasswords() {
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
		customers, err := fetchPasswords(mysqlDB, offset, batchSize)
		if err != nil {
			log.Fatalf("Ошибка получения данных из MySQL: %v", err)
		}

		// Если больше записей нет, выходим
		if len(customers) == 0 {
			log.Println("Все данные успешно перенесены!")
			break
		}

		// Вставляем батч данных в PostgreSQL
		err = insertPasswords(postgresDB, customers)
		if err != nil {
			log.Fatalf("Ошибка вставки данных в PostgreSQL: %v", err)
		}

		offset += batchSize
		log.Printf("Перенесено %d записей...", offset)
		time.Sleep(1 * time.Second)
	}
}

// fetchCustomers выбирает пакет покупателей из MySQL
func fetchPasswords(mysqlDB *sql.DB, offset, limit int) ([]Customer, error) {
	query := `
		SELECT u.ID AS id, u.user_pass AS password FROM wp_users u
		ORDER BY u.ID
		LIMIT ? OFFSET ?;`

	// Выполняем запрос с лимитом и смещением
	rows, err := mysqlDB.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	defer rows.Close()

	// Читаем данные пакета
	customers := make([]Customer, 0)
	for rows.Next() {
		var customer Customer

		err := rows.Scan(&customer.Id, &customer.Password)
		if err != nil {
			return nil, fmt.Errorf("ошибка чтения строки: %w", err)
		}

		customers = append(customers, customer)
	}

	return customers, rows.Err()
}

// insertCustomers выполняет батчевую вставку записей в PostgreSQL
func insertPasswords(postgresDB *sql.DB, customers []Customer) error {
	if len(customers) == 0 {
		return nil
	}

	// SQL-запрос для вставки
	query := `
		INSERT INTO accounts (customer_id, password)
		VALUES %s
		ON CONFLICT (customer_id) DO NOTHING;
	`

	// Подготавливаем placeholders для батчевой вставки
	placeholders := make([]string, len(customers))
	fieldsCnt := 2
	values := make([]interface{}, 0, len(customers)*fieldsCnt)

	for i, customer := range customers {
		placeholders[i] = fmt.Sprintf("($%d, $%d)",
			i*fieldsCnt+1, i*fieldsCnt+2)
		values = append(values, customer.Id, customer.Password)
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
