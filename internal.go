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
	capacity Size
}
