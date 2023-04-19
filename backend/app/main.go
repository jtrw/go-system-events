package main

import (
    //"fmt"
    "log"
    "github.com/jessevdk/go-flags"
    server "github.com/jtrw/go-events/v1/backend/app/server"
)

type Options struct {
   Host string `short:"h" long:"host" default:"localhost" description:"Host web server"`
   Port string `short:"p" long:"port" default:"3000" description:"Port web server"`
}

func main() {
    var opts Options
    parser := flags.NewParser(&opts, flags.Default)
    _, err := parser.Parse()
    if err != nil {
        log.Fatal(err)
    }

    srv := server.Server {
        Host: opts.Host,
        Port: opts.Port,
        PinSize: 1,
        WebRoot: "/",
        Version: "1.0",
    }

    if err := srv.Run(); err != nil {
        log.Printf("[ERROR] failed, %+v", err)
    }

    log.Printf("[INFO] Activate rest server")
    log.Printf("[INFO] Host: %s", opts.Host)
    log.Printf("[INFO] Port: %s", opts.Port)
}
