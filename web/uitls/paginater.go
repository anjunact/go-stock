package utils

const SHOW_NUM = 9

type Paginater struct {
	Total int
	Page  int
	Prev  int
	Next  int
	First int
	Last  int
	Pages []int
}

func New(Total int, Page int, pageSize int) *Paginater {
	// if Page < 0 {
	// 	Page = 1
	// }
	Prev := Page - 1
	Next := Page + 1
	First := 1
	Last := Total / pageSize
	if Last*pageSize < Total {
		Last++
	}
	Pages := []int{}
	for i := Page - SHOW_NUM/2; i < Page+SHOW_NUM/2; i++ {
		if i < 0 {
			i = 1
		}
		if i > Last {
			i = Last
		}
		Pages = append(Pages, i)
	}

	p := &Paginater{Total, Page, Prev, Next, First, Last, Pages}
	return p
}
