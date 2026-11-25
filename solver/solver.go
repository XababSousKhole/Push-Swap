// Package solver calcule une liste d'instructions qui trie la pile a.
//
//   - jusqu'à 6 nombres : recherche en largeur (BFS) -> solution OPTIMALE
//     (voir bfs.go) ;
//   - au-delà : "algorithme Turc", qui à chaque étape déplace le nombre
//     qui coûte le moins d'instructions (voir turk.go).
package solver

import (
	"sort"

	"push-swap/stack"
)

// bfsLimit : jusqu'à combien de nombres on cherche la solution parfaite.
// Au-delà, le BFS serait beaucoup trop lent.
const bfsLimit = 6

// Solve est la fonction principale, appelée par le programme push-swap.
// Elle reçoit les nombres et renvoie les instructions à afficher.
func Solve(values []int) []string {
	// Déjà trié ? Rien à faire : 0 instruction.
	if stack.IsAscending(values) {
		return nil
	}

	// On remplace les valeurs par leurs rangs (voir Normalize).
	ranks := Normalize(values)

	// Petite liste -> solution parfaite. Grande liste -> algo Turc.
	if len(ranks) <= bfsLimit {
		return bfs(ranks)
	}
	return turk(ranks)
}

// Normalize remplace chaque valeur par son rang (sa place une fois trié).
//
//	[42 -7 100]  ->  [1 0 2]
//
// -7 est le plus petit (rang 0), 42 le deuxième (rang 1), 100 le plus grand (rang 2).
// L'ordre entre les nombres ne change pas, donc trier les rangs revient à
// trier les vraies valeurs. Mais c'est plus simple : on sait que les nombres
// vont de 0 à n-1, et que le plus petit est toujours 0.
func Normalize(values []int) []int {
	// 1. On fait une copie triée : [-7 42 100]
	sorted := make([]int, len(values))
	copy(sorted, values)
	sort.Ints(sorted)

	// 2. On note le rang de chaque valeur : -7 -> 0, 42 -> 1, 100 -> 2
	rank := make(map[int]int, len(values))
	for i, v := range sorted {
		rank[v] = i
	}

	// 3. On remplace chaque valeur par son rang, dans l'ordre d'origine.
	ranks := make([]int, len(values))
	for i, v := range values {
		ranks[i] = rank[v]
	}
	return ranks
}

// recorder sert à deux choses en même temps :
//   - exécuter les instructions sur les piles (pour savoir où on en est) ;
//   - les mémoriser dans ops (c'est ce que push-swap affichera à la fin).
type recorder struct {
	s   *stack.Stacks
	ops []string
}

// do exécute l'instruction op, times fois de suite.
// Exemple : r.do("ra", 3) fait ra, ra, ra.
func (r *recorder) do(op string, times int) {
	for i := 0; i < times; i++ {
		_ = r.s.Exec(op) // pas d'erreur possible : c'est nous qui choisissons op
		r.ops = append(r.ops, op)
	}
}
