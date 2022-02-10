package form

import "github.com/quanxiang-cloud/form/internal/service"

// 需要在鉴权和不鉴权 的 都加上， 发消息的机制

type flow struct {
	permit service.Permission

	*comet

	*auth
}
