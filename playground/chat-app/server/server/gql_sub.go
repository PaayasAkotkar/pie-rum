package server

import (
	"app/server/server/graph/model"
	"context"
	"log"
)

func (r *IApp) ChessCoachReply(ctx context.Context, input *model.ChessStudentRequest) (<-chan *model.OnChessCoachReply, error) {
	log.Println("live update request triggered...")
	ch := r.cheetah.Subscribe(*input.ID)

	go func() {
		<-ctx.Done()
		r.cheetah.Unsubscribe(*input.ID, ch)
	}()

	return ch, nil
}
