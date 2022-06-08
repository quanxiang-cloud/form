package router

import (
	"fmt"
	"strings"
	"testing"
)

func TestName(t *testing.T) {
	s1 := "/xxxx/home/form/xxx/ddd"

	any := strings.LastIndex(s1, "/")

	str1 := s1[0:any]

	fmt.Sprint(str1)
	sprintf := fmt.Sprintf("%s/:id", str1)

	fmt.Println(sprintf)
}
