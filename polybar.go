// Package dnaprogress provides a DNA‐style progress bar with base complementing
package polybar

import (
	"fmt"
	"os"
	"strings"
)

const (
	// DNA‐style progress bar characters
	zipperChar = "┬"
	baseChar   = "┴"
	arrowText  = "===>"

	defaultSequence = "GCCAGTTTTGGGCTGGTTGGC" // first 21 letters of DNAPolA1 https://www.ncbi.nlm.nih.gov/nuccore/NG_016798.2?report=fasta&from=4972&to=308040
)

// ProgressBar represents a DNA‐style progress bar
type ProgressBar struct {
	width      int
	headerLine string
	topStrand  string
	complement string
	completed  int
	total      int
}

// New creates a new DNA progress bar.
// - topStrand: the DNA sequence to display (will be complemented on bottom).
//   If empty, defaults to the first 21 nucleotides of DNA polymerase I.
// - header: optional header text (if empty, uses topStrand length for width).
func New(topStrand string, header string) *ProgressBar {
	// If the caller passed an empty string, use our 21‐bp default.
	if strings.TrimSpace(topStrand) == "" {
		topStrand = defaultSequence
	}

	pb := &ProgressBar{
		topStrand: strings.ToUpper(topStrand),
		completed: 0,
	}

	// Generate the complement right away
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

// generateComplement creates DNA complement: A↔T, G↔C, others→N, dash→dash.
// It also flips literal '5'→'3' and '3'→'5' if someone manually includes them.
func generateComplement(sequence string) string {
	complement := make([]rune, len(sequence))

	for i, base := range sequence {
		switch base {
		case '5':
			complement[i] = '3'
		case '3':
			complement[i] = '5'
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

// padOrTruncate ensures string s is exactly the specified length.
// If shorter, pads with dashes. If longer, truncates.
func padOrTruncate(s string, length int) string {
	if len(s) == length {
		return s
	} else if len(s) < length {
		return s + strings.Repeat("-", length-len(s))
	}
	return s[:length]
}

// Start initializes the progress bar display (0 completed out of total).
func (pb *ProgressBar) Start(total int) {
	pb.total = total
	pb.completed = 0
	pb.render()
}

// Update increments progress by one step and refreshes the display.
func (pb *ProgressBar) Update() {
	pb.completed++
	pb.render()
}

// SetProgress sets the current progress value and refreshes the display.
func (pb *ProgressBar) SetProgress(completed int) {
	pb.completed = completed
	pb.render()
}

// Finish marks the progress bar as fully complete, then prints a newline.
func (pb *ProgressBar) Finish() {
	pb.completed = pb.total
	pb.render()
	fmt.Fprintln(os.Stderr) // Final newline
}

// render draws the DNA‐style progress bar with 5′→3′ labeling.
func (pb *ProgressBar) render() {
	// If total is zero, nothing to do.
	if pb.total == 0 {
		return
	}

	// Calculate how many bases to “fill in” based on completed/total.
	pos := pb.completed * pb.width / pb.total
	if pos > pb.width {
		pos = pb.width
	}

	// Line 1: prepend "5'" then zipper characters for full width
	lineZipper := "5'" + strings.Repeat(zipperChar, pb.width)

	// Line 2: two dashes + the first pos bases of topStrand
	var lineTop string
	if pos <= len(pb.topStrand) {
		lineTop = "--" + pb.topStrand[:pos]
	} else {
		lineTop = "--" + pb.topStrand
	}

	// Line 3: two dashes + the first pos bases of complement
	var lineComplement string
	if pos <= len(pb.complement) {
		lineComplement = "--" + pb.complement[:pos]
	} else {
		lineComplement = "--" + pb.complement
	}

	// Line 4: prepend "3'" then baseChar repeated pos times, followed by arrowText
	var linePrimer string
	if pos < pb.width {
		linePrimer = "3'" + strings.Repeat(baseChar, pos) + arrowText
	} else {
		linePrimer = "3'" + strings.Repeat(baseChar, pb.width) + arrowText
	}

	// Line 5: percentage and completed/total
	percent := float64(pb.completed) / float64(pb.total) * 100
	linePercent := fmt.Sprintf("%.1f%% (%d/%d)", percent, pb.completed, pb.total)

	// On subsequent frames, move cursor up 5 lines to overwrite
	if pb.completed > 0 {
		for i := 0; i < 5; i++ {
			fmt.Fprint(os.Stderr, "\033[F")
		}
	}

	// Print all five lines to stderr
	fmt.Fprintln(os.Stderr, lineZipper)
	fmt.Fprintln(os.Stderr, lineTop)
	fmt.Fprintln(os.Stderr, lineComplement)
	fmt.Fprintln(os.Stderr, linePrimer)
	fmt.Fprintln(os.Stderr, linePercent)
}

// Example usage function
func Example() {
	// 1) Create a progress bar with no explicit topStrand: will default to 21 bp of Pol I.
	pb := New("", "DNA‐POLYMERASE‐I")

	// 2) Start with total of 100 steps
	pb.Start(100)

	// 3) Simulate progress
	for i := 0; i < 100; i++ {
		// ... do some work ...
		pb.Update()
		// time.Sleep(50 * time.Millisecond)
	}

	// 4) Finish
	pb.Finish()
}
