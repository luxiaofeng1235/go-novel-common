package utils

//有层级关系的切片，通过父级id查找所有子级id数组
//parentId 父级id
//parentIndex 父级索引名称
//idIndex id索引名称
//id pid类型一定要相等 都是int64 否则无效
func FindSonByParentId(list []map[string]interface{}, parentId int64, parentIndex, idIndex string) []map[string]interface{}{
	newList := make([]map[string]interface{}, 0, len(list))
	for _, v := range list {
		if v[parentIndex] == parentId {
			newList = append(newList, v)
			fList := FindSonByParentId(list, v[idIndex].(int64), parentIndex, idIndex)
			newList = append(newList, fList...)
		}
	}
	return newList
}