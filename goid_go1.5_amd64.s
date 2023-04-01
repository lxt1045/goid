// Copyright 2016 Peter Mattis.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied. See the License for the specific language governing
// permissions and limitations under the License. See the AUTHORS file
// for names of contributors.

// Assembly to mimic runtime.getg.

// go:build (amd64 || amd64p32) && gc && go1.5
// +build amd64 amd64p32
// +build gc
// +build go1.5

#include "go_asm.h"
#include "textflag.h"

// func Get() int64
TEXT ·Get(SB),NOSPLIT,$0-8
	MOVQ (TLS), R14
	MOVQ g_goid(R14), R13
	MOVQ R13, ret+0(FP)
	RET

// func GetDefer() uintptr
TEXT ·GetDefer(SB),NOSPLIT,$0-8
	MOVQ (TLS), R14
	MOVQ g__defer(R14), R13
	MOVQ R13, ret+0(FP)
	RET

// func GetPC() uintptr
TEXT ·GetPC(SB),NOSPLIT,$0-8
	MOVQ retpc-8(FP), AX // MOVQ (SP), R13
	MOVQ AX, ret+0(FP)
	RET

// func GetPCBP() (uintptr,uintptr)
TEXT ·GetPCBP(SB),NOSPLIT,$0-16
	MOVQ retpc-8(FP), AX 
	MOVQ AX, ret+0(FP)
	MOVQ retbp-16(FP), AX 
	MOVQ AX, ret2+8(FP)
	RET


// func Retx()
TEXT ·Retx(SB), NOSPLIT, $0-8
	MOVQ	+8(BP), BX
	MOVQ	BP, AX
	ADDQ	$16,AX
	MOVQ	+0(BP), BP
	MOVQ	AX,SP
	JMP	BX



// func Ret1(p uintptr)
TEXT ·Ret1(SB), NOSPLIT, $0-16
	CMPQ	p+0(FP), $0
	JHI	unwind
	RET
unwind:
	MOVQ	BP, SP
	ADDQ	$16,SP
	MOVQ	+0(BP), BP
	CALL    runtime·deferreturn(SB)
	JMP	-8(SP)



// func Ret1(p uintptr)
TEXT ·Ret(SB), NOSPLIT, $0-16
	CMPQ	p+0(FP), $0
	JHI	unwind
	RET
unwind:
	MOVQ	BP, SP
	// ADDQ	$16,SP
	ADDQ	$16,SP
	MOVQ	+0(BP), BP
	//CALL runtime·deferreturn(SB)
	JMP	-8(SP)


TEXT ·Print(SB), $16-0
	MOVQ ·helloWorld+0(SB), AX
	MOVQ AX, 0(SP)
	MOVQ ·helloWorld+8(SB), BX
	MOVQ BX, 8(SP)
	CALL runtime·printstring(SB)
	CALL runtime·printnl(SB)
	RET


TEXT ·Getcallerpc(SB), $0-8
	CALL runtime·callers(SB)
	RET

// func GetG(p uintptr)
// TEXT ·GetG(SB), NOSPLIT|ABIInternal, $0-8
// 	CALL runtime·getg(SB)
// 	RET

//NewLog(code int, msg string) Log
TEXT    ·NewLog(SB), NOSPLIT, $0-56
    MOVQ 	retpc-8(FP), AX
	MOVQ 	AX, ret+24(FP)
	MOVQ 	code+0(FP), AX
	MOVQ 	AX, ret+32(FP)
	MOVQ 	msg+8(FP), AX
	MOVQ 	AX, ret+40(FP)
	MOVQ 	msg+16(FP), AX
	MOVQ 	AX, ret+48(FP)
	RET


//NewLine() Line
TEXT    ·NewLine(SB), NOSPLIT, $0-8
    MOVQ 	retpc-8(FP), AX
	MOVQ 	AX, ret+0(FP)
	RET
