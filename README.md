# OBBAA-477
This repository contains the code for the obbaa based TR-477 compliant VNFs.

## PPPoE IA VNF

The PPPoE IA VNF is deployed as a microservice and is responsible for receiving packets from the pOLT and process them. This includes appending new tags according to TR-101i2 and forward them to the PPPoE server either via control plane or data plane.

### Building PPPoE IA VNF
To build the PPPoE IA VNF docker image run 
```
docker build -t <image name> -f pppoe-relay-vnf/vnf.Dockerfile .
``` 
from the root of this project. Alternatively, see [Running PPPoE IA VNF](#running-pppoe-ia-vnf) to build and run the VNF via docker-compose.


### Running PPPoE IA VNF
To run PPPoE IA VNF, several docker-compose files are provided, to create different scenarios to run the VNF. To run any of the scenarios run 
```
docker-compose -f <name of docker-compose file> up -d --build
```
For example, to run the PPPoE directly connected to the OLT in an inband scenario run 
```
docker-compose -f docker-compose.yml up -d --build
```

The follwing environment variables can be used to customize certain aspects of the PPPoE VNF:

- KAFKA_HOST: ip address/hostname of the kafka container
- KAFKA_PORT: port of the kafka container
- MONGO_HOST: ip address/hostname of the mongo container
- MONGO_PORT: port of the mongo container
- VNF_MODE: defines the mode on which the VNF will operate. If "server", the vnf will wait for a connection on SOCKET_GRPC. If "client", the vnf will attempt to connect to 
the device on SOCKET_GRPC
- DB_NAME: name of the databse in the mongo container reserved for this VNF 
- SOCKET_GRPC: if VNF_MODE is "server", SOCKET_GRPC refers to the grpc socket the vnf is listening to. If VNF_MODE is "client", SOCKET_GRPC refers to the address the VNF will attempt to connect to.
- DISCARD_ON_ERROR: if "true", the VNF will discard packets that fails to process. If "false", the VNF will send the packets that fails to process unaltered
- VNF_NAME: name of the VNF. It is used as a prefix for the kafka topics and as entity/endpoint name in TR-477 messages

#

For more information refer to the OB-BAA public documentation.
