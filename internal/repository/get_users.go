package repository

import (
	"context"

	"github.com/dimayasha7123/quiz_service_tg_bot/internal/models"
)

func (r repository) GetUsers(ctx context.Context) ([]*models.User, error) {

	query := `
	select * from user_account;
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	var users []*models.User

	for rows.Next() {
		var user models.User
		err = rows.Scan(&user.ID, &user.Username, &user.TGID, &user.QSID)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return users, nil
}
