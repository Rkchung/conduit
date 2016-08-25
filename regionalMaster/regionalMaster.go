package main

import (
  "time"
  "math/rand"
  "github.com/patrickmn/go-cache"
  "log"
  "net"
  "net/rpc"
  "github/Conduit/common"
  "fmt"
  "os/exec"
)

type LocalMaster struct {
  r *regionalMaster
}

type Provider struct {
  r *regionalMaster
}

type Client struct {
  r *regionalMaster
}

type JobRequest struct {
  requestTime time.Time
  endTime     time.Time
  requestID   string
}

// Gets pings from local master and updates the current local master list
func (lm LocalMaster) registerBeats(Addr string) {
  // Get pings
  err := (*string)(nil)
  // Add local master to active local masters
  lm.r.activeLocalMasters.Set(Addr, err, cache.DefaultExpiration)
}

// internals

func main() {
  newRegionalMaster().run()
}

type regionalMaster struct{
  addr string
  activeLocalMasters *cache.Cache
  
  providers map[string]string // key is host addr
  clients map[string]string // key is host addr

  server *rpc.Client

  fail chan error
}

func random(min, max int) int {
    rand.Seed(time.Now().Unix())
    return rand.Intn(max - min) + min
}

func newRegionalMaster() *regionalMaster {
  rm := new(regionalMaster)
  // Active Local Masters is a cache where masters are removed after 30 seconds if no ping
  rm.activeLocalMasters = cache.New(5*time.Minute, 30*time.Second)
  // Set of logged in providers (possibly used to keep track of jobs done)
  rm.fail = make(chan error)
  
  serverConn, err := rpc.Dial("tcp", common.RegionalMasterListenerAddr)
  if err != nil {
    log.Fatalf("couldn't connect to conduit: %s", err)
  }
  rm.server = serverConn
  rm.addr = ":8010"
  return rm
}

func (r *regionalMaster) run() {
  go r.listenForProviders()
  go r.listenForClients()
  go r.listenForLocalMasters()
  go r.ping()
  log.Fatal(<-r.fail)
}

func (r *regionalMaster) listenForLocalMasters() {
  s := rpc.NewServer()
  s.Register(&LocalMaster{r})
  l, err := net.Listen("tcp", common.LocalMasterListenerAddr)
  if err != nil {
    r.fail <- err
  }
  s.Accept(l)
}

func (r *regionalMaster) listenForProviders() {
  // Make new server
  s := rpc.NewServer()
  s.Register(&Provider{r})
  l, err := net.Listen("tcp", common.ProviderListenerAddr)
  if err != nil {
    r.fail <- err
  }
  s.Accept(l)
}

func (r *regionalMaster) listenForClients() {
  s := rpc.NewServer()
  s.Register(&Client{r})
  l, err := net.Listen("tcp", common.ClientListenerAddr)
  if err != nil {
    r.fail <- err
  }
  s.Accept(l)
}

// Pings the conduit server so it knows that the regional master is active
func (rm *regionalMaster) ping() {
  args := common.PingArgs{rm.addr}
	err := rm.server.Call("ConduitServer.Ping", &args, nil)
	if err != nil {
		fmt.Errorf("Unable to Ping Server: %s", err)
	}
	time.Sleep(100 * time.Millisecond)
}

// Appends the JobRequest to the log and returns the requestID, request_time, and a set of local masters
func (c Client) makeNewRequest(provider_id string) (time.Time, string, string) {
  requestTime := time.Now()
  requestID := string(random(0, 2147483647))
  // TODO: check if collision with another request_id
  newRequest := JobRequest{requestTime: requestTime, requestID: requestID}
  c.r.appendNewRequest(newRequest)
  var master string
  for a := range c.r.activeLocalMasters.Items() {
    master = a
  }
  // TODO: send requestID, request_time and local masters to client
  return requestTime, requestID, master
}
  
// Appends start time to Request
func (p Provider) appendStartTime(requestID string, time time.Time) {
  // TODO: Append start time
}

// Appends end time to Request
func (p Provider) appendEndTime(requestID string, time time.Time) {
  // TODO: Append end of time
}

// Appends new request to log
func (r regionalMaster) appendNewRequest(newRequest JobRequest) {
  // TODO: Append request stuff
}

// Registers the provider and saves
func (p Provider) register(Addr string) {
  // Generate provider ID
  pID, err := exec.Command("uuidgen").Output()
  if err != nil {
      log.Fatal(err)
  }
  p.r.providers[Addr] = string(pID[:15])
}

// Registers the client and gives ID
func (c Client) register(publicKey string) {
  // Generate client ID
  pID, err := exec.Command("uuidgen").Output()
  if err != nil {
      log.Fatal(err)
  }
  // Saves the pID of client and publicKey
  c.r.clients[string(pID[:15])] = publicKey
}
  
// Returns request info from requestID
func (r regionalMaster) getRequest(requestID int) {
  
}
