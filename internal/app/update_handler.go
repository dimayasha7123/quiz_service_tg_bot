package app

import (
	"context"
	"errors"
	"fmt"
	pb "gitlab.ozon.dev/dimayasha7123/homework-2-dimayasha-7123/pkg/api"
	"google.golang.org/protobuf/types/known/emptypb"
	"hw2-tgbot/internal/models"
	"log"
	url2 "net/url"
	"strconv"
	"strings"
)

// TODO если quizserver вернет ошибку, то тут все упадет

func (b *bclient) updateHandler(ctx context.Context, update models.Update) error {

	messWords := strings.Split(update.Message.Text, "_")
	command := messWords[0]
	commandArgs := messWords[1:]

	log.Println(command)
	log.Println(commandArgs)

	var url string
	useDefaultUrl := true
	text := "This command not implemented (yet)"

	switch {

	case command == "/start":

		b.users.RLock()
		user, ok := b.users.M[update.Message.From.ID]
		b.users.RUnlock()

		var fmtText string

		if ok {

			fmtText = "Welcome, %s. Again"
			// TODO update user's username

		} else {

			fmtText = "Welcome, %s"
			user = models.NewUser(
				update.Message.From.ID,
				update.Message.From.Username,
			)

			qsID, err := b.quizClient.AddUser(ctx, &pb.User{Name: fmt.Sprintf("%d", user.TGID)})
			if err != nil {
				return err
			}
			user.QSID = qsID.ID

			err = b.repo.AddUser(ctx, user)
			if err != nil {
				return err
			}

			b.users.Lock()
			b.users.M[user.TGID] = user
			b.users.Unlock()

		}

		text = fmt.Sprintf(fmtText, user.Username)

	case command == "/help":

		text = "There should be help text, but there is no"

	case command == "/getquizes" || (command == "/startquiz" && len(commandArgs) == 0):

		quizes, err := b.quizClient.GetQuizList(ctx, &emptypb.Empty{})
		if err != nil {
			return err
		}

		text = ""
		for _, q := range quizes.QList {
			text += fmt.Sprintf("%s\nStart this quiz: /startquiz_%d\n\n", q.Name, q.ID)
		}

		// TODO add pagination

	case command == "/startquiz":

		if len(commandArgs) > 1 {
			return errors.New("too much arguments")
		}

		quizID, err := strconv.ParseInt(commandArgs[0], 10, 64)
		if err != nil {
			return errors.New("bad arguments")
		}

		b.users.RLock()
		user, ok := b.users.M[update.Message.From.ID]
		b.users.RUnlock()
		if !ok {
			return errors.New("user not found")
		}

		party, err := b.quizClient.StartQuizParty(ctx, &pb.QuizUserInfo{
			UserID: user.QSID,
			QuizID: quizID,
		})
		if err != nil {
			return errors.New("quiz service error")
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

		text, ok = user.GetQuestion(user.CurrentQuestion)
		if !ok {
			return errors.New("quiz has no quiestions")
		}

	case command == "/confirm":

		b.users.RLock()
		user, ok := b.users.M[update.Message.From.ID]
		b.users.RUnlock()
		if !ok {
			return errors.New("user not found")
		}

		if user.State != 1 {
			return errors.New("invalid state")
		}

		user.CurrentQuestion++
		if user.CurrentQuestion >= len(user.Questions) {

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
				return errors.New("quiz service error")
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

				sb.WriteString(fmt.Sprintf("%d. %s, with %d p.\n", r.Place, username, r.PointCount))
			}

			text = sb.String()

			break
		}

		text, ok = user.GetQuestion(user.CurrentQuestion)
		if !ok {
			return errors.New("quiz has no quiestions")
		}

	case command == "/pick":

		if len(commandArgs) != 1 {
			return errors.New("bad arguments")
		}

		b.users.RLock()
		user, ok := b.users.M[update.Message.From.ID]
		b.users.RUnlock()
		if !ok {
			return errors.New("user not found")
		}

		if user.State != 1 {
			return errors.New("invalid state")
		}

		ansNum, err := strconv.ParseInt(commandArgs[0], 10, 64)
		if err != nil {
			return errors.New("bad arguments")
		}

		user.Questions[user.CurrentQuestion].AnswerOptions[ansNum-1].Picked = true

		text, ok = user.GetQuestion(user.CurrentQuestion)
		if !ok {
			return errors.New("quiz has no quiestions")
		}

	case command == "/unpick":

		if len(commandArgs) != 1 {
			return errors.New("bad arguments")
		}

		b.users.RLock()
		user, ok := b.users.M[update.Message.From.ID]
		b.users.RUnlock()
		if !ok {
			return errors.New("user not found")
		}

		if user.State != 1 {
			return errors.New("invalid state")
		}

		ansNum, err := strconv.ParseInt(commandArgs[0], 10, 64)
		if err != nil {
			return errors.New("bad arguments")
		}

		user.Questions[user.CurrentQuestion].AnswerOptions[ansNum-1].Picked = false

		text, ok = user.GetQuestion(user.CurrentQuestion)
		if !ok {
			return errors.New("quiz has no quiestions")
		}

	case command == "/gettopbyquiz":

	default:

		text = "What?"

	}

	if useDefaultUrl {
		url = fmt.Sprintf(
			"https://api.telegram.org/bot%s/sendMessage?chat_id=%d&text=%s",
			b.apiKey,
			update.Message.Chat.ID,
			url2.PathEscape(text),
		)
	}

	// TODO может возвращать URL назад в RUN и перенести туда log MW
	_, err := b.httpClient.Post(url, "text/plain", nil)
	if err != nil {
		return err
	}

	log.Printf("Get <%#v> from %d, send him <%#v>", update.Message.Text, update.Message.From.ID, text)

	return nil
}

//https://api.telegram.org/bot5323294543:AAF-uuZS7_SR-j8XBuxTrYmpy_zN26u7qrA/sendPoll?chat_id=339069827&question=whatCarDoULike&options=[%22Ford%22,%20%22BMW%22,%20%22Fiat%22]&allows_multiple_answers=true
//https://api.telegram.org/bot5323294543:AAF-uuZS7_SR-j8XBuxTrYmpy_zN26u7qrA/sendMessage?chat_id=339069827&text=whatCarDoULike&reply_markup=%20{%22inline_keyboard%22:%20[[{%22text%22:%20%22A%22,%22callback_data%22:%20%22A1%22}],[{%22text%22:%20%22B%22,%22callback_data%22:%20%22C1%22}]]}
//https://api.telegram.org/bot5323294543:AAF-uuZS7_SR-j8XBuxTrYmpy_zN26u7qrA/answerCallbackQuery?callback_query_id=%221234%22&text=notificationText
