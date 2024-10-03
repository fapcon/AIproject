package storage

import (
	"context"

	"studentgit.kata.academy/quant/torque/internal/domain"
)

func (b *Builder) WithOrderBookRead() *Builder {
	b.buildCap(CapabilityOrderBookRead)
	return b
}

func (b *Builder) WithOrderBookWrite() *Builder {
	b.buildCap(CapabilityOrderBookWrite)
	return b
}

func (b *Builder) WithOrderBookReadWrite() *Builder {
	return b.WithOrderBookRead().WithOrderBookWrite()
}

func (s *Storage) GetOrderBook(_ context.Context, i domain.LocalInstrument) (domain.OrderBook, bool) {
	s.mustHave(CapabilityOrderBookRead)
	return s.cache.OrderBook.Get(i)
}

func (s *Storage) SaveOrderBook(_ context.Context, b domain.OrderBook) {
	s.mustHave(CapabilityOrderBookWrite)
	s.cache.OrderBook.Add(b)
}

func (s *Storage) DeleteOrderBook(_ context.Context, i domain.LocalInstrument) {
	s.mustHave(CapabilityOrderBookWrite)
	s.cache.OrderBook.Delete(i)
}

func (s *Storage) GetOrderBookList(_ context.Context) []domain.OrderBook {
	s.mustHave(CapabilityOrderBookRead)
	return s.cache.OrderBook.GetAll()
}
