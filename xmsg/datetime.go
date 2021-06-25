package xmsg

import (
	"database/sql/driver"
	"encoding/binary"
	"github.com/tinylib/msgp/msgp"
	"strconv"
	"time"
)

var (
	DateTime_EXT_TYPE int8 = -1
)

func init() {
	msgp.RegisterExtension(DateTime_EXT_TYPE, func() msgp.Extension { return new(DateTime) })
}

type DateTime time.Time

func (t DateTime) String() string {
	str, _ := time.Time(t).MarshalText()
	return string(str)
}

func (t DateTime) Value() (driver.Value, error) {
	return (time.Time)(t), nil
}

func (t *DateTime) Unix() int64 {
	return (*time.Time)(t).Unix()
}

func (t *DateTime) Nanosecond() int {
	return (*time.Time)(t).Nanosecond()
}

func (t *DateTime) UnixNano() int64 {
	return (*time.Time)(t).UnixNano()
}

func (r *DateTime) ExtensionType() int8 { return DateTime_EXT_TYPE }

func (r *DateTime) Len() int { return 12 }

func (r *DateTime) MarshalBinaryTo(b []byte) error {
	secs := uint64(r.Unix())
	nanoSecs := uint32(r.Nanosecond())
	binary.LittleEndian.PutUint64(b, secs)
	binary.LittleEndian.PutUint32(b[8:], nanoSecs)
	return nil
}

func (r *DateTime) UnmarshalBinary(b []byte) error {
	secs := int64(binary.LittleEndian.Uint64(b[0:8]))
	nanoSecs := int64(binary.LittleEndian.Uint32(b[8:12]))
	*r = DateTime(time.Unix(secs, nanoSecs))
	return nil
}

func (t DateTime) MarshalJSON() ([]byte, error) {
	str := strconv.FormatInt(t.UnixNano()/1000000, 10)
	return []byte(str), nil
}

func (t *DateTime) UnmarshalJSON(data []byte) error {
	i, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return err
	}
	*t = DateTime(time.Unix(int64(i/1000), (i%1000)*1000000))
	return nil
}
