package app

import (
	"context"

	"github.com/dimayasha7123/quiz_service_tg_bot/internal/models"
)

type repository interface {
	GetUsers(ctx context.Context) ([]*models.User, error)
	AddUser(ctx context.Context, user *models.User) error
	FindUser(ctx context.Context, tgID int64) (string, error)
}
