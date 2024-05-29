import grpc
import os
from protobuf import config_pb2
from protobuf import config_pb2_grpc


# 获取配置
# 从环境变量中获取 gRPC 服务器的地址
grpc_server_address = os.getenv('GRPC_SERVER_ADDRESS', '192.168.10.151:50051')
# 创建一个全局的 gRPC 连接
channel = grpc.insecure_channel(grpc_server_address)
# 创建一个全局的 stub
stub = config_pb2_grpc.ConfigServiceStub(channel)

def get_config(app_name, section, key):
    # 调用 GetConfig 方法
    response = stub.GetConfig(config_pb2.GetConfigRequest(section=f"{app_name}/{section}", key=key))
    return response.value

# 获取配置
app_name = 'DATA_PROCESS'
etcd_host = get_config(app_name, 'Etcd', 'ETCD_HOST')
etcd_port = get_config(app_name, 'Etcd', 'ETCD_PORT')

kafka_bootstrap_servers = get_config(app_name, 'Kafka', 'bootstrap.servers')
kafka_group_id = get_config(app_name, 'Kafka', 'group.id')
kafka_auto_offset_reset = get_config(app_name, 'Kafka', 'auto.offset.reset')
kafka_enable_auto_commit = get_config(app_name, 'Kafka', 'enable.auto.commit')

redis_ip_ports = get_config(app_name, 'Redis', 'REDIS_IP_PORTS')
redis_password = get_config(app_name, 'Redis', 'REDIS_PASSWORD')
redis_db = get_config(app_name, 'Redis', 'REDIS_DB')
redis_service_name = get_config(app_name, 'Redis', 'REDIS_SERVICE_NAME')
redis_kwargs = get_config(app_name, 'Redis', 'REDIS_KWARGS')

mysql_host = get_config(app_name, 'Mysql', 'MYSQL_HOST')
mysql_port = get_config(app_name, 'Mysql', 'MYSQL_PORT')
mysql_user = get_config(app_name, 'Mysql', 'MYSQL_USER')
mysql_passwd = get_config(app_name, 'Mysql', 'MYSQL_PASSWORD')
mysql_db = get_config(app_name, 'Mysql', 'MYSQL_DB')

# print(f"Etcd host: {etcd_host}, Etcd port: {etcd_port}")
# print(f"Kafka bootstrap_servers: {kafka_bootstrap_servers} group_id: {kafka_group_id} auto_offset_reset: {kafka_auto_offset_reset} enable_auto_commit: {kafka_enable_auto_commit}")
# print(f"Redis ip_ports: {redis_ip_ports} password: {redis_password} db: {redis_db} service_name: {redis_service_name} kwargs: {redis_kwargs}")
# print(f"MySQL host: {mysql_host} port: {mysql_port} user: {mysql_user} passwd: {mysql_passwd} db: {mysql_db}")

import concurrent.futures

# ...

def get_config_and_print(app_name, section, key):
    value = get_config(app_name, section, key)
    print(f"{app_name}/{section}/{key}: {value}")

# ...

with concurrent.futures.ThreadPoolExecutor(max_workers=2000) as executor:
    for i in range(100):
        executor.submit(get_config_and_print, app_name, 'Etcd', 'ETCD_HOST')
        executor.submit(get_config_and_print, app_name, 'Etcd', 'ETCD_PORT')
        executor.submit(get_config_and_print, app_name, 'Kafka', 'bootstrap.servers')
        executor.submit(get_config_and_print, app_name, 'Kafka', 'group.id')
        executor.submit(get_config_and_print, app_name, 'Kafka', 'auto.offset.reset')
        executor.submit(get_config_and_print, app_name, 'Kafka', 'enable.auto.commit')
        executor.submit(get_config_and_print, app_name, 'Redis', 'REDIS_IP_PORTS')
        executor.submit(get_config_and_print, app_name, 'Redis', 'REDIS_PASSWORD')
        executor.submit(get_config_and_print, app_name, 'Redis', 'REDIS_DB')
        executor.submit(get_config_and_print, app_name, 'Redis', 'REDIS_SERVICE_NAME')
        executor.submit(get_config_and_print, app_name, 'Redis', 'REDIS_KWARGS')
        executor.submit(get_config_and_print, app_name, 'Mysql', 'MYSQL_HOST')
        executor.submit(get_config_and_print, app_name, 'Mysql', 'MYSQL_PORT')
        executor.submit(get_config_and_print, app_name, 'Mysql', 'MYSQL_USER')
        executor.submit(get_config_and_print, app_name, 'Mysql', 'MYSQL_PASSWORD')
        executor.submit(get_config_and_print, app_name, 'Mysql', 'MYSQL_DB')