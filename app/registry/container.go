package registry

import (
	"github.com/sarulabs/di"

	"go-prj-skeleton/app/interface/persistence/postgre"
	"go-prj-skeleton/app/usecase"
)

type Container struct {
	ctn di.Container
}

func NewContainer() (*Container, error) {
	builder, err := di.NewBuilder()
	if err != nil {
		return nil, err
	}

	if err := builder.Add([]di.Def{
		{
			Name:  "user-usecase",
			Build: buildUserUsecase,
		},
	}...); err != nil {
		return nil, err
	}

	return &Container{
		ctn: builder.Build(),
	}, nil
}

func (c *Container) Resolve(name string) interface{} {
	return c.ctn.Get(name)
}

func (c *Container) Clean() error {
	return c.ctn.Clean()
}

func buildUserUsecase(ctn di.Container) (interface{}, error) {
	userRepo := postgre.NewUserRepo()
	accountRepo := postgre.NewAccountRepo()
	tranRepo := postgre.NewTransactionRepo()
	return usecase.NewUserUsecase(userRepo, accountRepo, tranRepo), nil
}
