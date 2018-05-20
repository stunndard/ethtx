# ethtx - a microservice to generate Ethereum addressess and sign Ethereum transactions

## What's the point?
`ethtx` is a small, stateless and portable microservice for generating Ethereum addresses and signing Ethereum transactions without the dedicated Ethereum node.
It comes extremely handy when you need to create Ethereum addresses (private/public key pair) or sign Ethereum transactions using a programming language that doesn't have the required crypto libraries.

Also it improves the security, because you don't need to call the remote Ethereum RPC node and send your private keys to it in order to just generate a new Ethereum address or sign your transaction.

That way the private keys won't ever leave the local host where your service or software is running.

`ethtx` consumes minimum of system resources and can be run on the local host together with your service or even desktop software.

## How it works?
It works as a HTTP service. You POST parameters and get a JSON response.

## What crypto libraries does ethtx use? Is it safe?
`ethtx` is using `geth`'s implentation directly by importing the corresponding source packages directly from the [go-ethereum](https://github.com/ethereum/go-ethereum) project.

## How it generates the private keys? Any implementation of such crucial part should be seriously reviewed.
`ethtx` derives key generation directly from [go-ethereum](https://github.com/ethereum/go-ethereum) sources. It doesn't add anything else to it. Please check the sources.

## Tell me more.
So far the following methods are supported:

### `generateAddress`
Generates a new Ethereum address (private/public keypair):

##### Input parameters:
`none`

##### Curl example:

```
curl -X POST \
  http://ethtx-address:8070/generateAddress
```

##### Response:
```JSON
{
    "address": "0x2cBD5C1D45DCcCD5147BED0314379c6b50c3e8a5",
    "privKey": "0xc3c4a14597d9ca779f730df5716c58321e46bc55d6559d8ba2932f0da04e1dbe",
    "pubKey": "0x7cbd2cd558251149285fc379fecbf030a92be3b6f2c793ae5c3a7dfefe4076da, 0xff5365976d4c8cddb531bf9dcb19963d8e03bebfcac713bca7017886ef5be7ab"
}
```

##### Output:

`address` : a new freshly generated Ethereum address.

`privKey` : the private key for this address. Be careful with it and NEVER save it unencrypted.

`pubKey`  : the public key for this address. Consists of two parts, if you need the public key, just concatenate this two parts.

### `signTx`
Signs a transaction and returns it in RLP coding, ready to be sent to the Ethereum network via a full node, or a 3rd party service like [Infura](https://infura.io/) or [Etherscan](https://etherscan.io/apis)

##### Input parameters:

`chainId` : (integer) your network chain ID. Use 1 for the production network and any other custom one for your own private networks.

`privKey` : (string), the private key used to sign this transaction. Do not use `0x` prefix for the private key.

`sendTo`  : (string), the Ethereum address to send this transaction to. Use `0x` prefix for the destination address.

`amount`  : (float), the amount of `ETH` to send. For example, 23.33323

`amountWei` : (integer), if you don't want to pass the amount in `ETH`, put `0` in `amount` above, and use the amount in `wei` here. For example, 30000000000000000

`gasLimit` : (integer), gas limit for this transaction, for example, 21000

`gasPrice` : (integer), gas price in `wei` for this transaction, for example, 20000000000

`nonce` : (integer), nonce value for this transaction. Use [eth_getTransactionCount](https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_gettransactioncount) to get the nonce before signing this transaction. Without the valid nonce, your transaction won't be relayed to the network and will be stuck.

##### Curl example:

```
curl -X POST \
  http://ethtx-address:8070/signTx \
  -H 'Content-Type: application/x-www-form-urlencoded' \
  -d 'chainId=88882&privKey=c3c4a14597d9ca779f730df5716c58321e46bc55d6559d8ba2932f0da04e1dbe&sendTo=0x3822e05392f097cf03ba3227cc16b73ceb0ee305&amount=0.1&gasLimit=21000&gasPrice=20000000000&nonce=1'
```

##### Output:

`result` : if it is `ok` then `signedTx` contains the signed transaction.

`signedTx` : signed transaction ready to be sent to the Ethereum network by any available means.

##### Response:
```JSON
{
    "result": "ok",
    "signedTx": "0xf86f018504a817c800825208943822e05392f097cf03ba3227cc16b73ceb0ee30588016345785d8a000080830300afa0ab96c871bf77e8f2a67b744fbec7138a45f7a706a3cdd758989f3b49da30299ba05eb9343eb000e620e62c31769d6a5e5ac8c191f3afa2f9c934d541520e4340f9"
}
```

## What platforms are supported?
Any that Go can compile for.

## What is required?
Just Go.

## How do I install ethtx?
The `go get` command will automatically fetch all dependencies required, compile the binary and place it in your $GOPATH/bin directory.

    go get github.com/stunndard/ethtx

## How do I configure it?
Read `ethtx.yml `. Tune it for your needs.
```YAML
---
  Gzip: true
  Other:
    listen: ":8070"
    redactLogs: True

```

`redactLogs` will hide the private keys from the logs. By default `ethtx` will write the logs to the standard output. You can hide that or redirect to files, as usual.

## How do I run it?
Just run the binary
```
    ./ethtx
```

You can run it inside `tmux` or `screen` or you can write a systemd config for it to run as daemon.