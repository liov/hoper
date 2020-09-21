package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
)

func MapReduce(mapper func(interface{}, chan interface{}),
	reducer func(chan interface{}, chan interface{}),
	input chan interface{},
	pool_size int) interface{} {
	reduce_input := make(chan interface{})
	reduce_output := make(chan interface{})
	worker_output := make(chan chan interface{}, pool_size)
	go reducer(reduce_input, reduce_output)
	go func() {
		for worker_chan := range worker_output {
			reduce_input <- <-worker_chan
		}
		close(reduce_input)
	}()
	go func() {
		for item := range input {
			my_chan := make(chan interface{})
			go mapper(item, my_chan)
			worker_output <- my_chan
		}
		close(worker_output)
	}()
	return <-reduce_output
}

func find_files(dirname string) chan interface{} {
	output := make(chan interface{})
	go func() {
		_find_files(dirname, output)
		close(output)
	}()
	return output
}

func _find_files(dirname string, output chan interface{}) {
	dir, _ := os.Open(dirname)
	dirnames, _ := dir.Readdirnames(-1)
	for i := 0; i < len(dirnames); i++ {
		fullpath := dirname + "/" + dirnames[i]
		file, _ := os.Stat(fullpath)
		if file.IsDir() {
			_find_files(fullpath, output)
		} else {
			output <- fullpath
		}
	}
}

func EachLine(filename string) chan string {
	output := make(chan string)
	go func() {
		file, err := os.Open(filename)
		if err != nil {
			return
		}
		defer file.Close()
		reader := bufio.NewReader(file)
		for {
			line, err := reader.ReadString('\n')
			output <- line
			if err == io.EOF {
				break
			}
		}
		close(output)
	}()
	return output
}

func mapper(filename interface{}, output chan interface{}) {
	results := map[string]int{}
	wordsRE := regexp.MustCompile(`[A-Za-z0-9_]*`)
	for line := range EachLine(filename.(string)) {
		for _, match := range wordsRE.FindAllString(line, -1) {
			results[match]++
		}
	}
	output <- results
}

func reducer(input chan interface{}, output chan interface{}) {
	results := map[string]int{}
	for new_matches := range input {
		for key, value := range new_matches.(map[string]int) {
			previous_count, exists := results[key]
			if !exists {
				results[key] = value
			} else {
				results[key] = previous_count + value
			}
		}
	}
	output <- results
}

func main() {
	fmt.Print(MapReduce(mapper, reducer, find_files("."), 20))
}
