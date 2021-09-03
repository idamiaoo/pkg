package http

// CallOption 请求设置
type CallOption interface {
	Before(*Request) error
	After(*Response) error
}

func combine(o1 []CallOption, o2 []CallOption) []CallOption {
	if len(o1) == 0 {
		return o2
	} else if len(o2) == 0 {
		return o1
	}
	ret := make([]CallOption, len(o1)+len(o2))
	copy(ret, o1)
	copy(ret[len(o1):], o2)
	return ret
}

type NopCallOption struct {
}

func (NopCallOption) Before(r *Request) error {
	return nil
}

func (NopCallOption) After(*Response) error {
	return nil
}
