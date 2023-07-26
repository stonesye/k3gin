package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"k3gin/app/config"
)

func SESSMiddleware() gin.HandlerFunc {
	// 创建基于cookie的存储引擎，SESSION_SECRET 参数是用于加密的密钥，可以随便填写
	cfg := config.C.SESSION

	store := cookie.NewStore([]byte(cfg.Secret))

	// 给Session的存储引擎设置范围和时间
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   7 * 86400,
		Secure:   false,
		HttpOnly: true,
	})
	// gin-session 就是session的名字，也是cookie的名字
	return sessions.Sessions("gin-session", store)
}
