package codeplug

import (
	"fmt"
	"os"
)

const (
	headerSize           = 0x100
	totalChannelsAddress = 0xF1
	modelOffset          = 0x09
	modelSize            = 10
	maxRadioIDs          = 10
)

type Codeplug struct {
	file *os.File
	path string
}

type Info struct {
	Model          string
	RadioIDs       []int
	RadioIDIndices []int
}

func Open(path string) (*Codeplug, error) {
	file, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	return &Codeplug{
		file: file,
		path: path,
	}, nil
}

func (cp *Codeplug) Close() error {
	return cp.file.Close()
}

func getSafeByteValue(data []byte, index int) byte {
	if index >= 0 && index < len(data) {
		return data[index]
	}
	return 0
}

func (cp *Codeplug) GetInfo() (*Info, error) {
	model := make([]byte, modelSize)
	if _, err := cp.file.ReadAt(model, modelOffset); err != nil {
		return nil, fmt.Errorf("failed to read model: %w", err)
	}

	radioIDs, err := cp.GetRadioIDs()
	if err != nil {
		return nil, fmt.Errorf("failed to get radio IDs: %w", err)
	}

	ids := make([]int, len(radioIDs))
	indices := make([]int, len(radioIDs))
	for i, entry := range radioIDs {
		ids[i] = entry.ID
		indices[i] = entry.Index
	}

	return &Info{
		Model:          string(model),
		RadioIDs:       ids,
		RadioIDIndices: indices,
	}, nil
}
