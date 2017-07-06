package main

import (
    "os"
    "fmt"
)

var (
    server = &Server{
      black_list: map[string]bool{},
      empty_links: map[string]*EmptyLink{},
    }
)

func init() {
}

func main() {
    if len(os.Args) != 2 {
        panic("please identify a config file")
        return
    }

    LoadCfg(os.Args[1])
    fmt.Print(cfg)

    fmt.Printf("listen on %s\n", cfg.Laddress)
    if server.Listen(cfg.Laddress) != nil {
        panic(fmt.Sprintf("ERROR: couldn't start listening"))
    }

    server.AcceptLoop()

}

