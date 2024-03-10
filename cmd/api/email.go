package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *application) sendEmailHandler(ctx *gin.Context) {
	var input struct {
		Email     string `json:"email" binding:"required,email"`
		Subject   string `json:"subject" binding:"required"`
		PlainText string `json:"plaintext"`
		HTML      string `json:"html"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil || (input.PlainText == "" && input.HTML == "") {
		app.badRequestResponse(ctx, err)
		return
	}

	app.background(func() {
		err := app.mailer.Send(input.Email, input.Subject, input.PlainText, input.HTML)
		if err != nil {
			app.logger.PrintError(err, nil)
		}
	})

	ctx.JSON(http.StatusAccepted, gin.H{"status": "accepted"})
}
