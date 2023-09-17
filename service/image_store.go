package service

import (
	"bytes"
	"fmt"
	"os"
	"sync"

	"github.com/google/uuid"
)

// ImageStore is an interface for store laptop images
type ImageStore interface {
	// Save saves the laptop images
	Save(laptopId string, imageType string, imageData bytes.Buffer) (string, error)
}

// DiskImageStore stores images on the disk and its info in memory
type DiskImageStore struct {
	mutex       sync.RWMutex
	imageFolder string
	images      map[string]*ImageInfo
}

// ImageInfo contains information about the laptop image
type ImageInfo struct {
	LaptopId string
	Type     string
	Path     string
}

// NewDiskImageStore returns a new DiskImageStore
func NewDiskImageStore(imageFolder string) *DiskImageStore {
	return &DiskImageStore{
		imageFolder: imageFolder,
		images:      make(map[string]*ImageInfo),
	}
}

// Save saves a new laptop image to the store
func (s *DiskImageStore) Save(laptopId string, imageType string, imageData bytes.Buffer) (string, error) {
	imageId, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("cannot generate image id: %w", err)
	}

	imagePath := fmt.Sprintf("%s/%s%s", s.imageFolder, imageId, imageType)

	file, err := os.Create(imagePath)
	if err != nil {
		return "", fmt.Errorf("cannot create image file: %w", err)
	}

	_, err = imageData.WriteTo(file)
	if err != nil {
		return "", fmt.Errorf("cannot write image to file: %w", err)
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.images[imageId.String()] = &ImageInfo{
		LaptopId: laptopId,
		Type:     imageType,
		Path:     imagePath,
	}

	return imageId.String(), nil
}
