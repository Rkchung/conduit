/*
package conduit

The primary Conduit server

Purpose: To keep track of responsive regional masters and send the address of an appropriate
  one when asked by the client/provider

Methods: Ping (gets ping from regional master and updates its aliveness), RequestRegionaMaster (gets request for regional master and sends back the address of the regional master)
*/

package main

import (
	"github.com/patrickmn/go-cache"
	"log"
	"net"
	"net/rpc"
	"time"
	"github/Conduit/common"
	"fmt"
)

// Conduit server connects with providers/clients and gives the IP of their regional master
type ProviderOrClient struct {
	c *ConduitServer
}

type RegionalMaster struct {
	rmAddr string
	c      *ConduitServer
}

// RequestRegionalMaster gets request for regional master and sends back the address of the regional master
func (o *ProviderOrClient) RequestRegionalMaster(p *common.RequestRegionalMasterReply) error {
	for a := range o.c.activeRegionalMasters.Items() { // random iteration TODO: Change to choosing by location
		p.Addr = a
		return nil
	}
	return fmt.Errorf("No Regional Masters Available")
}

// Ping gets ping from regional master and updates the active regional masters
func (r *RegionalMaster) Ping(a *common.PingArgs) {
	// Add regional master to active regional masters
    err := 0
	r.c.activeRegionalMasters.Set(a.Addr, err, cache.DefaultExpiration)
}

func main() {
	newConduitServer().run()
}

func newConduitServer() *ConduitServer {
	// Active Local Masters is a cache where masters are removed after 30 seconds if no ping
	cache := cache.New(5*time.Minute, 30*time.Second)
	return &ConduitServer{
	  activeRegionalMasters: cache,
	  fail:                  make(chan error),
	}
}

// ConduitServer the primary server struct
type ConduitServer struct {
	activeRegionalMasters *cache.Cache
	fail                  chan error
}

func (cs *ConduitServer) run() {
	go cs.listenForProviderClient()
	go cs.listenForRegionalMaster()
	log.Fatal(<-cs.fail)
}

// Accepts client or provider
func (cs *ConduitServer) listenForProviderClient() {
	s := rpc.NewServer()
	s.Register(&ProviderOrClient{cs})
	l, err := net.Listen("tcp", common.ProviderClientListenerAddr)
	if err != nil {
		cs.fail <- err
	}
	s.Accept(l)
}

func (cs *ConduitServer) listenForRegionalMaster() {
	s := rpc.NewServer()
	s.Register(&RegionalMaster{})
	l, err := net.Listen("tcp", common.RegionalMasterListenerAddr)
	if err != nil {
		cs.fail <- err
	}
	s.Accept(l)
}
