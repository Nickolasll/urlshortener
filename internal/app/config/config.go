package config

import (
	"flag"
	"os"
)

var (
	ServerEndpoint  = flag.String("a", "localhost:8080", "Server endpoint")
	SlugEndpoint    = flag.String("b", "http://localhost:8080", "Shorten endpoint")
	FileStoragePath = flag.String("f", "/tmp/short-url-db.json", "File storage path")
	DatabaseDSN     = flag.String("d", "postgresql://admin:admin@localhost:5432/postgres", "Database DSN")
	SlugSize        = 8
)

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
