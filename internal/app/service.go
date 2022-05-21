package app

import (
	pb "gitlab.ozon.dev/dimayasha7123/homework-2-dimayasha-7123/pkg/api"
	"hw2-tgbot/internal/models"
	"net/http"
)

type bclient struct {
	repo       repository
	apiKey     string
	quizClient pb.QuizServiceClient
	httpClient http.Client
	users      models.SyncMap
}

func New(repo repository, apiKey string, quizClient pb.QuizServiceClient) *bclient {

	bc := bclient{
		repo:       repo,
		apiKey:     apiKey,
		quizClient: quizClient,
		httpClient: http.Client{},
		users:      *models.NewSyncMap(),
	}

	return &bc
}
