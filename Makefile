NAME=habib-bet-backend
VERSION=0.0.1
DOCKER_DEVELOPMENT=docker compose --file docker/docker-compose-development.yml --project-name habib-bet-dev
DOCKER_PRODUCTION=docker compose --file docker/docker-compose-production.yml --project-name habib-bet

.PHONY: build
## build: Compile the packages.
build:
	@go build -o $(NAME)


.PHONY: build-docker
## build-docker: Build docker image.
build-docker:
	@docker build . -f docker/Dockerfile -t habib-bet/backend:latest

.PHONY: run-dev
## run-dev: Build and Run in development mode.
run-dev: build
	@./$(NAME) --profile development --command run_server

.PHONY: run-prod
## run-prod: Build and Run in production mode in docker.
run-prod: build-docker
	@$(DOCKER_PRODUCTION) up -d

.PHONY: create-contest-prd
## create-contest-prd: Build and Run create contest command in docker container.
create-contest-prd: build-docker
	@$(DOCKER_PRODUCTION) run migrate --profile production --command create_contest


.PHONY: clean-dev
## clean-dev: Clean project and previous builds and clear db
clean-dev:
	@rm -f $(NAME)
	@@$(DOCKER_DEVELOPMENT) down -v

.PHONY: deps-dev
## deps-dev: Download and sync modules and init db
deps-dev:
	@go mod tidy
	@go mod download
	@$(DOCKER_DEVELOPMENT) up -d

.PHONY: deps
## deps: Download modules
deps:
	@go mod download

.PHONY: help
all: help
# help: show this help message
help: Makefile
	@echo
	@echo " Choose a command to run in "$(APP_NAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
