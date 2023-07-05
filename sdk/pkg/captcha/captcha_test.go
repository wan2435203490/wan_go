package captcha

import "testing"

func TestDriverDigitFunc(t *testing.T) {
	tests := []struct {
		name     string
		wantId   string
		wantB64s string
		wantErr  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotId, gotB64s, err := DriverDigitFunc()
			if (err != nil) != tt.wantErr {
				t.Errorf("DriverDigitFunc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotId != tt.wantId {
				t.Errorf("DriverDigitFunc() gotId = %v, want %v", gotId, tt.wantId)
			}
			if gotB64s != tt.wantB64s {
				t.Errorf("DriverDigitFunc() gotB64s = %v, want %v", gotB64s, tt.wantB64s)
			}
		})
	}
}
