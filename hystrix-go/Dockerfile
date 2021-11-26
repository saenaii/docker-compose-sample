FROM golang:1.16-alpine AS builder

WORKDIR /app/
COPY . .

RUN apk add git && \
	go env -w GO111MODULE=off && \
	go env -w GOPROXY=https://goproxy.cn,https://gocenter.io,https://goproxy.io,direct && \
	go get -u github.com/afex/hystrix-go/hystrix

RUN go build -o main

FROM alpine:latest
COPY --from=builder /app/main .
CMD ["./main"]