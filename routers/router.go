package routers

import (
	"go-practice/controllers"

	"github.com/gin-gonic/gin"
)

var logic = controllers.NewLogic()

func HandleRouter(router *gin.Engine) *gin.Engine {
	AddDefaultRouter(router)
	AddApiRouter(router)

	return router
}
