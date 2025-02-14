// Copyright 2018 The go-ethereum Authors
// This file is part of go-puffscoin.
//
// go-puffscoin is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// go-puffscoin is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with go-puffscoin. If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/puffscoin/go-puffscoin/log"
	"github.com/puffscoin/go-puffscoin/node"
	"github.com/puffscoin/go-puffscoin/p2p/simulations"
	"github.com/puffscoin/go-puffscoin/p2p/simulations/adapters"
	"github.com/puffscoin/go-puffscoin/swarm/network"
	"github.com/puffscoin/go-puffscoin/swarm/network/simulation"
	cli "gopkg.in/urfave/cli.v1"
)

// create is used as the entry function for "create" app command.
func create(ctx *cli.Context) error {
	log.PrintOrigins(true)
	log.Root().SetHandler(log.LvlFilterHandler(log.Lvl(ctx.Int("verbosity")), log.StreamHandler(os.Stdout, log.TerminalFormat(true))))

	if len(ctx.Args()) < 1 {
		return errors.New("argument should be the filename to verify or write-to")
	}
	filename, err := touchPath(ctx.Args()[0])
	if err != nil {
		return err
	}
	return createSnapshot(filename, ctx.Int("nodes"), strings.Split(ctx.String("services"), ","))
}

// createSnapshot creates a new snapshot on filesystem with provided filename,
// number of nodes and service names.
func createSnapshot(filename string, nodes int, services []string) (err error) {
	log.Debug("create snapshot", "filename", filename, "nodes", nodes, "services", services)

	sim := simulation.New(map[string]simulation.ServiceFunc{
		"bzz": func(ctx *adapters.ServiceContext, bucket *sync.Map) (node.Service, func(), error) {
			addr := network.NewAddr(ctx.Config.Node())
			kad := network.NewKademlia(addr.Over(), network.NewKadParams())
			hp := network.NewHiveParams()
			hp.KeepAliveInterval = time.Duration(200) * time.Millisecond
			hp.Discovery = true // discovery must be enabled when creating a snapshot

			// store the kademlia in the bucket, needed later in the WaitTillHealthy function
			bucket.Store(simulation.BucketKeyKademlia, kad)

			config := &network.BzzConfig{
				OverlayAddr:  addr.Over(),
				UnderlayAddr: addr.Under(),
				HiveParams:   hp,
			}
			return network.NewBzz(config, kad, nil, nil, nil), nil, nil
		},
	})
	defer sim.Close()

	ids, err := sim.AddNodes(nodes)
	if err != nil {
		return fmt.Errorf("add nodes: %v", err)
	}

	err = sim.Net.ConnectNodesRing(ids)
	if err != nil {
		return fmt.Errorf("connect nodes: %v", err)
	}

	ctx, cancelSimRun := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancelSimRun()
	if _, err := sim.WaitTillHealthy(ctx); err != nil {
		return fmt.Errorf("wait for healthy kademlia: %v", err)
	}

	var snap *simulations.Snapshot
	if len(services) > 0 {
		// If service names are provided, include them in the snapshot.
		// But, check if "bzz" service is not among them to remove it
		// form the snapshot as it exists on snapshot creation.
		var removeServices []string
		var wantBzz bool
		for _, s := range services {
			if s == "bzz" {
				wantBzz = true
				break
			}
		}
		if !wantBzz {
			removeServices = []string{"bzz"}
		}
		snap, err = sim.Net.SnapshotWithServices(services, removeServices)
	} else {
		snap, err = sim.Net.Snapshot()
	}
	if err != nil {
		return fmt.Errorf("create snapshot: %v", err)
	}
	jsonsnapshot, err := json.Marshal(snap)
	if err != nil {
		return fmt.Errorf("json encode snapshot: %v", err)
	}
	return ioutil.WriteFile(filename, jsonsnapshot, 0666)
}

// touchPath creates an empty file and all subdirectories
// that are missing.
func touchPath(filename string) (string, error) {
	if path.IsAbs(filename) {
		if _, err := os.Stat(filename); err == nil {
			// path exists, overwrite
			return filename, nil
		}
	}

	d, f := path.Split(filename)
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return "", err
	}

	_, err = os.Stat(path.Join(dir, filename))
	if err == nil {
		// path exists, overwrite
		return filename, nil
	}

	dirPath := path.Join(dir, d)
	filePath := path.Join(dirPath, f)
	if d != "" {
		err = os.MkdirAll(dirPath, os.ModeDir)
		if err != nil {
			return "", err
		}
	}

	return filePath, nil
}
