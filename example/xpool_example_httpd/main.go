package main

import (
    "fmt"
    "net/http"
    "time"
    "github.com/tly1980/xpool"
    "io/ioutil"
    "regexp"
    "strconv"
    "flag"
)

var concurrent = flag.Int("c", 2, 
    "Concurrency number, or pool size. Default is 2")
var timeout = flag.Int("t", 10, 
    "Concurrency number, or pool size. Default is 2")


type url_input struct {
    url string
    sleep time.Duration
}

var pool = xpool.New(*concurrent, func ( i interface{}) interface{} {
    i2 := i.(*url_input)
    resp, err := http.Get(i2.url)

    if i2.sleep > 0 {
        fmt.Printf("Sleep: %v \n", i2.sleep)
        time.Sleep(i2.sleep)
    } else {
        fmt.Printf("No sleep\n")
    }

    if err != nil {
        return nil
    }
    return resp
})

var RE_URL, err = regexp.Compile("/(\\d+)/([^/]+)")

func handler(w http.ResponseWriter, r *http.Request) {
    var uinput *url_input
    result := RE_URL.FindStringSubmatch(r.URL.Path)

    if len(result) > 0 {
        val, _ := strconv.Atoi(result[1])
        uinput = &url_input {
            url: fmt.Sprintf("http://%s", result[2]),
            sleep: time.Second * time.Duration(val),
        }
    }else{
        uinput = &url_input {
            url: fmt.Sprintf("http://%s", r.URL.Path[1:]),
            sleep: 0,
        }
    }

    fmt.Printf("uinput is: %v\n", uinput)
    fu := pool.Run(uinput)
    timeout := time.Second * time.Duration(*timeout)
    r2, err := fu.Get(timeout)
    if err == nil {
        resp := r2.(*http.Response)
        defer resp.Body.Close()
        body, err := ioutil.ReadAll(resp.Body)
        fmt.Fprintf(w, "Hi there, I love %s - %v!", body, err)
    }else {
        fmt.Fprintf(w, "err: %v", err)
    }
}

func main() {
    http.HandleFunc("/", handler)
    http.ListenAndServe(":8080", nil)
}