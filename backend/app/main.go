package main

import (
    "fmt"
    "log"
    "net/http"
)

func homePage(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Home Page")
}


func setupRoutes() {
    manager := NewManager()

    http.HandleFunc("/", homePage)
    http.HandleFunc("/ws", manager.serveWS)
}

func main() {
    fmt.Println("Hello World")
    setupRoutes()
    log.Fatal(http.ListenAndServe(":8080", nil))
}
