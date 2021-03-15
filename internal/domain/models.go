package domain

type StoreRequest struct {
	Method string
	Path string
	Req string
}

type Request struct {
	Id int32
	Method string
	Path string
	Req string
}

type RequestShort struct {
	Id int32
	Method string
	Path string
}
