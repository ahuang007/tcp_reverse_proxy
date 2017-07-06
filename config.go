package main

import (
    "encoding/json"
    "os"
    "fmt"
)

type Config struct {
    Laddress string     `json:"Laddress"`
    Raddress string     `json:"Raddress"`
    MaxLen int          `json:"MaxLen"`
    PackFrequncy int    `json:"PackFrequncy"`
}

var cfg *Config = &Config{
  Laddress:":3000",             // 本地监听端口
  Raddress:"127.0.0.1:4000",    // 服务器监听端口
  MaxLen: 512,                  // 最大包长
  PackFrequncy: 3,              // 每分钟发包频率
}

func LoadCfg(path string) {
    file, err := os.Open(path)
    if err != nil {
        panic("fail to read config file: " + path)
        return
    }

    defer file.Close()

    fi, _ := file.Stat()

    buff := make([]byte, fi.Size())
    _, err = file.Read(buff)
    fmt.Println(buff)
    buff = []byte(os.ExpandEnv(string(buff)))

    err = json.Unmarshal(buff, cfg)
    if err != nil {
        panic("failed to unmarshal file")
        return
    }
}

