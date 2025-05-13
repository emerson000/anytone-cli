package codeplug

import (
	"fmt"
)

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

	NameOffset  int64
	NameLength  int
	TotalLength int
}

func (cp *Codeplug) readChannelMetadata(offset int64) (*Channel, error) {
	adjustedOffset := offset

	const nameOffset = 49
	header := make([]byte, nameOffset)
	if _, err := cp.file.ReadAt(header, adjustedOffset); err != nil {
		return nil, fmt.Errorf("failed to read channel header at offset %d: %w", adjustedOffset, err)
	}

	nameStartOffset := adjustedOffset + nameOffset
	nameBuf := make([]byte, 32)
	if _, err := cp.file.ReadAt(nameBuf, nameStartOffset); err != nil {
		return nil, fmt.Errorf("failed to read channel name at offset %d: %w", nameStartOffset, err)
	}

	nameLength := 0
	for i := 0; i < len(nameBuf); i++ {
		if nameBuf[i] == 0x00 {
			nameLength = i + 1
			break
		}
	}

	if nameLength == 0 {
		return nil, fmt.Errorf("invalid channel name at offset %d: no null terminator found", nameStartOffset)
	}

	trailingFieldsOffset := nameStartOffset + int64(nameLength)
	trailingFields := make([]byte, 27)

	if _, err := cp.file.ReadAt(trailingFields, trailingFieldsOffset); err != nil {
		return nil, fmt.Errorf("failed to read trailing fields at offset %d: %w", trailingFieldsOffset, err)
	}

	totalLength := nameOffset + nameLength + len(trailingFields)

	channel := &Channel{
		RxFreq:               uint32(header[3]) | uint32(header[4])<<8 | uint32(header[5])<<16 | uint32(header[6])<<24,
		TxFreqDirection:      header[7],
		TxFreq:               int32(header[8]) | int32(header[9])<<8 | int32(header[10])<<16 | int32(header[11])<<24,
		ChannelType:          header[12],
		TxPower:              header[13],
		Bandwidth:            header[14],
		PttProhibit:          header[16],
		CallConfirmation:     header[17],
		TalkAround:           header[18],
		CtcssDcsDecode:       header[19],
		CtcssDcsDecodeOption: header[20],
		CtcssDcsEncode:       header[23],
		CtcssDcsEncodeOption: header[24],
		Contact:              header[29],
		RadioId:              header[31],
		TxPermit:             header[33],
		SquelchMode:          header[34],
		ScanList:             int8(header[35]),
		ReceiveGroupList:     header[36],
		RxColorCode:          header[41],
		Slot:                 header[42],
		SlotSuit:             header[44],
		AprsRx:               header[45],
		AesEncryptionKey:     header[46],
		WorkAlone:            header[47],
		Name:                 string(nameBuf[:nameLength-1]),

		Ranging:            trailingFields[2],
		CorrectFreq:        int8(trailingFields[8]),
		SmsConfirmation:    trailingFields[11],
		ExcludeFromRoaming: trailingFields[12],
		MultipleKey:        trailingFields[15],
		RandomKey:          trailingFields[16],
		SmsForbid:          trailingFields[17],
		DataAckDisable:     trailingFields[18],
		AutoScan:           trailingFields[21],
		SendTalkerAlias:    getSafeByteValue(trailingFields, 22),
		ExtendEncryption:   getSafeByteValue(trailingFields, 27),

		NameOffset:  nameStartOffset,
		NameLength:  nameLength,
		TotalLength: totalLength,
	}

	return channel, nil
}
