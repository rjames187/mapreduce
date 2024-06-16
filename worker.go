package main

import (
	"fmt"
	"log"
	"net/rpc"
)

type Worker struct {
	masterAddr string
}

func (w *Worker) work() {
	w.CallRequestJob()
}

func (w *Worker) CallRequestJob() {
	args := RequestJobArgs{}
	reply := RequestJobReply{}
	ok := w.call("Coordinator.RequestJob", &args, &reply)
	if ok {
		fmt.Printf("Received a job: %v\n", *reply.Job)
	} else {
		fmt.Print("call failed!\n")
	}
}

func (w *Worker) call(rpcName string, args interface{}, reply interface{}) bool {
	c, err := rpc.DialHTTP("tcp", w.masterAddr)
	if err != nil {
		log.Fatal("dialing: ", err)
	}
	defer c.Close()

	err = c.Call(rpcName, args, reply)
	if err == nil {
		return true
	}
	fmt.Println(err)
	return false
}

// instantiate a new worker
func NewWorker(addr string) *Worker {
	return &Worker{masterAddr: addr}
}