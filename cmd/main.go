package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	_ "github.com/lib/pq"

	"proj1/internal/pkg/saving"
	"proj1/internal/pkg/server"
	"proj1/internal/pkg/storage"
)

const (
	file = "slice_storage.json"
)

func main() {
	filePath := os.Getenv("STORAGE_FILE_PATH")
	if filePath == "" {
		filePath = file
	}

	stor2, err := storage.NewSliceStorage(filePath)
	if err != nil {
		log.Fatal(err)
	}
	storageDB, err := saving.NewStorageDB(os.Getenv("POSTGRES"))
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	defer storageDB.Db.Close()
	var wg sync.WaitGroup
	closeChan := make(chan struct{})
	wg.Add(1)
	go func() {
		defer wg.Done()
		stor2.PeriodicClean(closeChan, 10*time.Minute, filePath)
	}()
	serverPort, ok := os.LookupEnv("BASIC_SERVER_PORT")
	if !ok {
		serverPort = "8090"
	}

	stor2.LoadFromFile(filePath)
	srv := server.New(":"+serverPort, &stor2)

	go func() {
		if err := srv.Start(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	close(closeChan)
	wg.Wait()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = srv.Shutdown(ctx)
	log.Fatalf("%s", err)
	if err != nil {
		log.Fatalf("Shutdown error: %s\n", err)
	}

	fmt.Println("Server exited")

}
