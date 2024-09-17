package snowflakeid

import (
	"fmt"
	"github.com/bwmarrin/snowflake"
	"log"
	"sync"
)

var node *snowflake.Node
var lock sync.RWMutex

func init() {
	innerNode, err := snowflake.NewNode(1)
	if err != nil {
		log.Fatalf("create snowflake id generate node failed: %v", err)
	}
	node = innerNode
}

func SetNode(nodeId int64) error {
	lock.Lock()
	defer lock.Unlock()
	innerNode, err := snowflake.NewNode(nodeId)
	if err != nil {
		return fmt.Errorf("create snowflake id generate node failed: %v", err)
	}
	node = innerNode
	return nil
}

func GetNode() *snowflake.Node {
	return node
}

func ID() int64 {
	lock.RLock()
	defer lock.RUnlock()
	return node.Generate().Int64()
}
