package main

import (
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

var (
	ErrKeyNotProvided = errors.New("Authorization key must be provided")
)

func (app *application) rateLimit() gin.HandlerFunc {
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	go func() {
		for {
			time.Sleep(time.Minute)

			mu.Lock()

			for ip, client := range clients {
				if time.Since(client.lastSeen) > time.Minute*3 {
					delete(clients, ip)
				}
			}

			mu.Unlock()
		}
	}()

	return func(ctx *gin.Context) {
		ip := ctx.RemoteIP()

		mu.Lock()

		if _, found := clients[ip]; !found {
			clients[ip] = &client{limiter: rate.NewLimiter(rate.Limit(app.config.limiter.rps), app.config.limiter.burst)}
		}

		clients[ip].lastSeen = time.Now()

		if !clients[ip].limiter.Allow() {
			mu.Unlock()
			app.rateLimitExceededResponse(ctx)
			ctx.Abort()
			return
		}

		mu.Unlock()

		ctx.Next()
	}
}

func (app *application) authenticate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		header := ctx.GetHeader("Authorization")

		split := strings.Split(header, " ")
		if len(split) != 2 || len(split[1]) == 0 {
			app.badRequestResponse(ctx, ErrKeyNotProvided)
			return
		}

		if split[1] != app.config.authKey {
			app.invalidTokenResponse(ctx)
			return
		}

		ctx.Next()
	}
}
