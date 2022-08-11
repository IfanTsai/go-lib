package promise

import "sync"

type Handler func(interface{}) interface{}

type callback struct {
	onFulfilled Handler
	onRejected  Handler
	resolve     Handler
	reject      Handler
}

type Thenable interface {
	Then(...Handler) Thenable
	Catch(Handler) Thenable
	Finally(func())
}

type Promise struct {
	callbacks []*callback
	state     TypeState
	value     interface{}
	mutex     sync.RWMutex
}

func New(fn func(resolve, reject Handler)) *Promise {
	p := &Promise{
		state: TypeStatePending,
	}

	fn(p.resolve, p.reject)

	return p
}

func (p *Promise) Then(handlers ...Handler) Thenable {
	var onFulfilled, onRejected Handler

	switch {
	case len(handlers) == 1:
		onFulfilled = handlers[0]
	case len(handlers) >= 2:
		onFulfilled, onRejected = handlers[0], handlers[1]
	}

	return New(func(resolve, reject Handler) {
		p.handle(&callback{
			onFulfilled: onFulfilled,
			onRejected:  onRejected,
			resolve:     resolve,
			reject:      reject,
		})
	})
}

func (p *Promise) Catch(onError Handler) Thenable {
	return p.Then(nil, onError)
}

func (p *Promise) Finally(onFinally func()) {
	if onFinally == nil {
		p.Then()

		return
	}

	doneHandler := func(interface{}) interface{} {
		onFinally()

		return nil
	}

	p.Then(doneHandler, doneHandler)
}

func (p *Promise) handle(cb *callback) {
	switch p.state {
	case TypeStatePending:
		p.mutex.Lock()
		defer p.mutex.Unlock()

		p.callbacks = append(p.callbacks, cb)
	case TypeStateFulfilled:
		if cb.onFulfilled != nil {
			cb.resolve(cb.onFulfilled(p.value))
		} else {
			cb.resolve(p.value)
		}
	case TypeStateRejected:
		if cb.onRejected != nil {
			cb.reject(cb.onRejected(p.value))
		} else {
			cb.reject(p.value)
		}
	}
}

func (p *Promise) resolve(value interface{}) interface{} {
	if newPromise, ok := value.(Thenable); ok {
		// p.resolve is promise bridge resolve
		// p.resolve will be called when newPromise is settled, then promise bridge's Then handler will be called
		newPromise.Then(p.resolve, p.reject)

		return nil
	}

	p.settle(TypeStateFulfilled, value)

	return nil
}

func (p *Promise) reject(value interface{}) interface{} {
	p.settle(TypeStateRejected, value)

	return nil
}

func (p *Promise) settle(state TypeState, value interface{}) {
	p.state = state
	p.value = value

	p.mutex.RLock()
	defer p.mutex.RUnlock()

	for _, cb := range p.callbacks {
		p.handle(cb)
	}
}
