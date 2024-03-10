# Include variables from the .envrc file
include .envrc

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## run/api: run the cmd/auth application
.PHONY: run/api
run/api:
	go run ./cmd/api \
		-limiter-enabled \
		-auth-key=${AUTH_KEY} \
		-smtp-host=${SMTP_HOST} \
		-smtp-port=${SMTP_PORT} \
		-smtp-username=${SMTP_USER} \
		-smtp-password=${SMTP_PASS} \
		-smtp-sender=${SMTP_FROM}
