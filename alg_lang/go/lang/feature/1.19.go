package main

type T[T any] struct{}

func (t T[T]) name() {

}
