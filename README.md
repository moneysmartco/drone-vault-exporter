# drone-vault-exporter

Drone plugin to export vault secret into dotenv file. For the usage information and a listing of the available options please take a look at [the docs](DOCS.md).

## Build

Build the binary with the following commands:

```
make install
make build
```

## Docker

Build the docker image with the following commands:

```
make linux_amd64 docker_image docker_deploy tag=X.X.X
```

## Usage

Execute from the working directory:

1) with vault token:
```sh
docker run --rm \
  -e PLUGIN_DEPLOY_ENV_PATH=./.deploy.env \
  -e PLUGIN_VAULT_KEY_PATH=secret/key/path \
  -e VAULT_ADDR=https://vault.server \
  -e VAULT_TOKEN=xxxxxxxxx \
  moneysmartco/drone-vault-exporter
```

2) with vault approle authentication:
```sh
docker run --rm \
  -e PLUGIN_DEPLOY_ENV_PATH=./.deploy.env \
  -e PLUGIN_VAULT_KEY_PATH=secret/key/path \
  -e VAULT_ADDR=https://vault.server \
  -e PLUGIN_VAULT_AUTH_METHOD=APPROLE \
  -e VAULT_ROLE_ID=xxx-xxx-xxx-xxx-xxx \
  -e VAULT_SECRET_ID=xxx-xxx-xxx-xxx-xxx \
  moneysmartco/drone-vault-exporter
```

3) output with helm-yaml format
```sh
docker run --rm \
  -e PLUGIN_OUTPUT_FORMAT=helm-yaml \
  -e PLUGIN_HELM_ENV_KEY=envs \
  -e PLUGIN_DEPLOY_ENV_PATH=./env.yaml \
  -e PLUGIN_VAULT_KEY_PATH=secret/key/path \
  -e PLUGIN_VAULT_AUTH_METHOD=APPROLE \
  -e VAULT_ADDR=https://vault.server \
  -e VAULT_ROLE_ID=xxx-xxx-xxx-xxx-xxx \
  -e VAULT_SECRET_ID=xxx-xxx-xxx-xxx-xxx \
  moneysmartco/drone-vault-exporter
```