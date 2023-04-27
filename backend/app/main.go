package main

import (
    //"fmt"
    "log"
    "github.com/jessevdk/go-flags"
    server "github.com/jtrw/go-events/v1/backend/app/server"
)

type Options struct {
   Listen string `short:"l" long:"listen" default:"0.0.0.0:3000" description:"Default 0.0.0.0:3000 for localhost"`
}

func main() {
    var opts Options
    parser := flags.NewParser(&opts, flags.Default)
    _, err := parser.Parse()
    if err != nil {
        log.Fatal(err)
    }

    srv := server.Server {
        Listen: opts.Listen,
        PinSize: 1,
        WebRoot: "/",
        Version: "1.0",
    }

    if err := srv.Run(); err != nil {
        log.Printf("[ERROR] failed, %+v", err)
    }
}
