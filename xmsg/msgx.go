package xmsg

type MsgMarshaler interface {
	MarshalMsg(b []byte) (o []byte, err error)
}
type MsgUnmarshaler interface {
	UnmarshalMsg(bts []byte) (o []byte, err error)
}
