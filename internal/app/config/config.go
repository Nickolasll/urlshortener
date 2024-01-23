package config

import (
	"encoding/json"
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
	// EnableHTTPS - Включение HTTPS в веб-сервере
	EnableHTTPS = flag.Bool("s", false, "Enable HTTPS")
	// ConfigPath - Путь к JSON файлу с конфигурацией приложения
	ConfigPath = flag.String("c", "", "Путь к JSON файлу с конфигурацией приложения")
)

// ParseFlags - Инициализирует конфигурацию сервиса, читая флаги и переменные окрудения
func ParseFlags() error {
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

	if envEnableHTTPS := os.Getenv("ENABLE_HTTPS"); envEnableHTTPS != "" {
		*EnableHTTPS = true
	}

	if envConfigPath, ok := os.LookupEnv("CONFIG"); ok {
		*ConfigPath = envConfigPath
	}

	if *ConfigPath != "" {
		err := loadConfigFromFile(*ConfigPath)
		if err != nil {
			return err
		}
	}

	return nil
}

// FileConfig - конфигурация приложения, прочитанная из файла
type FileConfig struct {
	// ServerEndpoint - адрес сокращателя ссылок
	ServerEndpoint string `json:"server_address"`
	// SlugEndpoint - адрес для получения полной ссылки по сокращенной
	SlugEndpoint string `json:"base_url"`
	// FileStoragePath - путь до репозитория-файла
	FileStoragePath string `json:"file_storage_path"`
	// DatabaseDSN - источник данных для подключения к postgres
	DatabaseDSN string `json:"database_dsn"`
	// EnableHTTPS - Включение HTTPS в веб-сервере
	EnableHTTPS bool `json:"enable_https"`
}

func loadConfigFromFile(path string) error {
	var cfg FileConfig
	file, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	err = json.Unmarshal(file, &cfg)
	if err != nil {
		return err
	}
	if *ServerEndpoint == "" {
		*ServerEndpoint = cfg.ServerEndpoint
	}
	if *SlugEndpoint == "" {
		*SlugEndpoint = cfg.SlugEndpoint
	}
	if *FileStoragePath == "" {
		*FileStoragePath = cfg.FileStoragePath
	}
	if *DatabaseDSN == "" {
		*DatabaseDSN = cfg.DatabaseDSN
	}
	if !*EnableHTTPS {
		*EnableHTTPS = cfg.EnableHTTPS
	}
	return nil
}
