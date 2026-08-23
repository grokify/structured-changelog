package gitlog

import (
	"reflect"
	"testing"
)

func TestExtractRMIs(t *testing.T) {
	tests := []struct {
		name    string
		message string
		want    []string
	}{
		{
			name:    "single trailer",
			message: "feat: add thing\n\nBody text.\n\nRefs: RMI-SCHANGELOG-101",
			want:    []string{"RMI-SCHANGELOG-101"},
		},
		{
			name:    "multiple trailers deduped in order",
			message: "feat: add thing\n\nRefs: RMI-SCHANGELOG-102\nRefs: RMI-SCHANGELOG-101\nRefs: RMI-SCHANGELOG-102",
			want:    []string{"RMI-SCHANGELOG-102", "RMI-SCHANGELOG-101"},
		},
		{
			name:    "multiple ids in one trailer",
			message: "fix: stuff\n\nRefs: RMI-FOO-1, RMI-BAR-22",
			want:    []string{"RMI-FOO-1", "RMI-BAR-22"},
		},
		{
			name:    "case-insensitive trailer key",
			message: "docs: x\n\nrefs: RMI-SCHANGELOG-104",
			want:    []string{"RMI-SCHANGELOG-104"},
		},
		{
			name:    "ignores RMI mentioned in prose, not on a Refs line",
			message: "feat: implement RMI-SCHANGELOG-101 support\n\nNo trailer here.",
			want:    nil,
		},
		{
			name:    "none present",
			message: "chore: tidy\n\nRefs: #123",
			want:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExtractRMIs(tt.message)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExtractRMIs() = %v, want %v", got, tt.want)
			}
		})
	}
}
