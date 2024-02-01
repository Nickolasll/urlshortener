package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "net/http/pprof"

	"github.com/Nickolasll/urlshortener/internal/app/config"
	"github.com/Nickolasll/urlshortener/internal/app/domain"
	"github.com/Nickolasll/urlshortener/internal/app/infrastructure/repositories"
	"github.com/Nickolasll/urlshortener/internal/app/presentation"
	"github.com/google/uuid"
)

var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
)

func benchmark() {
	for {
		URL := "www.testlongurl.com"
		userID := uuid.New().String()
		repository := repositories.RAMRepository{
			OriginalToShorts: map[string]domain.Short{},
			Cache:            map[string]string{},
		}
		for i := 0; i < 10000000; i++ {
			short := domain.Shorten(URL, userID)
			repository.Save(short)
		}
	}
}

func main() {
	var err error
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)
	err = config.ParseFlags()
	if err != nil {
		panic(err)
	}
	go benchmark()
	server := &http.Server{
		Addr:    *config.ServerEndpoint,
		Handler: presentation.ChiFactory(),
	}

	idleConnsClosed := make(chan struct{})

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	go func() {
		<-sigint

		if err := server.Shutdown(context.Background()); err != nil {
			log.Printf("Server shutdown: %v", err)
		}

		close(idleConnsClosed)
	}()

	if *config.EnableHTTPS {
		err = server.ListenAndServeTLS("server.crt", "server.key")
	} else {
		err = server.ListenAndServe()
	}
	if err != nil && err != http.ErrServerClosed {
		panic(err)
	}
	<-idleConnsClosed
	log.Println("Server shutdown gracefully")
}
