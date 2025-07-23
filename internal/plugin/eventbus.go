package plugin

import (
	"context"
	"sync"
)

// EventBus handles plugin event communication
type EventBus struct {
	subscribers map[string][]EventHandler
	mu          sync.RWMutex
}

// NewEventBus creates a new event bus
func NewEventBus() *EventBus {
	return &EventBus{
		subscribers: make(map[string][]EventHandler),
	}
}

// Initialize initializes the event bus
func (eb *EventBus) Initialize() error {
	return nil
}

// Shutdown shuts down the event bus
func (eb *EventBus) Shutdown() error {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	
	eb.subscribers = make(map[string][]EventHandler)
	return nil
}

// Subscribe adds an event handler for a specific event type
func (eb *EventBus) Subscribe(eventType string, handler EventHandler) {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	
	eb.subscribers[eventType] = append(eb.subscribers[eventType], handler)
}

// Unsubscribe removes an event handler
func (eb *EventBus) Unsubscribe(eventType string, plugin interface{}) {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	
	if handlers, exists := eb.subscribers[eventType]; exists {
		// Remove handlers for this plugin (simplified approach)
		// In a real implementation, you'd need a more sophisticated way to match handlers
		eb.subscribers[eventType] = handlers[:0]
	}
}

// Emit sends an event to all subscribers
func (eb *EventBus) Emit(ctx context.Context, event *PluginEvent) error {
	eb.mu.RLock()
	handlers := eb.subscribers[event.Type]
	eb.mu.RUnlock()
	
	for _, handler := range handlers {
		go func(h EventHandler) {
			_ = h(ctx, event)
		}(handler)
	}
	
	return nil
}