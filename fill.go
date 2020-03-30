package bitbuf

type Fill struct {
	*transportContainer
}

func NewFill(buf []byte) Fill {
	return Fill{&transportContainer{
		buf:    buf,
		length: 0,
		capacity: func(data []byte) Size {
			return Size(len(data)) * 8
		},
	}}
}

func NewCappedFill(buf []byte, cap Size) Fill {
	return Fill{&transportContainer{
		buf:    buf,
		length: 0,
		capacity: func(data []byte) Size {
			return cap
		},
	}}
}

func (f Fill) IntoInner() []byte {
	return f.buf
}

func (f Fill) FillFrom(from BitBuf) error {
	capacity := f.capacity(f.buf)
	to := NewBitSliceMut(f.buf)
	err := to.Advance(f.length)
	if err != nil {
		panic(err)
	}

	for true {
		if f.length < capacity {
			if from.Remaining() >= 8 && capacity-f.length >= 8 {
				val, err := from.ReadByte()
				if err != nil {
					panic(err)
				}
				err = to.WriteByte(val)
				if err != nil {
					panic(err)
				}
				f.length += 8
			} else {
				val, err := from.ReadBool()
				if err != nil {
					panic(err)
				}
				err = to.WriteBool(val)
				if err != nil {
					return InsufficientError{}
				}
				f.length++
			}
		} else {
			return nil
		}
	}
	return nil
}

func (f Fill) AsBuf() BitBuf {
	return NewBitSlice(f.buf)
}
