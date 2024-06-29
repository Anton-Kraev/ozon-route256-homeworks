//go:generate mockgen -package=wrap -source=./service.go -destination=./service_mocks.go
package wrap

import (
	"context"

	"gitlab.ozon.dev/antonkraeww/homeworks/hw-5/internal/models/domain/wrap"
)

type wrapRepository interface {
	GetWrapByName(ctx context.Context, name string) (*wrap.Wrap, error)
	AddWrap(ctx context.Context, wrap wrap.Wrap) error
}

type WrapService struct {
	wrapRepository wrapRepository
}

func NewWrapService(wrapRepository wrapRepository) WrapService {
	return WrapService{wrapRepository: wrapRepository}
}
