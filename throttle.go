package main

import "flag"
import "net"
import "log"
import "sync"
import "os"
import "time"

func clientHandler(listenSide net.Conn, connectPath string, latency int) {
    log.Printf("Client connected, connecting to %s", connectPath)
    defer listenSide.Close()

    // Connect to other socket
    connectSide, err := net.Dial("unix", connectPath)
    if err != nil {
        log.Fatal("dial error:", err)
    }
    defer connectSide.Close()

    var wg sync.WaitGroup

    // Read from listenSide and write to connectSide
    wg.Add(1)
    go func() {
        defer wg.Done()
        var buffers = make(chan []byte, 128)
        for {
            buf := make([]byte, 65536)
            n, err := listenSide.Read(buf)
            if err != nil {
                log.Printf("Error reading from listenSide: %s", err)
                return
            }
//             log.Printf("Read %d bytes from listenSide", n)
            buffers <- buf[:n]

            go func() {
                time.Sleep(time.Duration(latency) * time.Millisecond)
                buf := <-buffers
                _, err = connectSide.Write(buf)
//                 log.Printf("w: Read %d bytes from listenSide", len(buf))
                if err != nil {
                    log.Printf("Error writing to connectSide: %s", err)
                    return
                }
            }()
        }
    }()

    // Read from connectSide and write to listenSide
    wg.Add(1)
    go func() {
        defer wg.Done()
        var buffers = make(chan []byte, 128)
        for {
            buf := make([]byte, 65536)
            n, err := connectSide.Read(buf)
            if err != nil {
                log.Printf("Error reading from connectSide: %s", err)
                return
            }
//             log.Printf("Read %d bytes from connectSide", n)
            buffers <- buf[:n]

            go func() {
                time.Sleep(time.Duration(latency) * time.Millisecond)
                buf := <-buffers
                _, err = listenSide.Write(buf)
//                 log.Printf("w: Read %d bytes from listenSide", len(buf))
                if err != nil {
                    log.Printf("Error writing to listenSide: %s", err)
                    return
                }
            }()
        }
    }()

    wg.Wait()
    log.Printf("Client disconnected")
}

func main() {
    // Parse command line args
    listenPath := flag.String("listen", "", "Path to listen on")
    connectPath := flag.String("connect", "", "Path to connect to")
    latency := flag.Int("latency", 0, "Latency in milliseconds")
    flag.Parse()

    if *listenPath == "" {
        log.Fatal("Must specify --listen")
    }
    if *connectPath == "" {
        log.Fatal("Must specify --connect")
    }
    if *latency == 0 {
        log.Fatal("Must specify --latency")
    }

    // Create listener socket
    if err := os.RemoveAll(*listenPath); err != nil {
        log.Fatal(err)
    }

    l, err := net.Listen("unix", *listenPath)
    if err != nil {
        log.Fatal("listen error:", err)
    }
    defer l.Close()

    log.Printf("Ready")

    // Accept new connections, dispatching them to clientHandler in a goroutine.
    for {
        conn, err := l.Accept()
        if err != nil {
            log.Fatal("accept error:", err)
        }

        go clientHandler(conn, *connectPath, *latency)
    }
}
