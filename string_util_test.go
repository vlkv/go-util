package util

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"fmt"
)

func TestIsBase64(t *testing.T) {
	const base64 = `TWFuIGlzIGRpc3Rpbmd1aXNoZWQsIG5vdCBvbmx5IGJ5IGhpcyByZWFzb24sIGJ1dCBieSB0aGlzIHNpbmd1bGFyIHBhc3Npb24gZnJvbSBvdGhlciBhbmltYWxzLCB3aGljaCBpcyBhIGx1c3Qgb2YgdGhlIG1pbmQsIHRoYXQgYnkgYSBwZXJzZXZlcmFuY2Ugb2YgZGVsaWdodCBpbiB0aGUgY29udGludWVkIGFuZCBpbmRlZmF0aWdhYmxlIGdlbmVyYXRpb24gb2Yga25vd2xlZGdlLCBleGNlZWRzIHRoZSBzaG9ydCB2ZWhlbWVuY2Ugb2YgYW55IGNhcm5hbCBwbGVhc3VyZS4=`
	assert.True(t, IsBase64(base64))

	const notBase64 = `Hello World`
	assert.False(t, IsBase64(notBase64))
}

func TestStripMargin(t *testing.T) {
	var str = StripMargin(`1
	|234
	|    567
	|    |8910
    11`)

	fmt.Println(str)
	assert.Equal(t, "1\n234\n    567\n    |8910\n    11", str)
}
