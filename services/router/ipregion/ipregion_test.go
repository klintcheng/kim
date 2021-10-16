package ipregion

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIp2region_Search(t *testing.T) {
	region, err := NewIp2region("../ip2region.db")
	assert.Nil(t, err)

	got, err := region.Search("3.166.231.6")
	assert.Nil(t, err)
	t.Log(got)
}
