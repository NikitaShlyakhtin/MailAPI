# Include variables from the .envrc file
include .envrc

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## run/api: run the cmd/auth application
.PHONY: run/api
run/api:
	go run ./cmd/api 
