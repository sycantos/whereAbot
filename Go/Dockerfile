FROM golang:latest

RUN go get "github.com/go-sql-driver/mysql"
RUN go get "github.com/Krognol/go-wolfram"
RUN go get "github.com/christianrondeau/go-wit"
RUN go get "github.com/nlopes/slack"

RUN go build -o whereabot;
ENTRYPOINT ["./whereabot.go"]