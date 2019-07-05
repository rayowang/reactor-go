package mono_test

import (
	"context"
	"log"
	"testing"
	"time"

	rs "github.com/jjeffcaii/reactor-go"
	"github.com/jjeffcaii/reactor-go/mono"
	"github.com/jjeffcaii/reactor-go/scheduler"
	"github.com/stretchr/testify/assert"
)

func TestJust(t *testing.T) {
	var complete bool
	now := time.Now()
	mono.
		Just(now).
		Map(func(i interface{}) interface{} {
			return i.(time.Time).UnixNano()
		}).
		Map(func(i interface{}) interface{} {
			return i.(int64) * 2
		}).
		Subscribe(context.Background(),
			rs.OnNext(func(s rs.Subscription, v interface{}) {
				log.Println("next:", v)
				assert.Equal(t, now.UnixNano()*2, v, "bad result")
			}),
			rs.OnComplete(func() {
				log.Println("complete")
				complete = true
			}),
		)
	assert.True(t, complete, "not complete")
}

func TestMonoJust_FlatMap(t *testing.T) {
	v, err := mono.Just(1).
		FlatMap(func(i interface{}) mono.Mono {
			return mono.Just(i).
				Map(func(i interface{}) interface{} {
					time.Sleep(200 * time.Millisecond)
					return i.(int) * 2
				}).
				SubscribeOn(scheduler.Elastic())
		}).
		Block(context.Background())
	assert.NoError(t, err, "an error occurred")
	assert.Equal(t, 2, v, "bad result")
}