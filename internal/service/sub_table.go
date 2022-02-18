package service

type SubTable interface {
}

type subTable struct {
}

func NewSubTable() (SubTable, error) {
	return &subTable{}, nil
}
