package sanitize

import "testing"

func TestLogValue(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "clean alphanumeric",
			input: "order.created.v1",
			want:  "order.created.v1",
		},
		{
			name:  "clean URI source",
			input: "/default/sap.kyma/id",
			want:  "/default/sap.kyma/id",
		},
		{
			name:  "allowed special chars preserved",
			input: "app_name-1:thing@host",
			want:  "app_name-1:thing@host",
		},
		{
			name:  "newline injection replaced",
			input: "legit.type\nFAKE_LOG_ENTRY",
			want:  "legit.type_FAKE_LOG_ENTRY",
		},
		{
			name:  "carriage return injection replaced",
			input: "legit.type\r\nINFO forged log line",
			want:  "legit.type__INFO forged log line",
		},
		{
			name:  "tab replaced",
			input: "legit.type\tinjected",
			want:  "legit.type_injected",
		},
		{
			name:  "null byte replaced",
			input: "legit.type\x00injected",
			want:  "legit.type_injected",
		},
		{
			name:  "ANSI escape sequence replaced",
			input: "legit.type\x1b[31mREDTEXT\x1b[0m",
			want:  "legit.type__31mREDTEXT__0m",
		},
		{
			name:  "unicode control chars replaced",
			input: "legit\u200Btype",
			want:  "legit_type",
		},
		{
			name:  "empty string unchanged",
			input: "",
			want:  "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := LogValue(tc.input)
			if got != tc.want {
				t.Errorf("LogValue(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}
