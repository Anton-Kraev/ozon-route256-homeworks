//go:generate mockgen -package=wrap -source=./service.go -destination=./service_mocks.go

package wrap

import (
	"context"

	"github.com/Anton-Kraev/ozon-route256-homeworks/internal/models/domain/wrap"
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
