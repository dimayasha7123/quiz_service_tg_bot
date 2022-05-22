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

func (b *bclient) confirmHandler(ctx context.Context, update models.Update) (string, error) {

	b.users.RLock()
	user, ok := b.users.M[update.Message.From.ID]
	b.users.RUnlock()

	if !ok {
		return "", errors.New("user not found")
	}
	if user.State != 1 {
		return "", errors.New("invalid state")
	}

	var text string

	user.CurrentQuestion++
	if user.CurrentQuestion >= len(user.Questions) {
		user.State = 0

		answers := make([]*pb.QuestionRightAnswers, 0, 10)

		for _, q := range user.Questions {
			ansNums := make([]int32, 0, 6)
			for j, a := range q.AnswerOptions {
				if a.Picked {
					ansNums = append(ansNums, int32(j))
				}
			}
			answers = append(answers, &pb.QuestionRightAnswers{RightAnswerNumbers: ansNums})
		}

		ansPack := pb.AnswersPack{
			QuizPartyID: user.QuizPartyID,
			Answers:     answers,
		}

		sTop, err := b.quizClient.SendAnswers(ctx, &ansPack)
		if err != nil {
			return "", errors.New("quiz service error")
		}

		sb := strings.Builder{}
		sb.WriteString("Your results:\n")
		sb.WriteString(fmt.Sprintf("Points: %d\n", sTop.UserResults.PointCount))
		sb.WriteString(fmt.Sprintf("Place: %d\n", sTop.UserResults.Place))
		sb.WriteString("\nTop:\n")
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
		text = sb.String()

	} else {

		text, ok = user.GetQuestion(user.CurrentQuestion)
		if !ok {
			return "", errors.New("quiz has no quiestions")
		}

	}
	return text, nil
}
