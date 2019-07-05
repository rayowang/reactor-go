package mono

import (
	"context"

	rs "github.com/jjeffcaii/reactor-go"
)

type filterSubscriber struct {
	s rs.Subscriber
	f rs.Predicate
}

func (f filterSubscriber) OnComplete() {
	f.s.OnComplete()
}

func (f filterSubscriber) OnError(err error) {
	f.s.OnError(err)
}

func (f filterSubscriber) OnNext(s rs.Subscription, v interface{}) {
	if f.f(v) {
		f.s.OnNext(s, v)
	}
}

func (f filterSubscriber) OnSubscribe(s rs.Subscription) {
	f.s.OnSubscribe(s)
}

type monoFilter struct {
	*baseMono
	s Mono
	f rs.Predicate
}

func (m *monoFilter) Subscribe(ctx context.Context, options ...rs.SubscriberOption) {
	m.SubscribeWith(ctx, rs.NewSubscriber(options...))
}

func (m *monoFilter) SubscribeWith(ctx context.Context, s rs.Subscriber) {
	m.s.SubscribeWith(ctx, filterSubscriber{
		s: s,
		f: m.f,
	})
}

func newMonoFilter(s Mono, f rs.Predicate) Mono {
	m := &monoFilter{
		s: s,
		f: f,
	}
	m.baseMono = &baseMono{
		child: m,
	}
	return m
}