package repository

import "context"

func (r *repository) FindUser(ctx context.Context, tgID int64) (string, error) {

	query := `
	select username from user_account where tg_id = $1;
	`

	var username string
	err := r.pool.QueryRow(ctx, query, tgID).Scan(&username)
	if err != nil {
		return "", err
	}
	return username, nil
}
