package main

import (
	"context"
	"flag"
	"sync"
)

const (
	//MaxCrons max cron numbers
	MaxCrons = 10000
)

var port int

//Run base struct
type Run struct {
	Ctx      context.Context
	Cancel   context.CancelFunc
	MuxCrons struct {
		sync.RWMutex
		Crons chan *Cron
	}
}

func init() {
	flag.IntVar(&port, "port", 8888, "listen port")
}

//InitRun init run
func InitRun() *Run {
	ctx, cancel := context.WithCancel(context.Background())
	return &Run{
		Ctx:    ctx,
		Cancel: cancel,
		MuxCrons: struct {
			sync.RWMutex
			Crons chan *Cron
		}{
			Crons: make(chan *Cron, MaxCrons),
		},
	}
}

func main() {
	flag.Parse()

	r := InitRun()

	r.StartHTTP()

}
