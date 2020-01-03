package main

import (
	"strconv"
	"strings"

	"github.com/bcongdon/corral"
)

type wordCount struct{}

func (w wordCount) Map(key, value string, emitter corral.Emitter) {
	for _, word := range strings.Fields(value) {
		emitter.Emit(word, "")
	}
}

func (w wordCount) Reduce(key string, values corral.ValueIterator, emitter corral.Emitter) {
	count := 0
	for range values.Iter() {
		count++
	}
	emitter.Emit(key, strconv.Itoa(count))
}

func main() {
	wc := wordCount{}
	job := corral.NewJob(wc, wc)

	driver := corral.NewDriver(job)
	driver.Main()
}
