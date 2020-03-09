package main

type Ptr interface {
	Ptr()
}

type InterPtr struct {
	Ptr *Ptr
}

func main() {

}
