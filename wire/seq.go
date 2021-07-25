package wire

import (
	"math"
	"sync/atomic"
)

// Sequence Sequence
type sequence struct {
	num uint32
}

// Next return Next Seq id
func (s *sequence) Next() uint32 {
	next := atomic.AddUint32(&s.num, 1)
	if next == math.MaxUint32 {
		if atomic.CompareAndSwapUint32(&s.num, next, 1) {
			return 1
		}
		return s.Next()
	}
	return next
}

// Seq Seq
var Seq = sequence{num: 1}
