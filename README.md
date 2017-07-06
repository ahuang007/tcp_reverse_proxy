tcp反向代理，包含以下功能：

    空连接、发包过大、发包频率过快、包内容不对

1、编译

    make
    
2、配置 cfg.json

    Laddress:":3000",             // 本地监听端口

    Raddress:"127.0.0.1:4000",    // 服务器监听端口
    
    MaxLen: 512,                  // 最大包长
    
    PackFrequncy: 3,              // 每分钟发包频率

3、启动
    
    ./main cfg.json
