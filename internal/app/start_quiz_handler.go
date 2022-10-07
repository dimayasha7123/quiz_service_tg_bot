package app

import (
	"context"
	"errors"
	"strconv"

	pb "github.com/dimayasha7123/quiz_service/pkg/api"
	"github.com/dimayasha7123/quiz_service_tg_bot/internal/models"
)

func (b *bclient) startQuizHandler(ctx context.Context, update models.Update, commandArgs []string) (string, error) {

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

	party, err := b.quizClient.StartQuizParty(ctx, &pb.QuizUserInfo{
		UserID: user.QSID,
		QuizID: quizID,
	})
	if err != nil {
		return "", errors.New("quiz service error")
	}

	questions := make([]models.Question, 0, 10)
	for _, q := range party.Questions {
		answerOptions := make([]models.AnswerOption, 0, 6)
		for _, a := range q.AnswerOptions {
			answerOptions = append(answerOptions, models.AnswerOption{
				Title:  a,
				Picked: false,
			})
		}
		question := models.Question{
			Title:         q.Title,
			AnswerOptions: answerOptions,
		}
		questions = append(questions, question)
	}

	user.Questions = questions
	user.State = 1
	user.QuizPartyID = party.QuizPartyID
	user.CurrentQuestion = 0

	text, ok := user.GetQuestion(user.CurrentQuestion)
	if !ok {
		return "", errors.New("quiz has no quiestions")
	}

	return text, nil
}
