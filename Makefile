.PHONY: all clean chatserver vet

all: chatserver

chatserver:
	go build -o chatserver github.com/rodrigc/chat_application/cmd/chatserver

vet:
	go vet ./...

clean:
	go clean
	rm -f chatserver

