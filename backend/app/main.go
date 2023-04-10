package main

import (
    "fmt"
    "log"
    "net/http"
    "github.com/jessevdk/go-flags"
)

type Options struct {
   Host string `short:"h" long:"host" default:"127.0.0.1" description:"Host web server"`
   Port string `short:"p" long:"port" default:"8080" description:"Port web server"`
}

func homePage(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Home Page")
}


func setupRoutes() {
    manager := NewManager()

    http.HandleFunc("/", homePage)
    http.HandleFunc("/ws", manager.serveWS)
}

func main() {
    var opts Options
    parser := flags.NewParser(&opts, flags.Default)
    _, err := parser.Parse()
    if err != nil {
        log.Fatal(err)
    }


   // setupRoutes()
   // log.Fatal(http.ListenAndServe(":8080", nil))


     srv := Server {
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
