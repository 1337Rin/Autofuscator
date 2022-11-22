package main

import (
	"fmt"
	"github.com/1337Rin/Autofuscator/golang"
	//"example.com/autofuscator/powershell"
)

func main() {
	lines := golang.Preobfuscation("../go/countdown.go")
	for i:=0;i<len(lines);i++ {
		fmt.Println(lines[i])
	}

	fmt.Println(golang.Findvariables(lines))
	fmt.Println(golang.Findimports(lines))
	fmt.Println(golang.Findstrings(lines))



}