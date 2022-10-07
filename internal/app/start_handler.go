package app

import (
	"context"
	"fmt"

	pb "github.com/dimayasha7123/quiz_service/pkg/api"
	"github.com/dimayasha7123/quiz_service_tg_bot/internal/models"
)

func (b *bclient) startHandler(ctx context.Context, update models.Update) (string, error) {
	b.users.RLock()
	user, ok := b.users.M[update.Message.From.ID]
	b.users.RUnlock()

	var fmtText string

	if ok {

		fmtText = "Welcome, %s. Again"

	} else {

		fmtText = "Welcome, %s"
		user = models.NewUser(
			update.Message.From.ID,
			update.Message.From.Username,
		)

		qsID, err := b.quizClient.AddUser(ctx, &pb.User{Name: fmt.Sprintf("%d", user.TGID)})
		if err != nil {
			return "", err
		}
		user.QSID = qsID.ID

		err = b.repo.AddUser(ctx, user)
		if err != nil {
			return "", err
		}

		b.users.Lock()
		b.users.M[user.TGID] = user
		b.users.Unlock()

	}

	text := fmt.Sprintf(fmtText, user.Username)

	return text, nil
}
