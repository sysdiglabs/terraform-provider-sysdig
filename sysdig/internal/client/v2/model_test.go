package v2

import (
	"encoding/json"
	"testing"
)

func TestExpressionV2MarshalJSON(t *testing.T) {
	tests := []struct {
		name       string
		expr       ExpressionV2
		wantArray  bool   // true = expect "value" to be []interface{}
		wantLen    int    // expected array length (when wantArray)
		wantFirst  string // expected first element (when wantArray)
		wantAbsent bool   // true = expect "value" key to be absent
	}{
		{
			name:      "scalar Value field wraps to single-element array",
			expr:      ExpressionV2{Field: "agent.tag.key", Operator: "is not", Value: "asd"},
			wantArray: true,
			wantLen:   1,
			wantFirst: "asd",
		},
		{
			name:      "single-element Values stays as array",
			expr:      ExpressionV2{Field: "org", Operator: "in", Values: []string{"x"}},
			wantArray: true,
			wantLen:   1,
			wantFirst: "x",
		},
		{
			name:      "multi-element Values stays as array",
			expr:      ExpressionV2{Field: "org", Operator: "in", Values: []string{"a", "b"}},
			wantArray: true,
			wantLen:   2,
			wantFirst: "a",
		},
		{
			name:       "empty Value and nil Values omits key",
			expr:       ExpressionV2{Field: "f", Operator: "op"},
			wantAbsent: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := json.Marshal(tt.expr)
			if err != nil {
				t.Fatalf("marshal error: %v", err)
			}

			var got map[string]interface{}
			if err := json.Unmarshal(b, &got); err != nil {
				t.Fatalf("unmarshal marshalled json: %v", err)
			}

			v, ok := got["value"]

			if tt.wantAbsent {
				if ok {
					t.Fatalf("expected value key absent, got: %v", v)
				}
				return
			}

			if !ok {
				t.Fatalf("expected value key in marshalled json: %s", string(b))
			}

			arr, isArr := v.([]interface{})
			if !isArr {
				t.Fatalf("expected value to be array, got %T: %v (json: %s)", v, v, string(b))
			}
			if len(arr) != tt.wantLen {
				t.Fatalf("expected array length %d, got %d", tt.wantLen, len(arr))
			}
			if tt.wantFirst != "" && arr[0].(string) != tt.wantFirst {
				t.Fatalf("expected first element %q, got %q", tt.wantFirst, arr[0])
			}
		})
	}
}

func TestExpressionV2UnmarshalVariants(t *testing.T) {
	var e ExpressionV2

	// value as array
	if err := json.Unmarshal([]byte(`{"field":"f","operator":"op","value":["a","b"]}`), &e); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if len(e.Values) != 2 || e.Values[0] != "a" || e.Values[1] != "b" {
		t.Fatalf("unexpected values: %#v", e.Values)
	}

	// value as single string -> should populate Value (string)
	if err := json.Unmarshal([]byte(`{"field":"f","operator":"op","value":"x"}`), &e); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if e.Value != "x" {
		t.Fatalf("expected Value == 'x', got: %#v (Values: %#v)", e.Value, e.Values)
	}

	// values key
	if err := json.Unmarshal([]byte(`{"field":"f","operator":"op","values":["y"]}`), &e); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if len(e.Values) != 1 || e.Values[0] != "y" {
		t.Fatalf("unexpected values: %#v", e.Values)
	}
}
