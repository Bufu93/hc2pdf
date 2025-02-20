install:
    which modd || go install github.com/cortesi/modd/cmd/modd@latest

build:
    go mod download && CGO_ENABLED=0 GOOS=linux go build -o ./.bin/app ./cmd/main.go

run: build
    docker-compose up --remove-orphans app