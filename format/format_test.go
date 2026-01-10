package format

import (
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		input   string
		want    Format
		wantErr bool
	}{
		{"toon", TOON, false},
		{"", TOON, false},
		{"json", JSON, false},
		{"json-compact", JSONCompact, false},
		{"invalid", "", true},
		{"JSON", "", true}, // case-sensitive
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := Parse(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Parse(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestMarshal(t *testing.T) {
	type testStruct struct {
		Name  string `json:"name"`
		Count int    `json:"count"`
	}

	v := testStruct{Name: "test", Count: 42}

	t.Run("JSON", func(t *testing.T) {
		got, err := Marshal(v, JSON)
		if err != nil {
			t.Fatalf("Marshal() error = %v", err)
		}
		want := "{\n  \"name\": \"test\",\n  \"count\": 42\n}"
		if string(got) != want {
			t.Errorf("Marshal() = %q, want %q", string(got), want)
		}
	})

	t.Run("JSONCompact", func(t *testing.T) {
		got, err := Marshal(v, JSONCompact)
		if err != nil {
			t.Fatalf("Marshal() error = %v", err)
		}
		want := `{"name":"test","count":42}`
		if string(got) != want {
			t.Errorf("Marshal() = %q, want %q", string(got), want)
		}
	})

	t.Run("TOON", func(t *testing.T) {
		got, err := Marshal(v, TOON)
		if err != nil {
			t.Fatalf("Marshal() error = %v", err)
		}
		// TOON output should contain the field names and values
		// Note: TOON uses struct field names (capitalized), not json tags
		s := string(got)
		if !strings.Contains(s, "Name") || !strings.Contains(s, "test") {
			t.Errorf("Marshal() TOON output missing expected content: %q", s)
		}
		if !strings.Contains(s, "Count") || !strings.Contains(s, "42") {
			t.Errorf("Marshal() TOON output missing expected content: %q", s)
		}
	})
}

func TestMarshalArray(t *testing.T) {
	type item struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	items := []item{
		{ID: 1, Name: "first"},
		{ID: 2, Name: "second"},
	}

	t.Run("TOON tabular", func(t *testing.T) {
		got, err := Marshal(items, TOON)
		if err != nil {
			t.Fatalf("Marshal() error = %v", err)
		}
		s := string(got)
		// TOON should use tabular format for arrays of structs
		if !strings.Contains(s, "first") || !strings.Contains(s, "second") {
			t.Errorf("Marshal() TOON output missing expected content: %q", s)
		}
	})
}

func TestFormatString(t *testing.T) {
	tests := []struct {
		f    Format
		want string
	}{
		{TOON, "toon"},
		{JSON, "json"},
		{JSONCompact, "json-compact"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.f.String(); got != tt.want {
				t.Errorf("Format.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
