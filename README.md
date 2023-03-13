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

func (me *Observer) Watch(t *Target) {
    t.AddEventListener(Event.CLOSE, me.closeListener)
}

func (me *Observer) onClose(e *Event.Event) {
    t := e.Target().(*Target)
    me.logger.Infof("Target(%s) closed.", t.ID)
}
```

```go
func main() {
    // Learn how to create a logger here:
    //   https://github.com/oddengine/log
    var logger log.ILogger

    t := new(Target).Init("123", logger)

    o := new(Observer).Init(logger)
    o.Watch(t)

    t.Close()
}
```
