// main.go
package main

import (
	pb "config/protobuf"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
)

type server struct {
    pb.UnimplementedConfigServiceServer
    cli *clientv3.Client
}

func (s *server) GetConfig(ctx context.Context, in *pb.GetConfigRequest) (*pb.GetConfigResponse, error) {
    // 构造 etcd 的 key
    etcdKey := fmt.Sprintf("%s/%s", in.Section, in.Key)
    log.Printf("GetConfig key: %s", etcdKey)
    resp, err := s.cli.Get(ctx, etcdKey)
    if err != nil {
        return nil, err
    }
    if len(resp.Kvs) == 0 {
        return nil, fmt.Errorf("no value found for key %s", etcdKey)
    }
    value := resp.Kvs[0].Value
    log.Printf("GetConfig key: %s, value: %s", etcdKey, value)
    return &pb.GetConfigResponse{Value: string(value)}, nil
}

func (s *server) SetConfig(ctx context.Context, in *pb.SetConfigRequest) (*pb.SetConfigResponse, error) {
    // 构造 etcd 的 key
    etcdKey := fmt.Sprintf("%s/%s", in.Section, in.Key)
    _, err := s.cli.Put(ctx, etcdKey, in.Value)
    if err != nil {
        return &pb.SetConfigResponse{Success: false}, err
    }
    return &pb.SetConfigResponse{Success: true}, nil
}

func main() {
    // 连接到 etcd
    // 从环境变量中读取 etcd 的 endpoints
    endpoints := os.Getenv("ETCD_ENDPOINTS")
    if endpoints == "" {
        log.Fatalf("ETCD_ENDPOINTS environment variable is not set")
    }

    cli, err := clientv3.New(clientv3.Config{
        Endpoints:   strings.Split(endpoints, ","),
        DialTimeout: 5 * time.Second,
    })
    // 打印详细日志
    log.Printf("etcd client connected")
    if err != nil {
        log.Fatalf("Fail to connect to etcd: %v", err)
    }
    defer cli.Close()

    // 启动gRPC服务
    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }
    // 打印日志
    log.Printf("gRPC server listening at %v", lis.Addr())
    s := grpc.NewServer()
    pb.RegisterConfigServiceServer(s, &server{cli: cli})
    if err := s.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}