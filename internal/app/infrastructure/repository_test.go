package infrastructure

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRAMRepository(t *testing.T) {
	type args struct {
		url  string
		slug string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "save and get",
			args: args{
				url:  "http://test.com",
				slug: "/ABCDE",
			},
			want: "http://test.com",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RAMRepository.Save(tt.args.url, tt.args.slug)
			got, ok := RAMRepository.Get(tt.args.slug)
			assert.Equal(t, ok, true)
			assert.Equal(t, got, tt.want)
		})
	}
}
