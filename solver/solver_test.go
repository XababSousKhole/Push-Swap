package solver

import (
	"math/rand"
	"testing"

	"push-swap/stack"
)

// sorts vérifie que les instructions trient bien la liste (comme le checker).
func sorts(values []int, ops []string) bool {
	s := stack.New(values)
	for _, op := range ops {
		if s.Exec(op) != nil {
			return false
		}
	}
	return s.IsSorted()
}

func TestAlreadySorted(t *testing.T) {
	if ops := Solve([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}); len(ops) != 0 {
		t.Fatalf("liste déjà triée : %d instructions, want 0", len(ops))
	}
}

func TestSubjectExamples(t *testing.T) {
	cases := []struct {
		values []int
		max    int
	}{
		{[]int{2, 1, 3, 6, 5, 8}, 7},
		{[]int{4, 67, 3, 87, 23}, 6},
		{[]int{3, 2, 1, 0}, 6},
	}
	for _, c := range cases {
		ops := Solve(c.values)
		if !sorts(c.values, ops) {
			t.Errorf("%v n'est pas trié", c.values)
		}
		if len(ops) > c.max {
			t.Errorf("%v : %d instructions, want <= %d", c.values, len(ops), c.max)
		}
	}
}

// Toutes les permutations de 1 à 6 éléments.
func TestAllSmallPermutations(t *testing.T) {
	for n := 1; n <= 6; n++ {
		permute(seq(n), 0, func(p []int) {
			if ops := Solve(p); !sorts(p, ops) {
				t.Fatalf("%v non trié par %v", p, ops)
			}
		})
	}
}

func TestRandomSizes(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	for _, n := range []int{7, 8, 10, 20, 50, 100, 500} {
		for k := 0; k < 20; k++ {
			values := randomValues(rng, n)
			if ops := Solve(values); !sorts(values, ops) {
				t.Fatalf("n=%d : liste non triée", n)
			}
		}
	}
}

// Affiche le nombre moyen / max d'instructions (go test -v -run Score).
func TestScore(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	for _, n := range []int{3, 5, 100, 500} {
		total, worst := 0, 0
		const runs = 50
		for k := 0; k < runs; k++ {
			c := len(Solve(randomValues(rng, n)))
			total += c
			worst = max(worst, c)
		}
		t.Logf("n=%3d : moyenne %d, pire %d", n, total/runs, worst)
	}
}

func seq(n int) []int {
	s := make([]int, n)
	for i := range s {
		s[i] = i
	}
	return s
}

func permute(p []int, k int, f func([]int)) {
	if k == len(p) {
		cp := append([]int(nil), p...)
		f(cp)
		return
	}
	for i := k; i < len(p); i++ {
		p[k], p[i] = p[i], p[k]
		permute(p, k+1, f)
		p[k], p[i] = p[i], p[k]
	}
}

func randomValues(rng *rand.Rand, n int) []int {
	perm := rng.Perm(n * 10)[:n]
	for i := range perm {
		perm[i] -= n * 5 // inclut des négatifs
	}
	return perm
}
