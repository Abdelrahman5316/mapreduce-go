package mapreduce

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"log"
	"os"
	"sort"
)

func runMapTask(
	jobName string, // The name of the whole mapreduce job
	mapTaskIndex int, // The index of the map task
	inputFile string, // The path to the input file assigned to this task
	nReduce int, // The number of reduce tasks that will be run
	mapFn func(file string, contents string) []KeyValue, // The user-defined map function
) {
	
	data, err := ioutil.ReadFile(inputFile)

	if err != nil {
		fmt.Println(err)
	}
	file_data := string(data)

	r_map := mapFn(inputFile, file_data)

	var file_to_write_in uint32

	intermediateFiles := make([]*os.File, nReduce)
	encoders := make([]*json.Encoder, nReduce)
	for i := 0; i < nReduce; i++ {

		name := getIntermediateName(jobName, mapTaskIndex, i)
		intermediateFile, err := os.Create(name)
		if err != nil {
			log.Fatal(err)
		}
		defer intermediateFile.Close()
		intermediateFiles[i] = intermediateFile
		encoders[i] = json.NewEncoder(intermediateFile)
	}

	for _, pair := range r_map {

		file_to_write_in = hash32(pair.Key) % uint32(nReduce)
		err := encoders[file_to_write_in].Encode(&pair)
		if err != nil {
			log.Fatal(err)
		}

	}
	for i := 0; i < nReduce; i++ {
		err := intermediateFiles[i].Close()
		if err != nil {
			log.Fatal(err)
		}
	}

}

func hash32(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

func runReduceTask(
	jobName string, // the name of the whole MapReduce job
	reduceTaskIndex int, // the index of the reduce task
	nMap int, // the number of map tasks that were run
	reduceFn func(key string, values []string) string,
) {

	
	var sort_keys []string
	var intermediateKeyValues []KeyValue

	for i := 0; i < nMap; i++ {
		name := getIntermediateName(jobName, i, reduceTaskIndex)
		file, err := os.Open(name)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		decoder := json.NewDecoder(file)
		var kv KeyValue
		for decoder.More() {
			err := decoder.Decode(&kv)
			if err != nil {
				log.Fatal(err)
			}
			intermediateKeyValues = append(intermediateKeyValues, kv)
		}
	}
	grouping := make(map[string][]string)

	for _, i := range intermediateKeyValues {
		grouping[i.Key] = append(grouping[i.Key], i.Value)
	}
	for key := range grouping {
		sort_keys = append(sort_keys, key)
	}
	sort.Strings(sort_keys)
	outputFileName := getReduceOutName(jobName, reduceTaskIndex)
	outputFile, err := os.Create(outputFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()

	encoder := json.NewEncoder(outputFile)
	for _, key := range sort_keys {
		outputValue := reduceFn(key, grouping[key])
		err := encoder.Encode(KeyValue{Key: key, Value: outputValue})
		if err != nil {
			log.Fatal(err)
		}
	}
}
