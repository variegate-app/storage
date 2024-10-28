PID=/tmp/.${app}.pid
MAKEFLAGS += --silent

start-docker-dev: 
	@docker-compose -f ./.build/delve.docker-compose.yaml up -d --build

build-app: 
	@bash -c "go build -o ./.bin/${app} ./cmd/${app} && chmod +x ./.bin/${app}"
	
test:
	@go test -cover ./...