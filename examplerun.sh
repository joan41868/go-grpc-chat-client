#!/usr/bin/env bash


# when gateway is hosted on the same machine
# go run main.go --username=test --room=general


# when gateway is hosted on other machine
# go run main.go --username=test1 --room=general --gateway=GATEWAY_HOST --gatewayPort=GATEWAY_PORT

# when u want to see the rooms in a given gateway
# go run main.go --username=test --listRooms=true