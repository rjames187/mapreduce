package main

type RequestJobArgs struct{}

type RequestJobReply struct {
	Job *Job
}

type CompleteMapJobArgs struct {
	Id        int
	FilePaths map[int][]string
}

type CompleteMapJobReply struct{}

type CompleteReduceJobArgs struct {
	Id int
}

type CompleteReduceJobReply struct{}