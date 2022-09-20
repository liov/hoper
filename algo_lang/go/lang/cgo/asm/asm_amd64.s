// System V AMD64 ABI
// func asmCallCAdd(cfun uintptr, a, b int64) int64
TEXT ·AsmCallCAdd(SB),  $0
    MOVQ cfun+0(FP), AX // cfun
    MOVQ a+8(FP),    DI // a
    MOVQ b+16(FP),   SI // b
    CALL AX
    MOVQ AX, ret+24(FP)
    RET

// func SyscallWrite_Darwin(fd int, msg string) int
TEXT ·SyscallWrite_Darwin(SB),  $0
    MOVQ $(0x2000000+4), AX // #define SYS_write 4
    MOVQ fd+0(FP),       DI
    MOVQ msg_data+8(FP), SI
    MOVQ msg_len+16(FP), DX
    SYSCALL
    MOVQ AX, ret+0(FP)
    RET
