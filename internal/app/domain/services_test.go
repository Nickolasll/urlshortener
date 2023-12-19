package domain

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGenerateSlug(t *testing.T) {
	type args struct {
		size int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "zero size",
			args: args{
				size: 0,
			},
			want: 1,
		},
		{
			name: "simple success",
			args: args{
				size: 5,
			},
			want: 6,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateSlug(tt.args.size)
			assert.Equal(t, len(got), tt.want)
		})
	}
}

func BenchmarkShorteng(b *testing.B) {
	URL := "www.testlongurl.com"
	userID := uuid.New().String()
	for i := 0; i < b.N; i++ {
		Shorten(URL, userID)
	}
}
