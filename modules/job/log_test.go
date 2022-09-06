package job

import "testing"

func Test_isExpired1(t *testing.T) {
	path := "/mnt/data/zt-smartassistant/log"
	type args struct {
		str  []string
		path string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{
			name: "20220807",
			args: args{
				str:  []string{"smartassistant", "20220807", "log"},
				path: path,
			},
			want: true,
		},
		{
			name: "20220812",
			args: args{
				str:  []string{"smartassistant", "20220812", "log"},
				path: path,
			},
			want: true,
		},
		{
			name: "20220811",
			args: args{
				str:  []string{"smartassistant", "20220811", "log"},
				path: path,
			},
			want: true,
		},
		{
			name: "20220813",
			args: args{
				str:  []string{"smartassistant", "20220813", "log"},
				path: path,
			},
			want: false,
		},
		{
			name: "20220804",
			args: args{
				str:  []string{"smartassistant", "20220804", "txt"},
				path: path,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isExpired(tt.args.str); got != tt.want {
				t.Errorf("isExpired() = %v, want %v", got, tt.want)
			}
		})
	}
}
