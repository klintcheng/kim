FROM golang AS builder

# 为我们的镜像设置必要的环境变量
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY . .

# RUN go-wrapper install -ldflags "-linkmode external -extldflags -static"   
# RUN go build -a -ldflags '-linkmode external -extldflags "-static"' -o app ./services

RUN go build -o app ./services

FROM scratch

# 从builder镜像中把/build/app 拷贝到当前目录
COPY --from=builder /build/app /

EXPOSE 8080

# 需要运行的命令
ENTRYPOINT ["/app", "royal"]