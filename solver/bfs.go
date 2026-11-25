package solver

// ============================================================================
// BFS = Breadth-First Search = PARCOURS EN LARGEUR
//
// Idée : on imagine un arbre de toutes les possibilités.
//
//	niveau 0 : la liste de départ
//	niveau 1 : les 11 listes qu'on obtient avec 1 instruction
//	niveau 2 : les listes qu'on obtient avec 2 instructions
//	...
//
// On explore niveau par niveau. Dès qu'on tombe sur une liste triée,
// on s'arrête : comme tous les niveaux d'avant ont déjà été vérifiés,
// c'est FORCÉMENT la solution la plus courte possible.
//
// Limite : le nombre de cas explose. 6 nombres = 5 040 états possibles,
// 10 nombres = environ 40 millions. Donc on ne l'utilise que jusqu'à 6.
// ============================================================================

import "push-swap/stack"

// allOps : les 11 instructions qu'on essaie à chaque étape.
var allOps = []string{"sa", "sb", "ss", "pa", "pb", "ra", "rb", "rr", "rra", "rrb", "rrr"}

// node est une "case" de notre exploration : un état des deux piles.
type node struct {
	s      *stack.Stacks // l'état des piles a et b à ce moment-là
	parent int           // d'où on vient (index dans la file), -1 = départ
	op     string        // l'instruction qu'on a faite pour arriver ici
}

func bfs(ranks []int) []string {
	// La file d'attente : les états à examiner, dans l'ordre.
	// On commence avec un seul état : la liste de départ.
	queue := []node{{s: stack.New(ranks), parent: -1}}

	// visited retient les états déjà vus, pour ne pas les refaire
	// (ex : "sa" puis "sa" revient au point de départ, inutile de continuer).
	visited := map[string]bool{key(queue[0].s): true}

	// On parcourt la file. Elle grandit pendant qu'on avance :
	// c'est ce qui fait qu'on explore niveau par niveau.
	for i := 0; i < len(queue); i++ {
		cur := queue[i]

		// Trouvé ! On reconstruit le chemin et on s'arrête.
		if cur.s.IsSorted() {
			return path(queue, i)
		}

		// Sinon, on essaie les 11 instructions depuis cet état.
		for _, op := range allOps {
			next := cur.s.Clone() // copie pour ne pas abîmer l'état actuel
			_ = next.Exec(op)

			k := key(next)
			if visited[k] {
				continue // déjà vu : on ignore
			}
			visited[k] = true

			// Nouvel état : on l'ajoute en fin de file, en notant d'où il vient.
			queue = append(queue, node{s: next, parent: i, op: op})
		}
	}
	return nil // n'arrive jamais : un état trié est toujours atteignable
}

// key transforme les deux piles en une chaîne de caractères unique.
// Exemple : a = [2 0 1], b = [3]  ->  "\x02\x00\x01|\x03"
// Ça sert de "nom" à l'état, pour la map visited.
func key(s *stack.Stacks) string {
	b := make([]byte, 0, len(s.A)+len(s.B)+1)
	for _, v := range s.A {
		b = append(b, byte(v)) // les rangs vont de 0 à 5 : tiennent dans un octet
	}
	b = append(b, '|') // séparateur entre a et b
	for _, v := range s.B {
		b = append(b, byte(v))
	}
	return string(b)
}

// path reconstruit la liste d'instructions depuis l'état trouvé.
// On remonte de parent en parent jusqu'au départ, ce qui donne les
// instructions À L'ENVERS ; on inverse donc la liste à la fin.
func path(queue []node, i int) []string {
	var ops []string
	for ; queue[i].parent != -1; i = queue[i].parent {
		ops = append(ops, queue[i].op)
	}
	// inversion du slice (on échange le premier et le dernier, etc.)
	for l, r := 0, len(ops)-1; l < r; l, r = l+1, r-1 {
		ops[l], ops[r] = ops[r], ops[l]
	}
	return ops
}
