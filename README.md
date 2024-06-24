# Distributed MapReduce

## Getting Started

1. Make sure Go is installed (this project uses 1.19)
2. Make sure GNU Make is installed
3. If you are on a Windows machine, install Git Bash
4. Clone the repository to your machine
5. Inside the root directory, run `make build` or `go build` to compile the binary
6. Run `./run.sh` to execute a word count job with three workers
7. Run `make test_wc` and `make test_fault` to run tests for correctness and fault-tolerance

The test for fault-tolerance deliberately injects faults in the form of crashing or slow workers.

## What is MapReduce?

MapReduce is a programming model and framework for processing massive datasets in parallel. It was first developed by Google engineers in 2003 to perform huge compute jobs that could not happen on a single machine. MapReduce provides a simple interface that hides the complexities of ensuring fault-tolerance in a distributed system.

A MapReduce system consists of one master node and a bunch of worker nodes. The master node acts as a coordinator by assigning tasks to workers. The master node also handles machine failures by re-assigning jobs sent to workers that stopped responding.

## What is this project?

This project is an implementation of the MapReduce framework discussed in the original [2004 whitepaper](https://pdos.csail.mit.edu/6.824/papers/mapreduce.pdf) with a few differences:

- This project only implements the core features discussed in the paper
- The "user-defined" map and reduce functions are not truly user-defined as they are part of the source code in the plugins folder
- This implementation uses the operating system's filesystem as a shared filesystem between all processes rather than a true distributed filesystem like GFS or AWS S3
- Because of the above limitation, this project cannot run in a true distributed cluster, however, it could be extended to connect to a distributed filesystem

## Acknowledgements

This project was inspired by the MapReduce assignment in MIT graduate distributed systems course. While I was influenced by tiny portions of the starter code from the course, the source code in this repository is almost entirely my own.
