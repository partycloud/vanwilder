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


sudo docker run --name minecraft -p 25565:25565 -id --restart=always \
  -v minecraft-volume:/data --volume-driver=convoy partycloud/minecraft:ftb-infinity \
  java -Xms1024m -Xmx2048m -XX:PermSize=128m -jar server.jar nogui
