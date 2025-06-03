package main

import (
	"time"
	"github.com/William-Gardner-Biotech/polybar/polybar"
)

func main() {
	// Example 1: Basic usage with custom sequence
	println("Example 1: Basic DNA progress bar")
	pb1 := polybar.New("ATCG-NRZA-GCTA", "PROCESSING-DATA")
	pb1.Start(50)

	for i := 0; i < 50; i++ {
		time.Sleep(100 * time.Millisecond)
		pb1.Update()
	}
	pb1.Finish()

	time.Sleep(1 * time.Second)

	// Example 2: Using sequence length as width
	println("\nExample 2: Sequence-width progress bar")
	pb2 := polybar.New("ATCGATCGTTAACCGG", "")
	pb2.Start(100)

	// Simulate batch processing
	for i := 0; i < 10; i++ {
		time.Sleep(200 * time.Millisecond)
		pb2.SetProgress((i + 1) * 10)
	}
	pb2.Finish()
}
