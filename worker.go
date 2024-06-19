package main

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"log"
	"mapreduce/plugins"
	"net/rpc"
	"os"
	"time"
)

type Worker struct {
	masterAddr string
	plugin plugins.Plugin
}

func (w *Worker) work() {
	for {
		w.CallRequestJob()
	}
}

func (w *Worker) CallRequestJob() {
	args := RequestJobArgs{}
	reply := RequestJobReply{}
	ok := w.call("Coordinator.RequestJob", &args, &reply)
	if ok {
		// an empty job is a signal to wait
		if reply.Job == nil {
			duration, _ := time.ParseDuration("1s")
			time.Sleep(duration)
			return
		}
		fmt.Printf("Received a job: %v\n", *reply.Job)
		w.doTask(reply.Job)
	} else {
		fmt.Print("call failed!\n")
		log.Fatal("worker shutting down ...")
	}
}

func (w *Worker) CallCompleteMapJob(id int, filePaths map[int][]string) {
	args := CompleteMapJobArgs{Id: id}
	reply := CompleteMapJobReply{}
	ok := w.call("Coordinator.CompleteMapJob", &args, &reply)
	if ok {
		log.Printf("Worker successfully sent completion notice of map job #%d", id)
	} else {
		log.Fatal("Worker failed to send completion notice of map job")
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
		data, err := os.ReadFile(job.FilePaths[0])
		if err != nil {
			log.Fatal(err)
		}
		pairs := w.plugin.Map(string(data))
		partitions := w.partitionIntermediate(pairs, job.NReduce)
		filePaths := map[int][]string{}
		for p, part := range partitions {
			intFilePath := fmt.Sprintf("./mock_fs/m%d-%d.txt", job.Id, p)
			_, present := filePaths[p]
			if present {
				filePaths[p] = append(filePaths[p], intFilePath)
			} else {
				filePaths[p] = []string{intFilePath}
			}
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
		w.CallCompleteMapJob(job.Id, filePaths)
	} else if job.Type == "reduce" {
		intermediate_pairs := []*plugins.KeyValue{}
		for _, fp := range job.FilePaths {
			f, err := os.Open(fp)
			if err != nil {
				log.Fatal(err)
			}
			part_pairs := []*plugins.KeyValue{}
			err = json.NewDecoder(f).Decode(&part_pairs)
			if err != nil {
				log.Fatal(err)
			}
			intermediate_pairs = append(intermediate_pairs, part_pairs...)
		}
		final_pairs := w.plugin.Reduce(intermediate_pairs)
		f, err := os.Create(fmt.Sprintf("./mock_fs/o%d.txt", job.Id))
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		for _, pair := range final_pairs {
			f.Write([]byte(fmt.Sprintf("%v: %v\n", pair.Key, pair.Value)))
		}
		fmt.Printf("Reduce task #%d completed ...", job.Id)
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