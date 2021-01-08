package db

type ResultSet []interface{}

func (rs ResultSet) ForEach(cb func(p interface{})) {
	for _, r := range rs {
		cb(r)
	}
}

func (rs ResultSet) MapToStrList(cb func(p interface{}) string) []string {
	res := make([]string, 0)
	for _, r := range rs {
		sr := cb(r)
		res = append(res, sr)
	}
	return res
}

func (rs ResultSet) Filter(cb func(p interface{}) bool) ResultSet {
	res := make(ResultSet, 0)
	for _, r := range rs {
		if ok := cb(r); ok {
			res = append(res, r)
		}
	}
	return res
}
