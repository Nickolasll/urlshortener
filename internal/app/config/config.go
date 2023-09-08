package config

import (
	"flag"
)

var (
	ServerEndpoint = flag.String("a", "localhost:8080", "Server endpoint")
	SlugEndpoint   = flag.String("b", "http://localhost:8080", "Shorten endpoint")
)
