package sql

import (
	"testing"

	"github.com/aquasecurity/defsec/internal/types"

	"github.com/aquasecurity/defsec/pkg/state"

	"github.com/aquasecurity/defsec/pkg/providers/google/sql"
	"github.com/aquasecurity/defsec/pkg/scan"

	"github.com/stretchr/testify/assert"
)

func TestCheckEnableBackup(t *testing.T) {
	tests := []struct {
		name     string
		input    sql.SQL
		expected bool
	}{
		{
			name: "Database instance backups disabled",
			input: sql.SQL{
				Instances: []sql.DatabaseInstance{
					{
						Metadata:  types.NewTestMetadata(),
						IsReplica: types.Bool(false, types.NewTestMetadata()),
						Settings: sql.Settings{
							Metadata: types.NewTestMetadata(),
							Backups: sql.Backups{
								Metadata: types.NewTestMetadata(),
								Enabled:  types.Bool(false, types.NewTestMetadata()),
							},
						},
					},
				},
			},
			expected: true,
		},
		{
			name: "Database instance backups enabled",
			input: sql.SQL{
				Instances: []sql.DatabaseInstance{
					{
						Metadata:  types.NewTestMetadata(),
						IsReplica: types.Bool(false, types.NewTestMetadata()),
						Settings: sql.Settings{
							Metadata: types.NewTestMetadata(),
							Backups: sql.Backups{
								Metadata: types.NewTestMetadata(),
								Enabled:  types.Bool(true, types.NewTestMetadata()),
							},
						},
					},
				},
			},
			expected: false,
		},
		{
			name: "Read replica does not require backups",
			input: sql.SQL{
				Instances: []sql.DatabaseInstance{
					{
						Metadata:  types.NewTestMetadata(),
						IsReplica: types.Bool(true, types.NewTestMetadata()),
						Settings: sql.Settings{
							Metadata: types.NewTestMetadata(),
							Backups: sql.Backups{
								Metadata: types.NewTestMetadata(),
								Enabled:  types.Bool(false, types.NewTestMetadata()),
							},
						},
					},
				},
			},
			expected: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var testState state.State
			testState.Google.SQL = test.input
			results := CheckEnableBackup.Evaluate(&testState)
			var found bool
			for _, result := range results {
				if result.Status() == scan.StatusFailed && result.Rule().LongID() == CheckEnableBackup.Rule().LongID() {
					found = true
				}
			}
			if test.expected {
				assert.True(t, found, "Rule should have been found")
			} else {
				assert.False(t, found, "Rule should not have been found")
			}
		})
	}
}
