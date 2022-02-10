package form

import (
	"context"
	"errors"
)

type Components interface {
	GetTag() string
	SetValue(req *comReq)
	HandlerFunc(ctx context.Context, action string) error
}

var cs = []Components{
	&subTable{},
}

var (
	ErrNoComponents = errors.New("no ErrNoComponents like this")
)

// Component Container for components
type Component struct {
	com map[string]Components
}

// NewCom return component instance
func NewCom() *Component {
	c := &Component{
		com: make(map[string]Components, len(cs)),
	}
	for _, component := range cs {
		c.com[component.GetTag()] = component
	}
	return c
}

// GetCom build a component
func (c *Component) GetCom(tag string, req *comReq) (Components, error) {
	com, ok := c.com[tag] // 获取组件
	if !ok {
		return nil, ErrNoComponents
	}
	com.SetValue(req)
	return com, nil
}
