package mapreduce

func (mr *Master) schedule(phase jobPhase) {

	var ntasks int
	var numOtherPhase int
	switch phase {
	case mapPhase:
		ntasks = len(mr.files)     // number of map tasks
		numOtherPhase = mr.nReduce // number of reducers
	case reducePhase:
		ntasks = mr.nReduce           // number of reduce tasks
		numOtherPhase = len(mr.files) // number of map tasks
	}
	debug("Schedule: %v %v tasks (%d I/Os)\n", ntasks, phase, numOtherPhase)


	taskCompleted := make(chan bool)
	for i := 0; i < ntasks; i++ {
		go func(task_no int) {

			for {
				address := <-mr.registerChannel

				arg := RunTaskArgs{
					JobName:       mr.jobName,
					File:          mr.files[task_no],
					Phase:         phase,
					TaskNumber:    task_no,
					NumOtherPhase: numOtherPhase,
				}
				check_worker := call(address, "Worker.RunTask", &arg, new(struct{}))
				if check_worker {
					//fmt.Println("Worker Successed")
					taskCompleted <- true
					mr.registerChannel <- address
					break
				}

			}

		}(i)

	}

	for i := 0; i < ntasks; i++ {
		<-taskCompleted
	}

	debug("Schedule: %v phase done\n", phase)
}
