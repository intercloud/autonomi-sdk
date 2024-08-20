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
			name:                AdministrativeStateCreationPending.String(),
			administrativeState: AdministrativeStateCreationPending,
			expect:              "creation_pending",
		},
		{
			name:                AdministrativeStateCreationProceed.String(),
			administrativeState: AdministrativeStateCreationProceed,
			expect:              "creation_proceed",
		},
		{
			name:                AdministrativeStateCreationError.String(),
			administrativeState: AdministrativeStateCreationError,
			expect:              "creation_error",
		},
		{
			name:                AdministrativeStateDeployed.String(),
			administrativeState: AdministrativeStateDeployed,
			expect:              "deployed",
		},
		{
			name:                AdministrativeStateDeletePending.String(),
			administrativeState: AdministrativeStateDeletePending,
			expect:              "delete_pending",
		},
		{
			name:                AdministrativeStateDeleteProceed.String(),
			administrativeState: AdministrativeStateDeleteProceed,
			expect:              "delete_proceed",
		},
		{
			name:                AdministrativeStateDeleteError.String(),
			administrativeState: AdministrativeStateDeleteError,
			expect:              "delete_error",
		},
		{
			name:                AdministrativeStateDeleted.String(),
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
