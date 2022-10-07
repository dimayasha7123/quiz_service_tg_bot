package app

import (
	"net/http"

	pb "github.com/dimayasha7123/quiz_service/pkg/api"
	"github.com/dimayasha7123/quiz_service_tg_bot/internal/models"
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
