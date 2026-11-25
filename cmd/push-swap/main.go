// Programme push-swap : affiche la liste d'instructions qui trie les
// entiers reçus en argument.
//
//	./push-swap "2 1 3 6 5 8"
package main

import (
	"bufio"
	"fmt"
	"os"

	"push-swap/solver"
	"push-swap/stack"
)

func main() {
	// os.Args[0] est le nom du programme, les vrais arguments commencent à 1.
	values, err := stack.ParseArgs(os.Args[1:])
	if err != nil {
		// Le sujet demande "Error" sur la sortie d'ERREUR (stderr).
		fmt.Fprintln(os.Stderr, "Error")
		os.Exit(1)
	}
	// Pas d'argument -> values est vide -> Solve renvoie 0 instruction
	// -> rien n'est affiché, comme demandé.

	// bufio regroupe les écritures : beaucoup plus rapide pour afficher
	// des milliers de lignes (500 nombres = environ 5 000 instructions).
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush() // à la fin de main, on vide le tampon à l'écran

	// On calcule les instructions et on les affiche, une par ligne.
	for _, op := range solver.Solve(values) {
		fmt.Fprintln(out, op)
	}
}
