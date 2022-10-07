package app

import (
	"context"
	"fmt"
	"log"
	url2 "net/url"
	"strings"

	"github.com/dimayasha7123/quiz_service_tg_bot/internal/models"
)

func (b *bclient) updateRouter(ctx context.Context, update models.Update) (string, error) {

	messageWords := strings.Split(update.Message.Text, "_")
	command := messageWords[0]
	commandArgs := messageWords[1:]

	var url string
	text := "This command not implemented (yet)"
	var err error

	switch {
	case command == "/start":
		text, err = b.startHandler(ctx, update)
		if err != nil {
			return "", err
		}

	case command == "/help":
		text = "There should be help text, but there is no"

	case command == "/getquizes" ||
		(command == "/startquiz" && len(commandArgs) == 0) ||
		(command == "/gettopbyquiz" && len(commandArgs) == 0):
		text, err = b.getQuizesHandler(ctx)
		if err != nil {
			return "", err
		}

	case command == "/startquiz":
		text, err = b.startQuizHandler(ctx, update, commandArgs)
		if err != nil {
			return "", err
		}

	case command == "/confirm":
		text, err = b.confirmHandler(ctx, update)
		if err != nil {
			return "", err
		}

	case command == "/pick":
		text, err = b.pickUnpickHandler(ctx, update, commandArgs, true)
		if err != nil {
			return "", err
		}

	case command == "/unpick":
		text, err = b.pickUnpickHandler(ctx, update, commandArgs, false)
		if err != nil {
			return "", err
		}

	case command == "/gettopbyquiz":
		text, err = b.getTopByQuizHandler(ctx, update, commandArgs)
		if err != nil {
			return "", err
		}

	default:
		text = "What?"
	}

	url = fmt.Sprintf(
		"https://api.telegram.org/bot%s/sendMessage?chat_id=%d&text=%s",
		b.apiKey,
		update.Message.Chat.ID,
		url2.PathEscape(text),
	)

	log.Printf("Get <%#v> from %d, send him <%#v>", update.Message.Text, update.Message.From.ID, text)

	return url, nil
}

//https://api.telegram.org/bot5323294543:AAF-uuZS7_SR-j8XBuxTrYmpy_zN26u7qrA/sendPoll?chat_id=339069827&question=whatCarDoULike&options=[%22Ford%22,%20%22BMW%22,%20%22Fiat%22]&allows_multiple_answers=true
//https://api.telegram.org/bot5323294543:AAF-uuZS7_SR-j8XBuxTrYmpy_zN26u7qrA/sendMessage?chat_id=339069827&text=whatCarDoULike&reply_markup=%20{%22inline_keyboard%22:%20[[{%22text%22:%20%22A%22,%22callback_data%22:%20%22A1%22}],[{%22text%22:%20%22B%22,%22callback_data%22:%20%22C1%22}]]}
//https://api.telegram.org/bot5323294543:AAF-uuZS7_SR-j8XBuxTrYmpy_zN26u7qrA/answerCallbackQuery?callback_query_id=%221234%22&text=notificationText
