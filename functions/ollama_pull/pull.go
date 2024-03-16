package lib

import (
	"github.com/taubyte/go-sdk/event"
	"fmt"
)

//go:wasm-module env
//export _sleep
func Sleep(dur int64)

//go:wasm-module ollama
//export pull
func Pull(model string, id *uint64) uint32

//func Pull(status *byte, statusLen *uint32, err *byte, errLen *uint32) uint32

//go:wasm-module ollama
//export pull_status
func PullStatus(id uint64, status *byte, statusCap uint32, statusLen *uint32, total *int64, completed *int64, err *byte, errCap uint32, errLen *uint32) uint32

var (
	pstatus [256]byte
	perr    [256]byte
)


func printPullStatus(id uint64) (string,bool) {
	var pstatusLen, perrLen uint32
	var total, completed int64

	err := PullStatus(id, &pstatus[0], uint32(len(pstatus)), &pstatusLen, &total, &completed, &perr[0], uint32(len(perr)), &perrLen)
	if err != 0 {
		panic("failed to call pull_status")

	}

	status := string(pstatus[0:pstatusLen])

    var s string
	fmt.Sprintln(s,status)
	fmt.Println(s,completed, "/", total)
	fmt.Println(s,"ERR:", string(perr[0:perrLen]))
	fmt.Println(s,status == "success" || perrLen > 0)

	return s, status == "success" || perrLen > 0
}

//export generate
func generate(e event.Event) uint32 {
	h, err0 := e.HTTP()
	if err0 != nil {
		return 1
	}

	var id uint64
	err := Pull("gemma:2b-instruct", &id)
	if err != 0 {
		panic("failed to call pull")
	}

	for {
		s, done := printPullStatus(id)
		if done {
			break
		}
		h.Write([]byte(s))
		Sleep(10 * int64(time.Second))
	}

	return 0
}
