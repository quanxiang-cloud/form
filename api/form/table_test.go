package api

import (
	"fmt"
	"testing"

	id2 "github.com/quanxiang-cloud/cabin/id"
)

func TestName(t *testing.T) {
	uuid := id2.String(6)
	fmt.Println(uuid)

}
