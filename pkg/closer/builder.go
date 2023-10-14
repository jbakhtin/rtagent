package closer

type builder struct {
	err    error
	closer Closer
}

func New() *builder {
	return &builder{
		closer: Closer{},
		err:    nil,
	}
}

func (b *builder) WithFuncs(funcs ...Func) *builder {
	b.closer.funcs = append(b.closer.funcs, funcs...)
	return b
}

func (b *builder) Build() (*Closer, error) {
	if b.err != nil {
		return nil, b.err
	}
	return &b.closer, nil
}
