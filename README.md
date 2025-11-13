# OnLife

Simulateur minimaliste inspiré du jeu de la vie où différentes entités (feu, eau, roches, végétation, vie) interagissent sur une grille. Le moteur est écrit en Go et fournit un mode aléatoire ainsi qu'un système de scénarios JSON pour rejouer des situations précises.

## Règles du monde

| Symbole | Type   | Règles principales |
|---------|--------|--------------------|
| `.`     | Rock   | Case inerte servant de fond. Peut « renaître » en cellule de vie selon les règles de Conway. |
| `G`     | Grass  | Devient `Fire` si un feu est adjacent (4 directions). |
| `F`     | Fire   | Dispose d'une durée de vie finie. Chaque tick, la durée diminue, et le feu s'éteint si elle atteint 0 ou si de l'eau est adjacente. Propage le feu sur les `Grass`. (Affiché en 🔥 dans la console.) |
| `W`     | Water  | Éteint immédiatement les feux adjacents et bloque leurs déplacements. |
| `L`     | Life   | Suit strictement les règles du Jeu de la vie de Conway (voisinage à 8 directions). |

Résumé des règles de Conway (appliquées sur Rock/Life) :
- une cellule `Life` survit avec 2 ou 3 voisines vivantes, sinon redevient `Rock` ;
- une case `Rock` devient `Life` lorsqu’elle possède exactement 3 voisines vivantes.

## Exécution rapide

```bash
# Lancer 5 ticks sur une grille 5×5 aléatoire
go run .

# Charger un scénario JSON et afficher son évolution
go run . scenarios/forest_fire.json

# Charger un scénario puis sauvegarder l'état final dans un fichier JSON
go run . scenarios/life_oasis.json snapshots/oasis_final.json

# Démarrer le serveur WebSocket (HTTP) sur le port 8080
go run ./cmd/server -addr :8080 -scenario scenarios/ring_of_fire.json
```

Arguments de `main` :
1. chemin d'un fichier scénario (optionnel) ;
2. chemin de sortie pour sauvegarder un snapshot (optionnel).

### Serveur WebSocket

Le binaire `./cmd/server` expose un endpoint `ws://<host>/ws` qui partage la grille entre tous les clients connectés.

- À la connexion, le serveur envoie immédiatement un message `{ "type": "state", "tick": <n>, "grid": [[".", ...], ...] }`.
- Pour avancer d'un tour, le client envoie `{ "type": "tick" }`. Le serveur met à jour la simulation et diffuse le nouvel état à **tous** les clients.
- La commande `{ "type": "state" }` force l'envoi d'un snapshot sans faire avancer le temps.

Cette API suffit pour un premier frontend : établir un WebSocket, afficher les symboles reçus et déclencher des ticks manuels. Les scénarios chargés côté serveur sont les mêmes que décrits ci-dessous.

## Scénarios inclus

| Fichier | Description |
|---------|-------------|
| `scenarios/glider.json` | Petit planeur `Life` pour valider les règles de Conway. |
| `scenarios/forest_fire.json` | Prairie dense traversée par un canal d’eau servant de coupe-feu. |
| `scenarios/life_oasis.json` | Oasis d’eau avec un glider et un foyer de feu distant. |
| `scenarios/ring_of_fire.json` | Anneau de feu confiné dans un fossé d’eau entouré de cellules de vie. |

Ajoutez vos propres scénarios en copiant l’un de ces fichiers et en éditant le tableau `rows`.

## Format d'un scénario

```json
{
  "name": "example",
  "description": "Texte libre",
  "fireLifetime": 3, // optionnel, redéfinit la durée initiale des feux
  "rows": [
    "..GWF...",
    "..L.L..."
  ]
}
```

Règles du format :
- chaque ligne du tableau `rows` doit avoir la même longueur ;
- caractères acceptés : `.` (Rock), `G`, `W`, `L`, `F` ;
- `fireLifetime` est optionnel (valeur par défaut provenant du code Go) ;
- l’outil vérifie les erreurs de largeur ou de symbole lors du chargement.

## Développement

```bash
# Formater le code
gofmt -w ./...

# Lancer le programme principal
go run .
```

Le projet utilise Go modules (`go 1.25+`). Pensez à ignorer les binaires compilés (`OnLife`) pour éviter de polluer l’historique Git.
