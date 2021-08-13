package goToolsUtils

// 范例
//func main() {
//	var items goToolsUtils.ItemsInterface
//	var target Target
//
//	items = &target
//	target.Content = []map[string]float64{{"Value":0.3},{"Value":0.1},{"Value":0.2},{"Value":0.5}}
//	goToolsUtils.QuickSort2(items)
//	fmt.Println(target.Content)
//}
//
//type Target struct {
//	Content []map[string]float64
//	ContentInterface []interface{}
//}
//
//func (t *Target) OutputResult(items []interface{}) {
//	var output []map[string]float64
//	for _,item := range items{
//		output = append(output, item.(map[string]float64))
//	}
//	t.Content = output
//}
//
//func (t *Target) GetContentSlice() *[]interface{} {
//	if len(t.ContentInterface) ==0{
//		var Con []interface{}
//		for _,i :=range t.Content{
//			Con = append(Con, i )
//		}
//		t.ContentInterface = Con
//	}
//
//	return &t.ContentInterface
//}
//
//func (t *Target) GetReference(index int) float64 {
//	ret :=  t.ContentInterface[index].(map[string]float64)
//	return ret["Value"]
//}

type ItemsInterface interface {
	GetReference(index int) float64
	GetContentSlice() *[]interface{} //切片的指针指向切片
	OutputResult([]interface{}) // 将返回的[]interface{} 强转回去并赋值给 原属性
}

func QuickSort2(itemsInterface ItemsInterface)  {
	if len(*itemsInterface.GetContentSlice()) < 2 {
		return
	}
	var stack []int
	stack = append(stack, len(*itemsInterface.GetContentSlice())-1)
	stack = append(stack, 0)
	for len(stack) > 0 {
		l := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		r := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		index := partition2(itemsInterface, l, r)

		if l < index-1 {
			stack = append(stack, index-1)
			stack = append(stack, l)
		}
		if r > index+1 {
			stack = append(stack, r)
			stack = append(stack, index+1)
		}
	}
	itemsInterface.OutputResult(*itemsInterface.GetContentSlice())
}

func partition2(itemsInterface ItemsInterface, start int, end int) int {
	pivot := itemsInterface.GetReference(start)
	pivotItem := (*itemsInterface.GetContentSlice())[start]
	for start < end {
		for start < end && itemsInterface.GetReference(end) >= pivot {
			end -= 1
		}
		(*itemsInterface.GetContentSlice())[start] = (*itemsInterface.GetContentSlice())[end]
		for start < end && itemsInterface.GetReference(start)  <= pivot {
			start += 1
		}
		(*itemsInterface.GetContentSlice())[end] = (*itemsInterface.GetContentSlice())[start]
	}
	(*itemsInterface.GetContentSlice())[start] = pivotItem
	return start
}
