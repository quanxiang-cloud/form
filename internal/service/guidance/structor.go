package guidance

import (
	"context"
	"fmt"

	"github.com/quanxiang-cloud/form/internal/service/consensus"
)

type structor struct{}

func newStructor() (Guidance, error) {
	return &structor{}, nil
}

func (s *structor) Do(ctx context.Context, bus *consensus.Bus) (*consensus.Response, error) {
	// TODO
	fmt.Println("TODO structor")
	return nil, nil
}
