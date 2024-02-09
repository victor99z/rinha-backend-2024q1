FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

RUN go build -ldflags="-s -w" -o apiserver .

FROM scratch AS production

COPY --from=builder ["/build/apiserver", "/build/.env", "/"]

EXPOSE 3000

ENTRYPOINT ["/apiserver"]