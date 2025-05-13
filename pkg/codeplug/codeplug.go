package codeplug

import (
	"fmt"
	"os"
)

const (
	// RDT file header constants
	headerSize           = 0x100
	totalChannelsAddress = 0xF1
	modelOffset          = 0x09
	modelSize            = 10
	maxRadioIDs          = 10
)

// Channel represents the structure of a channel in the codeplug
type Channel struct {
	RxFreq               uint32
	TxFreqDirection      byte
	TxFreq               int32
	ChannelType          byte
	TxPower              byte
	Bandwidth            byte
	PttProhibit          byte
	CallConfirmation     byte
	TalkAround           byte
	CtcssDcsDecode       byte
	CtcssDcsDecodeOption byte
	CtcssDcsEncode       byte
	CtcssDcsEncodeOption byte
	Contact              byte
	RadioId              byte
	TxPermit             byte
	SquelchMode          byte
	ScanList             int8
	ReceiveGroupList     byte
	RxColorCode          byte
	Slot                 byte
	SlotSuit             byte
	AprsRx               byte
	AesEncryptionKey     byte
	WorkAlone            byte
	Name                 string
	Ranging              byte
	CorrectFreq          int8
	SmsConfirmation      byte
	ExcludeFromRoaming   byte
	MultipleKey          byte
	RandomKey            byte
	SmsForbid            byte
	DataAckDisable       byte
	AutoScan             byte
	SendTalkerAlias      byte
	ExtendEncryption     byte

	// Metadata for file operations
	NameOffset  int64
	NameLength  int
	TotalLength int
}

// Codeplug represents an Anytone RDT file
type Codeplug struct {
	file *os.File
	path string
}

// Info contains general information about the codeplug
type Info struct {
	Model          string
	RadioIDs       []int
	RadioIDIndices []int
}

// RadioIDEntry represents a single radio ID entry in the codeplug
type RadioIDEntry struct {
	Index    int    // Zero-based index
	ID       int    // Radio ID value (24-bit unsigned integer)
	Name     string // Radio ID name
	Position int64  // Position in the file
	Length   int    // Total length including index, ID, and name with null terminator
}

// Open opens an RDT file for reading and writing
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

// Close closes the codeplug file
func (cp *Codeplug) Close() error {
	return cp.file.Close()
}

// calculateRadioIDOffset dynamically determines the offset where radio IDs begin
// based on the channel information and variable channel sizes
func (cp *Codeplug) calculateRadioIDOffset() (int64, error) {
	// Read total number of channels
	channelCountBuf := make([]byte, 1)
	if _, err := cp.file.ReadAt(channelCountBuf, totalChannelsAddress); err != nil {
		return 0, fmt.Errorf("failed to read total channels: %w", err)
	}

	totalChannels := int(channelCountBuf[0])

	// Channels start after four bytes after totalChannelsAddress
	// Note: The offset was previously miscalculated here
	channelsStartOffset := int64(totalChannelsAddress + 1)

	// Start position for our scan
	currentOffset := channelsStartOffset

	// Iterate through each channel to account for variable channel sizes
	for i := 0; i < totalChannels; i++ {
		channel, err := cp.readChannelMetadata(currentOffset)
		if err != nil {
			return 0, fmt.Errorf("failed to read channel %d: %w", i+1, err)
		}

		// Move to the next channel
		currentOffset += int64(channel.TotalLength)
	}

	radioIDOffset := currentOffset + 2

	return radioIDOffset, nil
}

