# tic-tak-toe

Just simple investigations/attempts of tic-tak-toe game on
* Golang
* Websockets
* Redis
* Docker

## Features
* Multiplayer
* Horizontal scale

## Run server

```shell
cd ./.cloud/compose/
docker-compose up -d
```

### Server configuration

You can customize options via ENV variables. All app vars are starting with `TTT_`

|Variable name|Default|Description|
|---|---|---|
|TTT_ADDR|:8080|Service network address|
|TTT_REDIS_ADDR|redis:6379|Connection credentials for Redis. By default `redis` name is docker-compose service name. Change if you need|

## Run client

> Open `http://you_lan_ip:8080/` or `http://localhost:8080/` to start new game session. You have to have opponent
