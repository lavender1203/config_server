

首先生成Python3需要使用的代码

python3 -m grpc_tools.protoc -I protobuf/ --python_out=./protobuf --grpc_python_out=./protobuf protobuf/config.proto

生成Go需要使用的服务端代码
protoc --go_out=./protobuf --go_opt=paths=source_relative --go-grpc_out=./protobuf --go-grpc_opt=paths=source_relative protobuf/config.proto


注意：
需要修改config_pb2_grpc.py文件，以解决导入出错的问题
# import config_pb2 as config__pb2
from . import config_pb2 as config__pb2


# 部署
## 构建镜像
sudo docker build -t config_server:1.0.0 .
## 启动容器
sudo docker-compose up -d 


## 初始化配置命令
sudo docker exec -it config_server ./init