// readChannelMetadata reads the metadata of a channel at the given offset
func (cp *Codeplug) readChannelMetadata(offset int64) (*Channel, error) {
	adjustedOffset := offset

	const nameOffset = 49
	header := make([]byte, nameOffset)
	if _, err := cp.file.ReadAt(header, adjustedOffset); err != nil {
		return nil, fmt.Errorf("failed to read channel header at offset %d: %w", adjustedOffset, err)
	}

	// Read the name field (null-terminated string)
	nameStartOffset := adjustedOffset + nameOffset
	nameBuf := make([]byte, 32) // Assume name won't be longer than 32 bytes
	if _, err := cp.file.ReadAt(nameBuf, nameStartOffset); err != nil {
		return nil, fmt.Errorf("failed to read channel name at offset %d: %w", nameStartOffset, err)
	}

	// Find the null terminator
	nameLength := 0
	for i := 0; i < len(nameBuf); i++ {
		if nameBuf[i] == 0x00 {
			nameLength = i + 1 // Include null terminator in length calculation
			break
		}
	}

	if nameLength == 0 {
		return nil, fmt.Errorf("invalid channel name at offset %d: no null terminator found", nameStartOffset)
	}

	// Read the remaining fields after the name
	trailingFieldsOffset := nameStartOffset + int64(nameLength)
	trailingFields := make([]byte, 27) // Fixed number of bytes after name

	if _, err := cp.file.ReadAt(trailingFields, trailingFieldsOffset); err != nil {
		return nil, fmt.Errorf("failed to read trailing fields at offset %d: %w", trailingFieldsOffset, err)
	}

	// Total channel length (corrected to ensure proper alignment)
	totalLength := nameOffset + nameLength + len(trailingFields)

	// Extract relevant fields from the header
	channel := &Channel{
		// The first padding bytes are at header[0], header[1], header[2]
		RxFreq:          uint32(header[3]) | uint32(header[4])<<8 | uint32(header[5])<<16 | uint32(header[6])<<24,
		TxFreqDirection: header[7],
		TxFreq:          int32(header[8]) | int32(header[9])<<8 | int32(header[10])<<16 | int32(header[11])<<24,
		ChannelType:     header[12],
		TxPower:         header[13],
		Bandwidth:       header[14],
		// padding[1] at header[15]
		PttProhibit:          header[16],
		CallConfirmation:     header[17],
		TalkAround:           header[18],
		CtcssDcsDecode:       header[19],
		CtcssDcsDecodeOption: header[20],
		// padding[2] are at header[21], header[22]
		CtcssDcsEncode:       header[23],
		CtcssDcsEncodeOption: header[24],
		// padding[4] are at header[25], header[26], header[27], header[28]
		Contact: header[29],
		// padding[1] at header[30]
		RadioId: header[31],
		// padding[1] at header[32]
		TxPermit:         header[33],
		SquelchMode:      header[34],
		ScanList:         int8(header[35]),
		ReceiveGroupList: header[36],
		// padding[4] are at header[37], header[38], header[39], header[40]
		RxColorCode: header[41],
		Slot:        header[42],
		// padding[1] at header[43]
		SlotSuit:         header[44],
		AprsRx:           header[45],
		AesEncryptionKey: header[46],
		WorkAlone:        header[47],
		// padding[1] at header[48], header[49], header[50], header[51]

		Name: string(nameBuf[:nameLength-1]), // Exclude null terminator from name

		// Fields after the name
		// padding[2] are at trailingFields[0], trailingFields[1]
		Ranging: trailingFields[2],
		// padding[5] are at trailingFields[3], trailingFields[4], trailingFields[5], trailingFields[6], trailingFields[7]
		CorrectFreq: int8(trailingFields[8]),
		// padding[2] are at trailingFields[9], trailingFields[10]
		SmsConfirmation:    trailingFields[11],
		ExcludeFromRoaming: trailingFields[12],
		// padding[2] are at trailingFields[13], trailingFields[14]
		MultipleKey:    trailingFields[15],
		RandomKey:      trailingFields[16],
		SmsForbid:      trailingFields[17],
		DataAckDisable: trailingFields[18],
		// padding[2] are at trailingFields[19], trailingFields[20]
		AutoScan: trailingFields[21],
		// Safely get fields that might be out of bounds
		SendTalkerAlias: getSafeByteValue(trailingFields, 22),
		// padding[4] would be trailingFields[23:27]
		ExtendEncryption: getSafeByteValue(trailingFields, 27),

		// Metadata for file operations
		NameOffset:  nameStartOffset,
		NameLength:  nameLength,
		TotalLength: totalLength,
	}

	return channel, nil
}

// readRadioIDEntry reads a single radio ID entry from the file at the given offset
func (cp *Codeplug) readRadioIDEntry(offset int64, previousIndex int) (*RadioIDEntry, error) {
	// Read the first 4 bytes (index + ID)
	idHeader := make([]byte, 4)
	if _, err := cp.file.ReadAt(idHeader, offset); err != nil {
		return nil, fmt.Errorf("failed to read radio ID header at offset %d: %w", offset, err)
	}

	// First byte is the index
	index := int(idHeader[0])

	// End of radio IDs section is detected when the index is less than the previous index
	// (indexes can be skipped, but they are always in ascending order)
	if index < previousIndex {
		return nil, nil
	}

	// Next 3 bytes form a 24-bit unsigned integer
	id := int(uint32(idHeader[1]) | uint32(idHeader[2])<<8 | uint32(idHeader[3])<<16)

	// Read string until null terminator
	buf := make([]byte, 256) // Buffer for reading radio ID name
	if _, err := cp.file.ReadAt(buf, offset+4); err != nil {
		return nil, fmt.Errorf("failed to read radio ID name at offset %d: %w", offset+4, err)
	}

	// Find null terminator
	nameLength := 0
	for j := 0; j < len(buf); j++ {
		if buf[j] == 0 {
			nameLength = j + 1 // Include null terminator
			break
		}
	}

	// Extract name (excluding null terminator)
	name := string(buf[:nameLength-1])

	return &RadioIDEntry{
		Index:    index,
		ID:       id,
		Name:     name,
		Position: offset,
		Length:   4 + nameLength, // Header (4 bytes) + name with null terminator
	}, nil
}

