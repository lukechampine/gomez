package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/lukechampine/gomez/gomez"
)

func main() {
	m, err := gomez.New(os.Args[1])
	if err != nil {
		fmt.Println("Could not load maze:", err)
		return
	}
	t := time.Now()
	if !m.Solve() {
		fmt.Println("Maze has no solution!")
		return
	}
	dur := time.Now().Sub(t)
	if err = m.Save(strings.Trim(os.Args[1], ".gif") + "-solved.gif"); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Solved maze in", dur.Seconds(), "seconds")
	}
}
