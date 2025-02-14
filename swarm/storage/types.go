// Copyright 2016 The go-ethereum Authors
// This file is part of the go-puffscoin library.
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

package storage

import (
	"bytes"
	"context"
	"crypto"
	"crypto/rand"
	"encoding/binary"
	"io"

	"github.com/puffscoin/go-puffscoin/swarm/bmt"
	"github.com/puffscoin/go-puffscoin/swarm/chunk"
	"golang.org/x/crypto/sha3"
)

// MaxPO is the same as chunk.MaxPO for backward compatibility.
const MaxPO = chunk.MaxPO

// AddressLength is the same as chunk.AddressLength for backward compatibility.
const AddressLength = chunk.AddressLength

type SwarmHasher func() SwarmHash

// Address is an alias for chunk.Address for backward compatibility.
type Address = chunk.Address

// Proximity is the same as chunk.Proximity for backward compatibility.
var Proximity = chunk.Proximity

// ZeroAddr is the same as chunk.ZeroAddr for backward compatibility.
var ZeroAddr = chunk.ZeroAddr

func MakeHashFunc(hash string) SwarmHasher {
	switch hash {
	case "SHA256":
		return func() SwarmHash { return &HashWithLength{crypto.SHA256.New()} }
	case "SHA3":
		return func() SwarmHash { return &HashWithLength{sha3.NewLegacyKeccak256()} }
	case "BMT":
		return func() SwarmHash {
			hasher := sha3.NewLegacyKeccak256
			hasherSize := hasher().Size()
			segmentCount := chunk.DefaultSize / hasherSize
			pool := bmt.NewTreePool(hasher, segmentCount, bmt.PoolSize)
			return bmt.New(pool)
		}
	}
	return nil
}

type AddressCollection []Address

func NewAddressCollection(l int) AddressCollection {
	return make(AddressCollection, l)
}

func (c AddressCollection) Len() int {
	return len(c)
}

func (c AddressCollection) Less(i, j int) bool {
	return bytes.Compare(c[i], c[j]) == -1
}

func (c AddressCollection) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

// Chunk is an alias for chunk.Chunk for backward compatibility.
type Chunk = chunk.Chunk

// NewChunk is the same as chunk.NewChunk for backward compatibility.
var NewChunk = chunk.NewChunk

func GenerateRandomChunk(dataSize int64) Chunk {
	hasher := MakeHashFunc(DefaultHash)()
	sdata := make([]byte, dataSize+8)
	rand.Read(sdata[8:])
	binary.LittleEndian.PutUint64(sdata[:8], uint64(dataSize))
	hasher.ResetWithLength(sdata[:8])
	hasher.Write(sdata[8:])
	return NewChunk(hasher.Sum(nil), sdata)
}

func GenerateRandomChunks(dataSize int64, count int) (chunks []Chunk) {
	for i := 0; i < count; i++ {
		ch := GenerateRandomChunk(dataSize)
		chunks = append(chunks, ch)
	}
	return chunks
}

// Size, Seek, Read, ReadAt
type LazySectionReader interface {
	Context() context.Context
	Size(context.Context, chan bool) (int64, error)
	io.Seeker
	io.Reader
	io.ReaderAt
}

type LazyTestSectionReader struct {
	*io.SectionReader
}

func (r *LazyTestSectionReader) Size(context.Context, chan bool) (int64, error) {
	return r.SectionReader.Size(), nil
}

func (r *LazyTestSectionReader) Context() context.Context {
	return context.TODO()
}

type StoreParams struct {
	Hash          SwarmHasher `toml:"-"`
	DbCapacity    uint64
	CacheCapacity uint
	BaseKey       []byte
}

func NewDefaultStoreParams() *StoreParams {
	return NewStoreParams(defaultLDBCapacity, defaultCacheCapacity, nil, nil)
}

func NewStoreParams(ldbCap uint64, cacheCap uint, hash SwarmHasher, basekey []byte) *StoreParams {
	if basekey == nil {
		basekey = make([]byte, 32)
	}
	if hash == nil {
		hash = MakeHashFunc(DefaultHash)
	}
	return &StoreParams{
		Hash:          hash,
		DbCapacity:    ldbCap,
		CacheCapacity: cacheCap,
		BaseKey:       basekey,
	}
}

type ChunkData []byte

type Reference []byte

// Putter is responsible to store data and create a reference for it
type Putter interface {
	Put(context.Context, ChunkData) (Reference, error)
	// RefSize returns the length of the Reference created by this Putter
	RefSize() int64
	// Close is to indicate that no more chunk data will be Put on this Putter
	Close()
	// Wait returns if all data has been store and the Close() was called.
	Wait(context.Context) error
}

// Getter is an interface to retrieve a chunk's data by its reference
type Getter interface {
	Get(context.Context, Reference) (ChunkData, error)
}

// NOTE: this returns invalid data if chunk is encrypted
func (c ChunkData) Size() uint64 {
	return binary.LittleEndian.Uint64(c[:8])
}

type ChunkValidator interface {
	Validate(chunk Chunk) bool
}

// Provides method for validation of content address in chunks
// Holds the corresponding hasher to create the address
type ContentAddressValidator struct {
	Hasher SwarmHasher
}

// Constructor
func NewContentAddressValidator(hasher SwarmHasher) *ContentAddressValidator {
	return &ContentAddressValidator{
		Hasher: hasher,
	}
}

// Validate that the given key is a valid content address for the given data
func (v *ContentAddressValidator) Validate(ch Chunk) bool {
	data := ch.Data()
	if l := len(data); l < 9 || l > chunk.DefaultSize+8 {
		// log.Error("invalid chunk size", "chunk", addr.Hex(), "size", l)
		return false
	}

	hasher := v.Hasher()
	hasher.ResetWithLength(data[:8])
	hasher.Write(data[8:])
	hash := hasher.Sum(nil)

	return bytes.Equal(hash, ch.Address())
}

type ChunkStore interface {
	Put(ctx context.Context, ch Chunk) (err error)
	Get(rctx context.Context, ref Address) (ch Chunk, err error)
	Has(rctx context.Context, ref Address) bool
	Close()
}

// SyncChunkStore is a ChunkStore which supports syncing
type SyncChunkStore interface {
	ChunkStore
	BinIndex(po uint8) uint64
	Iterator(from uint64, to uint64, po uint8, f func(Address, uint64) bool) error
	FetchFunc(ctx context.Context, ref Address) func(context.Context) error
}

// FakeChunkStore doesn't store anything, just implements the ChunkStore interface
// It can be used to inject into a hasherStore if you don't want to actually store data just do the
// hashing
type FakeChunkStore struct {
}

// Put doesn't store anything it is just here to implement ChunkStore
func (f *FakeChunkStore) Put(_ context.Context, ch Chunk) error {
	return nil
}

// Has doesn't do anything it is just here to implement ChunkStore
func (f *FakeChunkStore) Has(_ context.Context, ref Address) bool {
	panic("FakeChunkStore doesn't support HasChunk")
}

// Get doesn't store anything it is just here to implement ChunkStore
func (f *FakeChunkStore) Get(_ context.Context, ref Address) (Chunk, error) {
	panic("FakeChunkStore doesn't support Get")
}

// Close doesn't store anything it is just here to implement ChunkStore
func (f *FakeChunkStore) Close() {
}
