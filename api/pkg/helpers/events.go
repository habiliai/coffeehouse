package helpers

import "context"

type EventListener struct {
	OnEndTransactions  []func(ctx context.Context)
	OnMessageCompleted []func(ctx context.Context)
}

type EventType int

const (
	EventTypeEndTransaction EventType = iota
	EventTypeCompletedAction
)

const (
	contextKeyEventListener contextKey = "_eventListener"
)

func (e *EventListener) Emit(ctx context.Context, eventType EventType) {
	switch eventType {
	case EventTypeEndTransaction:
		for _, fn := range e.OnEndTransactions {
			fn(ctx)
		}
	case EventTypeCompletedAction:
		for _, fn := range e.OnMessageCompleted {
			fn(ctx)
		}
	}
}

func WithEventListener(ctx context.Context, eventListener *EventListener) context.Context {
	return context.WithValue(ctx, contextKeyEventListener, eventListener)
}

func getEventListener(ctx context.Context) *EventListener {
	value, ok := ctx.Value(contextKeyEventListener).(*EventListener)
	if !ok {
		return nil
	}

	return value
}

func On(ctx context.Context, eventType EventType, fn func(ctx context.Context)) {
	eventListener := getEventListener(ctx)
	if eventListener != nil {
		switch eventType {
		case EventTypeEndTransaction:
			eventListener.OnEndTransactions = append(eventListener.OnEndTransactions, fn)
		case EventTypeCompletedAction:
			eventListener.OnMessageCompleted = append(eventListener.OnMessageCompleted, fn)
		}
	}
}

func NewEventListener() *EventListener {
	return &EventListener{}
}
