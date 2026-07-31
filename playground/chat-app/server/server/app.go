package server

import (
	rcrypt "app/server/crypt"
	"app/server/server/graph/model"
	"pie-rum-sdk/cheetah"
	"pie-rum-sdk/pie-rum/client"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/gin-gonic/gin"
)

type IApp struct {
	app     *gin.Engine
	cli     *client.PieRum
	srv     *handler.Server
	cheetah *cheetah.Cheetah[string, model.OnChessCoachReply]
}

func New(a *gin.Engine) *IApp {
	cli, err := client.New(rcrypt.PIERUMADRESS, nil)
	if err != nil {
		panic(err)
	}
	return &IApp{
		app:     a,
		cli:     cli,
		cheetah: cheetah.New[string, model.OnChessCoachReply](100),
	}
}
