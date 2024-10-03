package storage

import "studentgit.kata.academy/quant/torque/internal/domain"

func (s *Storage) GenerateULID() domain.ClientOrderID {
	return domain.ClientOrderID(s.ulidGenerator())
}
