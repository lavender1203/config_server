package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "go.etcd.io/etcd/client/v3"
    "gopkg.in/ini.v1"
)

func InitConfig(appName string, filePath string) {
    // 读取.env文件
    cfg, err := ini.Load(filePath)
    if err != nil {
        log.Fatalf("Fail to read file: %v", err)
    }

    // 连接到etcd
    etcdHost := cfg.Section("Etcd").Key("ETCD_HOST").MustString("localhost")
    etcdPort := cfg.Section("Etcd").Key("ETCD_PORT").MustInt(2379)

    cli, err := clientv3.New(clientv3.Config{
        Endpoints:   []string{fmt.Sprintf("%s:%d", etcdHost, etcdPort)},
        DialTimeout: 5 * time.Second,
    })
    if err != nil {
        log.Fatal(err)
    }
    defer cli.Close()

    // 将配置放到etcd
    for _, section := range cfg.Sections() {
        for _, key := range section.Keys() {
            _, err := cli.Put(context.Background(), fmt.Sprintf("%s/%s/%s", appName, section.Name(), key.Name()), key.String())
            if err != nil {
                log.Fatal(err)
            } else {
                log.Printf("Put %s/%s/%s = %s", appName, section.Name(), key.Name(), key.String())
            }
        }
    }
}

func main() {
    InitConfig("DATA_PROCESS", "init/data_process/.env")
    InitConfig("FLIGHT_TICKET", "init/flight_ticket/.env")
}