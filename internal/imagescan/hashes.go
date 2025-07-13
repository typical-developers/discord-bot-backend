package imagescan

import (
	"io/fs"
	"os"
	"path/filepath"
	"sync"

	"github.com/corona10/goimagehash"
	"github.com/typical-developers/discord-bot-backend/pkg/logger"
)

var (
	once   sync.Once
	hashes = map[string]goimagehash.ImageHash{}
)

func initHashes() {
	if _, err := os.Stat("assets/image-scan"); err != nil {
		logger.Log.Warn("Image scan assets directory not found. Hashes will not be preloaded.")
		return
	}

	// Preload all of the image hashes in the image-scan folder.
	// For deployments, this folder will have to be mounted to the container.
	err := filepath.Walk("assets/image-scan", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		hash, err := perceptionHash(file)
		if err != nil {
			return err
		}

		hashes[path] = *hash
		return nil
	})

	if err != nil {
		panic(err)
	}
}

func init() {
	once.Do(initHashes)
}
