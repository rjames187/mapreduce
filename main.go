package main

import "flag"

func main() {
	role := flag.String("r", "", "a role for the node")
	dir := flag.String("d", "", "the directory where the input files are stored")
	plugin := flag.String("p", "", "a string referring to the chosen map and reduce functions")

	flag.Parse()

	if *role == "sequential" {
		sequentialMapReduce(*dir, *plugin)
	}
}