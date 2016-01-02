# Vanwilder

Your Partycloud liaison.

## Protocol

JSON connection over websocket to API. Commands are sent from the api to the agent. Events are sent from the agent to the API.

Data volumes must be created before they can be used.

Requests made to the agent are in the format:
```
{
  "command": "$CMD",
  ...args...
}

```

The agent will publish updates as:
```
{
  "event": "$EVENT",
  ...properties...
}

## Commands

list-volumes
create-volume
delete-volume

start-game
  game: (minecraft/ftb-infinity/factorio/ark/terraria)
  volume: (data1)

```

## Test commands

echo -e '{"command":"start-game", "game":"ubuntu", "cmd-args": "nc -l 25565", "ports": ["25565/tcp"]}' | sudo nc -U /var/run/van.sock
echo -e '{"command":"stop-game", "id": ""}' | sudo nc -U /var/run/van.sock
