package dna

import (
	"fmt"

	"github.com/evolbioinfo/goalign/align"
)

type PDistModel struct {
	numSites      float64 // Number of selected sites (no gaps)
	selectedSites []bool  // true for selected sites
	removegaps    bool    // If true, we will remove posision with >=1 gaps
	// If 0, will not count as 1 mutation '-' to 'A"
	// If 1, will count as 1 mutation '-' to 'A"
	// If 2, will count as 1 mutation '-' to 'A", but only the internal
	// Default 0
	countgapmut int
}

func NewPDistModel(removegaps bool) *PDistModel {
	return &PDistModel{
		0,
		nil,
		removegaps,
		0,
	}
}

func (m *PDistModel) SetCountGapMutations(countgapmut int) (err error) {
	if countgapmut < 0 || countgapmut > 2 {
		err = fmt.Errorf("Gap count mode not available : %d", countgapmut)
	} else {
		m.countgapmut = countgapmut
	}
	return
}

/* computes p-distance between 2 sequences */
func (m *PDistModel) Distance(seq1 []rune, seq2 []rune, weights []float64) (diff float64, err error) {
	var total float64
	switch m.countgapmut {
	case GAP_COUNT_ALL:
		diff, total = countDiffsWithGaps(seq1, seq2, m.selectedSites, weights)
	case GAP_COUNT_INTERNAL:
		diff, total = countDiffsWithInternalGaps(seq1, seq2, m.selectedSites, weights)
	default:
		diff, total = countDiffs(seq1, seq2, m.selectedSites, weights)
	}
	diff = diff / total
	return
}

func (m *PDistModel) InitModel(al align.Alignment, weights []float64, gamma bool, alpha float64) (err error) {
	m.numSites, m.selectedSites = selectedSites(al, weights, m.removegaps)
	return
}
