package main

import (
	"net/http"

	cgin "github.com/actorgo-game/actorgo/components/gin"
	cherrySnowflake "github.com/actorgo-game/actorgo/extend/snowflake"
	"github.com/gin-gonic/gin"
)

type Test1Controller struct {
	cgin.BaseController
}

func (t *Test1Controller) Init() {
	t.GET("/", t.index)
	t.Engine.GET("/panic", t.panic)
	t.GET("/render_result", t.renderResult)

	cherrySnowflake.SetDefaultNode(1)
}

func (t *Test1Controller) index(c *cgin.Context) {
	c.RenderHTML("this is index... " + cherrySnowflake.Next().String())
}

func (t *Test1Controller) panic(c *gin.Context) {
	c.String(http.StatusOK, "test panic")
	panic("test panic!")
}

func (t *Test1Controller) renderResult(c *cgin.Context) {
	str := cherrySnowflake.Next().Base58()
	c.RenderJSON(str)
}
