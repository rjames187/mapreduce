package main

type RequestJobArgs struct{}

type RequestJobReply struct {
	Type     string
	FilePath string
}