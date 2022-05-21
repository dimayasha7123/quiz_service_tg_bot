package app

import (
	"context"
	"fmt"
	"hw2-tgbot/internal/models"
	"log"
	"strings"
)

func (b *bclient) updateHandler(ctx context.Context, update models.Update) error {

	messWords := strings.Split(update.Message.Text, " ")
	command := messWords[0]
	//commandArgs := messWords[1:]

	var url string
	useDefaultUrl := true
	text := "This command not implemented (yet)"

	switch command {

	case "/start":

		b.users.RLock()
		user, ok := b.users.M[update.Message.From.ID]
		b.users.RUnlock()

		var fmtText string

		if ok {
			fmtText = "Welcome, %s. Again"
			// тут можно апдейтить инфу о юзере, если понадобится
			// например, у юзера может поменяться username.
		} else {
			fmtText = "Welcome, %s"
			user = models.NewUser(
				update.Message.From.ID,
				update.Message.From.Username,
			)
			b.users.Lock()
			b.users.M[user.TGID] = user
			b.users.Unlock()
			err := b.repo.AddUser(ctx, user)
			if err != nil {
				return nil
			}
		}

		text = fmt.Sprintf(fmtText, user.Username)

	case "/help":

		text = "There should be help text, but there is no"

	case "/get_quizes":

	case "/start_quiz":

	case "/get_top_by_quiz":

	default:

		text = "What?"

	}

	if useDefaultUrl {
		url = fmt.Sprintf(
			"https://api.telegram.org/bot%s/sendMessage?chat_id=%d&text=%s",
			b.apiKey,
			update.Message.Chat.ID,
			text,
		)
	}

	// TODO может возвращать URL назад в RUN и перенести туда MW
	_, err := b.httpClient.Post(url, "text/plain", nil)
	if err != nil {
		return err
	}

	log.Printf("Get <%s> from %d, send him <%s>", update.Message.Text, update.Message.From.ID, text)

	return nil
}
