package main

import "fmt"

var A Wb
var B Wb

type Wb struct {
	Obj *int
}

func simpleSet(c *int) {
	A.Obj = nil
	B.Obj = c

	//if GC Begin
	A.Obj = c
	B.Obj = nil
	//scan B
}

//go build -gcflags "-N -l"
//go tool objdump -s 'main\.simpleSet' -S ./main
func main() {
	var c = 100
	simpleSet(&c)
	fmt.Println("GitCommit->", GitCommit)
	fmt.Println("Version->", Version)
	fmt.Println("aaa->", BuildDate)
}

//要想这里的变量初始化，再build的时候
//-ldflags "-s -w -X main.GitCommit=$(git rev-parse --short HEAD) -X main.Version=$(git rev-parse --abbrev-ref HEAD) -X main.BuildDate=$(date +%FT%T%z)"
var (
	GitCommit string
	Version   string
	BuildDate string
)

//GitCommit-> 246b577
//Version-> master
//aaa-> 2022-01-06T17:58:49+0800
