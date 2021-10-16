FROM golang:alpine AS builder

# 为我们的镜像设置必要的环境变量
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY . .

RUN go build -o app ./services

FROM scratch

# 从builder镜像中把/build/app 拷贝到当前目录
COPY --from=builder /build/app /
COPY --from=builder /build/services/router/data /data

EXPOSE 8000

# 需要运行的命令
ENTRYPOINT ["/app", "router", "-d", "/data"]