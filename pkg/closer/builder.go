package closer

type builder struct {
	closer closer
	err error
}

func New() (*builder) {
	return &builder{
		closer: closer{},
		err: nil,
	}
}

func (b *builder) WithFuncs(funcs ...Func) *builder {
	b.closer.funcs = append(b.closer.funcs, funcs...)
	return b
}

func (b *builder) Build() (*closer, error) {
	if b.err != nil {
		return nil, b.err
	}
	return &b.closer, nil
}