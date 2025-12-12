// internal/app/handler/common.go
package handler

import (
	"go-practice/internal/pkg/ginctx"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (h *Handler) TestGet(c *gin.Context) {
	h.Bind(c).Success(c, gin.H{"example": "for love"})
}

func (h *Handler) TestPost(c *gin.Context) {
	h.app.Service.GetUser(ctx, 7758258)

	h.app.Log.Info("get user", []zap.Field{zap.String("ip", ginctx.GetClientIP(c))}...)

	// ciphertext, _ := h.app.Crypto.Encrypt("AbcDefg8886")
	// plaintext, _ := h.app.Crypto.Decrypt(ciphertext)
	// h.app.Log.Info("plaintext", zap.String("plaintext", plaintext))
	// h.app.Log.Info("ciphertext", zap.String("ciphertext", ciphertext))

	// token, _ := h.app.Jwt.GenerateToken(7758258, 30*time.Minute)
	// claims, err := h.app.Jwt.ParseToken(token)
	// if err != nil {
	// 	h.Bind(c).Fail(c, response.CodeUnauthorized)
	// 	return
	// }
	// h.app.Log.Info("token", zap.String("token", token))
	// subject, _ := claims.GetSubject()
	// h.app.Log.Info("subject", zap.String("subject", subject))

	// var req dto.GetUserReq
	// if err := c.ShouldBindQuery(&req); err != nil {
	// 	h.Bind(c).Fail(c, response.CodeUnprocessableEntity)
	// 	return
	// }

	// id := cast.ToInt64(c.Param("id"))
	// if id == 0 {
	// 	id = 6
	// }
	// _, err := h.app.Repository.FindByID(c.Request.Context(), id)

	// if err != nil {
	// 	h.Bind(c).Error(c, err)
	// 	return
	// }

	h.Bind(c).Success(c, gin.H{"example": "for love"})
}
