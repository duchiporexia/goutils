package xutils

func If(cond bool, v0 string, v1 string) string {
	if cond {
		return v0
	}
	return v1
}
