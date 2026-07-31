package server

import (
	"app/server/server/graph"
)

type RagService struct{ *IApp }

func (r *IApp) Mutation() graph.MutationResolver { return &RagService{r} }

func (r *IApp) Query() graph.QueryResolver { return &RagService{r} }

func (r *IApp) Subscription() graph.SubscriptionResolver { return &RagService{r} }
