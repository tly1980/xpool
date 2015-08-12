package main

import (
    "fmt"
    "net/http"
    "time"
    "github.com/tly1980/xpool"
    "io/ioutil"
    "regexp"
    "strconv"
)

type url_input struct {
    url string
    sleep time.Duration
}

var pool = xpool.New(10, func ( i interface{}) interface{} {
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

var re_url, err = regexp.Compile("/(\\d+)/([^/]+)")

func handler(w http.ResponseWriter, r *http.Request) {
    result := re_url.FindStringSubmatch(r.URL.Path)
    var uinput *url_input

    fmt.Printf("Path: %v\n", r.URL.Path)
    fmt.Printf("result: %v\n", result)

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
    timeout := time.Second * time.Duration(3)
    r2, err := fu.Get(timeout)
    if err == nil {
        resp := r2.(*http.Response)
        defer resp.Body.Close()
        body, err := ioutil.ReadAll(resp.Body)
        fmt.Fprintf(w, "Hi there, I love %s - %v!", body, err)
    }else {
        fmt.Fprintf(w, "Timeout: %v", timeout)
    }
}

func main() {
    http.HandleFunc("/", handler)
    http.ListenAndServe(":8080", nil)
}