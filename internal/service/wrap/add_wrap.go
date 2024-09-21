package wrap

import (
	"context"

	errsdomain "github.com/Anton-Kraev/ozon-route256-homeworks/internal/models/domain/errors"
	"github.com/Anton-Kraev/ozon-route256-homeworks/internal/models/domain/wrap"
)

func (s WrapService) AddWrap(ctx context.Context, wrap wrap.Wrap) error {
	sameWrap, err := s.wrapRepository.GetWrapByName(ctx, wrap.Name)
	if err != nil {
		return err
	}

	if sameWrap != nil {
		return errsdomain.ErrWrapAlreadyExists
	}

	return s.wrapRepository.AddWrap(ctx, wrap)
}
