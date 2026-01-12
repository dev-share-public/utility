FROM golang:1.25

RUN go install github.com/air-verse/air@latest

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

CMD ["air", "-c", "docker/air.toml", "-build.poll=true"]