FROM alpine:latest AS root-certs

RUN apk add -U --no-cache ca-certificates
RUN addgroup -g 1001 app
RUN adduser app -u 1001 -D -G app /home/app

FROM golang:1.23 AS builder 

# Set the working directory
WORKDIR /source

# Copy and download dependencies
COPY --from=root-certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=vendor -o ./yt-server ./app/./...

FROM scratch AS final

WORKDIR /app

COPY --from=root-certs  /etc/passwd /etc/passwd
COPY --from=root-certs /etc/group /etc/group
COPY --chown=1001:1001 --from=root-certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --chown=1001:1001 --from=builder source/yt-server /yt-server

USER app


EXPOSE 10101

ENTRYPOINT [ "/yt-server" ]

