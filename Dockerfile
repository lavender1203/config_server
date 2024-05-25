# Dockerfile
# 使用官方的 Golang 镜像作为构建环境
FROM golang:1.20 as builder

# 设置工作目录
WORKDIR /

# 将 go.mod 和 go.sum 文件复制到工作目录
COPY go.mod go.sum ./

# 下载所有的依赖包 配置国内源
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go mod download

# 将源代码复制到工作目录
COPY . .

# 编译程序
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .
# init/main.go也需要编译
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o init_config init/main.go

# # 使用 scratch 镜像作为运行环境
FROM scratch
# 使用 alpine 镜像作为运行环境
# FROM alpine

# 从构建环境复制程序和 .env 文件到运行环境
COPY --from=builder /main /main
COPY --from=builder /init_config /init_config

# # .env文件在init/flight_ticket/.env目录中
# COPY --from=builder /init/flight_ticket/.env /init/flight_ticket/.env
# # .env文件在init/data_process/.env目录中
# COPY --from=builder /init/data_process/.env /init/data_process/.env

