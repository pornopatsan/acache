FROM golang:latest

LABEL maintainer="Hamlet Avetikyan <hamlet.avetikyn@gmail.com>"

WORKDIR /app

COPY . .

RUN mkdir -p "bin"
RUN go build -o ./bin ./cmd/acache

CMD ["./bin/acache"]
