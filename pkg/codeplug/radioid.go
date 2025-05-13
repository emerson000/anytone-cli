package codeplug

import (
	"fmt"
)

type RadioIDEntry struct {
	Index    int
	ID       int
	Name     string
	Position int64
	Length   int
}

func (cp *Codeplug) calculateRadioIDOffset() (int64, error) {
	channelCountBuf := make([]byte, 1)
	if _, err := cp.file.ReadAt(channelCountBuf, totalChannelsAddress); err != nil {
		return 0, fmt.Errorf("failed to read total channels: %w", err)
	}

	totalChannels := int(channelCountBuf[0])

	channelsStartOffset := int64(totalChannelsAddress + 1)

	currentOffset := channelsStartOffset

	for i := 0; i < totalChannels; i++ {
		channel, err := cp.readChannelMetadata(currentOffset)
		if err != nil {
			return 0, fmt.Errorf("failed to read channel %d: %w", i+1, err)
		}

		currentOffset += int64(channel.TotalLength)
	}

	radioIDOffset := currentOffset + 2

	return radioIDOffset, nil
}

func (cp *Codeplug) readRadioIDEntry(offset int64, previousIndex int) (*RadioIDEntry, error) {
	idHeader := make([]byte, 4)
	if _, err := cp.file.ReadAt(idHeader, offset); err != nil {
		return nil, fmt.Errorf("failed to read radio ID header at offset %d: %w", offset, err)
	}

	index := int(idHeader[0])

	if index < previousIndex {
		return nil, nil
	}

	id := int(uint32(idHeader[1]) | uint32(idHeader[2])<<8 | uint32(idHeader[3])<<16)

	buf := make([]byte, 256)
	if _, err := cp.file.ReadAt(buf, offset+4); err != nil {
		return nil, fmt.Errorf("failed to read radio ID name at offset %d: %w", offset+4, err)
	}

	nameLength := 0
	for j := 0; j < len(buf); j++ {
		if buf[j] == 0 {
			nameLength = j + 1
			break
		}
	}

	name := string(buf[:nameLength-1])

	return &RadioIDEntry{
		Index:    index,
		ID:       id,
		Name:     name,
		Position: offset,
		Length:   4 + nameLength,
	}, nil
}

func (cp *Codeplug) writeRadioIDEntry(entry *RadioIDEntry) error {
	totalLength := 4 + len(entry.Name) + 1
	buf := make([]byte, totalLength)

	buf[0] = byte(entry.Index)

	buf[1] = byte(entry.ID & 0xFF)
	buf[2] = byte((entry.ID >> 8) & 0xFF)
	buf[3] = byte((entry.ID >> 16) & 0xFF)

	copy(buf[4:], entry.Name)

	if _, err := cp.file.WriteAt(buf, entry.Position); err != nil {
		return fmt.Errorf("failed to write radio ID entry: %w", err)
	}

	return nil
}

func (cp *Codeplug) UpdateRadioID(index int, newID int) error {
	if index < 0 || index >= maxRadioIDs {
		return fmt.Errorf("invalid radio ID index: %d", index)
	}

	radioIDOffset, err := cp.calculateRadioIDOffset()
	if err != nil {
		return fmt.Errorf("failed to calculate radio ID offset: %w", err)
	}

	var currentOffset int64 = radioIDOffset
	var entry *RadioIDEntry
	var previousIndex = -1
	var entries []*RadioIDEntry

	for i := 0; i < maxRadioIDs; i++ {
		entry, err = cp.readRadioIDEntry(currentOffset, previousIndex)
		if err != nil {
			return err
		}

		if entry == nil {
			break
		}

		entries = append(entries, entry)
		previousIndex = entry.Index
		currentOffset += int64(entry.Length)
	}

	var targetEntry *RadioIDEntry
	for _, e := range entries {
		if e.Index == index {
			targetEntry = e
			break
		}
	}

	if targetEntry != nil {
		targetEntry.ID = newID
		return cp.writeRadioIDEntry(targetEntry)
	}

	var insertPosition int64 = radioIDOffset
	for _, e := range entries {
		if e.Index > index {
			break
		}
		insertPosition = e.Position + int64(e.Length)
	}

	newEntry := &RadioIDEntry{
		Index:    index,
		ID:       newID,
		Name:     fmt.Sprintf("Radio ID %d", index+1),
		Position: insertPosition,
		Length:   4 + len(fmt.Sprintf("Radio ID %d", index+1)) + 1,
	}

	return cp.writeRadioIDEntry(newEntry)
}
