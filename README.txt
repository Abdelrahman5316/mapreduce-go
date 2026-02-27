# Distributed MapReduce Framework in Go

## Overview

This project implements a simplified distributed **MapReduce framework** in Go. The system coordinates a master node and multiple worker nodes to execute parallel data processing tasks using the MapReduce programming model.

The implementation supports task scheduling, fault tolerance, worker coordination via RPC, and execution of user-defined map and reduce functions. A sample **word count application** is included to demonstrate the framework.

---

## Architecture

The system follows a classic MapReduce architecture:

* **Master**

  * Splits input data into map tasks
  * Schedules map and reduce jobs
  * Assigns tasks to available workers
  * Detects worker failures and reschedules unfinished tasks
  * Merges final reduce outputs

* **Workers**

  * Request tasks from the master via RPC
  * Execute map or reduce operations
  * Persist intermediate results
  * Report completion status back to the master

Communication between components is implemented using Go’s **net/rpc** package.

---

## Key Features

* Parallel execution of Map and Reduce phases
* Dynamic task scheduling
* Worker failure handling and task reassignment
* Intermediate data partitioning
* RPC-based distributed coordination
* Sequential and distributed execution modes
* Automated test suite validating correctness and fault tolerance

---

## Project Structure

```
mapreduce/
 ├── master.go              # Master coordination logic
 ├── worker.go              # Worker execution logic
 ├── schedule.go            # Task scheduling implementation
 ├── map_reduce.go          # MapReduce orchestration
 ├── common_rpc.go          # RPC definitions
 ├── master_splitmerge.go   # Output merging
 └── mr_test.go             # Framework tests

cmd/wordcount/
 └── word_count.go          # Example MapReduce application
```

---

## Example Application

The included **Word Count** program demonstrates how users can define custom:

* `Map()` — processes input data into key/value pairs
* `Reduce()` — aggregates values associated with each key

This showcases how the framework can be reused for large-scale data processing tasks.

---

## Technologies

* Go (Golang)
* Goroutines & Channels
* net/rpc
* Concurrent task scheduling
* Distributed systems concepts

---

## Testing

The implementation passes all provided tests, including:

* Sequential execution
* Parallel execution
* Worker failure scenarios
* Multiple failure recovery

Run tests using:

```bash
go test ./mapreduce -v
```

---

## Learning Objectives

This project demonstrates practical concepts in:

* Distributed systems design
* Fault tolerance
* Parallel computation
* RPC communication
* Concurrency management in Go

---

## Notes

This implementation is inspired by the MapReduce model introduced by Google and commonly used in distributed systems coursework.
