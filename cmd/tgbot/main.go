package main

import (
	"context"
	pb "gitlab.ozon.dev/dimayasha7123/homework-2-dimayasha-7123/pkg/api"
	"google.golang.org/grpc"
	"hw2-tgbot/config"
	"hw2-tgbot/internal/app"
	"hw2-tgbot/internal/db"
	"hw2-tgbot/internal/repository"
	"log"
	"os"
)

const (
	configPath = "./config/config.yaml"
)

func main() {
	b, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
	}

	cfg, err := config.ParseConfig(b)
	if err != nil {
		log.Fatal(err)
	}

	//log.Printf("Config = %+v\n", cfg)
	log.Println("Config unmarshalled")

	ctx := context.Background()

	adp, err := db.New(ctx, cfg.Dsn)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Get db adapter")

	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	log.Println("Create grpc client connection")

	bclient := app.New(repository.New(adp), cfg.ApiKeys.Telegram, pb.NewQuizServiceClient(conn))

	log.Println("Create botApiClient")
	log.Println("Client running!")

	err = bclient.Run(ctx)
	if err != nil {
		log.Fatal(err)
	}

}
