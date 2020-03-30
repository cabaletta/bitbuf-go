package bitbuf

type BitSliceMut struct {
	*sliceContainer
}

func NewBitSliceMut(data []byte) BitBufMut {
	return BitSliceMut{&sliceContainer{
		data:   data,
		offset: 0,
		length: 0,
	}}
}

func (s BitSliceMut) write(data byte, bits Size) error {
	bitsInv := 8 - bits
	low := data & (0xFF << bitsInv)
	high := data | (0xFF >> bitsInv)
	if s.offset == 0 {
		s.data[0] |= low
		s.data[0] &= high
	} else {
		offsetInv := 8 - s.offset
		s.data[0] |= low >> s.offset
		s.data[0] &= (high >> s.offset) | (0xFF << offsetInv)
		s.data[1] |= low >> offsetInv
		s.data[1] &= (high >> offsetInv) | (0xFF << s.offset)
	}
	return s.Advance(bits)
}

func (s BitSliceMut) Advance(bits Size) error {
	if s.Remaining() < bits {
		return InsufficientError{}
	}
	s.offset += byte(bits & 7)
	if s.offset >= 8 {
		s.offset -= 8
		s.data = s.data[bits/8+1:]
	} else {
		s.data = s.data[bits/8:]
	}
	s.length += bits
	return nil
}

func (s BitSliceMut) Write(data []byte, bits Size) (Size, error) {
	rem := s.Remaining()
	if bits > rem {
		bits = rem
	}
	if bits == 0 {
		return 0, nil
	}

	bytes := bits / 8
	length := len(data)
	if length == 0 || length*8 < int(bits) {
		return 0, OverflowError{}
	}

	for i := Size(0); i < bytes; i++ {
		err := s.WriteByte(data[i])
		if err != nil {
			panic(err)
		}
	}

	rem = bits & 7
	if rem != 0 {
		if length < int(bytes)+1 {
			return 0, OverflowError{}
		}
		val := data[bytes]
		err := s.write(val, rem)
		if err != nil {
			panic(err)
		}
	}

	return bits, nil
}

func (s BitSliceMut) WriteAll(data []byte, bits Size) error {
	if s.Remaining() < bits {
		return InsufficientError{}
	}

	_, err := s.Write(data, bits)
	if err != nil {
		return OverflowError{}
	}
	return nil
}

func (s BitSliceMut) WriteAligned(data []byte) Size {
	length, err := s.Write(data, Size(len(data))*8)
	if err != nil {
		panic(err)
	}
	return length
}

func (s BitSliceMut) WriteAlignedAll(data []byte) error {
	bits := Size(len(data))
	if bits > s.Remaining() {
		return InsufficientError{}
	}

	_, err := s.Write(data, bits)
	if err != nil {
		panic(err)
	}
	return nil
}

func (s BitSliceMut) WriteBool(val bool) error {
	if len(s.data) == 0 {
		return InsufficientError{}
	}
	data := &s.data[0]
	if val {
		*data |= 0x80 >> s.offset
	} else {
		*data &= 0xFF ^ (0x80 >> s.offset)
	}
	return s.Advance(1)
}

func (s BitSliceMut) WriteByte(val byte) error {
	if len(s.data) == 0 {
		return InsufficientError{}
	}
	if s.offset == 0 {
		s.data[0] = val
	} else {
		if len(s.data) == 1 {
			return InsufficientError{}
		}
		offsetInv := 8 - s.offset
		s.data[0] |= val >> s.offset
		s.data[0] &= (val >> s.offset) | (0xFF << offsetInv)
		s.data[1] |= val << offsetInv
		s.data[1] &= (val << offsetInv) | (0xFF << s.offset)
	}
	return s.Advance(8)
}

func (s BitSliceMut) Remaining() Size {
	return Size(len(s.data)*8 - int(s.offset))
}

func (s BitSliceMut) Len() Size {
	return s.length
}
