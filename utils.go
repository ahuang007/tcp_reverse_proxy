package main

import(
    "net"
)

func ChanFromConn(conn net.Conn) chan []byte {
    c := make(chan []byte)
    go func() {
        b := make([]byte, 1024)
        for {
            if n, err := conn.Read(b); err != nil {
                c <- nil
                break
            } else {
                c <- b[:n]
            }
        }
    }()
    return c
}
