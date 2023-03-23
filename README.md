# Tic-Tac-Toe

## The Spec

The game board grid looks as follows

    .-----------.
    | 0 | 1 | 2 |
    +---+---+---+
    | 3 | 4 | 5 |
    +---+---+---+
    | 6 | 7 | 8 |
    `-----------´

So, a board position

    .-----------.
    | X | O | - |
    +---+---+---+
    | - | X | - |
    +---+---+---+
    | - | O | X |
    `-----------´

translates to

    XO--X--OX
    012345678

See the accompanying swagger.yaml for the REST API documentation in Swagger
format (https://swagger.io).


## Game flow:

- The client (player) starts a game by making a POST request to /games.
  The POST request contains a representation of a game board, either empty
  (computer starts) or with the first move made (player starts).
  The player/computer can choose either noughts or crosses.

- The backend responds with the location URL of the started game.

- Client GETs the board state from the URL.

- Client PUTs the board state with a new move to the URL.

- Backend validates the move, makes it's own move and updates the game state.
  The updated game state is returned in the PUT response.

- And so on. The game is over once the computer or the player gets 3 noughts
  or crosses, horizontally, vertically or diagonally or there are no moves to
  be made.

## Application structure

The application is developed according to the principles of DDD.

    .
    ├── cmd
    │   └── srv
    ├── internal
    │   ├── api
    │   │   ├── protocol
    │   │   │   └── json
    │   │   └── transport
    │   │       └── http
    │   ├── app
    │   │   └── module
    │   │       └── game
    │   ├── domain
    │   │   ├── error
    │   │   └── repo
    │   ├── pkg
    │   │   └── postgres
    │   └── service
    │       ├── migration
    │       └── repo
    ├── docker-compose.yaml
    ├── Dockerfile
    ├── go.mod
    └── go.sum
    
## Run

The application is wrapped into Docker Compose.

```shell
docker-compose up
```

### Resources
 - UI: http://127.0.0.1:8081
 - API: http://127.0.0.1:8080
 - Database: `127.0.0.1:5432`, name `tictactoe`, user/password `postgres`
