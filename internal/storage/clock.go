package storage

import "time"

func (s *Storage) Now() time.Time {
	return s.clock.Now()
}

func (s *Storage) After(fireAfter time.Duration) <-chan time.Time {
	return s.clock.After(fireAfter)
}

func (s *Storage) NewTimer(fireAfter time.Duration) (<-chan time.Time, func() bool) {
	return s.clock.NewTimer(fireAfter)
}

func (s *Storage) Tick(fireEach time.Duration) <-chan time.Time {
	return s.clock.Tick(fireEach)
}

func (s *Storage) NewTicker(fireEach time.Duration) (<-chan time.Time, func()) {
	return s.clock.NewTicker(fireEach)
}
