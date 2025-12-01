// internal/app/controller/default.go
package controller

import (
	"go-practice/internal/dto"
	"go-practice/internal/pkg/ginctx"
	"go-practice/internal/pkg/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (l *Controller) TestGet(ctx *gin.Context) {
	l.Bind(ctx).Success(ctx, gin.H{"example": "for love"})
}

func (l *Controller) TestPost(c *gin.Context) {
	ctx := ginctx.New(c)
	l.app.Log.Info("get user", []zap.Field{zap.String("ip", ctx.GetClientIP())}...)

	// ciphertext, _ := l.app.Crypto.Encrypt("AbcDefg8886")
	// plaintext, _ := l.app.Crypto.Decrypt(ciphertext)
	// l.app.Log.Info("plaintext", zap.String("plaintext", plaintext))
	// l.app.Log.Info("ciphertext", zap.String("ciphertext", ciphertext))

	// token, _ := l.app.Jwt.GenerateToken(7758258, 30*time.Minute)
	// claims, err := l.app.Jwt.ParseToken(token)
	// if err != nil {
	// 	l.Bind(c).Fail(c, response.CodeUnauthorized)
	// 	return
	// }
	// l.app.Log.Info("token", zap.String("token", token))
	// subject, _ := claims.GetSubject()
	// l.app.Log.Info("subject", zap.String("subject", subject))

	var req dto.GetUserReq
	if err := c.ShouldBindQuery(&req); err != nil {
		l.Bind(c).Fail(c, response.CodeUnprocessableEntity)
		return
	}

	// id := cast.ToInt64(c.Param("id"))
	// if id == 0 {
	// 	id = 6
	// }
	// res, err := l.app.Repository.FindByID(c.Request.Context(), id)

	// if err != nil {
	// 	l.Bind(c).Error(c, err)
	// 	return
	// }

	l.Bind(c).Success(c, gin.H{"example": "for love"})
}
