package bitbuf

type Size uint

type BitBuf interface {
	Advance(bits Size) error
	Read(dst []byte, bits Size) (Size, error)
	ReadAll(dst []byte, bits Size) error
	ReadAligned(dst []byte) Size
	ReadAlignedAll(out []byte) error
	ReadBool() (bool, error)
	ReadByte() (byte, error)
	Remaining() Size
	Len() Size
}

type BitBufMut interface {
	Advance(bits Size) error
	Write(data []byte, bits Size) (Size, error)
	WriteAll(data []byte, bits Size) error
	WriteAligned(data []byte) Size
	WriteAlignedAll(data []byte) error
	WriteBool(data bool) error
	WriteByte(data byte) error
	Remaining() Size
	Len() Size
}
