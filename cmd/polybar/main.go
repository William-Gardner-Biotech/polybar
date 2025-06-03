package main

import (
    "flag"
    "time"
    "github.com/William-Gardner-Biotech/polybar/polybar"
)

func main() {
    seqPtr    := flag.String("seq", "", "DNA sequence (defaults to first 21 nt of Pol I)")
    headerPtr := flag.String("header", "", "Optional header above zipper")
    totalPtr  := flag.Int("total", 100, "Number of steps to reach 100%")
    interval  := flag.Duration("interval", 50*time.Millisecond, "Delay between Update() calls")

    flag.Parse()
    pb := polybar.New(*seqPtr, *headerPtr)
    pb.Start(*totalPtr)

    for i := 0; i < *totalPtr; i++ {
        time.Sleep(*interval)
        pb.Update()
    }
    pb.Finish()
}
