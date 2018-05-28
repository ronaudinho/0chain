package node

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"0chain.net/common"
	"0chain.net/config"
	"0chain.net/datastore"
	"0chain.net/encryption"
	"github.com/golang/snappy"
)

var (
	HeaderRequestTimeStamp      = "X-Request-Timestamp"
	HeaderRequestHashData       = "X-Request-Hashdata"
	HeaderRequestHash           = "X-Request-Hash"
	HeaderRequestRelayLength    = "X-Request-Relay-Length"
	HeaderRequestMaxRelayLength = "X-Request-Max-Relay-Length"
	HeaderRequestEntityName     = "X-Request-Entity-Name"
	HeaderRequestEntityID       = "X-Request-Entity-ID"
	HeaderRequestChainID        = "X-Chain-Id"

	HeaderNodeID               = "X-Node-Id"
	HeaderNodeRequestSignature = "X-Node-Request-Signature"
)

/*SendHandler is used to send any message to a given node */
type SendHandler func(n *Node) bool

/*EntitySendHandler is used to send an entity to a given node */
type EntitySendHandler func(entity datastore.Entity) SendHandler

type ReceiveEntityHandlerF func(ctx context.Context, entity *datastore.Entity) (interface{}, error)

/*SendAtleast - It tries to communicate to at least the given number of active nodes */
func (np *Pool) SendAtleast(numNodes int, handler SendHandler) {
	const THRESHOLD = 2
	nodes := np.shuffleNodes()

	validCount := 0
	allCount := 0
	for _, node := range nodes {
		if node.Status == NodeStatusInactive {
			continue
		}
		allCount++
		valid := handler(node)
		if valid {
			validCount++
			if validCount == numNodes {
				break
			}
		}
		if allCount >= numNodes+THRESHOLD {
			break
		}
	}
}

/*SetHeaders - sets the request headers */
func SetHeaders(req *http.Request, entity datastore.Entity, maxRelayLength int64, currentRelayLength int64) bool {
	ts := common.Now()
	hashdata := fmt.Sprintf("%v:%v:%v", Self.GetID(), ts, entity.GetKey())
	hash := encryption.Hash(hashdata)
	//TODO: Replace Self.privateKey with API from Ken
	signature, err := Self.Sign(hash)
	if err != nil {
		return false
	}
	req.Header.Set(HeaderRequestChainID, config.GetServerChainID())
	req.Header.Set(HeaderNodeID, Self.GetID())
	req.Header.Set(HeaderRequestTimeStamp, strconv.FormatInt(int64(ts), 10))
	req.Header.Set(HeaderRequestHashData, hashdata)
	req.Header.Set(HeaderRequestHash, hash)
	req.Header.Set(HeaderNodeRequestSignature, signature)

	if maxRelayLength > 0 {
		req.Header.Set(HeaderRequestMaxRelayLength, strconv.FormatInt(maxRelayLength, 10))
	}
	req.Header.Set(HeaderRequestRelayLength, strconv.FormatInt(currentRelayLength, 10))

	req.Header.Set(HeaderRequestEntityName, entity.GetEntityName())
	req.Header.Set(HeaderRequestEntityID, datastore.ToString(entity.GetKey()))
	return true
}

type SendOptions struct {
	MaxRelayLength     int64
	CurrentRelayLength int64
	Compress           bool
}

/*SendEntityHandler provides a client API to send an entity */
func SendEntityHandler(uri string, options SendOptions) EntitySendHandler {
	return func(entity datastore.Entity) SendHandler {
		return func(n *Node) bool {
			url := fmt.Sprintf("%v/%v", n.GetURLBase(), uri)
			client := &http.Client{Timeout: 500 * time.Millisecond}

			buffer := new(bytes.Buffer)
			json.NewEncoder(buffer).Encode(entity)
			if options.Compress {
				cbytes := snappy.Encode(nil, buffer.Bytes())
				buffer = bytes.NewBuffer(cbytes)
			}
			req, err := http.NewRequest("POST", url, buffer)
			if err != nil {
				return false
			}
			defer req.Body.Close()

			if options.Compress {
				req.Header.Set("Content-Encoding", "snappy")
			}
			req.Header.Set("Content-Type", "application/json; charset=utf-8")
			SetHeaders(req, entity, options.MaxRelayLength, options.CurrentRelayLength)
			resp, err := client.Do(req)
			if err != nil {
				return false
			}
			if resp.StatusCode != http.StatusOK {
				return false
			}
			return true
		}
	}
}

/*ToN2NReceiveEntityHandler - takes a handler that accepts an entity, processes and responds and converts it
* into somethign suitable for Node 2 Node communication
 */
func ToN2NReceiveEntityHandler(handler common.JSONEntityReqResponderF) common.ReqRespHandlerf {
	return func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("Content-type")
		if !strings.HasPrefix(contentType, "application/json") {
			http.Error(w, "Header Content-type=application/json not found", 400)
			return
		}

		chainID := r.Header.Get(HeaderRequestChainID)
		if config.GetServerChainID() != chainID {
			return
		}

		nodeID := r.Header.Get(HeaderNodeID)
		sender := GetNode(nodeID)
		if sender == nil {
			fmt.Printf("received request from unrecognized node %v\n", nodeID)
			return
		}

		//TODO: check the timestamp?
		//reqTS := r.Header.Get(HeaderRequestTimeStamp)
		reqHashdata := r.Header.Get(HeaderRequestHashData)
		reqHash := r.Header.Get(HeaderRequestHash)
		//TODO: Do we need this check?
		if reqHash != encryption.Hash(reqHashdata) {
			return
		}
		reqSignature := r.Header.Get(HeaderNodeRequestSignature)
		if ok, _ := sender.Verify(reqSignature, reqHash); !ok {
			return
		}

		entityName := r.Header.Get(HeaderRequestEntityName)
		if entityName == "" {
			return
		}
		entityProvider := datastore.GetProvider(entityName)
		if entityProvider == nil {
			return
		}
		var buffer io.Reader = r.Body
		if r.Header.Get("Content-Encoding") == "snappy" {
			cbuffer := new(bytes.Buffer)
			cbuffer.ReadFrom(r.Body)

			cbytes, err := snappy.Decode(nil, cbuffer.Bytes())
			if err != nil {
				return
			}
			buffer = bytes.NewReader(cbytes)
		}
		decoder := json.NewDecoder(buffer)
		entity := entityProvider()
		err := decoder.Decode(entity)
		if err != nil {
			http.Error(w, "Error decoding json", 500)
			return
		}
		ctx := r.Context()
		data, err := handler(ctx, entity)
		common.Respond(w, data, err)
	}
}

/*SetupN2NHandlers - Setup all the node 2 node communiations
 */
func SetupN2NHandlers() {
	http.HandleFunc("/v1/_n2n/entity/post", ToN2NReceiveEntityHandler(common.PrintEntityHandler))
}
