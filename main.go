package main

import (
	"fmt"

	"github.com/maxjustus/eudex-go/eudex"
)

func main() {
	fmt.Println(eudex.Eudex("Jeff Buckley").String())
	fmt.Println(eudex.Eudex("Tim Buckley").String())
	fmt.Println(eudex.Similar("Jonny", "Johnny"))
	fmt.Println(eudex.Similar("Jonny", "Jahnny"))
	fmt.Println(eudex.Similar("Jonny", "Jenny"))
	fmt.Println(eudex.StringDistance("Jonny", "Jenny"))
	fmt.Println(eudex.StringDistance("Jonny", "Jentny"))
	fmt.Println(eudex.Similar("Trimothy", "Tony"))
}
