package xpool;

import (
    "time"
    "errors"
    "log"
    "fmt"
)


type Worker struct {
    in chan interface{}
    out chan interface{}
}

func (w *Worker) loop(job func(interface{}) interface{}){
    for in_para := range w.in {
        w.out <- job(in_para)
    }
}

type Future struct{
    pool *XPool
    input interface{}
    worker *Worker
}

func (fu *Future) run() {
    fu.worker = fu.pool.Borrow()
    fu.worker.in <- fu.input
}

func (fu *Future) Get(duration time.Duration) (interface{}, error) {
    select {
        case r := <- fu.worker.out:
            fu.pool.Return(fu.worker)
            fu.worker = nil
            return r, nil
        case <-time.After(duration):
            // discard the timeout worker, and produce a new one instead
            msg := fmt.Sprintf("Timeout: %v", duration)
            log.Println(msg)
            fu.pool.produce_worker()
            fu.worker = nil
            return nil, errors.New(msg)
    }
}

type XPool struct{
    size int
    workers chan *Worker
    jobHandler func(interface{}) interface{}
}

func newWorker() *Worker {
    return &Worker {
        in: make(chan interface{}),
        out: make(chan interface{}),
    }
}


func New(size int, jobHandler func(interface{}) interface{}) *XPool{
    xpool := &XPool { 
        size: size, 
        workers: make(chan *Worker, size),
        jobHandler: jobHandler,
    }
    xpool.size = size;
    go xpool.start()
    return xpool;
}

func (xp *XPool) produce_worker(){
    w := newWorker()
    go w.loop(xp.jobHandler)
    xp.workers <- w
}

func (xp *XPool) start(){
    for i := 0; i < xp.size; i++ {
        xp.produce_worker()
    }
}

func (xp *XPool) Run(input interface{}) *Future{
    ret := &Future{
        pool: xp,
        input: input,
        worker: nil,
    }

    ret.run()

    return ret
}

func (xp *XPool) Borrow() *Worker{
    w := <-xp.workers
    return w
}

func (xp *XPool) Return(w *Worker) {
    xp.workers <- w
}
