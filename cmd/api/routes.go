package main

import (
	"github.com/gin-gonic/gin"
)

func (app *application) routes() *gin.Engine {
	router := gin.Default()

	if app.config.limiter.enabled {
		router.Use(app.rateLimit())
	}

	router.Use(app.authenticate())

	router.HandleMethodNotAllowed = true
	router.NoMethod(app.methodNotAllowedResponse)
	router.NoRoute(app.notFoundResponse)

	router.POST("/send-email", app.sendEmailHandler)

	return router
}
