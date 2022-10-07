package repository

import (
	"context"

	"github.com/dimayasha7123/quiz_service_tg_bot/internal/models"
)

func (r repository) AddUser(ctx context.Context, user *models.User) error {

	query := `
	insert into user_account (username, tg_id, qs_id)
	values ($1, $2, $3) returning id;
	`

	var ID int64
	err := r.pool.QueryRow(ctx, query, user.Username, user.TGID, user.QSID).Scan(&ID)
	if err != nil {
		return err
	}

	user.ID = ID

	return nil
}
