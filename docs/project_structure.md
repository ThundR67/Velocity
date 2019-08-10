# Project Structure
### The project structure of Velocity


# Global
#### This folder contains code which is accessed by each component of Velocity.

## Clients
##### This folder contains high level clients to each micro-service

## Config
##### This folder contains configurations of each component of Velocity.

## Logs
##### This folder contains logger.go which contains boiler plate code for creating a logger

# Gateway
#### This folder contains all the servers which will talke to clients. Currently it contains the auth server. Soon resource server will also be added.

## Auth Server
##### The folder handler contains code to handle all the routes. The folder router contains router.go which contains information about which routes are handled by which handler.

# Services
#### This folder contains all the microservices. Each service has its own sub folder. Each service has the handler folder which contains the service handler which is used by go-micro. Each service has proto folder which contains the proto defination of the service and compiled version of that proto file in go.