package main

import (
	"fmt"
	conf "github.com/gusleein/go-yaml-conf"
)

func main() {
	fmt.Println(conf.Local.GetUint("id"))
	fmt.Println(conf.Local.Get("tokenSecret"))
}
