// Code generated by "stringer -type=Source"; DO NOT EDIT

package refstate

import "fmt"

const _Source_name = "invalidSourceInternalFromDownstreamFromUpstream"

var _Source_index = [...]uint8{0, 13, 21, 35, 47}

func (i Source) String() string {
	if i < 0 || i >= Source(len(_Source_index)-1) {
		return fmt.Sprintf("Source(%d)", i)
	}
	return _Source_name[_Source_index[i]:_Source_index[i+1]]
}
