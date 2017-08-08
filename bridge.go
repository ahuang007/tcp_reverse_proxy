package main

import(
    "net"
    "fmt"
    "time"
)

type Bridge struct {
    local net.Conn
    remote net.Conn
    buf []byte
    buf_i int
    ip string
    last_time int64
    last_count int
}

func (this *Bridge) Init() {
    this.buf = make([]byte, 1024, 1024)
    this.buf_i = 0
    this.last_time = 0
    this.last_count = 0
}

func (this *Bridge) SetLocal(conn net.Conn, ip string) {
    this.local = conn
    this.ip = ip
}

func (this *Bridge) CheckFrequency() bool {
    now := time.Now().Unix()
    if(now/60 == this.last_time/60) {
        this.last_count++
        fmt.Println("frequency",this.last_count)
        if this.last_count > 20 {
            server.AddBlackIP(this.ip)
            fmt.Println("发包频率过快, 加入黑名单",this.ip)
            return false
        }
    } else {
        this.last_time = now
        this.last_count = 1
    }

    return true
}

func (this *Bridge) Start() {
    defer this.Close()
    fmt.Println("读取第一个包")
    if !this.readFirstPack() {
        fmt.Println("读取第一个包失败")
        return
    }

    fmt.Println("检查地一个包")
    check, n := this.CheckPack()
    if !check || n == 0 {
        fmt.Println("检测第一个包失败")
        return
    }
    fmt.Println("链接服务器")
    if !this.ConnRemote() {
        fmt.Println("链接服务器失败")
        return
    }
    this.Pipe()
}

func (this *Bridge) ConnRemote() bool {
    s, err := net.ResolveTCPAddr("tcp", cfg.Raddress)
    conn, err := net.DialTCP("tcp", nil, s)
    if err != nil {
        fmt.Println("ERROR: Dial failed:", err.Error())
        return false
    }
    this.remote = conn
    // 第一个包发送给服务器
    this.remote.Write(this.buf[:this.buf_i])
    this.buf_i = 0
    fmt.Println("链接服务器，客户端ip",this.ip)
    return true
}

func (this *Bridge) readFirstPack() bool {
    this.local.SetReadDeadline(time.Now().Add(time.Duration(10) * time.Second))
    n, err := this.local.Read(this.buf[:])
    if err != nil || n == 0{
        server.AddEmptyLink(this.ip)
        fmt.Println(n, err)
        return false
    }
    this.buf_i = n

    this.local.SetReadDeadline(time.Time{})
    return true
}

func (this *Bridge) CheckPack() (bool, int) {
    /*start := 0
    for {
        check,n := this.CheckOnePack(start)
        fmt.Println(check, n)
        if !check {
            return false,0
        }
        if n <= 0 {
            return true,start
        }
        start += n
        if start == this.buf_i {
            break
        }
    }*/
    start := this.buf_i

    return true, start
}

func (this *Bridge) AddData(data []byte) bool{
    n := len(data)
    if this.buf_i + n > 1024 {
        server.AddBlackIP(this.ip)
        return false
    }

    for i:=0;i<n;i++{
        this.buf[this.buf_i + i] = data[i]
    }

    this.buf_i += n
    return true
}

func (this *Bridge) CheckOnePack(start int) (bool, int) {
    if this.buf_i <= start + 2 {
        return false, 0
    }

    n := (uint)(this.buf[start + 0])*256 + (uint)(this.buf[start + 1])
    if n > uint(cfg.MaxLen) {
        return false, 0
    }

    if n + 2 > uint(this.buf_i) {
        return true, 0
    }

    // 检查发包频率
    if !this.CheckFrequency() {
        return false, 0
    }

    // 检查包体内容是否合法

    return true, int(n+2)
}

func (this *Bridge) Pipe() {
    local_chan := ChanFromConn(this.local)
    remote_chan := ChanFromConn(this.remote)
    for {
        select {
        case b1 := <-local_chan:
            if b1 == nil {
                fmt.Println("客户端断开连接")
                return
            }
            if !this.AddData(b1) {
                return
            }
            check,n := this.CheckPack()
            if !check {
                fmt.Println("检查客户端包失败")
                return
            }
            if n > 0 {
                this.remote.Write(this.buf[:n])
                for i:=0; i<this.buf_i-n;i++{
                    this.buf[i] = this.buf[n+i]
                }
                this.buf_i -= n
            }
        case b2 := <-remote_chan:
            if b2 == nil {
                fmt.Println("服务器关闭连接")
                return
            }
            fmt.Println("转发服务器消息")
            this.local.Write(b2)
        }
    }
}

func (this *Bridge) Close() {
    fmt.Println("关闭桥", this.ip)
    if this.local != nil {
        this.local.Close()
        this.local = nil
    }

    if this.remote != nil {
        this.remote.Close()
        this.remote = nil
    }
}
