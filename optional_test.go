package buildkite

import (
	"encoding/json"
	"testing"
)

func TestOptionalMarshalJSON(t *testing.T) {
	t.Parallel()

	type payload struct {
		Name    Optional[string]            `json:"name,omitzero"`
		Tags    Optional[[]string]          `json:"tags,omitzero"`
		Env     Optional[map[string]string] `json:"env,omitzero"`
		Enabled Optional[bool]              `json:"enabled,omitzero"`
	}

	tests := []struct {
		name string
		in   payload
		want string
	}{
		{
			name: "unset fields omitted",
			in:   payload{},
			want: `{}`,
		},
		{
			name: "empty string is sent",
			in: payload{
				Name: Some(""),
			},
			want: `{"name":""}`,
		},
		{
			name: "false is sent",
			in: payload{
				Enabled: Some(false),
			},
			want: `{"enabled":false}`,
		},
		{
			name: "empty slice is sent",
			in: payload{
				Tags: Some([]string{}),
			},
			want: `{"tags":[]}`,
		},
		{
			name: "empty map is sent",
			in: payload{
				Env: Some(map[string]string{}),
			},
			want: `{"env":{}}`,
		},
		{
			name: "nil slice is sent as null when set",
			in: payload{
				Tags: Some([]string(nil)),
			},
			want: `{"tags":null}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := json.Marshal(tt.in)
			if err != nil {
				t.Fatalf("Marshal returned error: %v", err)
			}
			assertJSONEqual(t, string(got), tt.want)
		})
	}
}

func TestOptionalUnmarshalJSON(t *testing.T) {
	t.Parallel()

	type payload struct {
		Name Optional[string]   `json:"name,omitzero"`
		Tags Optional[[]string] `json:"tags,omitzero"`
	}

	var got payload
	if err := json.Unmarshal([]byte(`{"name":null,"tags":[]}`), &got); err != nil {
		t.Fatalf("Unmarshal returned error: %v", err)
	}

	name, nameOK := got.Name.Value()
	if !nameOK {
		t.Fatal("Name was not marked set")
	}
	if name != "" {
		t.Errorf("Name = %q, want empty string", name)
	}

	tags, tagsOK := got.Tags.Value()
	if !tagsOK {
		t.Fatal("Tags was not marked set")
	}
	if len(tags) != 0 {
		t.Errorf("Tags length = %d, want 0", len(tags))
	}
}
