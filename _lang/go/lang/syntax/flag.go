package main

import (
	"github.com/spf13/pflag"
	"log"
	"os"
)

func main() {
	pflag.CommandLine = pflag.NewFlagSet(os.Args[0], pflag.ContinueOnError)
	var a int
	pflag.IntVarP(&a, "aa", "a", 0, "a")
	pflag.Parse()
	log.Println(a)
	var b int
	pflag.IntVarP(&b, "bb", "b", 0, "a")
	pflag.Parse()
	log.Println(b)
}
