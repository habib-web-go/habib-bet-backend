NAME=habib-bet-backend
VERSION=0.0.1
DOCKER_DEVELOPMENT=docker compose --file docker/docker-compose-development.yml --project-name habib-bet

.PHONY: build
## build: Compile the packages.
build:
	@go build -o $(NAME)

.PHONY: run
## run: Build and Run in development mode.
run: build
	@./$(NAME) --profile development

.PHONY: run-prod
## run-prod: Build and Run in production mode.
run-prod: build
	@./$(NAME) --profile production

.PHONY: clean
## clean: Clean project and previous builds.
clean:
	@rm -f $(NAME)
	@@$(DOCKER_DEVELOPMENT) down -v

.PHONY: deps
## deps: Download modules
deps:
	@go mod tidy
	@go mod download
	@$(DOCKER_DEVELOPMENT) up -d

.PHONY: help
all: help
# help: show this help message
help: Makefile
	@echo
	@echo " Choose a command to run in "$(APP_NAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
