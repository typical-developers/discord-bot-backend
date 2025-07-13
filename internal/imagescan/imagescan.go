package imagescan

import (
	"fmt"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"sort"

	"github.com/corona10/goimagehash"
)

func perceptionHash(file *os.File) (*goimagehash.ImageHash, error) {
	var hash *goimagehash.ImageHash

	ext := filepath.Ext(file.Name())

	switch ext {
	case ".jpg", ".jpeg":
		image, err := jpeg.Decode(file)
		if err != nil {
			return nil, err
		}

		hash, err = goimagehash.PerceptionHash(image)
		if err != nil {
			return nil, err
		}
	case ".png":
		image, err := png.Decode(file)
		if err != nil {
			return nil, err
		}

		hash, err = goimagehash.PerceptionHash(image)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unsupported image type: %s", ext)
	}

	return hash, nil
}

// Compare an image to predefined hashes.
// This will return an array of results, with the most similar hash being first.
func Compare(file *os.File) (results ComparisonResults, err error) {
	compareHash, err := perceptionHash(file)
	if err != nil {
		return results, err
	}

	for _, hash := range hashes {
		distance, err := compareHash.Distance(&hash)
		if err != nil {
			return results, err
		}

		results = append(results, ComparisonResult{
			InputHash:       hash.GetHash(),
			Hash:            compareHash.GetHash(),
			HammingDistance: distance,
		})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].HammingDistance < results[j].HammingDistance
	})

	return results, err
}
