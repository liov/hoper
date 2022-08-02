package main

type T[T any] struct{}

func (t *T[T]) method() {

}

type B[T any] struct{}

func (t *B[T]) method() {

}

type C struct{}

func (t *C[C]) method() {

}
