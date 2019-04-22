## go-PUFFScoin

Official Golang implementation of the PUFFScoin protocol.

[![API Reference]](https://godoc.org/github.com/puffscoin/go-puffscoin)
[![Go Report Card](https://goreportcard.com/badge/github.com/puffscoin/go-puffscoin)](https://goreportcard.com/report/github.com/puffscoin/go-puffscoin)
[![Travis](https://travis-ci.org/puffscoin/go-puffscoin.svg?branch=master)](https://travis-ci.org/puffscoin/go-puffscoin)
[![Discord](https://img.shields.io/badge/discord-join%20chat-blue.svg)](https://discord.gg/A5nNsZF)

Automated builds are available for stable releases and the unstable master branch.
Links to binary archives are published at http://puffscoin.leafycauldronapothecary.com/downloads/.

## Building the source

For prerequisites and detailed build instructions please read the
[Installation Instructions](http://puffscoin.leafycauldronapothecary.com/puffwiki/the-basics/building-puffscoin/)
on the puffswiki.

Building gpuffs requires both a Go (version 1.10 or later) and a C compiler.
You can install them using your favourite package manager.
Once the dependencies are installed, run

    make gpuffs

or, to build the full suite of utilities:

    make all

## Executables

The go-puffscoin project comes with several wrappers/executables found in the `cmd` directory.

| Command    | Description |
|:----------:|-------------|
| **`gpuffs`** | The main PUFFScoin CLI client. It is the entry point into the PUFFScoin network, capable of running as a full node (default), archive node (retaining all historical state) or a light node (retrieving data live). It can be used by other processes as a gateway into the PUFFScoin network via JSON RPC endpoints exposed on top of HTTP, WebSocket and/or IPC transports. `gpuffs --help` and the [CLI puffswiki page](http://puffscoin.leafycauldronapothecary.com/puffwiki/the-basics/command-line-options/) for command line options. |
| `abigen` | Source code generator to convert PUFFScoin contract definitions into easy to use, compile-time type-safe Go packages. It operates on plain [Solidity contract ABIs](https://solidity.readthedocs.io/en/develop/abi-spec.html) with expanded functionality if the contract bytecode is also available. However, it also accepts Solidity source files, making development much more streamlined. Please see our [Native DApps](http://puffscoin.leafycauldronapothecary.com/native-dapps-go-bindings-to-puffscoin-contracts/) puffswiki entry for details. |
| `bootnode` | Stripped down version of our PUFFScoin client implementation that only takes part in the network node discovery protocol, but does not run any of the higher level application protocols. It can be used as a lightweight bootstrap node to aid in finding peers in private networks. |
| `evm` | Developer utility version of the PUFFScoin EVM (Ethereum Virtual Machine) that is capable of running bytecode snippets within a configurable environment and execution mode. Its purpose is to allow isolated, fine-grained debugging of EVM opcodes (e.g. `evm --code 60ff60ff --debug`). |
| `gethrpctest` | Developer utility tool to support our [ethereum/rpc-test](https://github.com/puffscoin/rpc-tests) test suite which validates baseline conformity to the [Ethereum JSON RPC](http://puffscoin.leafycauldronapothecary.com/puffwiki/blockchain-protocols/json-rpc-api/) specs.  |
| `rlpdump` | Developer utility tool to convert binary RLP ([Recursive Length Prefix](https://github.com/ethereum/wiki/wiki/RLP)) dumps (data encoding used by the PUFFScoin protocol both network as well as consensus wise) to user-friendlier hierarchical representation (e.g. `rlpdump --hex CE0183FFFFFFC4C304050583616263`). |
| `swarm`    | Swarm daemon and tools. This is the entry point for the Swarm network. `swarm --help` for command line options and subcommands. See [Swarm](http://puffscoin.leafycauldronapothecary.com/services/swarm/) for more information. |
| `puppeth`    | a CLI wizard that aids in creating a new PUFFScoin-compliant network. |

## Running gpuffs

Going through all the possible command line flags is out of scope here (please consult our
[CLI Wiki page](http://puffscoin.leafycauldronapothecary.com/puffwiki/the-basics/command-line-options/)), but we've
enumerated a few common parameter combos to get you up to speed quickly on how you can run your
own gpuffs instance.

### Full node on the main PUFFScoin network

By far the most common scenario is people wanting to simply interact with the PUFFScoin network:
create accounts; transfer funds; deploy and interact with contracts. For this particular use-case
the user doesn't care about years-old historical data, so we can fast-sync quickly to the current
state of the network. To do so:

```
$ gpuffs console
```

This command will:

 * Start gpuffs in fast sync mode (default, can be changed with the `--syncmode` flag), causing it to
   download more data in exchange for avoiding processing the entire history of the PUFFScoin network,
   which may become very CPU intensive as he blockchain grows, as evidenced form Ethereum, Bitcoin and
   other projects with blockchains well in excess of 1GB.
 * Start up gpuffs' built-in interactive [JavaScript console](http://puffscoin.leafycauldronapothecary.com/javascript-console/),
   (via the trailing `console` subcommand) through which you can invoke all official [`web3` methods](http://puffscoin.leafycauldronapothecary.com/javascript-api/)
   as well as gpuffs' own [management APIs](http://puffscoin.leafycauldronapothecary.com/puffwiki/blockchain-protocols/management-apis/).
   This tool is optional and if you leave it out you can always attach to an already running gpuffs instance
   with `gpuffs attach`.



### Configuration

As an alternative to passing the numerous flags to the `gpuffs` binary, you can also pass a configuration file via:

```
$ gpuffs --config /path/to/your_config.toml
```

To get an idea how the file should look like you can use the `dumpconfig` subcommand to export your existing configuration:

```
$ gpuffs --your-favourite-flags dumpconfig
```


#### Docker quick start

One of the quickest ways to get PUFFScoin up and running on your machine is by using Docker:

```
docker run -d --name puffscoin-node -v /Users/alice/puffscoin:/root \
           -p 11363:11363 -p 31313:31313 \
           puffscoin/client-go
```

This will start gpuffs in fast-sync mode with a DB memory allowance of 1GB just as the above command does.  It will also create a persistent volume in your home directory for saving your blockchain as well as map the default ports. There is also an `alpine` tag available for a slim version of the image.

Do not forget `--rpcaddr 0.0.0.0`, if you want to access RPC from other containers and/or hosts. By default, `gpuffs` binds to the local interface and RPC endpoints is not accessible from the outside.

### Programmatically interfacing gpuffs nodes

As a developer, sooner rather than later you'll want to start interacting with gpuffs and the PUFFScoin
network via your own programs and not manually through the console. To aid this, gpuffs has built-in
support for a JSON-RPC based APIs ([standard APIs](http://puffscoin.leafycauldronapothecary.com/puffwiki/blockchain-protocols/json-rpc-api/) and
[Geth specific APIs](http://puffscoin.leafycauldronapothecary.com/puffwiki/blockchain-protocols/management-apis/)). These can be
exposed via HTTP, WebSockets and IPC (UNIX sockets on UNIX based platforms, and named pipes on Windows).

The IPC interface is enabled by default and exposes all the APIs supported by gpuffs, whereas the HTTP
and WS interfaces need to manually be enabled and only expose a subset of APIs due to security reasons.
These can be turned on/off and configured as you'd expect.

HTTP based JSON-RPC API options:

  * `--rpc` Enable the HTTP-RPC server
  * `--rpcaddr` HTTP-RPC server listening interface (default: "localhost")
  * `--rpcport` HTTP-RPC server listening port (default: 8545)
  * `--rpcapi` API's offered over the HTTP-RPC interface (default: "eth,net,web3")
  * `--rpccorsdomain` Comma separated list of domains from which to accept cross origin requests (browser enforced)
  * `--ws` Enable the WS-RPC server
  * `--wsaddr` WS-RPC server listening interface (default: "localhost")
  * `--wsport` WS-RPC server listening port (default: 11363)
  * `--wsapi` API's offered over the WS-RPC interface (default: "eth,net,web3")
  * `--wsorigins` Origins from which to accept websockets requests
  * `--ipcdisable` Disable the IPC-RPC server
  * `--ipcapi` API's offered over the IPC-RPC interface (default: "admin,debug,eth,miner,net,personal,shh,txpool,web3")
  * `--ipcpath` Filename for IPC socket/pipe within the datadir (explicit paths escape it)

You'll need to use your own programming environments' capabilities (libraries, tools, etc) to connect
via HTTP, WS or IPC to a Geth node configured with the above flags and you'll need to speak [JSON-RPC](https://www.jsonrpc.org/specification)
on all transports. You can reuse the same connection for multiple requests!

**Note: Please understand the security implications of opening up an HTTP/WS based transport before
doing so! Hackers omay try to subvert PUFFScoin nodes with exposed APIs!
Further, all browser tabs can access locally running web servers, so malicious web pages could try to
subvert locally available APIs!**

### Operating a private network

Maintaining your own private network is more involved as a lot of configurations taken for granted in
the official networks need to be manually set up.

#### Defining the private genesis state

First, you'll need to create the genesis state of your networks, which all nodes need to be aware of
and agree upon. This consists of a small JSON file (e.g. call it `genesis.json`):

```json
{
  "config": {
        "chainId": 0,
        "homesteadBlock": 0,
        "eip155Block": 0,
        "eip158Block": 0
    },
  "alloc"      : {},
  "coinbase"   : "0x0000000000000000000000000000000000000000",
  "difficulty" : "0x20000",
  "extraData"  : "",
  "gasLimit"   : "0x2fefd8",
  "nonce"      : "0x0000000000000042",
  "mixhash"    : "0x0000000000000000000000000000000000000000000000000000000000000000",
  "parentHash" : "0x0000000000000000000000000000000000000000000000000000000000000000",
  "timestamp"  : "0x00"
}
```

The above fields should be fine for most purposes, although we'd recommend changing the `nonce` to
some random value so you prevent unknown remote nodes from being able to connect to you. If you'd
like to pre-fund some accounts for easier testing, you can populate the `alloc` field with account
configs:

```json
"alloc": {
  "0x0000000000000000000000000000000000000001": {"balance": "111111111"},
  "0x0000000000000000000000000000000000000002": {"balance": "222222222"}
}
```

With the genesis state defined in the above JSON file, you'll need to initialize **every** Geth node
with it prior to starting it up to ensure all blockchain parameters are correctly set:

```
$ geth init path/to/genesis.json
```

#### Creating the rendezvous point

With all nodes that you want to run initialized to the desired genesis state, you'll need to start a
bootstrap node that others can use to find each other in your network and/or over the internet. The
clean way is to configure and run a dedicated bootnode:

```
$ bootnode --genkey=boot.key
$ bootnode --nodekey=boot.key
```

With the bootnode online, it will display an [`enode` URL](https://github.com/ethereum/wiki/wiki/enode-url-format)
that other nodes can use to connect to it and exchange peer information. Make sure to replace the
displayed IP address information (most probably `[::]`) with your externally accessible IP to get the
actual `enode` URL.

*Note: You could also use a full-fledged gpuffs node as a bootnode, but it's the less recommended way.*

#### Starting up your member nodes

With the bootnode operational and externally reachable (you can try `telnet <ip> <port>` to ensure
it's indeed reachable), start every subsequent gpuffs node pointed to the bootnode for peer discovery
via the `--bootnodes` flag. It will probably also be desirable to keep the data directory of your
private network separated, so do also specify a custom `--datadir` flag.

```
$ geth --datadir=path/to/custom/data/folder --bootnodes=<bootnode-enode-url-from-above>
```

*Note: Since your network will be completely cut off from the main and test networks, you'll also
need to configure a miner to process transactions and create new blocks for you.*

#### Running a private miner

Mining on the public PUFFScoin network is a complex task as it's only feasible using GPUs, requiring
an OpenCL or CUDA enabled `ethminer` instance. For information on such a setup, please consult the
[EtherMining subreddit](https://www.reddit.com/r/EtherMining/) and the [Genoil miner](https://github.com/Genoil/cpp-ethereum)
repository. 

In a private network setting, however a single CPU miner instance is more than enough for practical
purposes as it can produce a stable stream of blocks at the correct intervals without needing heavy
resources (consider running on a single thread, no need for multiple ones either). To start a Geth
instance for mining, run it with all your usual flags, extended by:

```
$ geth <usual-flags> --mine --minerthreads=1 --etherbase=0x0000000000000000000000000000000000000000
```

Which will start mining blocks and transactions on a single CPU thread, crediting all proceedings to
the account specified by `--etherbase`. You can further tune the mining by changing the default gas
limit blocks converge to (`--targetgaslimit`) and the price transactions are accepted at (`--gasprice`).



Please see the [Developers' Guide](https://github.com/ethereum/go-ethereum/wiki/Developers'-Guide)
for more details on configuring your environment, managing project dependencies, and testing procedures.

## License

The go-ethereum library (i.e. all code outside of the `cmd` directory) is licensed under the
[GNU Lesser General Public License v3.0](https://www.gnu.org/licenses/lgpl-3.0.en.html), also
included in our repository in the `COPYING.LESSER` file.

The go-ethereum binaries (i.e. all code inside of the `cmd` directory) is licensed under the
[GNU General Public License v3.0](https://www.gnu.org/licenses/gpl-3.0.en.html), also included
in our repository in the `COPYING` file.

go-puffsconi is a project initiated and maintained by the PUFFScoin Core Development Team. PUFFScoin is a subsidiary service of The Leafy Cauldron Apothecary, LLC. (Nova Scotia)
