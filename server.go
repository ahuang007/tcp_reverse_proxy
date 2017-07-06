package main

import (
    "net"
    "fmt"
    "log"
    "time"
)

type EmptyLink struct{
    time int64
    count int
}

type Server struct {
    l *net.TCPListener
    black_list map[string]bool
    empty_links map[string]*EmptyLink
}

func (this *Server) AddBlackIP(ip string) {
    this.black_list[ip] = true
}

func (this *Server) AddEmptyLink(ip string) {
    _,ok := this.black_list[ip]
    if ok {
        return
    }

    now := time.Now().Unix()
    c,err1 := this.empty_links[ip]
    if !err1 {
        this.empty_links[ip] = &EmptyLink{time: now, count: 1}
    } else {
        if(now/60 == c.time/60) {
            c.count++
            if c.count > 20 {
                this.black_list[ip] = true
                log.Printf("ip:%s add to black list, reason: empty link", ip)
                return
            }
        } else {
            c.time = now
            c.count = 1
        }
    }
}

func (this *Server) Listen(addr string)  error {
    s, err := net.ResolveTCPAddr("tcp", addr)

    if err != nil {
        panic(fmt.Sprintf("ResolveTCPAddr failed:%v", err))
        return err
    }

    l, err := net.ListenTCP("tcp", s)

    if err != nil {
        panic(fmt.Sprintf("can't listen on %s,%v", addr, err))
        return err
    }

    this.l = l
    return nil
}

func (this *Server) AcceptLoop() {
    for {
        if c, err := this.l.Accept(); err == nil && c != nil {
            // 是否在黑名单中
            ip, _, _ := net.SplitHostPort(c.RemoteAddr().String())
            fmt.Printf("新连接 addr=%s\n", ip)
            _,ok := this.black_list[ip]
            if ok {
                fmt.Printf("连接在黑名单中，关闭")
                c.Close()
                continue
            }
            go this.HandleConn(c, ip)
        } else {
            fmt.Printf("ERROR: couldn't accept: %v", err)
        }
    }
}

func (this *Server) HandleConn(c net.Conn, ip string) {
    b := &Bridge{}
    b.Init()
    b.SetLocal(c, ip)
    b.Start()
}
