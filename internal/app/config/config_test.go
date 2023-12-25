package config_test

import "github.com/Nickolasll/urlshortener/internal/app/config"

func Example() {
	config.ParseFlags()
	myDSN := *config.DatabaseDSN
	if myDSN != "" {
		// Получаем из конфига значение переменной DatabaseDSN
	}
}
