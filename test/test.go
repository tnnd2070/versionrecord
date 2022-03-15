package main

import (
	"fmt"
	"strings"
)

func main() {
	filename:="hxuat_pkgversion_2022012017.csv"
	//还可以是这种方式打开
	var envName = strings.Split(filename, "_")[1]
	fmt.Println(envName)
}