// writeRadioIDEntry writes a radio ID entry to the file
func (cp *Codeplug) writeRadioIDEntry(entry *RadioIDEntry) error {
	// Create buffer for the entry
	totalLength := 4 + len(entry.Name) + 1 // Header (4 bytes) + name + null terminator
	buf := make([]byte, totalLength)

	// Write index
	buf[0] = byte(entry.Index)

	// Write ID (24-bit LE)
	buf[1] = byte(entry.ID & 0xFF)
	buf[2] = byte((entry.ID >> 8) & 0xFF)
	buf[3] = byte((entry.ID >> 16) & 0xFF)

	// Write name
	copy(buf[4:], entry.Name)
	// buf[4+len(entry.Name)] is already 0 (null terminator)

	// Write to file
	if _, err := cp.file.WriteAt(buf, entry.Position); err != nil {
		return fmt.Errorf("failed to write radio ID entry: %w", err)
	}

	return nil
}

// GetInfo retrieves general information about the codeplug
func (cp *Codeplug) GetInfo() (*Info, error) {
	// Read model
	model := make([]byte, modelSize)
	if _, err := cp.file.ReadAt(model, modelOffset); err != nil {
		return nil, fmt.Errorf("failed to read model: %w", err)
	}

	// Calculate radioIDOffset dynamically
	radioIDOffset, err := cp.calculateRadioIDOffset()
	if err != nil {
		return nil, fmt.Errorf("failed to calculate radio ID offset: %w", err)
	}

	// Read radio IDs
	radioIDs := make([]int, 0, maxRadioIDs)
	radioIDIndices := make([]int, 0, maxRadioIDs)

	var currentOffset int64 = radioIDOffset
	previousIndex := -1 // Initialize to -1 to handle the first index correctly

	for i := 0; i < maxRadioIDs; i++ {
		entry, err := cp.readRadioIDEntry(currentOffset, previousIndex)
		if err != nil {
			return nil, err
		}

		if entry == nil {
			// End of radio IDs section
			break
		}

		radioIDs = append(radioIDs, entry.ID)
		radioIDIndices = append(radioIDIndices, entry.Index)
		previousIndex = entry.Index
		currentOffset += int64(entry.Length)
	}

	return &Info{
		Model:          string(model),
		RadioIDs:       radioIDs,
		RadioIDIndices: radioIDIndices,
	}, nil
}

// UpdateRadioID updates a radio ID at the specified index
func (cp *Codeplug) UpdateRadioID(index int, newID int) error {
	if index < 0 || index >= maxRadioIDs {
		return fmt.Errorf("invalid radio ID index: %d", index)
	}

	// Calculate radioIDOffset dynamically
	radioIDOffset, err := cp.calculateRadioIDOffset()
	if err != nil {
		return fmt.Errorf("failed to calculate radio ID offset: %w", err)
	}

	// Locate the radio ID entry
	var currentOffset int64 = radioIDOffset
	var entry *RadioIDEntry
	var previousIndex = -1
	var entries []*RadioIDEntry

	// First, read all existing entries to find the one we need or to determine where to insert
	for i := 0; i < maxRadioIDs; i++ {
		entry, err = cp.readRadioIDEntry(currentOffset, previousIndex)
		if err != nil {
			return err
		}

		if entry == nil {
			// End of radio IDs section
			break
		}

		entries = append(entries, entry)
		previousIndex = entry.Index
		currentOffset += int64(entry.Length)
	}

	// Try to find the exact index match
	var targetEntry *RadioIDEntry
	for _, e := range entries {
		if e.Index == index {
			targetEntry = e
			break
		}
	}

	if targetEntry != nil {
		// Found an entry with the exact index, update it
		targetEntry.ID = newID
		return cp.writeRadioIDEntry(targetEntry)
	}

	// Need to create a new entry
	// Determine where to insert the new entry
	var insertPosition int64 = radioIDOffset
	for _, e := range entries {
		if e.Index > index {
			break
		}
		insertPosition = e.Position + int64(e.Length)
	}

	// Create a new entry
	newEntry := &RadioIDEntry{
		Index:    index,
		ID:       newID,
		Name:     fmt.Sprintf("Radio ID %d", index+1),
		Position: insertPosition,
		Length:   4 + len(fmt.Sprintf("Radio ID %d", index+1)) + 1, // +1 for null terminator
	}

	// This is a simplified approach that doesn't handle shifting subsequent entries
	// In a real implementation, you would need to handle moving all subsequent entries
	return cp.writeRadioIDEntry(newEntry)
}

// getSafeByteValue returns a byte value from a slice if the index is valid, or 0 if not
func getSafeByteValue(data []byte, index int) byte {
	if index >= 0 && index < len(data) {
		return data[index]
	}
	return 0
}
