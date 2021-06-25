package xfiber

//go:generate msgp --tests=false
type ErrResponse struct {
	Code int    `msg:"code"`
	Msg  string `msg:"msg"`
}
