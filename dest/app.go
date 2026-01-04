package dest

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Server struct {
	Port int
	Addr string
}

func SpinServers(logger *log.Logger) ([]*Server, error) {
	servers := []*Server{
		{Addr: "localhost", Port: 8080},
		{Addr: "localhost", Port: 8081},
		{Addr: "localhost", Port: 8082},
		{Addr: "localhost", Port: 8083},
	}

	for _, srv := range servers {
		createApp(srv, logger)
	}

	return servers, nil
}

func createApp(srv *Server, logger *log.Logger) {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = logger.Writer()
	gin.DefaultErrorWriter = logger.Writer()
	router := gin.Default()
	logger.Println("SPINNING UP DEST INSTANCE ON :%d", srv.Port)

	router.GET("/ping", func(c *gin.Context) {
		time.Sleep(100 * time.Millisecond)
		logger.Println("RECIEVED", c.Request.Method)
		c.String(http.StatusOK, "pong")
	})

	go func() {
		router.Run(fmt.Sprintf(":%d", srv.Port))
	}()
}
