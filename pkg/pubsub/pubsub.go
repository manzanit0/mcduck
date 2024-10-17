package pubsub

import (
	"context"
	"errors"
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

const DefaultStreamName = "mcduck"

func MarshalProto(event proto.Message) ([]byte, string, error) {
	eventType := string(proto.MessageName(event))

	eventBytes, err := proto.Marshal(event)
	if err != nil {
		return nil, eventType, fmt.Errorf("failed to marshal event: %w", err)
	}

	return eventBytes, eventType, nil
}

func UnmarshalProto[T proto.Message](raw []byte, event T) (T, error) {
	eventType := proto.MessageName(event)

	mt, err := protoregistry.GlobalTypes.FindMessageByName(protoreflect.FullName(eventType))
	if err != nil {
		return event, fmt.Errorf("failed to get proto from registry: %w", err)
	}

	msg := mt.New().Interface()
	err = proto.Unmarshal(raw, msg)
	if err != nil {
		return event, fmt.Errorf("failed to unmarshal message: %w", err)
	}

	typed, ok := msg.(T)
	if !ok {
		return event, fmt.Errorf("type mismatch: payload and event type don't match")
	}

	return typed, nil
}

func NewStream(ctx context.Context, url, name, subject string) (jetstream.JetStream, jetstream.Stream, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, nil, fmt.Errorf("connect to NATS: %v", err)
	}

	js, err := jetstream.New(nc)
	if err != nil {
		return nil, nil, fmt.Errorf("create new stream: %v", err)
	}

	stream, err := js.CreateStream(ctx, jetstream.StreamConfig{
		Name:     name,
		Subjects: []string{subject},
	})
	if err != nil {
		if !errors.Is(err, jetstream.ErrStreamNameAlreadyInUse) {
			return nil, nil, err
		}
		var jsErr error
		_, jsErr = js.Stream(ctx, name)
		if jsErr != nil {
			return nil, nil, err
		}
	}

	return js, stream, nil
}
