package wav

import (
	"encoding/binary"
	"encoding/hex"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ToInt(t *testing.T) {
	// buf := []byte{255, 255, 127, 0}
	var x int32 = -_MAX24_SAMPLE

	// buf32 := bytes.NewBuffer(nil)
	// if err := binary.Write(buf32, binary.LittleEndian, x); err != nil {
	// 	t.Fatal(err)
	// }

	exp := make([]byte, 4)
	binary.LittleEndian.PutUint32(exp, uint32(x))

	actual := make([]byte, 4)
	putUint24(actual, x)
	assert.Equal(t, hex.EncodeToString(exp), hex.EncodeToString(actual))

}

func Test_ToFloat(t *testing.T) {
	val := []byte{255, 255, 127, 0}
	exp := binary.LittleEndian.Uint32(val)
	act := (fromBit24(val))
	assert.Equal(t, exp, uint32(act))

	e32 := math.Float32frombits(exp)
	a32 := math.Float32frombits(uint32(act))
	assert.Equal(t, e32, a32)

	e64 := math.Float64frombits(uint64(exp))
	a64 := math.Float64frombits(uint64(act))
	assert.Equal(t, e64, a64)

	u64 := math.Float64bits(a64)
	assert.Equal(t, act, int32(u64))
}

func Test_binary(t *testing.T) {
	// var n = _MAX24_SAMPLE
	x := int16(32760)
	y := int16(-32760)

	act1 := float64(x) / _MAX16_SAMPLE
	act2 := float64(y) / _MAX16_SAMPLE

	xx := int32(act1 * _MAX24_SAMPLE)
	yy := int32(act2 * _MAX24_SAMPLE)

	assert.Equal(t, x, int16(xx>>8))
	assert.Equal(t, y, int16(yy>>8))

}
