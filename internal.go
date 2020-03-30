package bitbuf

// hake

type sliceContainer struct {
	data   []byte
	offset byte
	length Size
}

type transportContainer struct {
	buf      []byte
	length   Size
	capacity func([]byte) Size
}

type bitSlice struct {
	*sliceContainer
}

type bitSliceMut struct {
	*sliceContainer
}

type drain struct {
	*transportContainer
}

type fill struct {
	*transportContainer
}
