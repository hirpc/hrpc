# 生成全局唯一ID

示例：

```go
func TestGetUniqueID(t *testing.T) {
	err := Register()
	if err != nil {
		t.Error(err)
	}

	nodeMap := make(map[int64]int)
	for i := 1; i < 10000; i++ {
		uniqueID := New().Number()
		if id, ok := nodeMap[uniqueID]; ok {
			t.Error(id)
		} else {
			nodeMap[uniqueID] = i
		}
	}
  t.Log()
}
```





## todo

当前生成的上一个数字形式的全局id，以后待补充一个UUID

实现UUID可参考https://github.com/google/uuid

北京时间2021年10月08日00:06:42，UUID暂时不考虑做入这个包内，也许会直接使用Google的包

## ref. Copied from houzhenkai