# Distributed MapReduce Implementation

## Getting Started

1. Make sure Go is installed (this project uses 1.19)
2. Make sure GNU Make is installed
3. If you are on a Windows machine, install Git Bash
4. Clone the repository to your machine
5. Inside the root directory, run `make build` or `go build` to compile the binary
6. Run `./run.sh` to execute a word count job with three workers
8. Run `make test_wc` and `make test_fault` to run tests for correctness and fault-tolerance

The test for fault-tolerance deliberately injects faults in the form of crashing or slow workers.

### Command-Line Interface

In this project, a MapReduce job is ran as one or more compiled binaries that are given access to a directory for file reading and writing.

The binary can be executed as one of three possible roles: **sequential**, **worker**, or **master**. Role can be specified with the `-r` flag.
- The sequential role is meant to run a job as a single program (without concurrency). The output of a sequential execution is used as a benchmark for testing the distributed execution.
- The master role, used for distributed execution, coordinates the workers. A MapReduce job should only involve one master. The master program always needs to be started before the workers.
- The programs with the worker role read and write files and perform the core data processing. Workers will continually ask the master for a task before attempting to perform it.

Master and sequential programs must be passed the `-d` flag which specifies the filepath of the directory where input files are stored and output files will be stored. Workers do not need to directly be told the filepath because the master will tell them.

Using the `-p` flag when invoking the master or sequential role specifies the plugin, or which map and reduce functions to use. I have implemented a word count plugin that can be used by passing in the value `wc`.

The `-a` flag specifies the IP address and port number of the master node. It is required for both the master and worker roles. Currently, the project only supports local execution so the IP address should be 127.0.0.1 (localhost).

The `-nr` flag specifies to the master the number of files to output from the reduce phase.

### Extending with Custom Map and Reduce Functions

You can add custom map and reduce functions by implementing the `Plugin` interface in `plugins/plugin.go`. In the same file, you would need to create a new entry in the `Plugins` map to enable access to your functions via command-line invocation of the program. 

## What is MapReduce?

MapReduce is a programming model and framework for processing massive datasets in parallel. A user supplies a **map** function that performs filtering and/or computation and a **reduce** function that determines how results are grouped and aggregated. A MapReduce system consists of one master node and a bunch of worker nodes. The master node acts as a coordinator by assigning tasks to workers. The master node also handles machine failures by re-assigning tasks sent to workers that stopped responding.

The first MapReduce framework was developed by Google engineers in 2003 to provide a simple interface for harnessing the scale and parallel processing power of distributed systems. A few years later, engineers at Yahoo created an open-source implementation of MapReduce called Hadoop that became widely used in industry. Although MapReduce has become obsolete due to the domination of SQL-based interfaces, performance of OLAP database systems, and deprecation of batch processing, many of the core ideas such as scalability, fault tolerance, and shared-disk architecture are still relevant. 

## What is this project?

This project is an implementation of the MapReduce framework discussed in the original [2004 whitepaper](https://pdos.csail.mit.edu/6.824/papers/mapreduce.pdf) with a few differences:

- This project only implements some core features (programming model, distribution, fault-tolerance) discussed in the paper
- The "user-defined" map and reduce functions are not truly user-defined as they are part of the source code in the plugins folder
- This implementation uses the operating system's filesystem as a shared filesystem between all processes rather than a true distributed filesystem like HDFS or S3
- Because of the above limitation, this project cannot run in a true multi-node cluster distributed over a network, however, it could be extended to connect to a distributed filesystem

## Acknowledgements

This project was inspired by the MapReduce assignment in MIT graduate distributed systems course. While I was influenced by tiny portions of the starter code from the course, the source code in this repository is almost entirely my own.
