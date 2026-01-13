FROM golang:1.25.5-alpine

WORKDIR /app

RUN apk --no-cache add git

RUN go install github.com/air-verse/air@latest

COPY go.mod go.sum* ./
RUN go mod download

COPY . .

CMD ["air", "-c", ".air.toml"]
