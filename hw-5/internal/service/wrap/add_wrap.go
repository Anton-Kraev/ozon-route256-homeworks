package wrap

import (
	"context"

	errsdomain "gitlab.ozon.dev/antonkraeww/homeworks/hw-5/internal/models/domain/errors"
	"gitlab.ozon.dev/antonkraeww/homeworks/hw-5/internal/models/domain/wrap"
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
