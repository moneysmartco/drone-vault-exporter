---
date: 2017-01-30T00:00:00+00:00
title: vault-exporter
author: dhoeric
tags: [ secret, vault ]
repo: moneysmartco/drone-vault-exporter
image: moneysmartco/drone-vault-exporter
---


Use the Vault exporter plugin to export env var saved in Vault server, in dotenv file format. The below pipeline configuration demonstrates simple usage:

```yaml
pipeline:
  export_envvar:
    image: moneysmartco/drone-vault-exporter:0.0.1
    vault_key_path: secret/app_env/staging/api_server
    deploy_env_path: .deploy.env
    secrets:
      - vault_addr
      - vault_token
```

Example configuration in your `.drone.yml` file for using [AppRole](https://www.vaultproject.io/docs/auth/approle.html) auth method:

```diff
pipeline:
  export_envvar:
    image: moneysmartco/drone-vault-exporter:0.0.1
    vault_key_path: secret/app_env/staging/api_server
+   vault_auth_method: APPROLE
    deploy_env_path: .deploy.env
    secrets:
      - vault_addr
-     - vault_token
+     - vault_role_id
+     - vault_secret_id
```

Example configuration in your `.drone.yml` file for exporting in helm-yaml format:

```diff
pipeline:
  export_envvar:
    image: moneysmartco/drone-vault-exporter:0.0.1
    vault_key_path: secret/app_env/staging/api_server
+   output_format: helm-yaml
+   deploy_env_path: env.yaml
    vault_auth_method: APPROLE
    secrets:
      - vault_addr
      - vault_token
      - vault_role_id
      - vault_secret_id
```

# Parameter Reference

vault_key_path
: key path for secrets in vault server

deploy_env_path
: filename to be saved in the workspace (default: `.deploy.env`)

vault_auth_method (optional)
: auth method to be use on connecting vault server

output_format (optional)
: fileformat to be saved, helm-yaml / dotenv (default)

helm_env_key (optional)
: Helm key for env vars (default: `envs`)


# Secret Reference

vault_addr
: vault server address

vault_token
: vault token for access

vault_role_id
: vault role_id for generating vault token when `vault_auth_method = APPROLE`

vault_secret_id
: vault secret_id for generating vault token when `vault_auth_method = APPROLE`

