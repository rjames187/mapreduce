package main

import (
	"flag"
)

func main() {
	role := flag.String("r", "", "a role for the node")
	dir := flag.String("d", "", "the directory where the input files are stored")
	plugin := flag.String("p", "", "a string referring to the chosen map and reduce functions")
	addr := flag.String("a", "", "the ip address and port number of the master node")
	nReduce := flag.Int("nr", 2, "the number of files to output")

	flag.Parse()
	if *role == "master" {
		coordinator := NewCoordinator(*nReduce)
		coordinator.coordinate(*dir, *addr)
	}
	if *role == "worker" {
		worker := NewWorker(*addr, *plugin)
		worker.work()
	}
	if *role == "sequential" {
		sequentialMapReduce(*dir, *plugin)
	}
}