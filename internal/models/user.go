package models

import (
	"fmt"
	"strings"
)

type User struct {
	ID              int64
	TGID            int64
	QSID            int64
	Username        string
	State           int64 // 0 - in lobby, 1 - in quiz
	Questions       []Question
	CurrentQuestion int
	QuizPartyID     int64
}

type Question struct {
	Title         string
	AnswerOptions []AnswerOption
}

type AnswerOption struct {
	Title  string
	Picked bool
}

func NewUser(chatID int64, username string) *User {
	u := &User{
		TGID:     chatID,
		Username: username,
	}
	return u
}

func (u *User) GetQuestion(n int) (string, bool) {
	if n >= len(u.Questions) {
		return "", false
	}

	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("Question #%d:\n", n+1))
	sb.WriteString(fmt.Sprintf("%s\n", u.Questions[n].Title))
	sb.WriteString(fmt.Sprintf("\nOptions:\n"))
	for j, ansOpt := range u.Questions[n].AnswerOptions {
		pickSign := '❌'
		if ansOpt.Picked {
			pickSign = '✅'
		}
		sb.WriteString(fmt.Sprintf(
			"%c  %s     /pick_%d     /unpick_%d\n",
			pickSign,
			ansOpt.Title,
			j+1,
			j+1,
		))
	}
	sb.WriteString("\nConfirm question: /confirm")

	return sb.String(), true
}
