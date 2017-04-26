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

func NewPaginater(Total int, Page int, pageSize int) (p *Paginater ){
	First := 1
	Last := Total / pageSize
	if Last*pageSize < Total {
		Last++
	}
	Prev := Page - 1
	if Page == 1{
		Prev =1
	}

	Next := Page + 1
	if Page ==  Last{
		Next = Last
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
		if i==Last{
			break
		}
	}
	p = &Paginater{Total, Page, Prev, Next, First, Last, Pages}
	return
}
