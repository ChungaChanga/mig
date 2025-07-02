package brand

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// downloadImage загружает изображение из URL и сохраняет его в указанную папку
func DownloadImage(url, folder string) {

	// Получение имени файла из URL
	fileName := filepath.Base(url)
	filePath := filepath.Join(folder, fileName)

	// Создание HTTP-запроса
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Ошибка загрузки %s: %v\n", url, err)
		return
	}
	defer resp.Body.Close()

	// Проверка статуса HTTP-ответа
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Не удалось загрузить %s: HTTP %d\n", url, resp.StatusCode)
		return
	}

	// Создание локального файла для записи
	out, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("Ошибка при создании файла %s: %v\n", filePath, err)
		return
	}
	defer out.Close()

	// Копирование содержимого HTTP-ответа в файл
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Printf("Ошибка при сохранении %s: %v\n", filePath, err)
		return
	}

	fmt.Printf("Успешно загружено: %s\n", filePath)
}
