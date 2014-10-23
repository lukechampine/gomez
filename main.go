package main

import (
	"fmt"
	"time"

	"github.com/lukechampine/gomez/gomez"
)

func main() {
	m, err := gomez.New("maze.gif")
	if err != nil {
		fmt.Println(err)
		return
	}
	t := time.Now()
	if !m.Solve() {
		fmt.Println("Maze has no solution!")
		return
	}
	dur := time.Now().Sub(t)
	if err = m.Save("maze-solved.gif"); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Solved maze in", dur.Seconds(), "seconds")
	}
}
