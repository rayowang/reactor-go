package flux_test

import (
	"context"
	"log"
	"testing"
	"time"

	rs "github.com/jjeffcaii/reactor-go"
	"github.com/jjeffcaii/reactor-go/flux"
)

func TestUnicastProcessor(t *testing.T) {
	p := flux.NewUnicastProcessor()

	time.AfterFunc(100*time.Millisecond, func() {
		p.Next(1)
	})

	time.AfterFunc(150*time.Millisecond, func() {
		p.Next(2)
	})
	time.AfterFunc(200*time.Millisecond, func() {
		p.Next(2)
		p.Complete()
	})

	done := make(chan struct{})
	var su rs.Subscription
	p.DoOnNext(func(v interface{}) {
		log.Println("onNext:", v)
		su.Request(1)
	}).DoOnRequest(func(n int) {
		log.Println("request:", n)
	}).DoOnComplete(func() {
		log.Println("complete")
		close(done)
	}).Subscribe(context.Background(), rs.OnSubscribe(func(s rs.Subscription) {
		su = s
		su.Request(1)
	}))
	<-done
}

func TestDoOnSubscribe(t *testing.T) {
	pc := flux.NewUnicastProcessor()
	time.AfterFunc(1*time.Second, func() {
		pc.Next(111)
		pc.Next(222)
		pc.Complete()
	})
	done := make(chan struct{})
	pc.DoFinally(func(s rs.SignalType) {
		close(done)
	}).DoOnSubscribe(func(su rs.Subscription) {
		log.Println("doOnSubscribe")
	}).BlockLast(context.Background())

	<-done

}
