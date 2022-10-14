package logic

import (
	"context"
	"github.com/eine-doodka/twoStage/cache"
	"github.com/eine-doodka/twoStage/customerrors"
	"github.com/google/uuid"
	"math/rand"
)

var (
	digits = []rune("0123456789")
)

type Logic interface {
	HandleInit(ctx context.Context) (uuid.UUID, string, error)
	HandleCommit(ctx context.Context, uid uuid.UUID, code string) error
}

type Impl struct {
	cache      cache.Cache
	codeLength int
}

func NewImpl(cache cache.Cache, codeLength int) *Impl {
	return &Impl{
		cache:      cache,
		codeLength: codeLength,
	}
}

func (l *Impl) randDigitStr(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = digits[rand.Intn(len(digits))]
	}
	return string(b)
}

func (l *Impl) HandleInit(ctx context.Context) (uuid.UUID, string, error) {
	randn := l.randDigitStr(l.codeLength)
	opid := uuid.New()
	err := l.cache.Set(ctx, opid, randn)
	return opid, randn, err
}

func (l *Impl) HandleCommit(ctx context.Context, uid uuid.UUID, code string) error {
	result, err := l.cache.Get(ctx, uid)
	if err != nil {
		return err
	}
	if result != code {
		return customerrors.ErrNotMatch
	}
	return nil
}
