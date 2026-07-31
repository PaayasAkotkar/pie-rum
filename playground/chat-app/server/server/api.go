package server

import (
	"log"

	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
)

func (a *IApp) Monitor() {
	a.app.GET("/ask-chess-coach", func(ctx *gin.Context) {
		log.Println("get")
		a.srv.ServeHTTP(ctx.Writer, ctx.Request)
	})
	a.app.POST("/ask-chess-coach", func(ctx *gin.Context) {
		log.Println("post")
		a.srv.ServeHTTP(ctx.Writer, ctx.Request)
	})
}

func (a *IApp) Playground() {
	a.app.GET("/playground", func(ctx *gin.Context) {
		h := playground.Handler("GraphQL", "/query")
		h.ServeHTTP(ctx.Writer, ctx.Request)
	})
}

func (a *IApp) Listen(addr string) error {
	return a.app.Run(addr)
}
