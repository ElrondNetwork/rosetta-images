# Rosetta images for Docker

**Elrond Network runs on a sharded architecture** - transaction, data and network sharding are leveraged. In the Rosetta implementation, we've decided to abstract away the particularities and complexity of our architecture, and provide a single-chain perspective to the API consumer.

This is achieved through our [Elrond Proxy](https://github.com/ElrondNetwork/elrond-proxy-go), the piece of software that, among others, implements the Rosetta specification and acts as an entry point into the Network, through a set of Observer Nodes. In other words, the Proxy starts as a gateway that resolves the impedance mismatch between the Elrond API (exposed by the Observer Nodes) and the Rosetta API.

The Rosetta deployment for Elrond takes the shape of two Docker images (Elrond Proxy and Elrond Observer) - and a Docker Compose definition to orchestrate `4+1` containers: one Elrond Proxy and four Elrond Observers. This `4+1` setup is usually referred to as an **Elrond Observing Squad**.


**Note:** An Elrond Observing Squad is defined as _a set of N Observer Nodes (one for each Shard, including the Metachain) plus the Elrond Proxy instance (which connects to these Observers and delegates requests towards them)._ Currently the Elrond Mainnet has 3 Shards, plus the Metachain. Therefore, the Observing Squad is composed of 4 Observers and one Proxy instance.

## Prerequisites

### Clone the repository (repositories)

```
cd $HOME
git clone https://github.com/ElrondNetwork/rosetta-images.git
```

### Build the images

```
cd $HOME/rosetta-images

docker image build . -t proxy:latest -f ./proxy/proxy.dockerfile

docker image build . -t rosetta-observer-testnet:latest -f ./observer/testnet.dockerfile
docker image build . -t rosetta-observer-devnet:latest -f ./observer/devnet.dockerfile
docker image build . -t rosetta-observer-mainnet:latest -f ./observer/mainnet.dockerfile
```

### Prepare folders on host

The following script prepares the required folder structure on host:

```
cd $HOME/rosetta-images

./prepare_host.sh ${HOME}/rosetta
```

### Generate keys for observers

The following script generates the node keys, required by the observers:

```
cd $HOME/rosetta-images

./generate_keys.sh ${HOME}/rosetta/keys
```

## Run rosetta

## Run on testnet

```
cd $HOME/rosetta-images

export PROXY_IMAGE=proxy:latest
export OFFLINE_CONFIG=offline_testnet.toml
export OBSERVER_IMAGE=rosetta-observer-testnet:latest
export DATA_FOLDER=${HOME}/rosetta/testnet
export KEYS_FOLDER=${HOME}/rosetta/keys

docker compose --file ./docker-compose.yml up --detach
```

## Run on devnet

```
cd $HOME/rosetta-images

export PROXY_IMAGE=proxy:latest
export OFFLINE_CONFIG=offline_devnet.toml
export OBSERVER_IMAGE=rosetta-observer-devnet:latest
export DATA_FOLDER=${HOME}/rosetta/devnet
export KEYS_FOLDER=${HOME}/rosetta/keys

docker compose --file ./docker-compose.yml up --detach
```

## Run on mainnet

```
cd $HOME/rosetta-images

export PROXY_IMAGE=proxy:latest
export OFFLINE_CONFIG=offline_mainnet.toml
export OBSERVER_IMAGE=rosetta-observer-mainnet:latest
export DATA_FOLDER=${HOME}/rosetta/mainnet
export KEYS_FOLDER=${HOME}/rosetta/keys

docker compose --file ./docker-compose.yml up --detach
```

## Update rosetta

Update the repository (repositories):

```
cd $HOME/rosetta-images
git pull origin
```

Stop the running containers:

```
docker stop rosetta-images-observer-0-1
docker stop rosetta-images-observer-1-1
docker stop rosetta-images-observer-2-1
docker stop rosetta-images-observer-metachain-1
docker stop rosetta-images-proxy-1
docker stop rosetta-images-proxy-rosetta-1
docker stop rosetta-images-proxy-rosetta-offline-1
```

Re-build the images as described above.

Optionally, remove the attached volumes (for testnet):

```
TBD
```

Run the containers as described above.
