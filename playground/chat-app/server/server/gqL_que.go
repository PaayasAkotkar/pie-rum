package server

import (
	"app/server/server/graph/model"
	"context"
	"log"
)

func (r *IApp) InitCoachSuggestions(ctx context.Context, input *model.ChessStudentRequest) ([]*model.OnChessCoachReply, error) {
	log.Println("Rag init request made...")
	return []*model.OnChessCoachReply{}, nil
}
