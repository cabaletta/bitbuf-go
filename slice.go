package bitbuf

type BitSlice struct {
	*sliceContainer
}

func NewBitSlice(data []byte) BitSlice {
	return BitSlice{&sliceContainer{
		data:   data,
		offset: 0,
		length: 0,
	}}
}

func (s BitSlice) dataAtOffset(offset, size Size) (byte, error) {
	length := s.Remaining()
	offset += Size(s.offset)
	if offset == 0 {
		if length == 0 {
			return 0, InsufficientError{}
		}
		return s.data[0], nil
	} else if length < size {
		return 0, InsufficientError{}
	} else {
		offsetBytes := offset / 8
		offsetRem := offset & 7
		if offsetRem == 0 {
			return s.data[offsetBytes], nil
		} else {
			data := (s.data[offsetBytes] & (0xFF >> offsetRem)) << offsetRem
			if size+offsetRem <= 8 {
				return data, nil
			} else {
				offsetRem = 8 - offsetRem
				return data + (s.data[(offsetBytes)+1]&(0xFF<<offsetRem))>>offsetRem, nil
			}
		}
	}
}

func (s BitSlice) byteAtOffset(offset Size) (byte, error) {
	return s.dataAtOffset(offset, 8)
}

func (s BitSlice) Advance(bits Size) (e error) {
	if bits > s.Remaining() {
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
	return
}

func (s BitSlice) Read(dst []byte, bits Size) (Size, error) {
	rem := s.Remaining()
	if bits > rem {
		bits = rem
	}

	bytes := bits / 8
	length := len(dst)
	if length*8 < int(bits) {
		return 0, OverflowError{}
	}

	for i := Size(0); i < bytes; i++ {
		data, err := s.byteAtOffset(i * 8)
		if err != nil {
			panic(InsufficientError{})
		}
		dst[i] = data
	}

	rem = bits & 7
	if rem != 0 {
		if length < int(bytes)+1 {
			return 0, OverflowError{}
		}
		data, err := s.dataAtOffset(bytes*8, rem)
		if err != nil {
			panic(InsufficientError{})
		}
		data &= 0xFF << (8 - rem)
		dst[bytes] |= data
		dst[bytes] &= data
	}

	err := s.Advance(bits)
	if err != nil {
		panic(err)
	}
	return bits, nil
}

func (s BitSlice) ReadAll(dst []byte, bits Size) error {
	if s.Remaining() < bits {
		return InsufficientError{}
	}

	_, err := s.Read(dst, bits)
	if err != nil {
		return OverflowError{}
	}
	return nil
}

func (s BitSlice) ReadAligned(dst []byte) Size {
	rem := s.Remaining()
	length := Size(len(dst))
	if length*8 > rem {
		length = rem
	}

	if length&7 != 0 {
		length, err := s.Read(dst, length*8)
		if err != nil {
			panic(err)
		}
		return length
	} else {
		for i := 0; i < len(dst); i++ {
			data, err := s.byteAtOffset(Size(i * 8))
			if err != nil {
				panic(err)
			}
			dst[i] = data
		}
	}
	return length
}

func (s BitSlice) ReadAlignedAll(dst []byte) error {
	err := s.ReadAll(dst, Size(len(dst))*8)
	if err != nil {
		switch err.(type) {
		case InsufficientError:
			{
				return err
			}
		case OverflowError:
			{
				panic(err)
			}
		}
	}
	return nil
}

func (s BitSlice) ReadBool() (bool, error) {
	data, err := s.dataAtOffset(0, 1)
	if err != nil {
		return false, err
	}
	err = s.Advance(1)
	if err != nil {
		return false, err
	}
	return data&0x80 != 0, nil
}

func (s BitSlice) ReadByte() (byte, error) {
	data, err := s.byteAtOffset(0)
	if err != nil {
		return 0, err
	}
	err = s.Advance(8)
	if err != nil {
		return 0, err
	}
	return data, nil
}

func (s BitSlice) Remaining() Size {
	return Size(len(s.data)*8 - int(s.offset))
}

func (s BitSlice) Len() Size {
	return s.length
}
