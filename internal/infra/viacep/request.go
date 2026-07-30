package viacep

type Request struct {
	Zipcode string
}

func NewRequest(zipcode string) *Request {
	return &Request{
		Zipcode: zipcode,
	}
}
