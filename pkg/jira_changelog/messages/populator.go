package messages

import "context"

type Messager interface {
	Message() string
}

type Populator interface {
	Populate(ctx context.Context) ([]Messager, error)
}

type NoopPopulator struct{}

func (e *NoopPopulator) Populate(ctx context.Context) ([]Messager, error) {
	return []Messager{}, nil
}
