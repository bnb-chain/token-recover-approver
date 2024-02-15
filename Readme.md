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

## Migrate Local Data to SQL Store

```bash
## sqlite
./build/bin/approver tool migration-from-local-to-sql --config ./configs/sqlite.config.yaml --proof_path ./example/store/merkle_proofs.json
## mysql
./build/bin/approver tool migration-from-local-to-sql --config ./configs/mysql.config.yaml --proof_path ./example/store/merkle_proofs.json
## pgsql
./build/bin/approver tool migration-from-local-to-sql --config ./configs/pgsql.config.yaml --proof_path ./example/store/merkle_proofs.json
```

## How To Get Approval
```bash
curl -X 'POST' http://localhost:8080/approve -d '{"token_symbol": "BNB","owner_pub_key": "0x036d5d41cd7da2e96d39bcbd0390bfed461a86382f7a2923436ff16c65cabc7719","owner_signature": "0x5f5391ba7f2b002b4746025f7e803a43e57a397ea66f3939d05302eb7851bbbc0773cda87aae0fbb1e2a29367b606209ed47dc5cba6d1a83f6b79cb70e56efdb","claim_address": "0x2e9247B67ae885a8dcfBf77Eb6d0e93A32bea24C"}'
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
| http.read_timeout | HTTP_READ_TIMEOUT | time.Duration | | HTTP read timeout | `"5s"` |
| http.read_header_timeout | HTTP_READ_HEADER_TIMEOUT | time.Duration | | HTTP read header timeout | `"5s"` |
| http.write_timeout | HTTP_WRITE_TIMEOUT | time.Duration | | HTTP write timeout | `"10s"` |
| http.idle_timeout | HTTP_IDLE_TIMEOUT | time.Duration | | HTTP idle timeout | `"5s"` |
| http.max_header_bytes | HTTP_MAX_HEADER_BYTES | int | | HTTP max header bytes | `1 << 20` |
|---|---|---|---|---|---|
| metrics.enable | METRICS_ENABLE | bool | | Whether to enable metrics router | `true` |
| metrics.pprof | METRICS_PPROF | bool | | Whether to enable pprof router | `false` |
| metrics.path | METRICS_PATH | string | | Metrics router path | `/metrics` |
| metrics.addr | METRICS_ADDR | string | | Metrics address | `"0.0.0.0"` |
| metrics.port | METRICS_PORT | uint16 | | Metrics port | `6060` |
| metrics.read_timeout | METRICS_READ_TIMEOUT | time.Duration | | Metrics read timeout | `"5s"` |
| metrics.read_header_timeout | METRICS_READ_HEADER_TIMEOUT | time.Duration | | Metrics read header timeout | `"5s"` |
| metrics.write_timeout | METRICS_WRITE_TIMEOUT | time.Duration | | Metrics write timeout | `"10s"` |
| metrics.idle_timeout | METRICS_IDLE_TIMEOUT | time.Duration | | Metrics idle timeout | `"5s"` |
| metrics.max_header_bytes | METRICS_MAX_HEADER_BYTES | int | | Metrics max header bytes | `1 << 20` |
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