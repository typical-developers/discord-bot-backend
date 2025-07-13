package imagescan

type ComparisonResult struct {
	InputHash       uint64
	Hash            uint64
	HammingDistance int
}

// Convert the hamming distance to a percentage for more simplified comparison.
func (r *ComparisonResult) Percentage() int {
	return (64 - r.HammingDistance) * 100 / 64
}

type ComparisonResults []ComparisonResult
