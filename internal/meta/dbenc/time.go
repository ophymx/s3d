package dbenc

import (
	"encoding/binary"
	"fmt"
	"time"
)

// implements MsgPack Timestamp Extension as defined in spec:
//    https://github.com/msgpack/msgpack/blob/master/spec.md#formats-timestamp
type timeExtension struct {
	sec  uint64
	nsec uint32
}

const nsecMask = 0x00000003ffffffff

func newTime(t time.Time) *timeExtension {
	return &timeExtension{
		sec:  uint64(t.Unix()),
		nsec: uint32(t.Nanosecond()),
	}
}

func (e *timeExtension) ExtensionType() int8 {
	return -1
}

func (e *timeExtension) Len() int {
	if (e.sec >> 34) == 0 {
		if e.nsec == 0 {
			return 4
		}
		return 8
	}
	return 12
}

func (e *timeExtension) MarshalBinaryTo(b []byte) error {
	switch len(b) {
	case 4:
		binary.BigEndian.PutUint32(b, uint32(e.sec))
	case 8:
		binary.BigEndian.PutUint64(b, (uint64(e.nsec)<<34)|e.sec)
	case 12:
		binary.BigEndian.PutUint32(b, e.nsec)
		binary.BigEndian.PutUint64(b[4:], e.sec)
	default:
		return fmt.Errorf("invalid length for time extension: %d", len(b))
	}
	return nil
}

func (e *timeExtension) UnmarshalBinary(b []byte) error {
	switch len(b) {
	case 4:
		e.nsec = 0
		e.sec = uint64(binary.BigEndian.Uint32(b))
	case 8:
		data := binary.BigEndian.Uint64(b)
		e.nsec = uint32(data >> 34)
		e.sec = data & nsecMask
	case 12:
		e.nsec = binary.BigEndian.Uint32(b)
		e.sec = binary.BigEndian.Uint64(b[4:])
	default:
		return fmt.Errorf("invalid length of time extension: %d", len(b))
	}
	return nil
}

func (e *timeExtension) UTC() time.Time {
	return time.Unix(int64(e.sec), int64(e.nsec)).UTC()
}
