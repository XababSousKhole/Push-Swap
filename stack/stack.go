// Package stack contient les deux piles (a et b) et les 11 instructions
// du langage push-swap.
//
// Convention : l'index 0 d'un slice est le SOMMET de la pile.
package stack

import "errors"

// ErrUnknownInstruction est renvoyée quand une instruction n'existe pas.
var ErrUnknownInstruction = errors.New("unknown instruction")

// Stacks regroupe les deux piles.
type Stacks struct {
	A []int
	B []int
}

// New crée les piles : a contient une copie de values, b est vide.
func New(values []int) *Stacks {
	a := make([]int, len(values))
	copy(a, values)
	return &Stacks{A: a, B: []int{}}
}

// Clone renvoie une copie indépendante des deux piles.
// Utilisé par le BFS : on essaie une instruction sur la copie
// sans abîmer l'état d'origine.
func (s *Stacks) Clone() *Stacks {
	c := &Stacks{A: make([]int, len(s.A)), B: make([]int, len(s.B))}
	copy(c.A, s.A)
	copy(c.B, s.B)
	return c
}

// Exec exécute une instruction ("sa", "pb", "rrr", ...).
//
// Les 11 instructions se rangent en 4 familles :
//
//	s  = swap           : échange les 2 éléments du dessus
//	p  = push           : déplace le sommet d'une pile vers l'autre
//	r  = rotate         : le sommet descend tout en bas
//	rr = reverse rotate : le bas remonte tout en haut
//
// La dernière lettre dit sur quelle pile on agit :
// a, b, ou les deux (s / r / r répétée : ss, rr, rrr).
//
// Dans les exemples ci-dessous, la pile est écrite de gauche (sommet)
// à droite (bas) : [1 2 3] veut dire 1 en haut, 3 en bas.
func (s *Stacks) Exec(op string) error {
	switch op {

	// ---------- SWAP : échanger les 2 premiers ----------

	case "sa": // swap a : a = [1 2 3] -> [2 1 3]
		swap(s.A)
	case "sb": // swap b : b = [4 5 6] -> [5 4 6]
		swap(s.B)
	case "ss": // sa + sb en même temps (1 seule instruction comptée)
		swap(s.A)
		swap(s.B)

	// ---------- PUSH : passer le sommet d'une pile à l'autre ----------

	case "pa": // push a : le sommet de b va sur a
		// a = [1 2]  b = [3 4]   ->   a = [3 1 2]  b = [4]
		push(&s.B, &s.A)
	case "pb": // push b : le sommet de a va sur b
		// a = [1 2]  b = [3 4]   ->   a = [2]  b = [1 3 4]
		push(&s.A, &s.B)

	// ---------- ROTATE : le premier devient le dernier ----------

	case "ra": // rotate a : a = [1 2 3] -> [2 3 1]
		rotate(s.A)
	case "rb": // rotate b : b = [4 5 6] -> [5 6 4]
		rotate(s.B)
	case "rr": // ra + rb en même temps (1 seule instruction comptée)
		rotate(s.A)
		rotate(s.B)

	// ---------- REVERSE ROTATE : le dernier devient le premier ----------

	case "rra": // reverse rotate a : a = [1 2 3] -> [3 1 2]
		reverseRotate(s.A)
	case "rrb": // reverse rotate b : b = [4 5 6] -> [6 4 5]
		reverseRotate(s.B)
	case "rrr": // rra + rrb en même temps (1 seule instruction comptée)
		reverseRotate(s.A)
		reverseRotate(s.B)

	// N'importe quoi d'autre ("SA", "sa ", "foo"...) est refusé :
	// le checker affichera alors "Error".
	default:
		return ErrUnknownInstruction
	}
	return nil
}

// IsSorted indique si a est trié par ordre croissant (plus petit en haut)
// et si b est vide.
func (s *Stacks) IsSorted() bool {
	return len(s.B) == 0 && IsAscending(s.A)
}

// IsAscending indique si le slice est trié par ordre croissant.
func IsAscending(values []int) bool {
	for i := 1; i < len(values); i++ {
		if values[i-1] > values[i] { // un nombre plus grand que le suivant
			return false
		}
	}
	return true
}

// swap échange les deux premiers éléments (ne fait rien s'il y en a moins de 2).
func swap(s []int) {
	if len(s) >= 2 {
		s[0], s[1] = s[1], s[0] // échange en une ligne grâce à Go
	}
}

// push prend le sommet de from et le pose au sommet de to.
// On passe des pointeurs (*[]int) car la TAILLE des deux slices change :
// sans pointeur, la modification serait perdue en sortant de la fonction.
func push(from, to *[]int) {
	if len(*from) == 0 {
		return // pile vide : rien à pousser
	}
	top := (*from)[0]                // on récupère le sommet
	*from = (*from)[1:]              // on l'enlève de from
	*to = append([]int{top}, *to...) // on le met devant (= au sommet) de to
}

// rotate : le premier élément devient le dernier.
func rotate(s []int) {
	if len(s) >= 2 {
		first := s[0]       // on garde le premier de côté   [1 2 3]
		copy(s, s[1:])      // on décale tout vers le haut   [2 3 3]
		s[len(s)-1] = first // on remet le premier à la fin  [2 3 1]
	}
}

// reverseRotate : le dernier élément devient le premier.
func reverseRotate(s []int) {
	if len(s) >= 2 {
		last := s[len(s)-1]       // on garde le dernier de côté  [1 2 3]
		copy(s[1:], s[:len(s)-1]) // on décale tout vers le bas   [1 1 2]
		s[0] = last               // on met le dernier en haut    [3 1 2]
	}
}
