package collection

type GoSlice []interface{}

func (rs GoSlice) ForEach(cb func(p interface{})) {
	for _, r := range rs {
		cb(r)
	}
}

func (rs GoSlice) MapToStrList(cb func(p interface{}) string) []string {
	res := make([]string, 0)
	for _, r := range rs {
		sr := cb(r)
		res = append(res, sr)
	}
	return res
}

func (rs GoSlice) Filter(cb func(p interface{}) bool) GoSlice {
	res := make(GoSlice, 0)
	for _, r := range rs {
		if ok := cb(r); ok {
			res = append(res, r)
		}
	}
	return res
}

func (rs GoSlice) FindIndex(cb func(p interface{}) bool) int {
	for idx, r := range rs {
		if ok := cb(r); ok {
			return idx
		}
	}
	return -1
}

func ConvStrList2GoSlice(sl []string) GoSlice {
	goSlice := make([]interface{}, len(sl))
	for i, v := range sl {
		goSlice[i] = v
	}
	return goSlice
}
