package node

type Http struct {
	host string
}

func NewHttp(host string) *Http {
	return &Http{host}
}
