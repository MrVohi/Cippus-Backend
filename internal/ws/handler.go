package ws

import (
	"cippus-backend/internal/services"

	"github.com/gin-gonic/gin"
)

func ServeWS(hub *Hub, secret string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.Query("token")
		if token == "" {
			ctx.AbortWithStatus(401)
			return
		}

		claims, err := services.ValidateAccessToken(token, secret)
		if err != nil {
			ctx.AbortWithStatus(401)
			return
		}
		userID := uint(claims["userID"].(float64))

		conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
		if err != nil {
			return
		}

		client := &Client{
			Hub:    hub,
			UserID: userID,
			Conn:   conn,
			Send:   make(chan []byte, 256),
		}
		hub.Register <- client

		go client.writePump()
		go client.readPump()
	}
}
