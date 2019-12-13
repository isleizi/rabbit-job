package str

import "fmt"

// 统一转string
func ToStr(val interface{}) string {
	return fmt.Sprintf("%v", val)
}
