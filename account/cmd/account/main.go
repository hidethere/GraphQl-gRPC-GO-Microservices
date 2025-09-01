package main

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/hidethere/GraphQl-gRPC-GO-Microservices/account"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/tinrab/retry"
)

type Config struct {
	DatabaseURL string `envconfig:"ACCOUNT_DATABASE_URL"`
}

func main() {
	cwd, _ := os.Getwd()

	envPath := filepath.Join(cwd, "..", "..", "..", ".env.local")

	if err := godotenv.Load(envPath); err != nil {
		log.Printf("No .env.local file found at %s, skipping\n", envPath)
	}
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	var r account.Repository
	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {
		r, err = account.NewPostgresRepository(cfg.DatabaseURL)
		if err != nil {
			log.Println(err)
		}
		return
	})
	defer r.Close()
	log.Println("Listening on port 8081...")
	s := account.NewService(r)
	log.Fatal(account.ListenGRPC(s, 8081))

}
