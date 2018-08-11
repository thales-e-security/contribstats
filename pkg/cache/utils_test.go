package cache

import "testing"

func Test_stringInSlice(t *testing.T) {
	type args struct {
		a    string
		list []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "True",
			args: args{
				a:    "A",
				list: []string{"A", "B"},
			},
			want: true,
		},
		{
			name: "False",
			args: args{
				a:    "c",
				list: []string{"A", "B"},
			},
			want: false,
		},
	}

	// Test I
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stringInSlice(tt.args.a, tt.args.list); got != tt.want {
				t.Errorf("stringInSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}
