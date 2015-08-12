# xpool - minimalistic implementation for worker pool in Golang

## Highlights 

1. Only channels / goroutines / select are used. ** NOT using "sync" package at all **.
2. Minimalistic and simple code.
3. Simple "Future" pattern implemented.


## Easy to use


```GO
var pool = xpool.New(10, handler_function)
...
...

fu := pool.Run(input)
result, err := fu.Get(time.Second * 5) // Timeout in 5 seconds

//
if err == nil {
	// do your stull with result
}else{
	// Timeout 
	log.Printf("err:", err)
}

```