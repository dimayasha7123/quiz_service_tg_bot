package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	url2 "net/url"
	"time"

	"github.com/dimayasha7123/quiz_service_tg_bot/internal/models"
)

func (b *bclient) Run(ctx context.Context) error {

	users, err := b.repo.GetUsers(ctx)
	if err != nil {
		return err
	}

	b.users.Lock()
	for _, user := range users {
		_, ok := b.users.M[user.TGID]
		if ok {
			return errors.New("few users with equal ID gotten from repo")
		}
		b.users.M[user.TGID] = user
	}
	b.users.Unlock()

	var lastUpdateId int64

	for {

		url := fmt.Sprintf(
			"https://api.telegram.org/bot%s/getUpdates?offset=%d",
			b.apiKey,
			lastUpdateId+1,
		)

		resp, err := b.httpClient.Get(url)
		if err != nil {
			return err
		}

		bytes, err := ioutil.ReadAll(resp.Body)

		updates := models.Updates{}
		err = json.Unmarshal(bytes, &updates)
		if err != nil {
			return err
		}

		if updates.Ok {
			for _, update := range updates.Result {

				postUrl, err := b.updateRouter(ctx, update)
				if err != nil {
					postUrl = fmt.Sprintf(
						"https://api.telegram.org/bot%s/sendMessage?chat_id=%d&text=%s",
						b.apiKey,
						update.Message.Chat.ID,
						url2.PathEscape(fmt.Sprintf("Ooops, something was wrong.\nError: %v", err)),
					)
				}
				_, err = b.httpClient.Post(postUrl, "text/plain", nil)
				if err != nil {
					return err
				}

			}
		}

		if len(updates.Result) != 0 {
			lastUpdateId = updates.Result[len(updates.Result)-1].UpdateID
		}

		err = resp.Body.Close()
		if err != nil {
			return err
		}

		time.Sleep(50 * time.Millisecond)
	}
}
