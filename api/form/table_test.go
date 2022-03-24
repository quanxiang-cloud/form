package api

import (
	"fmt"
	id2 "github.com/quanxiang-cloud/cabin/id"
	"testing"
)

func TestName(t *testing.T) {
	uuid := id2.String(6)
	fmt.Println(uuid)

}
