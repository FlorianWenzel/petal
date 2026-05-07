package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand/v2"

	"github.com/jackc/pgx/v5"
)

var flowerAdjectives = []string{
	"dewy", "soft", "gentle", "quiet", "sunlit", "mossy", "wild", "sweet",
	"hushed", "golden", "velvet", "misty", "twilit", "dappled", "tender",
	"rosy", "amber", "pale", "sleepy", "warm", "shady", "still", "breezy",
	"fragrant", "cozy", "lush", "silken", "starlit", "creamy", "meadow",
}

var flowers = []string{
	"rose", "lily", "peony", "daisy", "iris", "tulip", "dahlia", "poppy",
	"aster", "lotus", "jasmine", "violet", "orchid", "marigold", "zinnia",
	"magnolia", "camellia", "hyacinth", "primrose", "sunflower", "bluebell",
	"foxglove", "lilac", "clover", "fern", "sage", "thistle", "heather",
	"crocus", "freesia",
}

func generateUsername() string {
	a := flowerAdjectives[rand.IntN(len(flowerAdjectives))]
	f := flowers[rand.IntN(len(flowers))]
	n := 1000 + rand.IntN(9000)
	return fmt.Sprintf("%s-%s-%d", a, f, n)
}

func (s *server) allocateUsername(ctx context.Context) (string, error) {
	for i := 0; i < 25; i++ {
		candidate := generateUsername()
		var id string
		err := s.pool.QueryRow(ctx, `select id from users where username = $1`, candidate).Scan(&id)
		if errors.Is(err, pgx.ErrNoRows) {
			return candidate, nil
		}
		if err != nil {
			return "", err
		}
	}
	return "", errors.New("failed to allocate unique username")
}
