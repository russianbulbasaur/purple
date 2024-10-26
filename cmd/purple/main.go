package main

import (
	"fmt"
	myResp "github.com/russianbulbasaur/my-resp"
)

func main() {
	resp := myResp.Init()
	fmt.Println("%#v", resp)
}
