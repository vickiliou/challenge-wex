FROM golang:1.21-alpine

WORKDIR /app

RUN apk --no-cache add gcc musl-dev sqlite

ENV CGO_ENABLED=1

COPY . .

RUN go build cmd/main.go

EXPOSE 8082

ENTRYPOINT [ "./main"]