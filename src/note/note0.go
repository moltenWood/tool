package note

import (
	"fmt"
	"sort"
)

type info struct {
	Name string
	Time int
}

func demoForSort() {
	var list newlist = newlist{info{"456",456}, info{"789",789}, info{"123",123}}
	sort.Sort(list)  //调用标准库的sort.Sort必须要先实现Len(),Less(),Swap() 三个方法.
	fmt.Println(list)
}

type newlist []info

func (I newlist) Len() int {
	return len(I)
}
func (I newlist) Less(i, j int) bool {
	return I[i].Time < I[j].Time
}
func (I newlist) Swap(i, j int) {
	I[i], I[j] = I[j], I[i]
}
