package example

import "testing"

func TestTestDataType_Value(t *testing.T) {
	tests := []struct {
		name string
		x    TestDataType
		want string
	}{
		{
			name: "test1",
			x:    TestDataType_One,
			want: "1111",
		},
		{
			name: "test2",
			x:    TestDataType_Two,
			want: "2222",
		},
		{
			name: "test3",
			x:    TestDataType_Three,
			want: "333",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.x.Value(); got != tt.want {
				t.Errorf("Value() = %v, want %v", got, tt.want)
			}
		})
	}
}
