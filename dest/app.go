package dest

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func StartApp() {
	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		fmt.Println("RECIEVED", c.Request.Method)
		c.String(http.StatusOK, "pong")
	})

	router.Run()
}
