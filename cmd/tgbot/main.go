package main

import (
	"context"
	"fmt"
	"log"
	"os"

	pb "github.com/dimayasha7123/quiz_service/pkg/api"
	"github.com/dimayasha7123/quiz_service_tg_bot/config"
	"github.com/dimayasha7123/quiz_service_tg_bot/internal/app"
	"github.com/dimayasha7123/quiz_service_tg_bot/internal/db"
	"github.com/dimayasha7123/quiz_service_tg_bot/internal/repository"
	"google.golang.org/grpc"
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

	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", cfg.Socket.Host, cfg.Socket.GrpcPort), grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	log.Println("Create grpc client connection")

	bclient := app.New(repository.New(adp), cfg.TelegramAPIKey, pb.NewQuizServiceClient(conn))

	log.Println("Create botApiClient")
	log.Println("Client running!")

	err = bclient.Run(ctx)
	if err != nil {
		log.Fatal(err)
	}

}
