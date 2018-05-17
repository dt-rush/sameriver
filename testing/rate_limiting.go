package main

import (
    "sync"
    "fmt"
    "time"
)


func writer_a (c chan(bool), stopWriting chan(bool)) {

writeloop:
    for {
        select {
        case <-stopWriting:
            fmt.Println ("[stop writing]")
            break writeloop
        default:
            fmt.Println (len(c))
            if len(c) == 32 {
                fmt.Println ("full")
            } else {
                c <- true
            }
            time.Sleep(5 * time.Millisecond)
        }
    }
}

type SignalRateLimiter struct {
    mutex sync.Mutex
    guard sync.Once
    out chan(bool)
    delay time.Duration
}

func (r *SignalRateLimiter) send (signal bool) {
    r.mutex.Lock()
    r.guard.Do(func() {
        r.out <- signal
        go func() {
            time.Sleep(r.delay)
            r.mutex.Lock()
            r.guard = sync.Once{}
            r.mutex.Unlock()
        }()
    })
    r.mutex.Unlock()
}

func writer_b (c chan(bool), stopWriting chan(bool)) {

    signalRateLimiter := SignalRateLimiter{
        guard: sync.Once{},
        out: c,
        delay: time.Millisecond * 200}

writeloop:
    for {
        select {
        case <-stopWriting:
            fmt.Println ("[stop writing]")
            break writeloop
        default:
            fmt.Println (len(c))
            signalRateLimiter.send (true)
            time.Sleep(50 * time.Millisecond)
        }
    }
}



func reader_a (c chan(bool), stopReading chan(bool)) {

    rate_limiter := sync.Once{}

    react := func() {
        fmt.Println("[react] yo!")
        go func() {
            time.Sleep(1 * time.Second)
            rate_limiter = sync.Once{}
        }()
    }

readloop:
    for {
        select {
        case <-stopReading:
            fmt.Println ("[stop reading]")
            break readloop
        case <-c:
            fmt.Println("[read] !")
            rate_limiter.Do(react)
        default:
            fmt.Println("[read] _")
            time.Sleep(20 * time.Millisecond)
        }
    }
}

func reader_b (c chan(bool), stopReading chan(bool)) {


    react := func() {
        fmt.Println("[react] yo!")
        time.Sleep(1 * time.Second)
        n := len (c)
        for i := 0; i < n; i++ {
            <-c
        }
    }

readloop:
    for {
        select {
        case <-stopReading:
            fmt.Println ("[stop reading]")
            break readloop
        case <-c:
            fmt.Println("[read] !")
            react()
        default:
            fmt.Println("[read] _")
            time.Sleep(20 * time.Millisecond)
        }
    }

}

func main() {
    c := make(chan(bool), 32)
    stop := make(chan(bool),2)
    go writer_a(c, stop)
    go reader_b(c, stop)
    time.Sleep(time.Second * 10)
    stop <- true
    stop <- true
}
