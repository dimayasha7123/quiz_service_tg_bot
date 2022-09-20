package app

import (
	"context"
	"fmt"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (b *bclient) getQuizesHandler(ctx context.Context) (string, error) {

	quizes, err := b.quizClient.GetQuizList(ctx, &emptypb.Empty{})
	if err != nil {
		return "", err
	}
	text := ""
	for _, q := range quizes.QList {
		text += fmt.Sprintf(
			"%s\nStart: /startquiz_%d\nTop: /gettopbyquiz_%d\n\n",
			q.Name,
			q.ID,
			q.ID,
		)
	}

	return text, nil
}
