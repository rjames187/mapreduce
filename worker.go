package main

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"log"
	"mapreduce/plugins"
	"net/rpc"
	"os"
)

type Worker struct {
	masterAddr string
	plugin plugins.Plugin
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
		w.doTask(reply.Job)
	} else {
		fmt.Print("call failed!\n")
		log.Fatal("worker shutting down ...")
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

func (w *Worker) doTask(job *Job) {
	if job.Type == "map" {
		data, err := os.ReadFile(job.FilePath)
		if err != nil {
			log.Fatal(err)
		}
		pairs := w.plugin.Map(string(data))
		partitions := w.partitionIntermediate(pairs, job.NReduce)
		for p, part := range partitions {
			intFilePath := fmt.Sprintf("./mock_fs/m%d-%d.txt", job.Num, p)
			f, err := os.Create(intFilePath)
			if err != nil {
				log.Fatal("error creating new file: ", err)
			}
			defer f.Close()
			err = json.NewEncoder(f).Encode(part)
			if err != nil {
				log.Fatal("error encoding intermediate pairs: ", err)
			}
		}
	}
}

func (w *Worker) partitionIntermediate(pairs []*plugins.KeyValue, nReduce int) map[int][]*plugins.KeyValue {
	res := map[int][]*plugins.KeyValue{}
	for _, p := range pairs {
		h := fnv.New32a()
		h.Write([]byte(p.Key))
		hash := int(h.Sum32() & 0x7fffffff)
		num := hash % nReduce
		res[num] = append(res[num], p)
	}
	return res
}

// instantiate a new worker
func NewWorker(addr string, plugin string) *Worker {
	if plugin == "" {
		log.Fatal("a worker requires a plugin ...")
	}
	return &Worker{masterAddr: addr, plugin: plugins.Plugins[plugin]}
}