# bifrost-relayer

![bifrost-relayer](./bifrost_relayer.gif)

This is a bifrost relayer that transmits that continuously watches tezos contract storage and relays it to bifrost zone.

# Installation

```
$ git clone github.com/sap200/bifrost-relayer

```

```
$ go install 
```

```
$ bifrost-relayer
```

This will start the relayer 

# Important ports

7009 - verification engine

7010 - Operation Engine

# Operations

```
mint
```

Mint mints a new FA12 Token using tezos-client

```
burn
```

burn unlocks the locked tezos in bifrost contract

```
verify
```

verify verifies the storage and checks if the corresponding txs exists in the bifrost contract storage

This happens automatically if tezos-client is configured, and bifrost zone is running.

### This for now works in development mode



