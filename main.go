package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime/pprof"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	cpuProfile := flag.String("cpuprofile", "", "write cpu profile to file")
	flag.Parse()
	if *cpuProfile != "" {
		f, err := os.Create(*cpuProfile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	rand.Seed(time.Now().Unix())
	w, h := 1024, 768
	s, err := newSquads(w, h)
	if err != nil {
		fmt.Printf("setup system: %v\n", err)
		os.Exit(1)
	}

	ebiten.SetWindowSize(w, h)
	if err := ebiten.RunGame(s); err == errExitGame {
		fmt.Println("See you next time.")
	} else if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
}

var errExitGame = errors.New("game has completed")
