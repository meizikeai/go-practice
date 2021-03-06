package routers

import "github.com/gin-gonic/gin"

func HandleRouter(router *gin.Engine) *gin.Engine {
	AddApiRouter(router)
	AddWebRouter(router)

	return router
}
