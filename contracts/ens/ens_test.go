// Copyright 2016 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package ens

import (
	"math/big"
	"testing"

	"github.com/puffscoin/go-puffscoin/accounts/abi/bind"
	"github.com/puffscoin/go-puffscoin/accounts/abi/bind/backends"
	"github.com/puffscoin/go-puffscoin/common"
	"github.com/puffscoin/go-puffscoin/contracts/ens/contract"
	"github.com/puffscoin/go-puffscoin/contracts/ens/fallback_contract"
	"github.com/puffscoin/go-puffscoin/core"
	"github.com/puffscoin/go-puffscoin/crypto"
)

var (
	key, _       = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
	name         = "my name on ENS"
	hash         = crypto.Keccak256Hash([]byte("my content"))
	fallbackHash = crypto.Keccak256Hash([]byte("my content hash"))
	addr         = crypto.PubkeyToAddress(key.PublicKey)
	testAddr     = common.HexToAddress("0x1234123412341234123412341234123412341234")
)

func TestENS(t *testing.T) {
	contractBackend := backends.NewSimulatedBackend(core.GenesisAlloc{addr: {Balance: big.NewInt(1000000000)}}, 10000000)
	transactOpts := bind.NewKeyedTransactor(key)

	ensAddr, ens, err := DeployENS(transactOpts, contractBackend)
	if err != nil {
		t.Fatalf("can't deploy root registry: %v", err)
	}
	contractBackend.Commit()

	// Set ourself as the owner of the name.
	if _, err := ens.Register(name); err != nil {
		t.Fatalf("can't register: %v", err)
	}
	contractBackend.Commit()

	// Deploy a resolver and make it responsible for the name.
	resolverAddr, _, _, err := contract.DeployPublicResolver(transactOpts, contractBackend, ensAddr)
	if err != nil {
		t.Fatalf("can't deploy resolver: %v", err)
	}

	if _, err := ens.SetResolver(EnsNode(name), resolverAddr); err != nil {
		t.Fatalf("can't set resolver: %v", err)
	}
	contractBackend.Commit()

	// Set the content hash for the name.
	cid, err := EncodeSwarmHash(hash)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = ens.SetContentHash(name, cid); err != nil {
		t.Fatalf("can't set content hash: %v", err)
	}
	contractBackend.Commit()

	// Try to resolve the name.
	resolvedHash, err := ens.Resolve(name)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resolvedHash.Hex() != hash.Hex() {
		t.Fatalf("resolve error, expected %v, got %v", hash.Hex(), resolvedHash.Hex())
	}

	// set the address for the name
	if _, err = ens.SetAddr(name, testAddr); err != nil {
		t.Fatalf("can't set address: %v", err)
	}
	contractBackend.Commit()

	// Try to resolve the name to an address
	recoveredAddr, err := ens.Addr(name)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if testAddr.Hex() != recoveredAddr.Hex() {
		t.Fatalf("resolve error, expected %v, got %v", testAddr.Hex(), recoveredAddr.Hex())
	}

	// deploy the fallback contract and see that the fallback mechanism works
	fallbackResolverAddr, _, _, err := fallback_contract.DeployPublicResolver(transactOpts, contractBackend, ensAddr)
	if err != nil {
		t.Fatalf("can't deploy resolver: %v", err)
	}
	if _, err := ens.SetResolver(EnsNode(name), fallbackResolverAddr); err != nil {
		t.Fatalf("can't set resolver: %v", err)
	}
	contractBackend.Commit()

	// Set the content hash for the name.
	if _, err = ens.SetContentHash(name, fallbackHash.Bytes()); err != nil {
		t.Fatalf("can't set content hash: %v", err)
	}
	contractBackend.Commit()

	// Try to resolve the name.
	fallbackResolvedHash, err := ens.Resolve(name)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if fallbackResolvedHash.Hex() != fallbackHash.Hex() {
		t.Fatalf("resolve error, expected %v, got %v", hash.Hex(), resolvedHash.Hex())
	}
}
