package events

import (
	"reflect"
)

const (
	REMOVE_ALL_LISTENERS = "removeAllListeners"
	NEW_LISTENER         = "newListeners"
)

var (
	defaultMaxListeners = 10
)

type Listener func(args ...interface{})
type Listeners []Listener

type EventEmitter struct {
	events       map[string]Listeners
	eventsCount  int
	MaxListeners int
}

func NewEventEmitter() (e *EventEmitter) {
	e = new(EventEmitter)
	e.events = make(map[string]Listeners)
	e.eventsCount = 0
	e.MaxListeners = defaultMaxListeners
	return e
}

func (this *EventEmitter) Emit(eventName string, args ...interface{}) (ok bool) {
	events := this.events
	if nil == events {
		return false
	}
	listeners := events[eventName]
	if nil == listeners || len(listeners) <= 0 {
		return false
	}
	for _, handle := range listeners {
		if nil != handle {
			handle(args...)
		}
	}
	return true
}

func addListener(target *EventEmitter, eventName string, listener Listener, prepend bool) (ok bool) {
	if "" == eventName || nil == listener {
		return false
	}
	if nil == target.events {
		target.events = make(map[string]Listeners)
		target.eventsCount = 0
	}
	events := target.events
	if _, ok := events[NEW_LISTENER]; ok {
		target.Emit(NEW_LISTENER, eventName, listener)
		events = target.events
	}
	existing, ok := events[eventName]
	if ok {
		if len(existing) >= target.MaxListeners {
			return false
		}
		if prepend {
			existing = append(Listeners{listener}, existing...)
		} else {
			existing = append(existing, listener)
		}
		events[eventName] = existing
	} else {
		existing = append(make(Listeners, 0), listener)
		target.events[eventName] = existing
		target.eventsCount++
	}

	return true
}

func (this *EventEmitter) AddListener(eventName string, listener Listener) (ok bool) {
	return addListener(this, eventName, listener, false)
}

func (this *EventEmitter) On(eventName string, listener Listener) (ok bool) {
	return this.AddListener(eventName, listener)
}

func (this *EventEmitter) PrependListener(eventName string, listener Listener) (ok bool) {
	return addListener(this, eventName, listener, true)
}

type onceWapper struct {
	fired     bool
	target    *EventEmitter
	eventName string
	listener  Listener
	wrapperFn Listener
}

func (ow onceWapper) onceWrapper(args ...interface{}) {
	if !ow.fired {
		ow.target.RemoveListener(ow.eventName, ow.wrapperFn)
		ow.fired = true
		ow.listener(args...)
	}
}

func onceWrap(target *EventEmitter, eventName string, listener Listener) func(args ...interface{}) {
	var ow onceWapper
	wrapperFn := func(args ...interface{}) {
		ow.onceWrapper(args...)
	}
	ow.fired = false
	ow.target = target
	ow.eventName = eventName
	ow.listener = listener
	ow.wrapperFn = wrapperFn
	return wrapperFn
}

func (this *EventEmitter) Once(eventName string, listener Listener) {
	this.On(eventName, onceWrap(this, eventName, listener))
}

func (this *EventEmitter) PrependOnceListener(eventName string, listener Listener) {
	this.PrependListener(eventName, onceWrap(this, eventName, listener))
}

func (this *EventEmitter) RemoveListener(eventName string, listener Listener) *EventEmitter {
	events := this.events
	if nil == events {
		return this
	}
	list := events[eventName]
	if nil == list || len(list) <= 0 {
		return this
	}
	position := -1
	var originListener Listener
	for i := range list {
		if reflect.ValueOf(listener) == reflect.ValueOf(list[i]) {
			position = i
			originListener = list[i]
			break
		}
	}
	if position < 0 {
		return this
	}
	list = append(list[:position], list[position+1:]...)
	events[eventName] = list
	if _, ok := events[REMOVE_ALL_LISTENERS]; ok {
		this.Emit(REMOVE_ALL_LISTENERS, eventName, originListener)
	}
	return this
}

func (this *EventEmitter) RemoveAllListener(eventName string) *EventEmitter {
	events := this.events
	if nil == events {
		return this
	}
	_, ok := events[REMOVE_ALL_LISTENERS]
	// not listening for removeListener, no need to emit
	if !ok {
		if "" == eventName {
			this.events = make(map[string]Listeners)
			this.eventsCount = 0
		} else if _, ok := events[eventName]; ok {
			if this.eventsCount > 0 {
				this.eventsCount--
				delete(this.events, eventName)
			}
		}
		return this
	}
	// emit removeListener for all listeners on all events
	if "" == eventName {
		for key := range events {
			if key == REMOVE_ALL_LISTENERS {
				continue
			}
			this.RemoveAllListener(key)
		}
		this.RemoveAllListener(REMOVE_ALL_LISTENERS)
	} else {
		if listeners, ok := events[eventName]; ok {
			for _, listener := range listeners {
				this.RemoveListener(eventName, listener)
			}
		}
	}
	return this
}

func (this *EventEmitter) Listeners(eventName string) (listener Listeners) {
	events := this.events
	if nil == events {
		return nil
	}
	listener = events[eventName]
	return listener
}

func (this *EventEmitter) ListenerCount(eventName string) (count int) {
	count = 0
	events := this.events
	if nil == events {
		return count
	}
	if "" != eventName {
		for _, listeners := range events {
			count += len(listeners)
		}
	} else {
		if listeners, ok := events[eventName]; ok {
			count = len(listeners)
		}
	}
	return count
}

func (this *EventEmitter) EventNames() (names []string) {
	if nil == this.events {
		return nil
	}
	length := len(this.events)
	names = make([]string, length)
	for key := range this.events {
		names = append(names, key)
	}
	return names
}
