package main

import (
	"clifford/pkg/clifford"
	"log"
	"math"
	"math/rand"
	"os"
	"runtime"
	"runtime/pprof"
	"time"
)

const (
	DotsPerCm = 119 // 300 DPI
	Size      = 50 * DotsPerCm
	// Size  = 512
	Steps = 1e8
)

func main() {
	fCpuProf, err := os.Create("./cpu.prof")
	if err != nil {
		log.Fatal("could not create CPU profile: ", err)
	}
	defer fCpuProf.Close()
	if err := pprof.StartCPUProfile(fCpuProf); err != nil {
		log.Fatal("could not start CPU profile: ", err)
	}
	defer pprof.StopCPUProfile()

	rand.Seed(time.Now().UnixNano())

	log.Println("fetching random gradient...")
	grad, err := randomGradient()
	if err != nil {
		log.Fatalf("could not get random gradient: %s", err)
	}

	// a, b, c, d := -1.4, 1.6, 1.0, 0.7
	// a, b, c, d := 1.6, -0.6, -1.2, 1.6
	// a, b, c, d := 1.7, 1.7, 0.6, 1.2
	// a, b, c, d := 1.5, -1.8, 1.6, 0.9
	// a, b, c, d := -1.7, 1.3, -0.1, -1.2
	// a, b, c, d := -1.7, 1.8, -1.9, -0.4
	a, b, c, d := -1.8, -2.0, -0.5, -0.9
	att := clifford.NewAttractor(a, b, c, d)

	log.Println("stabilizing attractor...")
	for i := 0; i < 128; i++ {
		att.Advance()
	}

	log.Println("building histogram...")
	hist := clifford.NewHistogram(Size, math.Phi/2, att)
	for i := 0; i < Steps; i++ {
		att.Advance()
		hist.Inc(att.X, att.Y)
	}

	log.Println("rendering histogram...")
	img := clifford.RenderHistogram(hist, Size, grad)

	log.Println("writing output image...")
	if err := writeImage("./output.jpg", img); err != nil {
		log.Fatalf("could not write image: %s", err)
	}

	fMemProf, err := os.Create("./mem.prof")
	if err != nil {
		log.Fatal("could not create memory profile: ", err)
	}
	defer fMemProf.Close()
	runtime.GC()
	if err := pprof.WriteHeapProfile(fMemProf); err != nil {
		log.Fatal("could not write memory profile: ", err)
	}
}
