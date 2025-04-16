package closer

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
)

var globalCloser *Closer
var once sync.Once

func Add(f ...func() error) {
	globalCloser.Add(f...)
}

func Wait() {
	globalCloser.Wait()
}

func CloseAll(ctx context.Context) {
	globalCloser.CloseAll(ctx)
}

type Closer struct {
	mu    sync.Mutex
	once  sync.Once
	done  chan struct{}
	funcs []func() error
}

func New(sig ...os.Signal) {
	once.Do(func() {
		c := &Closer{done: make(chan struct{})}
		if len(sig) > 0 {
			go func() {
				ch := make(chan os.Signal, 1)
				signal.Notify(ch, sig...)
				<-ch
				signal.Stop(ch)
				c.CloseAll(context.Background())
			}()
		}

		globalCloser = c
	})
}

func (c *Closer) Add(f ...func() error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.funcs = append(c.funcs, f...)
}

func (c *Closer) Wait() {
	<-c.done
}

func (c *Closer) CloseAll(ctx context.Context) {
	ctx = context.WithoutCancel(ctx)

	c.once.Do(func() {
		defer close(c.done)
		c.mu.Lock()
		funcs := c.funcs
		c.funcs = nil
		c.mu.Unlock()

		// call all closer funcs async
		errs := make(chan error, len(funcs))
		for _, f := range funcs {
			go func(f func() error) {
				errs <- f()
			}(f)
		}

		for i := 0; i < len(funcs); i++ {
			if err := <-errs; err != nil {
				fmt.Println("ошибка при закрытии:", err)
			}
		}
	})
}
