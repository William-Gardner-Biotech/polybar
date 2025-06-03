// Package dnaprogress provides a DNA-style progress bar with base complementing
package polybar

import (
	"fmt"
	"os"
	"strings"
)

const (
	// DNA-style progress bar characters
	zipperChar = "┬"
	baseChar   = "┴"
	arrowText  = "===>"
)

// ProgressBar represents a DNA-style progress bar
type ProgressBar struct {
	width       int
	headerLine  string
	topStrand   string
	complement  string
	completed   int
	total       int
}

// New creates a new DNA progress bar
// topStrand: the DNA sequence to display (will be complemented on bottom)
// header: optional header text (if empty, uses topStrand length for width)
func New(topStrand string, header string) *ProgressBar {
	pb := &ProgressBar{
		topStrand: strings.ToUpper(topStrand),
		completed: 0,
	}

	// Generate complement
	pb.complement = generateComplement(pb.topStrand)

	// Set width and header
	if header != "" {
		pb.headerLine = header
		pb.width = len(header)
		// Pad or truncate strands to match header width
		pb.topStrand = padOrTruncate(pb.topStrand, pb.width)
		pb.complement = padOrTruncate(pb.complement, pb.width)
	} else {
		pb.width = len(pb.topStrand)
		pb.headerLine = pb.topStrand
	}

	return pb
}

// generateComplement creates DNA complement: A↔T, G↔C, others→N, dash→dash
func generateComplement(sequence string) string {
	complement := make([]rune, len(sequence))

	for i, base := range sequence {
		switch base {
		case 'A':
			complement[i] = 'T'
		case 'T':
			complement[i] = 'A'
		case 'G':
			complement[i] = 'C'
		case 'C':
			complement[i] = 'G'
		case '-':
			complement[i] = '-'
		default:
			complement[i] = 'N'
		}
	}

	return string(complement)
}

// padOrTruncate ensures string is exactly the specified length
func padOrTruncate(s string, length int) string {
	if len(s) == length {
		return s
	} else if len(s) < length {
		return s + strings.Repeat("-", length-len(s))
	}
	return s[:length]
}

// Start initializes the progress bar display
func (pb *ProgressBar) Start(total int) {
	pb.total = total
	pb.completed = 0

	// Print header
	fmt.Fprintln(os.Stderr, pb.headerLine)
	pb.render()
}

// Update increments progress and refreshes the display
func (pb *ProgressBar) Update() {
	pb.completed++
	pb.render()
}

// SetProgress sets the current progress value
func (pb *ProgressBar) SetProgress(completed int) {
	pb.completed = completed
	pb.render()
}

// Finish completes the progress bar
func (pb *ProgressBar) Finish() {
	pb.completed = pb.total
	pb.render()
	fmt.Fprintln(os.Stderr) // Add final newline
}

// render draws the DNA-style progress bar
func (pb *ProgressBar) render() {
	if pb.total == 0 {
		return
	}

	// Calculate position (how many bases to fill in)
	pos := pb.completed * pb.width / pb.total
	if pos > pb.width {
		pos = pb.width
	}

	// Line 1: zipper of '┬' characters
	lineZipper := strings.Repeat(zipperChar, pb.width)

	// Line 2: top strand (first 'pos' characters)
	lineTop := ""
	if pos <= len(pb.topStrand) {
		lineTop = pb.topStrand[:pos]
	} else {
		lineTop = pb.topStrand
	}

	// Line 3: complement strand (first 'pos' characters)
	lineComplement := ""
	if pos <= len(pb.complement) {
		lineComplement = pb.complement[:pos]
	} else {
		lineComplement = pb.complement
	}

	// Line 4: primer + arrow
	var linePrimer string
	if pos < pb.width {
		linePrimer = strings.Repeat(baseChar, pos) + arrowText
	} else {
		linePrimer = strings.Repeat(baseChar, pb.width) + arrowText
	}

	// Line 5: percentage and X/N
	percent := float64(pb.completed) / float64(pb.total) * 100
	linePercent := fmt.Sprintf("%.1f%% (%d/%d)", percent, pb.completed, pb.total)

	// On subsequent frames, move cursor up to overwrite
	if pb.completed > 0 {
		for i := 0; i < 5; i++ {
			fmt.Fprint(os.Stderr, "\033[F")
		}
	}

	// Print all lines to stderr
	fmt.Fprintln(os.Stderr, lineZipper)
	fmt.Fprintln(os.Stderr, lineTop)
	fmt.Fprintln(os.Stderr, lineComplement)
	fmt.Fprintln(os.Stderr, linePrimer)
	fmt.Fprintln(os.Stderr, linePercent)
}

// Example usage function
func Example() {
	// Create progress bar with custom DNA sequence
	pb := New("ATCG-NNTA-GCTA", "DNA-SEQUENCING")

	// Start with total of 100 steps
	pb.Start(100)

	// Simulate progress
	for i := 0; i < 100; i++ {
		// Do some work...
		pb.Update()
		// Small delay for demo
		// time.Sleep(50 * time.Millisecond)
	}

	pb.Finish()
}
