package main

func main() {
	//runtime.GOMAXPROCS(1) //First
	exit := make(chan int)
	go func() {
		close(exit)
		for {
			if true {
				println("Looping!") //Second
			}
		}
	}()
	<-exit
	println("Am I printed?")
}
