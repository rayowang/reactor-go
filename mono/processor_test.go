package mono_test

import (
	"context"
	"testing"
	"time"

	"github.com/jjeffcaii/reactor-go"
	"github.com/jjeffcaii/reactor-go/mono"
	"github.com/stretchr/testify/assert"
)

func TestProcessor(t *testing.T) {
	p := mono.CreateProcessor()

	time.AfterFunc(100*time.Millisecond, func() {
		p.Success(333)
	})

	v, err := p.
		Map(func(i interface{}) interface{} {
			return i.(int) * 2
		}).
		Block(context.Background())
	assert.NoError(t, err, "block failed")
	assert.Equal(t, 666, v.(int), "bad result")

	var actual int
	p.
		DoOnNext(func(v interface{}) {
			actual = v.(int)
		}).
		Subscribe(context.Background())
	assert.Equal(t, 333, actual, "bad result")
}

func TestProcessor_Context(t *testing.T) {
	p := mono.CreateProcessor()
	ctx, cancel := context.WithTimeout(context.Background(), 22*time.Millisecond)
	defer cancel()
	time.AfterFunc(33*time.Millisecond, func() {
		p.Success(77778888)
	})
	done := make(chan struct{})
	p.
		DoOnError(func(e error) {
			assert.Equal(t, rs.ErrSubscribeCancelled, e, "bad error")
		}).
		DoFinally(func(signal rs.SignalType) {
			close(done)
			assert.Equal(t, rs.SignalTypeError, signal)
		}).
		Subscribe(ctx)
	<-done
}
