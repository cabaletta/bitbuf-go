package bitbuf

func NewDrain(buf []byte) Drain {
	return drain{&transportContainer{
		buf:    buf,
		length: 0,
		capacity: func(data []byte) Size {
			return Size(len(data)) * 8
		},
	}}
}

func NewCappedDrain(buf []byte, cap Size) Drain {
	return drain{&transportContainer{
		buf:    buf,
		length: 0,
		capacity: func(data []byte) Size {
			return cap
		},
	}}
}

func (d drain) IntoInner() []byte {
	return d.buf
}

func (d drain) DrainInto(to BitBufMut) error {
	capacity := d.capacity(d.buf)
	from := NewBitSlice(d.buf)
	err := from.Advance(d.length)
	if err != nil {
		panic(err)
	}

	for true {
		if d.length < capacity {
			if to.Remaining() >= 8 && capacity-d.length >= 8 {
				val, err := from.ReadByte()
				if err != nil {
					panic(err)
				}
				err = to.WriteByte(val)
				if err != nil {
					panic(err)
				}
				d.length += 8
			} else {
				val, err := from.ReadBool()
				if err != nil {
					panic(err)
				}
				err = to.WriteBool(val)
				if err != nil {
					return InsufficientError{}
				}
				d.length++
			}
		} else {
			return nil
		}
	}
	return nil
}

func (d drain) AsBuf() BitBuf {
	return NewBitSlice(d.buf)
}
