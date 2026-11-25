// Programme checker : lit des instructions sur l'entrée standard, les
// exécute sur les entiers reçus en argument, puis affiche :
//
//   - OK si a est trié et b est vide ;
//
//   - KO sinon ;
//
//   - Error si un argument ou une instruction est invalide.
//
//     ./push-swap "3 2 1" | ./checker "3 2 1"
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"push-swap/stack"
)

func main() {
	// 1. Lire et vérifier les nombres (même fonction que push-swap).
	values, err := stack.ParseArgs(os.Args[1:])
	if err != nil {
		fail()
	}
	if len(values) == 0 {
		return // pas d'argument : on n'affiche rien
	}

	// 2. Lire toutes les instructions, une par ligne, sur l'entrée standard
	//    (le clavier, ou ce qu'envoie push-swap grâce au "|").
	var ops []string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		// On enlève un éventuel "\r" (fins de ligne Windows).
		line := strings.TrimSuffix(scanner.Text(), "\r")
		if line == "" {
			continue // tolère la ligne vide finale ajoutée par `echo -e "...\n"`
		}
		ops = append(ops, line)
	}
	if scanner.Err() != nil {
		fail()
	}

	// 3. Exécuter les instructions une par une sur les piles.
	s := stack.New(values)
	for _, op := range ops {
		if s.Exec(op) != nil {
			fail() // instruction inconnue ("foo", "SA", "sa "...)
		}
	}

	// 4. Verdict.
	if s.IsSorted() {
		fmt.Println("OK")
	} else {
		fmt.Println("KO")
	}
}

// fail affiche "Error" sur la sortie d'erreur et quitte le programme.
func fail() {
	fmt.Fprintln(os.Stderr, "Error")
	os.Exit(1)
}
