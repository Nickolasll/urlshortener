package config

import (
	"flag"
	"os"
)

// Переменные конфигурации для сервиса
var (
	// ServerEndpoint - адрес сокращателя ссылок
	ServerEndpoint = flag.String("a", "localhost:8080", "Server endpoint")
	// SlugEndpoint - адрес для получения полной ссылки по сокращенной
	SlugEndpoint = flag.String("b", "http://localhost:8080", "Shorten endpoint")
	// FileStoragePath - путь до репозитория-файла
	FileStoragePath = flag.String("f", "/tmp/short-url-db.json", "File storage path")
	// DatabaseDSN - источник данных для подключения к postgres
	DatabaseDSN = flag.String("d", "", "Database DSN")
	// SlugSize - Размер сокращенной ссылки
	SlugSize = 8
	// TokenExp - Время жизни токена в секундах
	TokenExp = 3600
	// SecretKey - Ключ шифрования токена
	SecretKey = "supersecretkey"
)

// Инициализирует конфигурацию сервиса, читая флаги и переменные окрудения
func ParseFlags() {
	flag.Parse()

	if envServerAddr := os.Getenv("SERVER_ADDRESS"); envServerAddr != "" {
		*ServerEndpoint = envServerAddr
	}

	if envBaseURL := os.Getenv("BASE_URL"); envBaseURL != "" {
		*SlugEndpoint = envBaseURL
	}

	if envFileStoragePath, ok := os.LookupEnv("FILE_STORAGE_PATH"); ok {
		*FileStoragePath = envFileStoragePath
	}

	if envDatabaseDSN, ok := os.LookupEnv("DATABASE_DSN"); ok {
		*DatabaseDSN = envDatabaseDSN
	}
}
