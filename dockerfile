FROM golang:1.21-alpine AS base


RUN apk add -U tzdata
RUN apk --update add ca-certificates

WORKDIR /app
COPY . .

RUN go mod download
RUN go mod verify

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

RUN go build -ldflags="-s -w" -o /server .

FROM scratch

COPY --from=base /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=base /etc/passwd /etc/passwd
COPY --from=base /etc/group /etc/group

COPY --from=base /server .

EXPOSE 3000

CMD [ "/server" ]