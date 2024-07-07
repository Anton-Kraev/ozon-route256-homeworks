package wrap

import (
	"context"

	"gitlab.ozon.dev/antonkraeww/homeworks/hw-6/internal/models/domain/wrap"
)

type wrapService interface {
	AddWrap(ctx context.Context, wrap wrap.Wrap) error
}

type WrapController struct {
	wrapService wrapService
}

func NewWrapController(wrapService wrapService) WrapController {
	return WrapController{wrapService: wrapService}
}
