# push-swap

Deux programmes en Go (bibliothèque standard uniquement) :

- **push-swap** : affiche les instructions qui trient la pile `a` (le plus petit en haut).
- **checker** : lit des instructions sur l'entrée standard, les exécute, puis affiche `OK` si c'est trié, `KO` sinon, et `Error` si l'entrée est invalide.

## Organisation

```
stack/          les 2 piles, les 11 instructions, la lecture des arguments
solver/         l'algorithme : BFS (≤ 6 nombres) et algorithme Turc (> 6)
cmd/push-swap/  programme push-swap
cmd/checker/    programme checker
```

## Algorithme

1. **Normalisation** : chaque nombre est remplacé par son rang ([42 -7 100] → [1 0 2]).
2. **≤ 6 nombres → BFS** : on teste toutes les instructions niveau par niveau (1 coup, 2 coups…).
   La première solution trouvée est donc la plus courte possible.
3. **> 6 nombres → algorithme Turc (glouton)** :
   - on pousse tout dans `b` sauf 3 nombres, en choisissant toujours le nombre le moins cher à placer ;
   - on trie les 3 restants (2 coups max) ;
   - on ramène chaque nombre de `b` à sa place dans `a` ;
   - on tourne `a` pour remettre le minimum en haut.
   - `rr` / `rrr` font tourner les deux piles en un seul coup.

Résultats : 100 nombres ≈ 550 instructions, 500 nombres ≈ 5 100.

---

# Commandes pour la présentation

Toutes les commandes se lancent depuis le dossier du projet.

## 0. Préparation

```bash
cd ~/Documents/YNOV/Push-Swap
go build -o push-swap ./cmd/push-swap
go build -o checker ./cmd/checker
clear
```

## 1. Démo de base

```bash
./push-swap "2 1 3 6 5 8"
```
→ affiche les instructions, une par ligne (exactement l'exemple du sujet)

```bash
ARG="4 67 3 87 23"; ./push-swap "$ARG" | ./checker "$ARG"
```
→ `OK`

```bash
ARG="4 67 3 87 23"; ./push-swap "$ARG" | wc -l
```
→ `6` (comme demandé dans le sujet)

## 2. Exemples du sujet pour le checker

```bash
echo -e "rra\npb\nsa\nrra\npa" | ./checker "3 2 1 0"
```
→ `OK`

```bash
echo -e "sa\nrra\npb" | ./checker "3 2 1 0"
```
→ `KO` (les instructions ne trient pas)

## 3. Gestion des erreurs

```bash
./push-swap "0 one 2 3"
```
→ `Error` (pas un nombre)

```bash
./push-swap "1 2 1"
```
→ `Error` (doublon)

```bash
./push-swap "2147483648"
```
→ `Error` (trop grand pour un int)

```bash
echo -e "rra\npb\nsa\n" | ./checker "3 2 one 0"
```
→ `Error` (argument invalide)

```bash
echo -e "sa\nfoo" | ./checker "2 1"
```
→ `Error` (instruction inconnue)

```bash
./push-swap
./checker
```
→ rien (pas d'argument)

```bash
./push-swap "1 2 3 4 5"
```
→ rien (déjà trié : 0 instruction)

## 4. Performances (nombres au hasard)

```bash
ARG=$(seq 1 10000 | sort -R | head -n 100 | tr '\n' ' '); ./push-swap "$ARG" | wc -l; ./push-swap "$ARG" | ./checker "$ARG"
```
→ environ `550`, puis `OK`

```bash
ARG=$(seq 1 10000 | sort -R | head -n 500 | tr '\n' ' '); ./push-swap "$ARG" | wc -l; ./push-swap "$ARG" | ./checker "$ARG"
```
→ environ `5100`, puis `OK`

```bash
ARG=$(seq 1 10000 | sort -R | head -n 500 | tr '\n' ' '); time ./push-swap "$ARG" > /dev/null
```
→ quelques millisecondes

## 5. Tests unitaires

```bash
go test ./...
```
→ `ok` pour `push-swap/solver` et `push-swap/stack`

```bash
go test -v -run Score ./solver
```
→ nombre moyen et pire nombre d'instructions pour 3, 5, 100 et 500 nombres

```bash
go vet ./...
gofmt -l .
```
→ rien affiché = aucun problème détecté, code bien formaté

## 6. Montrer le code

```bash
open .
open -a "Visual Studio Code" .
```
→ ouvre le dossier dans le Finder / dans VS Code

Fichiers à montrer : `stack/stack.go` (les 11 instructions), `solver/bfs.go` (BFS), `solver/turk.go` (algorithme Turc).
