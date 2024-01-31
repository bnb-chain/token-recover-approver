# TOKEN RECOVER APPROVER

<p>
  <a href="https://github.com/bnb-chain/token-recover-approver/blob/develop/COPYING">
    <img src="https://img.shields.io/github/license/bnb-chain/token-recover-approver?style=flat-square&color=blue">
  </a>
  <img src="https://img.shields.io/github/go-mod/go-version/bnb-chain/token-recover-approver?style=flat-square">
  <a href="https://pkg.go.dev/github.com/bnb-chain/token-recover-approver">
    <img src="https://img.shields.io/badge/Go-reference-blue?style=flat-square">
  </a>
</p>


## Introduction

Prior to the scheduled discontinuation of operations on the Beacon Chain (BC), it is strongly advised that users promptly initiate cross-chain transactions to transfer their assets to other networks. Following the cessation of BC operations, the community team will capture a snapshot of users' assets on BC, which will be publicly released and acknowledged by the community.

Subsequently, a Merkle tree will be generated based on users' balances in the snapshot. The root of this tree will be securely stored in the system contract of the BNB Smart Chain (BSC).

Users wishing to prove ownership of the original tokens on BC can do so by providing a Merkle proof and their BC account's signature. Once ownership and token information are verified, the system contract on BSC will unlock the corresponding amount of tokens from the token hub and allocate them to the user's account on BSC.

The token recovery process can be initiated through a user-friendly web app or a command-line interface (CLI), providing flexibility for users to choose their preferred method.

![Token Recovery Process](./assets/token-recover.png)

The `TOKEN RECOVER APPROVER` system has been designed with security in mind, ensuring that recover requests are approved by a designated approver account. The approver account holds the authority to approve recover requests and will be specified in the system contract on BSC. The community team retains the ability to change the approver account as needed.

## Build Binary

```bash
make build
```

## Build Docker Image

```bash
make build-image
```

## Run Locally

```bash
make run
```

## How To Get Approval
```bash
curl -X 'POST' http://localhost:8080/approve -d '{"token_symbol": "BNB","owner_pub_key": "0x02dcd743516b78366a217a1bf2aa562ec5accd07163db3332d924fa48e643875a6","owner_signature": "0xcd32af98a3cf4b66deaba53dc81c7cf8c810a83eb2fa23bf1a555a718826e2f03d47e3711a10e1ae72031fcd3faabac51325999d74c0cff31d554b4d657dbc64","claim_address": "0x5b38da6a701c568545dcfcb03fcb875f56beddc4"}'
```

## Configuration

| Name | Env | Type | Option | Description | Default |
|------|-----|------|--------|-------------|---------|
| chain_id | CHAIN_ID | string | | Chain ID | `"Binance-Chain-Ganges"` |
| merkle_root | MERKLE_ROOT | string | | Merkle root | `"0x0000000000000000000000000000000000000000000000000000000000000000"` |
| account_white_list | ACCOUNT_WHITE_LIST | []string | | Account white list | `[]` |
|---|---|---|---|---|---|
| logger.level | LOGGER_LEVEL | string | `ERROR`, `WARN`, `INFO`, `DEBUG`, `TRACE` | Log level | `"INFO"` |
| logger.format | LOGGER_FORMAT | string | `console`, `json` | Log format | `console` |
|---|---|---|---|---|---|
| http.addr | HTTP_ADDR | string | | HTTP address | `"0.0.0.0"` |
| http.port | HTTP_PORT | uint16 | | HTTP port | `8080` |
|---|---|---|---|---|---|
| metrics.enable | METRICS_ENABLE | bool | | Whether to enable metrics router | `true` |
| metrics.pprof | METRICS_PPROF | bool | | Whether to enable pprof router | `false` |
| metrics.path | METRICS_PATH | string | | Metrics router path | `/metrics` |
| metrics.addr | METRICS_ADDR | string | | Metrics address | `"0.0.0.0"` |
| metrics.port | METRICS_PORT | uint16 | | Metrics port | `6060` |
|---|---|---|---|---|---|
| secret.type | SECRET_TYPE | string | | Secret type | `"local"` |
| secret.local_secret.private_key | SECRET_LOCAL_SECRET_PRIVATE_KEY | string | | Local secret private key | `""` |
| secret.aws_secret_manager.region | SECRET_AWS_SECRET_MANAGER_REGION | string | | AWS Secret Manager region | `""` |
| secret.aws_secret_manager.secret_name | SECRET_AWS_SECRET_MANAGER_SECRET_NAME | string | | AWS Secret Manager secret name | `""` |
|---|---|---|---|---|---|
| store.driver | STORE_DRIVER | string | | Store driver | `"memory"` |
|---|---|---|---|---|---|
| store.memory_store.merkle_proofs | STORE_MEMORY_STORE_MERKLE_PROOFS | string | | Memory store Merkle proofs file | `"./example/merkle_proofs.json"` |
|---|---|---|---|---|---|
| store.sql_store.sql_driver | STORE_SQL_STORE_SQL_DRIVER | string | `mysql`, `postgres`, `sqlite` | SQL store driver | `"mysql"` |
| store.sql_store.host | STORE_SQL_STORE_HOST | string | | SQL store host | `"localhost"` |
| store.sql_store.port | STORE_SQL_STORE_PORT | uint | | SQL store port | `3306` |
| store.sql_store.dbname | STORE_SQL_STORE_DBNAME | string | | SQL store database name | `"approver"` |
| store.sql_store.user | STORE_SQL_STORE_USER | string | | SQL store user | `"root"` |
| store.sql_store.password | STORE_SQL_STORE_PASSWORD | string | | SQL store password | `""` |
| store.sql_store.connect_timeout | STORE_SQL_STORE_CONNECT_TIMEOUT | string | | SQL store connect timeout | `"10s"` |
| store.sql_store.read_timeout | STORE_SQL_STORE_READ_TIMEOUT | string | | SQL store read timeout | `"30s"` |
| store.sql_store.write_timeout | STORE_SQL_STORE_WRITE_TIMEOUT | string | | SQL store write timeout | `"30s"` |
| store.sql_store.dial_timeout | STORE_SQL_STORE_DIAL_TIMEOUT | time.Duration | | SQL store dial timeout | `"10s"` |
| store.sql_store.max_idletime | STORE_SQL_STORE_MAX_IDLETIME | time.Duration | | SQL store max idle time | `"1h"` |
| store.sql_store.max_lifetime | STORE_SQL_STORE_MAX_LIFETIME | time.Duration | | SQL store max lifetime | `"1h"` |
| store.sql_store.max_idle_conn | STORE_SQL_STORE_MAX_IDLE_CONN | int | | SQL store max idle connections | `2` |
| store.sql_store.max_open_conn | STORE_SQL_STORE_MAX_OPEN_CONN | int | | SQL store max open connections | `5` |
| store.sql_store.ssl_mode | STORE_SQL_STORE_SSL_MODE | bool | | SQL store SSL mode | `false` |