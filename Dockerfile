FROM golang:1.22-alpine3.19 AS builder

WORKDIR /go/src
COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -ldflags="-s -w" -o main .

FROM alpine:3.19

RUN apk update && apk upgrade && \
    apk --no-cache add ca-certificates && \
	apk add --no-cache tzdata && \
	ln -snf "/usr/share/zoneinfo/$TZ" /etc/localtime && echo "$TZ" > /etc/timezone
ENV TZ=Asia/Bangkok
COPY --from=builder /go/src/main /
COPY --from=builder /go/src/certs /certs
ENTRYPOINT ["/main"]