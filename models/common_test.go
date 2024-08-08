package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdministrativeStateString(t *testing.T) {
	tests := []struct {
		name                string
		administrativeState AdministrativeState
		expect              string
	}{
		{
			name:                NodeTypeAccess,
			administrativeState: AdministrativeStateCreationPending,
			expect:              "creation_pending",
		},
		{
			name:                NodeTypeBridge,
			administrativeState: AdministrativeStateCreationProceed,
			expect:              "creation_proceed",
		},
		{
			name:                NodeTypeCloud,
			administrativeState: AdministrativeStateCreationError,
			expect:              "creation_error",
		},
		{
			name:                NodeTypeRouter,
			administrativeState: AdministrativeStateDeployed,
			expect:              "deployed",
		},
		{
			name:                NodeTypeRouter,
			administrativeState: AdministrativeStateDeletePending,
			expect:              "delete_pending",
		},
		{
			name:                NodeTypeRouter,
			administrativeState: AdministrativeStateDeleteProceed,
			expect:              "delete_proceed",
		},
		{
			name:                NodeTypeRouter,
			administrativeState: AdministrativeStateDeleteError,
			expect:              "delete_error",
		},
		{
			name:                NodeTypeRouter,
			administrativeState: AdministrativeStateDeleted,
			expect:              "deleted",
		},
	}

	for _, tc := range tests {
		t.Log(tc.name)
		tc := tc
		assert.Equal(t, tc.expect, tc.administrativeState.String())
	}
}
