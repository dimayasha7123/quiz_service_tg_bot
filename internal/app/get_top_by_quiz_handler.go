package app

import (
	"context"
	"errors"
	"fmt"
	pb "gitlab.ozon.dev/dimayasha7123/homework-2-dimayasha-7123/pkg/api"
	"hw2-tgbot/internal/models"
	"strconv"
	"strings"
)

func (b *bclient) getTopByQuizHandler(ctx context.Context, update models.Update, commandArgs []string) (string, error) {

	if len(commandArgs) > 1 {
		return "", errors.New("too much arguments")
	}

	quizID, err := strconv.ParseInt(commandArgs[0], 10, 64)
	if err != nil {
		return "", errors.New("bad arguments")
	}

	b.users.RLock()
	user, ok := b.users.M[update.Message.From.ID]
	b.users.RUnlock()
	if !ok {
		return "", errors.New("user not found")
	}

	sTop, err := b.quizClient.GetQuizTop(ctx, &pb.QuizUserInfo{
		UserID: user.QSID,
		QuizID: quizID,
	})
	if err != nil {
		return "", err
	}

	sb := strings.Builder{}

	if sTop.UserResults.Place == 0 && sTop.UserResults.PointCount == 0 {
		sb.WriteString("You didn't take part\n\n")
	} else {
		sb.WriteString("Your results:\n")
		sb.WriteString(fmt.Sprintf("Points: %d\n", sTop.UserResults.PointCount))
		sb.WriteString(fmt.Sprintf("Place: %d\n", sTop.UserResults.Place))
	}

	if len(sTop.QuizTop.Results) == 0 {
		sb.WriteString("No one has participated yet =(\nBut you can be the first!\n")
	} else {
		sb.WriteString("\nTop:\n")
	}

	for _, r := range sTop.QuizTop.Results {

		username := r.Name

		tgID, err := strconv.ParseInt(r.Name, 10, 64)
		if err == nil {
			nameFromDB, err := b.repo.FindUser(ctx, tgID)
			if err == nil && nameFromDB != "" {
				username = nameFromDB
			}
		}

		sb.WriteString(fmt.Sprintf("%d. @%s, with %d p.\n", r.Place, username, r.PointCount))
	}

	return sb.String(), nil
}
