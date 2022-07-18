package server

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/eviltomorrow/rogue/lib/httpmiddleware"
	"github.com/gin-gonic/gin"
)

var (
	Port int

	server *http.Server
	Router = &gin.Engine{}
)

func StartupHTTP() error {
	Router.Use(gin.Recovery())
	Router.Use(httpmiddleware.NewLogger())

	server = &http.Server{
		Addr:    fmt.Sprintf(":%d", Port),
		Handler: Router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Printf("[F] HTTP Server startup failure, nest error: %v", err)
			os.Exit(1)
		}
	}()

	return nil
}

func ShutdownHTTP() error {
	if server != nil {
		server.Close()
	}
	return nil
}
