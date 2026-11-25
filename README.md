# push-swap

Deux programmes en Go (bibliothèque standard uniquement) :

- `push-swap` : affiche les instructions qui trient la pile `a`.
- `checker` : lit des instructions sur l'entrée standard, les exécute, affiche `OK` ou `KO`.

## Compiler et lancer

```console
$ go build -o push-swap ./cmd/push-swap
$ go build -o checker ./cmd/checker
$ ARG="4 67 3 87 23"; ./push-swap "$ARG" | ./checker "$ARG"
OK
```

## Tests

```console
$ go test ./...
$ go test -v -run Score ./solver   # nombre moyen d'instructions
```

## Organisation

```
stack/          les deux piles, les 11 instructions, le parsing des arguments
solver/         l'algorithme : BFS (≤ 6 nombres) et algorithme Turc (> 6)
cmd/push-swap/  programme push-swap
cmd/checker/    programme checker
```

## Algorithme

1. **Normalisation** : chaque valeur est remplacée par son rang (0 à n-1).
2. **≤ 6 nombres → BFS** : on explore tous les états possibles, niveau par niveau.
   Le premier état trié trouvé utilise le nombre minimal d'instructions.
3. **> 6 nombres → algorithme Turc** :
   - on pousse tout dans `b` sauf 3 nombres, en choisissant à chaque fois le nombre
     de `a` le moins cher à placer au bon endroit dans `b` (trié en décroissant) ;
   - on trie les 3 nombres restants ;
   - on ramène chaque nombre de `b` au bon endroit dans `a` ;
   - on tourne `a` pour remettre le minimum en haut.

Résultats moyens : 100 nombres ≈ 550 instructions, 500 nombres ≈ 5 100.
