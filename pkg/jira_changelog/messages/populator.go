package messages

import "context"

type Message interface {
	Message() string
}

type Populator interface {
	Populate(ctx context.Context) ([]Message, error)
}
