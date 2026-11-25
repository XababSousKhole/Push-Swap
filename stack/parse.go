package stack

import (
	"errors"
	"strconv"
	"strings"
)

// ErrInvalidArgs est renvoyée si un argument n'est pas un entier
// ou s'il y a des doublons. Le programme affichera alors "Error".
var ErrInvalidArgs = errors.New("invalid arguments")

// ParseArgs transforme les arguments du programme en liste d'entiers.
//
//	./push-swap "2 1 3"   -> args = ["2 1 3"]       -> [2 1 3]
//	./push-swap 2 1 3     -> args = ["2", "1", "3"] -> [2 1 3]
//
// Les deux écritures marchent. Le premier nombre sera le sommet de a.
func ParseArgs(args []string) ([]int, error) {
	// On recolle tous les arguments avec des espaces, puis on découpe
	// sur les espaces (Fields ignore aussi les espaces en trop).
	fields := strings.Fields(strings.Join(args, " "))

	values := make([]int, 0, len(fields))

	// seen = les nombres déjà rencontrés, pour repérer les doublons.
	seen := make(map[int]bool, len(fields))

	for _, f := range fields {
		// Conversion texte -> nombre.
		// "one" ou "1.5" échouent ici -> Error.
		// Le 32 refuse les nombres trop grands pour un int classique
		// (au-delà de 2147483647 ou en dessous de -2147483648).
		n, err := strconv.ParseInt(f, 10, 32)
		if err != nil {
			return nil, ErrInvalidArgs
		}
		v := int(n)

		// Déjà vu ? C'est un doublon -> Error.
		if seen[v] {
			return nil, ErrInvalidArgs
		}
		seen[v] = true

		values = append(values, v)
	}
	return values, nil
}
