package form

import (
	"context"
	"errors"
)

type components interface {
	getTag() string
	setValue(req *comReq)
	handlerFunc(ctx context.Context, action string) error
}

var cs = []components{
	&subTable{},
}

var (
	ErrNoComponents = errors.New("no ErrNoComponents like this")
)

// Component Container for components
type component struct {
	com map[string]components
}

// NewCom return component instance
func newFormComponent() *component {
	c := &component{
		com: make(map[string]components, len(cs)),
	}
	for _, component := range cs {
		c.com[component.getTag()] = component
	}
	return c
}

// GetCom build a component
func (c *component) getCom(tag string, req *comReq) (components, error) {
	com, ok := c.com[tag] // 获取组件
	if !ok {
		return nil, ErrNoComponents
	}
	com.setValue(req)
	return com, nil
}
