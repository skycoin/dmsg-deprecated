package setup

import (
	"encoding/json"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewProtocol(t *testing.T) {
	connA, connB := net.Pipe()
	protoA := NewProtocol(connA)
	protoB := NewProtocol(connB)

	cases := []struct {
		Type PacketType
		Data string
	}{
		{PacketType(0), "this is a test!"},
		{PacketType(255), "this is another test!"},
	}

	for _, c := range cases {
		errChan := make(chan error, 1)
		go func() {
			errChan <- protoA.WritePacket(c.Type, []byte(c.Data))
		}()

		pt, data, err := protoB.ReadPacket()

		var decoded []byte
		require.NoError(t, json.Unmarshal(data, &decoded))

		assert.NoError(t, err)
		assert.Equal(t, c.Type, pt)
		assert.Equal(t, c.Data, string(decoded))

		assert.NoError(t, <-errChan)
	}
}