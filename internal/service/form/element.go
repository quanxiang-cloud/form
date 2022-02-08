package form

import "context"

type subTable struct {
}

func (s *subTable) GetTag() string {
	return "sub_table"
}

func (s *subTable) Create(ctx context.Context) error {
	return nil
}

func (s *subTable) Update(ctx context.Context) error {
	return nil
}

func (s *subTable) FindOne(ctx context.Context) error {
	return nil
}
