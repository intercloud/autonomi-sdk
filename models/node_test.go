package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNodeTypeString(t *testing.T) {
	tests := []struct {
		name     string
		nodeType NodeType
		expect   string
	}{
		{
			name:     NodeTypeAccess,
			nodeType: NodeTypeAccess,
			expect:   "access",
		},
		{
			name:     NodeTypeBridge,
			nodeType: NodeTypeBridge,
			expect:   "bridge",
		},
		{
			name:     NodeTypeCloud,
			nodeType: NodeTypeCloud,
			expect:   "cloud",
		},
		{
			name:     NodeTypeRouter,
			nodeType: NodeTypeRouter,
			expect:   "router",
		},
	}

	for _, tc := range tests {
		t.Log(tc.name)
		tc := tc
		assert.Equal(t, tc.expect, tc.nodeType.String())
	}
}
