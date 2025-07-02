package brand

/*
устанавливаются дефолтные значения имен для юзеров без имен
если телефон > 20 устанавливает в nil
*/
import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

const (
	batchSize = 10000
)

type Brand struct {
	Name            string
	Description     string
	Slug            string
	Image           string
	MetaTitle       string
	MetaDescription string
	MetaKeywords    string
	Active          bool
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

	// Получение батча данных из MySQL
	brands, err := fetchBrands(mysqlDB)
	if err != nil {
		log.Fatalf("Ошибка получения данных из MySQL: %v", err)
	}

	// Вставляем батч данных в PostgreSQL
	err = insertBrands(postgresDB, brands)
	if err != nil {
		log.Fatalf("Ошибка вставки данных в PostgreSQL: %v", err)
	}
}

// fetchCustomers выбирает пакет покупателей из MySQL
func fetchBrands(mysqlDB *sql.DB) ([]Brand, error) {
	query := `select wp_terms.name, wp_terms.slug, wp_term_taxonomy.description, imgs.img from wp_term_taxonomy
         join wp_terms on wp_terms.term_id = wp_term_taxonomy.term_id
         join (
    select distinct(wp_termmeta.term_id) as term_id, wp_postmeta.meta_value as img from wp_term_taxonomy
    join wp_termmeta on wp_termmeta.term_id = wp_term_taxonomy.term_id
    join wp_postmeta on  wp_termmeta.meta_value=wp_postmeta.post_id
    where wp_term_taxonomy.taxonomy='product_brand'
      and wp_termmeta.meta_key='thumbnail_id'
      and wp_termmeta.meta_value != '' and wp_termmeta.meta_value is not null
      and wp_postmeta.meta_key='_wp_attached_file'
    order by wp_termmeta.term_id DESC
    ) as imgs on wp_term_taxonomy.term_id = imgs.term_id
         where wp_term_taxonomy.taxonomy='product_brand'
           and wp_term_taxonomy.description != ''
         and imgs.img not like 'https://www.directliquidation.com/contents/uploads%'`

	// Выполняем запрос с лимитом и смещением
	rows, err := mysqlDB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	defer rows.Close()

	// Читаем данные пакета
	brands := make([]Brand, 0)
	imageRe := regexp.MustCompile(`^\d{4}/\d{2}/`)
	for rows.Next() {
		var brand Brand

		err := rows.Scan(&brand.Name, &brand.Slug, &brand.Description, &brand.Image)
		if err != nil {
			return nil, fmt.Errorf("ошибка чтения строки: %w", err)
		}

		if brand.Image != "" {
			//DownloadImage("https://files.directliquidation.com/directliquidation/"+brand.Image, "/home/andrey/projects/blog-images")
			brand.Image = imageRe.ReplaceAllString(brand.Image, "")
			brand.Image = fmt.Sprintf("https://www.directliquidation.com/media/cms/brand/%s", brand.Image)
		}
		brand.MetaTitle = fmt.Sprintf("%s Liquidation Auctions - Bulk Wholesale Lots - DirectLiquidation", brand.Name)
		brand.MetaDescription = fmt.Sprintf("Liquidation auctions w/ %s surplus inventory in bulk wholesale lots by box, pallet or truckload. Source high quality goods from a top US retailer.", brand.Name)
		brand.MetaKeywords = brand.Name
		brand.Active = true

		brands = append(brands, brand)
	}

	return brands, rows.Err()
}

// insertCustomers выполняет батчевую вставку записей в PostgreSQL
func insertBrands(postgresDB *sql.DB, brands []Brand) error {
	if len(brands) == 0 {
		return nil
	}

	// SQL-запрос для вставки
	query := `
		INSERT INTO content_pages (title, content, slug, image, meta_title, meta_description, meta_keywords, active)
		VALUES %s ON CONFLICT (slug) DO NOTHING;
	`

	// Подготавливаем placeholders для батчевой вставки
	placeholders := make([]string, len(brands))
	fieldsCnt := 8
	values := make([]interface{}, 0, len(brands)*fieldsCnt)

	for i, brand := range brands {
		placeholders[i] = fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
			i*fieldsCnt+1, i*fieldsCnt+2, i*fieldsCnt+3, i*fieldsCnt+4, i*fieldsCnt+5, i*fieldsCnt+6, i*fieldsCnt+7, i*fieldsCnt+8)
		values = append(values, brand.Name, brand.Description, brand.Slug, brand.Image, brand.MetaTitle, brand.MetaDescription, brand.MetaKeywords, brand.Active)
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
