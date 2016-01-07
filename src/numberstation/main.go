package main

import (
    "net/http"
    "flag"
    "encoding/json"
    "net"
    "io"
    "fmt"
    "log"
    "math/rand"
    "encoding/binary"

    "github.com/bemasher/rtltcp"
)

var (
    httpListen = flag.String("http-listen", ":8080", "HTTP listen string")
    rtlTcp = flag.String("rtl-tcp", "127.0.0.1:1234", "IP:Port for rtl_tcp daemon")
)

type (
    GenerateRandom struct {}
)

func get_random(amount int, sdr *rtltcp.SDR, ch chan []byte) {

    random := make([]byte, amount)

    val, err := io.ReadFull(sdr, random)
    if err != nil {
        log.Print(val)
        log.Fatal("Error reading samples:", err)
    }

    ch <- random
}

func (h *GenerateRandom) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    sdr := new(rtltcp.SDR)

    addr, err := net.ResolveTCPAddr("tcp4", *rtlTcp)
    if err != nil {
        log.Fatal("invalid address", err)
    }

    sdr.Connect(addr)
    defer sdr.Close()

    sdr.SetCenterFreq(1420405751)
    //sdr.SetSampleRate(5e6)

    ch := make(chan []byte)
    go get_random(32, sdr, ch)
    hash := <-ch

    num := int64(binary.BigEndian.Uint64(hash))
    rand.Seed(num)

    w.Header().Add("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{ "hash": fmt.Sprintf("%x", rand.Int63())})
}


func main() {
    flag.Parse();

    http.Handle("/", new(GenerateRandom))

    http.ListenAndServe(*httpListen, nil)
}

