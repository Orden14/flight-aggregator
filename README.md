# Flight Aggregator

## Groupe
- Thomas L.
- David W.
- Antoine H.
- Eden O.

## Projet

Fork de [flight-aggregator par Romain Chenard](https://github.com/RomainC75/flight-aggregator)  
Projet réalisé dans le cadre du cours de Golang en suivant les consignes du repo original.

## Installation

1. Cloner le repository
```bash
git clone https://github.com/Orden14/flight-aggregator
cd flight-aggregator
```

2. Build et run le projet
```bash
docker-compose up -d
```

## Tests

1. Installer les dépendances
```bash
cd server
go mod tidy
```

2. lancer les tests
``` bash
go test ./test -v
```

## Structure du projet

- `/main.go` : Point d'entrée de l'application
- `/config/` : contient `config.go` pour la gestion de la configuration de l'application
- `/httpserver/` : contient `router.go` pour la gestion des routes HTTP (equivalent d'un controleur)
- `/handler/` : contient les handlers pour la gestion des requêtes HTTP
- `/domain/` : contient les structures de données internes à l'application
- `/model/` : contient les structures de données des vols en fonction du schema de donnée des deux serveurs JSON
- `/repository/` : contient les repositories pour la gestion des appels aux serveurs JSON
- `/service/` : contient `flight_service.go` pour la logique métier
- `/sorter/` : contient les fonctions de tri des vols

## Utilisation

### A. Accès aux serveurs

1. Serveur principal : http://localhost:3001/
2. Serveurs JSON database 1 : http://localhost:4001/
3. Serveurs JSON database 2 : http://localhost:4002/

(customisable dans le [.env](.env))

### B. Endpoints pour le serveur principal

1. [GET] `/health` : Vérifie l'état de santé du serveur
2. [GET] `/flights` : Récupère tous les vols (triés par prix par défaut)

### C. Paramètres pour la route /flight

- `sort` : Critère de tri (price, travel_time). Par défaut : price
- `from` : Code IATA de l'aéroport de départ (ex: CDG)
- `to` : Code IATA de l'aéroport d'arrivée (ex: HND)
