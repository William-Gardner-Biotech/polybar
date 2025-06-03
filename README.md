# DNA Progress Bar

A Go package that provides a DNA-style progress bar with base complementing, designed to look like a DNA polymerase in action.

## Features

- **DNA Complementing**: Automatically generates complement strand (A↔T, G↔C, others→N)
- **Custom Sequences**: Use your own DNA sequence as the top strand
- **Visual Design**: Looks like DNA replication with zipper, strands, and primer
- **Thread Safe**: Safe to use from multiple goroutines
- **Customizable Width**: Match header text or use sequence length

## Installation

```bash
go get github.com/yourusername/polybar
```

## Usage

### Basic Usage

```go
package main

import (
    "time"
    "github.com/William-Gardner-Biotech/polybar"
)

func main() {
    // Create progress bar with custom DNA sequence
    pb := polybar.New("ATCG-NNTA-GCTA", "DNA-SEQUENCING")

    // Start with total of 100 steps
    pb.Start(100)

    // Simulate work with progress updates
    for i := 0; i < 100; i++ {
        // Do some work...
        time.Sleep(50 * time.Millisecond)
        pb.Update()
    }

    pb.Finish()
}
```

### Advanced Usage

```go
// Use sequence length as width (no header padding)
pb := polybar.New("ATCGATCGATCG", "")

// Set progress directly instead of incrementing
pb.Start(1000)
pb.SetProgress(250)  // 25% complete
pb.SetProgress(500)  // 50% complete
pb.Finish()
```

## Output Example

```
DNA-SEQUENCING
┬┬┬┬┬┬┬┬┬┬┬┬┬┬
ATCG-NNTA-G
TAGC-NNAT-C
┴┴┴┴┴┴┴┴┴┴┴====>
75.0% (75/100)
```

## DNA Complement Rules

- **A** ↔ **T** (Adenine ↔ Thymine)
- **G** ↔ **C** (Guanine ↔ Cytosine)
- **-** → **-** (Gap remains gap)
- **Any other character** → **N** (Unknown base)

## API Reference

### Functions

#### `New(topStrand string, header string) *ProgressBar`
Creates a new DNA progress bar.
- `topStrand`: DNA sequence for the top strand (will be complemented)
- `header`: Optional header text. If provided, strands are padded/truncated to match width

#### Methods

- `Start(total int)`: Initialize progress bar with total steps
- `Update()`: Increment progress by 1 and refresh display
- `SetProgress(completed int)`: Set current progress value
- `Finish()`: Complete progress bar and add final newline

## License

MIT License
