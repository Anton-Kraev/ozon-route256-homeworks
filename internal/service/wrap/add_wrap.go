package wrap

import (
	"context"

	errsdomain "gitlab.ozon.dev/antonkraeww/ozon-route256-homeworks/internal/models/domain/errors"
	"gitlab.ozon.dev/antonkraeww/ozon-route256-homeworks/internal/models/domain/wrap"
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
