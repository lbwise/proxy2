package dest

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type Server struct {
	Port int
	Addr string
}

func SpinServers(ctx context.Context, logger *log.Logger) ([]*Server, error) {
	servers := []*Server{
		{Addr: "localhost", Port: 8080},
		//{Addr: "localhost", Port: 8081},
		//{Addr: "localhost", Port: 8082},
		//{Addr: "localhost", Port: 8083},
	}

	var wg sync.WaitGroup
	for _, srv := range servers {
		go func() {
			defer wg.Done()
			err := createApp(ctx, srv, logger)
			if err != nil {
				logger.Fatalf("Error creating instance: %s", err.Error())
			}
		}()
	}

	wg.Wait()
	return servers, nil
}

func createApp(ctx context.Context, srv *Server, logger *log.Logger) error {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = logger.Writer()
	gin.DefaultErrorWriter = logger.Writer()
	router := gin.Default()

	httpSrv := &http.Server{
		Addr:    "localhost:8080",
		Handler: router,
	}
	logger.Printf("Spinning up instance on: %d\n", srv.Port)

	router.GET("/ping", func(c *gin.Context) {
		logger.Println("RECIEVED", c.Request.Method)
		time.Sleep(1 * time.Second)
		c.String(http.StatusOK, "pong")
	})

	go func() {
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			// we got an actual error here
			logger.Fatalf("Error starting HTTP server: %s", err.Error())
		}
	}()

	// So ctx gets shutdown signal from root
	// But create new timeout ctx in case Shutdown takes too long to sort itself out
	<-ctx.Done()
	logger.Println("Shutting down the server...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return httpSrv.Shutdown(shutdownCtx)
}
