FROM golang:1.16-alpine AS builder
WORKDIR /app/
COPY . .
RUN go env -w GO111MODULE=on \
&& go env -w GOPROXY=https://goproxy.cn,https://gocenter.io,https://goproxy.io,direct \
&& go mod download
RUN go build -o main

FROM saenaii/saenai-image:latest
WORKDIR /cmd/
COPY --from=builder /app/main .
CMD ["./main"]