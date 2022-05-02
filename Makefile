.PHONY:
.SILENT:

build:
	go build -o ./.bin/bot.exe cmd/bot/main.go

run: build
	./.bin/bot.exe

build-image:
	docker build -t go-pocket-bot:v0.1 .

start-container:
	docker run --name pocket-bot -p 80:80 --env-file .env go-pocket-bot:v0.1