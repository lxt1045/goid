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

//go:build (amd64 || amd64p32) && gc && go1.5
// +build amd64 amd64p32
// +build gc
// +build go1.5

package goid

import (
	"runtime"
	"strconv"
	"sync/atomic"
	"unsafe"
	_ "unsafe"
)

func Get() int64

func GetDefer() uintptr

func GetPC() uintptr

func GetPCBP() (uintptr, uintptr)

func Ret(p interface{})

//go:linkname deferreturn runtime.deferreturn
func deferreturn(s string)

//go:linkname printnl runtime.printnl
func printnl()

//go:linkname printstring runtime.printstring
func printstring(s string)

var helloWorld = "你好，世界"

func Print()

//go:linkname callers runtime.callers
func callers(skip int, pcbuf []uintptr) int

func Getcallerpc() uintptr

func GetG() uintptr

type Log struct {
	PC   uintptr
	Code int64
	Msg  string
}

func NewLog(code int, msg string) Log

var (
	mCache3 unsafe.Pointer = func() unsafe.Pointer {
		m := make(map[uintptr]string)
		// for i := 0; i < 1000; i++ {
		// 	m[uintptr(i)] = fmt.Sprintf("%d", i)
		// }
		return unsafe.Pointer(&m)
	}()
)

func (l Log) LineNO() (line string) {
	mPCs := *(*map[uintptr]string)(atomic.LoadPointer(&mCache3))
	line, ok := mPCs[l.PC]
	if !ok {
		file, n := runtime.FuncForPC(l.PC).FileLine(l.PC)
		line = file + ":" + strconv.Itoa(n)
		mPCs2 := make(map[uintptr]string, len(mPCs)+10)
		mPCs2[l.PC] = line
		for {
			p := atomic.LoadPointer(&mCache3)
			mPCs = *(*map[uintptr]string)(p)
			for k, v := range mPCs {
				mPCs2[k] = v
			}
			swapped := atomic.CompareAndSwapPointer(&mCache3, p, unsafe.Pointer(&mPCs2))
			if swapped {
				break
			}
		}
	}
	return
}

type Line uintptr

func NewLine() Line

var (
	mCache4 unsafe.Pointer = func() unsafe.Pointer {
		m := make(map[Line]string)
		return unsafe.Pointer(&m)
	}()
)

func (l Line) LineNO() (line string) {
	mPCs := *(*map[Line]string)(atomic.LoadPointer(&mCache3))
	line, ok := mPCs[l]
	if !ok {
		file, n := runtime.FuncForPC(uintptr(l)).FileLine(uintptr(l))
		line = file + ":" + strconv.Itoa(n)
		mPCs2 := make(map[Line]string, len(mPCs)+10)
		mPCs2[l] = line
		for {
			p := atomic.LoadPointer(&mCache3)
			mPCs = *(*map[Line]string)(p)
			for k, v := range mPCs {
				mPCs2[k] = v
			}
			swapped := atomic.CompareAndSwapPointer(&mCache3, p, unsafe.Pointer(&mPCs2))
			if swapped {
				break
			}
		}
	}
	return
}

// /usr/local/go/src/runtime/traceback.go:187
// frame.fp = frame.sp + uintptr(funcspdelta(f, frame.pc, &cache))
// 表明: funcspdelta 可以用于查找stack中sp和fb间的距离
// 同时要注意,很多内联的函数,的特殊操作 ,,, ,,,

// pc++
// /usr/local/go/src/runtime/traceback.go
/*
// Normally, pc is a return address. In that case, we want to look up
// file/line information using pc-1, because that is the pc of the
// call instruction (more precisely, the last byte of the call instruction).
// Callers expect the pc buffer to contain return addresses and do the
// same -1 themselves, so we keep pc unchanged.
// When the pc is from a signal (e.g. profiler or segv) then we want
// to look up file/line information using pc, and we store pc+1 in the
// pc buffer so callers can unconditionally subtract 1 before looking up.
// See issue 34123.
// The pc can be at function entry when the frame is initialized without
// actually running code, like runtime.mstart.
*/
