package wrap

import (
	"context"

	"github.com/Anton-Kraev/ozon-route256-homeworks/internal/models/domain/wrap"
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
