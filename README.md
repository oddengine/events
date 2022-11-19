# events

Golang event-driven module with reentrant mutex.

## Example

```go
type Target struct {
    events.EventTarget

    ID string
}

func (me *Target) Init(id string, logger log.ILogger) *Target {
    me.EventTarget.Init(logger)
    me.ID = id
    return me
}

func (me *Target) Close() error {
    me.DispatchEvent(Event.New(Event.CLOSE, me))
}
```

```go
type Observer struct {
    logger        log.ILogger
    closeListener *events.EventListener
}

func (me *Observer) Init(logger log.ILogger) *Observer {
    me.logger = logger
    me.closeListener = events.NewEventListener(me.onClose)
    return me
}

func (me *Observer) Attach(t *Target) {
    t.AddEventListener(Event.CLOSE, me.closeListener)
}

func (me *Observer) onClose(e *Event.Event) {
    t := e.Target().(*Target)
    me.logger.Infof("Target(%s) closed.", t.ID)
}
```

```go
var logger log.ILogger

func onClose(e *Event.Event) {
    t := e.Target().(*Target)
    logger.Infof("Target(%s) closed.", t.ID)
}

func main() {
    listener := events.NewEventListener(onClose)

    t := new(Target).Init("123", logger)
    t.AddEventListener(Event.CLOSE, listener)

    o := new(Observer).Init(logger)
    o.Attach(t)

    t.Close()
}
```
