package taskmanager

type builder struct {
	err    error
	closer taskmanager
}

func New() *builder {
	return &builder{
		closer: taskmanager{},
		err:    nil,
	}
}

func (b *builder) WithFuncs(funcs ...Func) *builder {
	b.closer.funcs = append(b.closer.funcs, funcs...)
	return b
}

func (b *builder) Build() (*taskmanager, error) {
	if b.err != nil {
		return nil, b.err
	}
	return &b.closer, nil
}
