package server

import (
	"app/server/admin"
	"app/server/server/graph/model"
	"context"
	"log"
	rumrpc "pie-rum-sdk/misc/rum"
)

func (a *IApp) Config(ctx context.Context, input *model.IConfigRequest) (*model.IConfigResp, error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Panic recovered in Config: %v", r)
		}
	}()
	log.Println("Config made ...")
	x := &rumrpc.IConfigRequest{
		ProfileName:    input.ProfileName,
		Config:         conGQLToRPC(input.Config),
		KitName:        input.KitName,
		ServiceName:    input.ServiceName,
		DispatcherName: input.DispatcherName,
		EventName:      input.EventName,
	}
	cli := a.cli

	var err error
	var resp *rumrpc.IConfigResponse
	switch input.Depth {
	case 1:
		resp, err = cli.UpdateProfileConfig(ctx, x)
	case 2:
		resp, err = cli.UpdateKitConfig(ctx, x)
	case 4:
		resp, err = cli.UpdateServiceConfig(ctx, x)
	case 5:
		resp, err = cli.UpdateDispatcherConfig(ctx, x)
	case 6:
		resp, err = cli.UpdateEventConfig(ctx, x)
	}
	if err != nil {
		return &model.IConfigResp{
			Error:   resp.Error,
			Succeed: resp.Succeed,
		}, err
	}

	return &model.IConfigResp{Succeed: true, Error: ""}, nil
}

func (a *IApp) AskChessCoach(ctx context.Context, input *model.ChessStudentRequest) (*model.ChessCoachPayload, error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Panic recovered in AskChessCoach: %v", r)
		}
	}()
	log.Println("RAG request made...")
	cli := a.cli

	reqx := admin.Req{
		Student: input,
		//Store:   r.store,
	}

	cli.POST(ctx, &rumrpc.IPostRequest{
		Post: []*rumrpc.IPost{
			{
				Profile: &rumrpc.ISequence{
					Profile: "nvidea",
					Input:   reqx.Pack(),
				},
				Push: true,
			},
		},
	})
	status := int32(200)
	msg := "coaching thinking..."
	return &model.ChessCoachPayload{Status: &status, Message: &msg}, nil
}
