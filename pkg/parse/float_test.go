package parse

import "testing"

func TestParseExplicitFloat(t *testing.T) {
	config := Config{Tokens: []string{"float_elapsed"}, Patterns: []string{"duration=float_elapsed"}}
	for _, tt := range []struct {
		input string
		want  float64
	}{
		{"42", 42},
		{"-42", -42},
		{"0", 0},
		{"1e3", 1000},
		{"2E-2", 0.02},
		{"1.5", 1.5},
	} {
		t.Run(tt.input, func(t *testing.T) {
			got, err := Parse("duration="+tt.input, &config)
			if err != nil {
				t.Fatal(err)
			}
			if len(got) != 1 {
				t.Fatalf("got %d records, want 1", len(got))
			}
			if got[0].Pattern != "duration=float_elapsed" {
				t.Errorf("pattern = %q, want duration=float_elapsed", got[0].Pattern)
			}
			if got[0].Params["float_elapsed"] != (Param{Value: tt.want, Type: "float"}) {
				t.Errorf("param = %+v, want float %v", got[0].Params["float_elapsed"], tt.want)
			}
		})
	}
	for _, input := range []string{"invalid", "NaN", "Inf", "+Inf", "-Inf", "1e309"} {
		t.Run(input, func(t *testing.T) {
			got, err := Parse("duration="+input, &config)
			if err != nil {
				t.Fatal(err)
			}
			if len(got) != 1 {
				t.Fatalf("got %d records, want 1", len(got))
			}
			if param, ok := got[0].Params["float_elapsed"]; ok {
				t.Errorf("invalid float was retained: %+v", param)
			}
		})
	}
}
