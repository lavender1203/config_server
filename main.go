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
	"sync"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedConfigServiceServer
	cli    *clientv3.Client
	cache  map[string]string
	mu     sync.RWMutex
	ctx    context.Context
	cancel context.CancelFunc
}

func (s *server) GetConfig(ctx context.Context, in *pb.GetConfigRequest) (*pb.GetConfigResponse, error) {
	// 构造 etcd 的 key
	etcdKey := fmt.Sprintf("%s/%s", in.Section, in.Key)
	log.Printf("GetConfig key: %s", etcdKey)

	// 先从缓存中获取
	s.mu.RLock()
	value, ok := s.cache[etcdKey]
	s.mu.RUnlock()
	if ok {
		log.Printf("GetConfig key: %s, value: %s (from cache)", etcdKey, value)
		return &pb.GetConfigResponse{Value: value}, nil
	}

	// 如果缓存中没有，再从 etcd 中获取
	resp, err := s.cli.Get(ctx, etcdKey)
	if err != nil {
		return nil, err
	}
	if len(resp.Kvs) == 0 {
		return nil, fmt.Errorf("no value found for key %s", etcdKey)
	}
	value = string(resp.Kvs[0].Value)

	// 将获取到的值存入缓存
	s.mu.Lock()
	s.cache[etcdKey] = value
	s.mu.Unlock()

	log.Printf("GetConfig key: %s, value: %s", etcdKey, value)
	return &pb.GetConfigResponse{Value: value}, nil
}

func (s *server) watchEtcd() {
	rch := s.cli.Watch(s.ctx, "", clientv3.WithPrefix())
	for wresp := range rch {
		for _, ev := range wresp.Events {
			etcdKey := string(ev.Kv.Key)
			s.mu.Lock()
			if ev.Type == clientv3.EventTypeDelete {
				delete(s.cache, etcdKey)
			} else {
				s.cache[etcdKey] = string(ev.Kv.Value)
			}
			s.mu.Unlock()
		}
	}
}

func (s *server) SetConfig(ctx context.Context, in *pb.SetConfigRequest) (*pb.SetConfigResponse, error) {
	// 构造 etcd 的 key
	etcdKey := fmt.Sprintf("%s/%s", in.Section, in.Key)
	_, err := s.cli.Put(ctx, etcdKey, in.Value)
	if err != nil {
		return &pb.SetConfigResponse{Success: false}, err
	}

	// 更新本地缓存
	s.mu.Lock()
	s.cache[etcdKey] = in.Value
	s.mu.Unlock()

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
	ctx, cancel := context.WithCancel(context.Background())
	s := &server{
		cli:    cli,
		cache:  make(map[string]string), // 初始化 cache
		ctx:    ctx,
		cancel: cancel,
	}
	go s.watchEtcd()

	// 从环境变量获取端口号
	port := os.Getenv("GRPC_SERVER_LISTEN_PORT")
	// 如果没有设置端口号，则默认为50051
	if port == "" {
		port = "50051"
	}
	// 启动gRPC服务
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// 打印日志
	log.Printf("gRPC server listening at %v", lis.Addr())
	grpcServer := grpc.NewServer()
	pb.RegisterConfigServiceServer(grpcServer, s) // 使用已经创建并初始化的 server 实例
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
