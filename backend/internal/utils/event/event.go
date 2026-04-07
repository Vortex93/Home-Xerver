package event

import "sync"

// ========================================
// Types
// ========================================
type Handler func(payload any)

type listener struct {
    id   uint64
    fn   Handler
    once bool
}

type emitter struct {
    mu        sync.Mutex
    nextID    uint64
    listeners map[string][]listener
    closed    bool
}

// ========================================
// State
// ========================================
var global = emitter{
    listeners: make(map[string][]listener),
}

// ========================================
// Functions
// ========================================

// On registers a handler for a specific event and returns a unique listener ID and an unsubscription function.
func On(event string, handler Handler) (uint64, func()) {
    global.mu.Lock()
    defer global.mu.Unlock()

    if global.closed {
        return 0, func() {}
    }

    global.nextID++

    listenerID := global.nextID

    global.listeners[event] = append(global.listeners[event], listener{
        id:   listenerID,
        fn:   handler,
        once: false,
    })

    return listenerID, func() {
        Off(event, listenerID)
    }
}

// Once registers a handler for a specific event that will be called at most once. It returns a unique listener ID and an unsubscription function.
func Once(event string, handler Handler) (uint64, func()) {
    global.mu.Lock()
    defer global.mu.Unlock()

    if global.closed {
        return 0, func() {}
    }

    global.nextID++

    listenerID := global.nextID

    global.listeners[event] = append(global.listeners[event], listener{
        id:   listenerID,
        fn:   handler,
        once: true,
    })

    return listenerID, func() {
        Off(event, listenerID)
    }
}

// Off unregisters a handler for a specific event using its unique listener ID.
func Off(event string, listenerID uint64) {
    global.mu.Lock()
    defer global.mu.Unlock()

    listeners, exists := global.listeners[event]
    if !exists {
        return
    }

    filtered := listeners[:0]

    for _, currentListener := range listeners {
        if currentListener.id != listenerID {
            filtered = append(filtered, currentListener)
        }
    }

    if len(filtered) == 0 {
        delete(global.listeners, event)
        return
    }

    global.listeners[event] = filtered
}

// Emit triggers all handlers registered for a specific event, passing the provided payload to each handler. Handlers registered with Once will be automatically unregistered after being called.
func Emit(event string, payload any) {
    global.mu.Lock()

    if global.closed {
        global.mu.Unlock()
        return
    }

    listeners, exists := global.listeners[event]
    if !exists || len(listeners) == 0 {
        global.mu.Unlock()
        return
    }

    handlers := make([]Handler, 0, len(listeners))
    persistentListeners := make([]listener, 0, len(listeners))

    for _, currentListener := range listeners {
        handlers = append(handlers, currentListener.fn)

        if !currentListener.once {
            persistentListeners = append(persistentListeners, currentListener)
        }
    }

    if len(persistentListeners) == 0 {
        delete(global.listeners, event)
    } else {
        global.listeners[event] = persistentListeners
    }

    global.mu.Unlock()

    for _, handler := range handlers {
        handler(payload)
    }
}

// Close disables the emitter and clears all registered listeners. After calling Close, no new listeners can be registered and no events will be emitted.
func Close() {
    global.mu.Lock()
    defer global.mu.Unlock()

    global.closed = true
    global.listeners = make(map[string][]listener)
}

// Reset clears all registered listeners and resets the emitter to its initial state. This allows new listeners to be registered and events to be emitted again after a Close.
func Reset() {
    global.mu.Lock()
    defer global.mu.Unlock()

    global.nextID = 0
    global.closed = false
    global.listeners = make(map[string][]listener)
}
