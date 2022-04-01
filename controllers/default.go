package controllers

import (
	"net/http"

	"go-practice/libs/jwt"
	"go-practice/libs/tool"
	"go-practice/models"

	"github.com/gin-gonic/gin"

	log "github.com/sirupsen/logrus"
)

func Home(ctx *gin.Context) {
	// kafka producer
	// tool.SendKafkaProducerMessage("broker", "topic", "sync", "test")

	// Server Api Host
	// con := tool.GetZookeeperServerConfig()
	// log.Info(con["send"])

	// EncryptToken
	j := jwt.NewJWT()

	custom := jwt.Custom{
		Uid:      113,
		UserName: "love",
	}

	etoken, _ := j.EncryptToken(custom)
	log.Info(etoken)

	// eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOjExMywidXNlcm5hbWUiOiJsb3ZlIiwiZXhwIjoxNjE0Nzg4MTk4LCJpc3MiOiJtZWl6aWtlYWlAMTYzLmNvbSJ9.koGpHgG1ukECOyTLgmOgTvH5eFPI-ZET_k53-ffO8VQ

	dtoken, _ := j.DecryptToken(etoken)
	log.Info(string(tool.MarshalJson(dtoken)))

	utoken, _ := j.UpdateToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOjExMywidXNlcm5hbWUiOiJsb3ZlIiwiZXhwIjoxNjE0Nzg4MTk4LCJpc3MiOiJtZWl6aWtlYWlAMTYzLmNvbSJ9.koGpHgG1ukECOyTLgmOgTvH5eFPI-ZET_k53-ffO8VQ")
	log.Info(utoken)

	// test redis
	// redis, _ := models.GetUserName()
	// log.Info(redis)

	// test mysql
	// email := "admin@bank.com"
	// mysql, _ := models.GetPerson(email)
	// log.Info(string(tool.MarshalJson(mysql)))

	// mysql, _ := models.GetMySQL()
	// log.Info(string(tool.MarshalJson(mysql)))

	// test log
	log.Error("this is error test")

	// test get
	// getparams := map[string]string{"type": "1"}
	// get, _ := fetch.GET("https://www.test.com/api/test", getparams, nil)
	// log.Info(string(get))

	// test post
	// postparams := map[string]string{"type": "1"}
	// postbody := map[string]interface{}{"uid": 113}

	// post, _ := fetch.POST("https://www.test.com/api/test", postbody, postparams, nil)
	// log.Info(string(post))

	ctx.HTML(http.StatusOK, "index.tmpl", gin.H{
		"title": "GoLang",
	})
}

func NotFound(ctx *gin.Context) {
	ctype := ctx.Request.Header.Get("Content-Type")
	// test := regexp.MustCompile(`^application\/json$`)

	if ctype == "application/json" {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"status":  403,
			"message": "Forbidden",
		})
	} else {
		ctx.HTML(http.StatusOK, "error.tmpl", gin.H{
			"title": "404 page",
		})
	}
}

// ?email=test10@bank.com
func ApiAddPerson(ctx *gin.Context) {
	name := ctx.DefaultQuery("name", "guest")
	email := ctx.Query("email")

	if email != "" {
		person := []string{email, name, "汉族", "男", "11010199812187756", "13412345678", "北京市朝阳区百子湾路苹果社区B区", "100000"}
		lastId, _ := models.AddPerson(person)

		ctx.JSON(http.StatusOK, gin.H{
			"status":  200,
			"lastId":  lastId,
			"message": "Added successfully",
		})
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"status":  400,
			"message": "Add failed",
		})
	}
}

func ApiTest(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "test",
	})
}
