package node

import (
	"os"
	"testing"

	"github.com/absolute8511/redcon"
	"github.com/stretchr/testify/assert"
)

func TestKVNode_jsonCommand(t *testing.T) {
	nd, dataDir, stopC := getTestKVNode(t)
	testKey := []byte("default:test:1")
	testJSONField := []byte("1")
	testJSONFieldValue := []byte("1")

	tests := []struct {
		name string
		args redcon.Command
	}{
		{"json.get", buildCommand([][]byte{[]byte("json.get"), testKey, testJSONField})},
		{"json.get", buildCommand([][]byte{[]byte("json.get"), testKey})},
		{"json.keyexists", buildCommand([][]byte{[]byte("json.keyexists"), testKey})},
		{"json.mkget", buildCommand([][]byte{[]byte("json.mkget"), testKey, testJSONField})},
		{"json.type", buildCommand([][]byte{[]byte("json.type"), testKey})},
		{"json.type", buildCommand([][]byte{[]byte("json.type"), testKey, testJSONField})},
		{"json.arrlen", buildCommand([][]byte{[]byte("json.arrlen"), testKey, testJSONField})},
		{"json.objkeys", buildCommand([][]byte{[]byte("json.objkeys"), testKey})},
		{"json.objlen", buildCommand([][]byte{[]byte("json.objlen"), testKey})},
		{"json.set", buildCommand([][]byte{[]byte("json.set"), testKey, testJSONField, testJSONFieldValue})},
		{"json.del", buildCommand([][]byte{[]byte("json.del"), testKey, testJSONField})},
		{"json.arrappend", buildCommand([][]byte{[]byte("json.arrappend"), testKey, testJSONField, testJSONFieldValue})},
		{"json.arrpop", buildCommand([][]byte{[]byte("json.arrpop"), testKey, testJSONField})},
		{"json.del", buildCommand([][]byte{[]byte("json.del"), testKey})},
	}
	defer os.RemoveAll(dataDir)
	defer nd.Stop()
	defer close(stopC)
	c := &fakeRedisConn{}
	for _, cmd := range tests {
		c.Reset()
		origCmd := append([]byte{}, cmd.args.Raw...)
		handler, ok := nd.router.GetCmdHandler(cmd.name)
		if ok {
			handler(c, cmd.args)
			assert.Nil(t, c.GetError())
		} else {
			whandler, _ := nd.router.GetWCmdHandler(cmd.name)
			rsp, err := whandler(cmd.args)
			assert.Nil(t, err)
			_, ok := rsp.(error)
			assert.True(t, !ok)
		}
		assert.Equal(t, origCmd, cmd.args.Raw)
	}
}
