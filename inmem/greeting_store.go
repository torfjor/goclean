package inmem

import (
	"context"
	"math/rand"
)

type Store struct {
	Templates []string
}

func (s *Store) RandomGreetingTemplate(ctx context.Context) (string, error) {
	return s.Templates[rand.Intn(len(s.Templates))], nil
}
