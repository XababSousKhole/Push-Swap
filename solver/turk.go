package solver

import "push-swap/stack"

// turk trie les grandes listes en 4 phases :
//
//  1. on pousse tout dans b sauf 3 nombres. À chaque fois on choisit le
//     nombre de a le moins cher à placer au bon endroit dans b, pour que b
//     reste trié en ordre DÉCROISSANT (à une rotation près) ;
//  2. on trie les 3 nombres restants dans a ;
//  3. on ramène chaque nombre de b dans a, juste au-dessus du plus petit
//     nombre de a qui lui est supérieur ;
//  4. on tourne a pour remettre le minimum en haut.
func turk(ranks []int) []string {
	// r contient les piles + la liste des instructions faites (voir solver.go)
	r := &recorder{s: stack.New(ranks)}

	// Phase 1 : vider a dans b (sauf 3 nombres)
	for len(r.s.A) > 3 {
		if len(r.s.B) < 2 {
			r.do("pb", 1) // les deux premiers partent sans calcul
			continue
		}
		r.apply(r.cheapestMove()) // ra/rb/rr/rra/rrb/rrr pour bien placer les 2 piles
		r.do("pb", 1)             // puis on pousse le nombre choisi dans b
	}

	// Phase 2 : trier les 3 nombres restants dans a
	r.sortThree()

	// Phase 3 : tout ramener dans a, chacun à sa place
	for len(r.s.B) > 0 {
		r.rotateA(targetInA(r.s.A, r.s.B[0])) // on prépare la bonne place dans a
		r.do("pa", 1)                         // on y pose le sommet de b
	}

	// Phase 4 : a est trié mais peut-être "tourné" ([3 4 1 2]) :
	// on remet le plus petit en haut.
	r.rotateA(indexOfMin(r.s.A))
	return r.ops
}

// move = un "plan de rotations" à faire avant un pb.
// Il dit combien de ra, rb, rra et rrb il faut pour que :
//   - le nombre choisi soit en haut de a ;
//   - sa place soit en haut de b.
type move struct {
	ra, rb, rra, rrb int
}

// cost = nombre d'instructions que coûte ce plan.
// On prend le max et pas la somme car rr fait ra + rb en un seul coup :
// 5 ra + 3 rb = 3 rr + 2 ra = 5 instructions (et pas 8).
func (m move) cost() int {
	return max(m.ra, m.rb) + max(m.rra, m.rrb)
}

// cheapestMove : LE CŒUR DE L'ALGO.
// Pour CHAQUE nombre de a, on calcule combien il coûterait de le pousser
// au bon endroit dans b. On garde le moins cher.
func (r *recorder) cheapestMove() move {
	a, b := r.s.A, r.s.B
	var best move
	bestCost := -1 // -1 = "pas encore de meilleur trouvé"

	for i, x := range a { // i = position du nombre dans a, x = sa valeur
		// où x doit aller dans b ?
		j := targetInB(b, x)
		// combien de rotations pour mettre x en haut de a ET sa place en haut de b ?
		m := bestRotations(i, len(a), j, len(b))

		// moins cher que le meilleur actuel ? on le garde.
		if c := m.cost(); bestCost == -1 || c < bestCost {
			best, bestCost = m, c
		}
	}
	return best
}

// bestRotations : pour amener l'index i de a ET l'index j de b au sommet,
// il y a 4 façons de faire (haut/haut, bas/bas, haut/bas, bas/haut).
// On garde la moins chère.
//
// Rappel : pour amener l'index i en haut, soit on fait i fois ra (il monte),
// soit lenA-i fois rra (il passe par le bas et revient en haut).
func bestRotations(i, lenA, j, lenB int) move {
	options := []move{
		{ra: i, rb: j},                 // on monte les deux -> rr
		{rra: lenA - i, rrb: lenB - j}, // on descend les deux -> rrr
		{ra: i, rrb: lenB - j},         // a monte, b descend
		{rra: lenA - i, rb: j},         // a descend, b monte
	}
	best := options[0]
	for _, o := range options[1:] {
		if o.cost() < best.cost() {
			best = o
		}
	}
	return best
}

