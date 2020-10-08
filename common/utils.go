package common

import "math"

type Page struct {
	Start int64
	End   int64
}

// n 为页数
func SplitPage(total int64, page int) []*Page {

	if total == 0 || page == 0 {
		return nil
	}

	// 总数量还没有页码多, 比如total=5, page=30
	// 那就分成5页即可
	ret := make([]*Page, 0)
	if total < int64(page) {
		for i := int64(0); i < total; i++ {
			s, e := i, i+1
			if e > total {
				e = total
			}
			ret = append(ret, &Page{
				Start: s,
				End:   e,
			})
		}
		return ret
	}

	// 应该每页是多少?
	// total:47 page:30 => pageSize:2
	pageSize := int64(math.Ceil(float64(total) / float64(page)))

	// 每页的大小, 比总量还大
	// total:47 page:1 => pageSize:47
	if pageSize >= total {
		ret = append(ret, &Page{
			Start: 0,
			End:   total,
		})
		return ret
	}

	//
	for i := 0; i < page; i++ {
		s, e := int64(i)*pageSize, int64(i+1)*pageSize
		if e > total {
			e = total
		}
		if s >= total {
			break
		}
		ret = append(ret, &Page{
			Start: s,
			End:   e,
		})
	}
	return ret
}
