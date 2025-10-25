package qrencode

import (
	"errors"
)

func stringContentBits(content string, ecLevel ECLevel) (*BitVector, versionNumber, error) {
	if !supportedECLevel(ecLevel) {
		return nil, versionNumber(0), errors.New("Unrecognized ECLevel")
	}
	headerBits := &BitVector{}
	mode := getMode(content)
	if mode == modeByte {
		headerBits.Append(int(modeECI), 4)
		headerBits.Append(26, 8) // UTF-8
	}
	headerBits.Append(int(mode), 4)
	return contentBits([]byte(content), ecLevel, mode, headerBits)
}

func binaryContentBits(content []byte, ecLevel ECLevel) (*BitVector, versionNumber, error) {
	if !supportedECLevel(ecLevel) {
		return nil, versionNumber(0), errors.New("Unrecognized ECLevel")
	}
	headerBits := &BitVector{}
	headerBits.Append(int(modeByte), 4)
	return contentBits(content, ecLevel, modeByte, headerBits)
}

func contentBits(content []byte, ecLevel ECLevel, mode modeIndicator, headerBits *BitVector) (*BitVector, versionNumber, error) {
	dataBits := BitVector{}
	appendContent(content, mode, &dataBits)

	bitsNeeded := headerBits.Length() + dataBits.Length() + mode.characterCountBits(versionNumber(40))
	version, err := chooseVersion(bitsNeeded, ecLevel)
	if err != nil {
		return nil, version, err
	}

	headerAndDataBits := &BitVector{}
	headerAndDataBits.AppendBits(*headerBits)
	headerAndDataBits.Append(len(content), mode.characterCountBits(version))
	headerAndDataBits.AppendBits(dataBits)

	appendTerminator(version.totalCodewords()-ecBlocks[version][ecLevel].totalECCodewords(), headerAndDataBits)
	return headerAndDataBits, version, nil
}

var (
	invalidAlphanumericByte = errors.New("Invalid Alphanumeric Byte")
	alphanumericTable       = [256]int8{
		-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
		-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
		36, -1, -1, -1, 37, 38, -1, -1, -1, -1, 39, 40, -1, 41, 42, 43,
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 44, -1, -1, -1, -1, -1,
		-1, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24,
		25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, -1, -1, -1, -1, -1,
		-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
		-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
		-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
		-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
		-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
		-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
		-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
		-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
		-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
		-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
	}
)

func alphanumericCode(b byte) (byte, error) {
	code := alphanumericTable[b]
	if code < 0 {
		return 0, invalidAlphanumericByte
	}
	return byte(code), nil
}

func appendContent(content []byte, mode modeIndicator, bits *BitVector) {
	switch mode {
	case modeNumeric:
		for i := 0; i+2 < len(content); i += 3 {
			n1, err := alphanumericCode(content[i])
			if n1 > 9 || err != nil {
				panic("Invalid numeric mode content")
			}
			n2, err := alphanumericCode(content[i+1])
			if n2 > 9 || err != nil {
				panic("Invalid numeric mode content")
			}
			n3, err := alphanumericCode(content[i+2])
			if n3 > 9 || err != nil {
				panic("Invalid numeric mode content")
			}
			bits.Append(int(n1)*100+int(n2)*10+int(n3), 10)
		}
		switch len(content) % 3 {
		case 1:
			n1, err := alphanumericCode(content[len(content)-1])
			if n1 > 9 || err != nil {
				panic("Invalid numeric mode content")
			}
			bits.Append(int(n1), 4)
		case 2:
			n1, err := alphanumericCode(content[len(content)-2])
			if n1 > 9 || err != nil {
				panic("Invalid numeric mode content")
			}
			n2, err := alphanumericCode(content[len(content)-1])
			if n2 > 9 || err != nil {
				panic("Invalid numeric mode content")
			}
			bits.Append(int(n1)*10+int(n2), 7)
		}
	case modeAlphanumeric:
		for i := 0; i+1 < len(content); i += 2 {
			n1, err := alphanumericCode(content[i])
			if err != nil {
				panic("Invalid alphanumeric mode content")
			}
			n2, err := alphanumericCode(content[i+1])
			if err != nil {
				panic("Invalid alphanumeric mode content")
			}
			bits.Append(int(n1)*45+int(n2), 11)
		}
		if len(content)%2 == 1 {
			n1, err := alphanumericCode(content[len(content)-1])
			if err != nil {
				panic("Invalid alphanumeric mode content")
			}
			bits.Append(int(n1), 6)
		}
	case modeByte:
		for _, b := range content {
			bits.Append(int(b), 8)
		}
	default:
		panic("Unsupported mode")
	}
}

func appendTerminator(capacityBytes int, bits *BitVector) {
	capacity := capacityBytes * 8
	if bits.Length() > capacity {
		panic("bits.Length() > capacity")
	}
	for i := 0; i < 4 && bits.Length() < capacity; i++ {
		bits.AppendBit(false)
	}
	if bits.Length()%8 != 0 {
		for i := bits.Length() % 8; i < 8; i++ {
			bits.AppendBit(false)
		}
	}
	for {
		if bits.Length() >= capacity {
			break
		}
		bits.Append(0xec, 8)
		if bits.Length() >= capacity {
			break
		}
		bits.Append(0x11, 8)
	}
	if bits.Length() != capacity {
		panic("bits.Length() != capacity")
	}
}