// apply exécute un move en regroupant au maximum en rr / rrr.
//
// Exemple : il faut 5 ra et 3 rb.
// Au lieu de 8 instructions, on fait 3 rr (qui font ra + rb d'un coup)
// puis les 2 ra qui restent : 5 instructions au total.
func (r *recorder) apply(m move) {
	// partie "on monte" : autant de rr que possible, puis le reste
	both := min(m.ra, m.rb)
	r.do("rr", both)      // ra + rb ensemble
	r.do("ra", m.ra-both) // les ra en trop
	r.do("rb", m.rb-both) // les rb en trop

	// partie "on descend" : autant de rrr que possible, puis le reste
	both = min(m.rra, m.rrb)
	r.do("rrr", both)       // rra + rrb ensemble
	r.do("rra", m.rra-both) // les rra en trop
	r.do("rrb", m.rrb-both) // les rrb en trop
}

// rotateA amène l'index i de a au sommet par le chemin le plus court.
//
// Si l'élément est dans la moitié haute, on fait des ra (il remonte).
// S'il est dans la moitié basse, des rra sont plus rapides
// (il passe par le bas et revient en haut).
func (r *recorder) rotateA(i int) {
	if i <= len(r.s.A)/2 {
		r.do("ra", i) // moitié haute : i fois ra
	} else {
		r.do("rra", len(r.s.A)-i) // moitié basse : (taille - i) fois rra
	}
}

// sortThree trie a quand il contient au plus 3 nombres (2 instructions max).
//
// Étape 1 : on envoie le plus grand tout en bas.
//
//	max en haut   [3 1 2] -> ra  -> [1 2 3]
//	max au milieu [1 3 2] -> rra -> [2 1 3]
//	max en bas    : déjà bien placé, rien à faire
//
// Étape 2 : si les 2 du dessus sont à l'envers, sa.
//
//	[2 1 3] -> sa -> [1 2 3]
func (r *recorder) sortThree() {
	a := r.s.A
	if len(a) == 3 {
		// étape 1 : le plus grand part en bas
		switch indexOfMax(a) {
		case 0: // le max est en haut -> ra l'envoie en bas
			r.do("ra", 1)
		case 1: // le max est au milieu -> rra fait remonter le bas, le max descend
			r.do("rra", 1)
		}
	}
	// étape 2 : si les deux du haut sont dans le mauvais ordre, on les échange
	if len(r.s.A) >= 2 && r.s.A[0] > r.s.A[1] {
		r.do("sa", 1)
	}
}

// targetInB : x doit se poser sur le plus grand nombre de b qui lui est
// inférieur. S'il n'y en a pas, x est le nouveau minimum : il va sur le max.
//
// Pourquoi ? b est rangé en ordre DÉCROISSANT (grand en haut).
// Exemple : b = [8 5 2], on veut pousser x = 6.
// Le plus grand nombre de b plus petit que 6 est 5 -> on tourne b pour
// avoir [5 2 8], puis pb donne [6 5 2 8] : toujours décroissant (en boucle).
//
// Renvoie la POSITION (l'index) de la cible dans b.
func targetInB(b []int, x int) int {
	target := -1
	for i, v := range b {
		if v < x && (target == -1 || v > b[target]) {
			target = i
		}
	}
	if target == -1 {
		return indexOfMax(b)
	}
	return target
}

// targetInA : y doit se poser sur le plus petit nombre de a qui lui est
// supérieur. S'il n'y en a pas, y est le nouveau maximum : il va sur le min.
//
// C'est le miroir de targetInB, car a est rangé en ordre CROISSANT.
// Exemple : a = [2 7 9], on ramène y = 5 depuis b.
// Le plus petit nombre de a plus grand que 5 est 7 -> on tourne a pour
// avoir [7 9 2], puis pa donne [5 7 9 2] : toujours croissant (en boucle).
func targetInA(a []int, y int) int {
	target := -1
	for i, v := range a {
		if v > y && (target == -1 || v < a[target]) {
			target = i
		}
	}
	if target == -1 {
		return indexOfMin(a)
	}
	return target
}

// indexOfMin renvoie la position du plus petit nombre du slice.
func indexOfMin(s []int) int {
	idx := 0
	for i, v := range s {
		if v < s[idx] {
			idx = i
		}
	}
	return idx
}

// indexOfMax renvoie la position du plus grand nombre du slice.
func indexOfMax(s []int) int {
	idx := 0
	for i, v := range s {
		if v > s[idx] {
			idx = i
		}
	}
	return idx
}
