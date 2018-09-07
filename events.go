package events

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

func (this *EventEmitter) Emit(eventName string, args ...interface{}) {

}

func (this *EventEmitter) AddListener(eventName string, listen Listener) {

}

func (this *EventEmitter) On(eventName string, listen Listener) {
	this.AddListener(eventName, listen)
}

func (this *EventEmitter) PrependListener(eventName string, listener Listener) {

}

func (this *EventEmitter) Once(eventName string, listener Listener) {

}

func (this *EventEmitter) PrependOnceListener(eventName string, listener Listener) {

}

func (this *EventEmitter) RemoveListener(eventName string, listener Listener) {

}

func (this *EventEmitter) RemoveAllListener(eventName string) {

}

func (this *EventEmitter) Listeners(eventName string) {

}

func (this *EventEmitter) ListenerCount(eventName string) {

}

func (this *EventEmitter) EventNames() {

}
