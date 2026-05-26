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
		{
			name: "nil map is sent as null when set",
			in: payload{
				Env: Some(map[string]string(nil)),
			},
			want: `{"env":null}`,
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
