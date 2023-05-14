package pkg

import "testing"

func Test_convertTimeOnIce(t *testing.T) {
	tests := []struct {
		name      string
		timeOnIce string
		games     string
		want      string
	}{
		{
			name:      "converts time on ice to seconds",
			timeOnIce: "10:10",
			games:     "1",
			want:      "610.00",
		},
		{
			name:      "converts time on ice to seconds, multiple games",
			timeOnIce: "10:10",
			games:     "2",
			want:      "305.00",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convertTimeOnIce(tt.timeOnIce, tt.games); got != tt.want {
				t.Errorf("convertTimeOnIce() = %v, want %v", got, tt.want)
			}
		})
	}
}
