package config

import (
	"flag"
	"os"
)

var (
	ServerEndpoint = flag.String("a", "localhost:8080", "Server endpoint")
	SlugEndpoint   = flag.String("b", "http://localhost:8080", "Shorten endpoint")
	SlugSize       = 8
)

func ParseFlags() {
	flag.Parse()

	if envServerAddr := os.Getenv("SERVER_ADDRESS"); envServerAddr != "" {
		*ServerEndpoint = envServerAddr
	}

	if envBaseURL := os.Getenv("BASE_URL"); envBaseURL != "" {
		*SlugEndpoint = envBaseURL
	}
}