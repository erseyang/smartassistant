package wangpan

import "testing"

func TestGetUcs2Md4(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "123456",
			args: args{"123456"},
			want: "32ED87BDB5FDC5E9CBA88547376818D4",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetUcs2Md4(tt.args.str); got != tt.want {
				t.Errorf("GetUcs2Md4() = %v, want %v", got, tt.want)
			}
		})
	}
}
