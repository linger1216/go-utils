package common

import "math"

type Page struct {
	Start int64
	End   int64
}

func SplitPage(total int64, n int) []*Page {
	pageSize := int64(math.Ceil(float64(total) / float64(n)))
	ret := make([]*Page, 0)
	if pageSize >= total {
		ret = append(ret, &Page{
			Start: 0,
			End:   total,
		})
	} else {
		for i := 0; i < n-1; i++ {
			ret = append(ret, &Page{
				Start: int64(i) * pageSize,
				End:   int64(i+1) * pageSize,
			})
		}
		ret = append(ret, &Page{
			Start: int64(n-1) * pageSize,
			End:   total,
		})
	}
	return ret
}
