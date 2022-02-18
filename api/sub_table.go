package api

import "github.com/gin-gonic/gin"

type SubTable struct {
}

func NewSubTable() (*SubTable, error) {
	return &SubTable{}, nil
}

func (s *SubTable) CreateSubTable(c *gin.Context) {

}
