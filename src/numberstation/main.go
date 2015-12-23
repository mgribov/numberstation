package main

import (
    "net/http"
    "flag"
    "encoding/json"
    "net"
    "io"
    "log"
    "fmt"
    //"strconv"

    "github.com/bemasher/rtltcp"
)

var (
    httpListen = flag.String("http-listen", ":8080", "HTTP listen string")
    rtlTcp = flag.String("rtl-tcp", "127.0.0.1:1234", "IP:Port for rtl_tcp daemon")
)

type (
    GenerateRandom struct {}
)

func get_random(amount int, rtl *rtltcp.SDR, ch chan []byte) {
    random := make([]byte, amount)

    _, err := io.ReadFull(rtl, random)
    if err != nil {
        log.Fatal("Error reading samples:", err)
    }

    log.Print(random)

    ch <- random
}

func (h *GenerateRandom) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    amount := 10

    /*
    args := r.URL.Query()
    l := args["l"][0]
    if l != "" {
        amount, err := strconv.Atoi(l)
    } else {
    }
    */

    sdr := new(rtltcp.SDR)

    addr, err := net.ResolveTCPAddr("tcp4", *rtlTcp)
    if err != nil {
        log.Fatal("invalid address", err)
    }

    sdr.Connect(addr)
    defer sdr.Close()

    sdr.SetCenterFreq(1420e6)
    //sdr.SetSampleRate(5e6)

    ch := make(chan []byte)
    go get_random(amount, sdr, ch)
    hash := <-ch

    w.Header().Add("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{ "hash": fmt.Sprintf("%x", hash)})
}


func main() {
    flag.Parse();

    http.Handle("/", new(GenerateRandom))

    http.ListenAndServe(*httpListen, nil)
}

