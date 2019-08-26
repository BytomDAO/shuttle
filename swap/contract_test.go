package swap

import (
	"reflect"
	"testing"
)

var (
	assetRequested  = "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"
	seller          = "00145dd7b82556226d563b6e7d573fe61d23bd461c1f"
	cancelKey       = "3e5d7d52d334964eef173021ef6a04dc0807ac8c41700fe718f5a80c2109f79e"
	amountRequested = uint64(1000000000)
)

func TestCompileLockContract(t *testing.T) {
	var tests = []struct {
		assetRequested  string
		seller          string
		cancelKey       string
		amountRequested uint64
		want            string
	}{
		{
			assetRequested:  "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
			seller:          "00145dd7b82556226d563b6e7d573fe61d23bd461c1f",
			cancelKey:       "3e5d7d52d334964eef173021ef6a04dc0807ac8c41700fe718f5a80c2109f79e",
			amountRequested: uint64(1000000000),
			want:            "203e5d7d52d334964eef173021ef6a04dc0807ac8c41700fe718f5a80c2109f79e1600145dd7b82556226d563b6e7d573fe61d23bd461c1f0400ca9a3b20ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff741a547a6413000000007b7b51547ac1631a000000547a547aae7cac00c0",
		},
	}

	for i, tt := range tests {
		got := CompileLockContract(tt.assetRequested, tt.seller, tt.cancelKey, tt.amountRequested).Program

		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("CompileLockContract %d: content mismatch: have %x, want %x", i, got, tt.want)
		}
	}
}
