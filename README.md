# Interview task for Kraken Systems, 3 days deadline


## Zombie Game
To play the Zombie Game

https://github.com/thehowl/claws

### Starting the game
To start the game, position inside src/ and run:

'go run main.go'

OR

'go build; ./src;'

#### Setting the port
The server starts on port 8080 by default, you can set the port by using the "addr" flag:

'./src -addr=9000'

#### Connecting to the server
To connect to the server, you can use any client that can use websockets.

One client you can use is https://github.com/thehowl/claws, install it and run:

'./claws ws://localhost:8080'

## Playing the game
To play the game, you need to enter text commands

### Start a new game
"START {playerName}", e.g. "START zombie"

### Join an existing game
If someone else created a game and you want to join:

"JOIN {hostingPlayer} {yourName}", e.g. "JOIN zombie leon"

### Shooting
"SHOOT {row} {column}", e.g. "SHOOT 2 3"

## Rules
First to shoot the zombie wins!
