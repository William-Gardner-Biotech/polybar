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

	// Default = first 21 nt of DNA polymerase I (NCBI: NG_016798.2, positions 4972–308040)
	defaultSequence = "GCCAGTTTTGGGCTGGTTGGC"
)

// ProgressBar represents a DNA-style progress bar
type ProgressBar struct {
	width      int    // number of bases across
	headerLine string // if non-empty, print this above zipper
	topStrand  string // uppercase DNA (template)
	complement string // computed complement of topStrand
	completed  int    // how many “steps” done so far
	total      int    // total number of “steps”
}

// New creates a new DNA progress bar.
//   • topStrand: the DNA sequence to display (will be complemented on bottom).
//     If empty, defaults to defaultSequence (21 nt).
//   • header:    optional header text. If non-empty, printed above zipper;
//                if empty, we set headerLine="" (so nothing prints there).
func New(topStrand, header string) *ProgressBar {
	// 1) If caller did not provide any sequence, use defaultSequence.
	if strings.TrimSpace(topStrand) == "" {
		topStrand = defaultSequence
	}

	pb := &ProgressBar{
		topStrand:  strings.ToUpper(topStrand),
		completed:  0,
		headerLine: header, // may be "" if caller wants no header
	}

	// 2) Generate the complement once
	pb.complement = generateComplement(pb.topStrand)

	// 3) Decide width: if header is non-empty, use its length; else use length of topStrand
	if header != "" {
		pb.width = len(header)
		// Pad or truncate both strands so their printed width = len(header)
		pb.topStrand = padOrTruncate(pb.topStrand, pb.width)
		pb.complement = padOrTruncate(pb.complement, pb.width)
	} else {
		pb.width = len(pb.topStrand)
		// leave topStrand, complement as-is
	}

	return pb
}

// generateComplement returns the complement of a DNA sequence.
// A↔T, G↔C; digits '5' ↔ '3'; dash→dash; others→'N'.
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

// padOrTruncate returns s padded with dashes or truncated so its length == length.
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

// Update increments progress by one step and refreshes.
func (pb *ProgressBar) Update() {
	pb.completed++
	pb.render()
}

// SetProgress jumps to a given “completed” count and refreshes.
func (pb *ProgressBar) SetProgress(completed int) {
	pb.completed = completed
	pb.render()
}

// Finish marks the bar fully complete, then prints a newline.
func (pb *ProgressBar) Finish() {
	pb.completed = pb.total
	pb.render()
	fmt.Fprintln(os.Stderr)
}

// render draws five lines to stderr (overwriting previous five if not first frame).
// 1) If headerLine != "", print headerLine (alone).
// 2) Zipper line (“3′” + zipper characters spanning pb.width).
// 3) Top strand: “--” + first pos bases of template.
// 4) Complement: “--” + first pos bases of complement.
// 5) Primer line: “5′” + `┴` repeated pos times + “===>”.
// 6) Percentage line “xx.x% (c/t)”.
func (pb *ProgressBar) render() {
	if pb.total == 0 {
		return
	}

	// 1) Calculate how many bases to “fill in” (pos), scaled to width.
	pos := pb.completed * pb.width / pb.total
	if pos > pb.width {
		pos = pb.width
	}

	// 2) Build zipper line with “3′” label.
	lineZipper := "3'" + strings.Repeat(zipperChar, pb.width)

	// 3) Build top-strand (template) showing only the first pos bases, with “--” in front.
	var lineTop string
	if pos <= len(pb.topStrand) {
		lineTop = "--" + pb.topStrand[:pos]
	} else {
		lineTop = "--" + pb.topStrand
	}

	// 4) Build complement line similarly.
	var lineComplement string
	if pos <= len(pb.complement) {
		lineComplement = "--" + pb.complement[:pos]
	} else {
		lineComplement = "--" + pb.complement
	}

	// 5) Build primer line (“5′” + baseChar × pos + arrow).
	var linePrimer string
	if pos < pb.width {
		linePrimer = "5'" + strings.Repeat(baseChar, pos) + arrowText
	} else {
		linePrimer = "5'" + strings.Repeat(baseChar, pb.width) + arrowText
	}

	// 6) Percentage line
	percent := float64(pb.completed) / float64(pb.total) * 100
	linePercent := fmt.Sprintf("%.1f%% (%d/%d)", percent, pb.completed, pb.total)

	// 7) If not the very first frame (completed > 0), move cursor up 5 lines to overwrite.
	if pb.completed > 0 {
		for i := 0; i < 5+(func() int {
			if pb.headerLine != "" {
				return 1
			}
			return 0
		}()); i++ {
			// If headerLine exists, that's one extra line to overwrite.
			fmt.Fprint(os.Stderr, "\033[F")
		}
	}

	// 8) Actually print:
	if pb.headerLine != "" {
		fmt.Fprintln(os.Stderr, pb.headerLine)
	}
	fmt.Fprintln(os.Stderr, lineZipper)
	fmt.Fprintln(os.Stderr, lineTop)
	fmt.Fprintln(os.Stderr, lineComplement)
	fmt.Fprintln(os.Stderr, linePrimer)
	fmt.Fprintln(os.Stderr, linePercent)
}
