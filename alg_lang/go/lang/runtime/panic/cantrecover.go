package main

//不能recover的panic bulkBarrierPreWrite: unaligned arguments

//产生场景，fmt.Println((*type)(unsafe.Pointer(uintptr(address)))) address不是整字节
//type为基本类型时(string不算基本类型)不会panic

//unexpected fault address 0x
func main() {

}
