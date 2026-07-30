package bgr555

import "testing"

func TestDecode(t *testing.T) {
	tests := []struct {
		name string
		in   uint16
		want Color
	}{
		{"black", 0x0000, Color{}},
		{"red", 0x001F, Color{R: 255}},
		{"green", 0x03E0, Color{G: 255}},
		{"blue", 0x7C00, Color{B: 255}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Decode(tt.in); got != tt.want {
				t.Fatalf("Decode(%#04x)=%+v, want %+v", tt.in, got, tt.want)
			}
		})
	}
}
