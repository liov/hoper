package asm

func SyscallWrite_Darwin(fd int, msg string) int

func AsmCallCAdd(cfun uintptr, a, b int64) int64
