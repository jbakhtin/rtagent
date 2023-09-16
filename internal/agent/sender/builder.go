package sender

type Builder struct {
	sender sender
	err error
}

func New() *Builder {
	return &Builder{}
}

func (b *Builder) WithConfig(cfg Configer) *Builder {
	b.sender.cfg = cfg
	return b
}

func (b *Builder) Build() (*sender, error) {
	if b.err != nil {
		return nil, b.err
	}

	return &b.sender, b.err
}
