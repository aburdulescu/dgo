package main

import (
	"fmt"
	"os"

	"github.com/aburdulescu/dgo"
)

func main() {
	g, err := dgo.Parse(os.Stdin)
	if err != nil {
		panic(err)
	}
	fmt.Println(g.Write())
}
