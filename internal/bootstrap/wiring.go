package bootstrap

import (
	"fmt"

	"vyshkat/internal/application"
	"vyshkat/internal/di"
	"vyshkat/internal/inspection"
	"vyshkat/internal/storage"
)

const (
	InspectorKey  di.Key = "inspector"
	RepositoryKey di.Key = "repository"
	FleetKey      di.Key = "fleet"
)

func NewContainer() *di.Container {
	container := di.New()

	must(container.Provide(InspectorKey, func(*di.Container) (any, error) {
		return application.Inspector(inspection.ServiceCenter{}), nil
	}))

	must(container.Provide(RepositoryKey, func(*di.Container) (any, error) {
		return application.Repository(storage.NewMemoryRepository()), nil
	}))

	must(container.Provide(FleetKey, func(deps *di.Container) (any, error) {
		repoValue, err := deps.Resolve(RepositoryKey)
		if err != nil {
			return nil, err
		}

		inspectorValue, err := deps.Resolve(InspectorKey)
		if err != nil {
			return nil, err
		}

		repo, ok := repoValue.(application.Repository)
		if !ok {
			return nil, fmt.Errorf("неверный тип repository: %T", repoValue)
		}

		inspector, ok := inspectorValue.(application.Inspector)
		if !ok {
			return nil, fmt.Errorf("неверный тип inspector: %T", inspectorValue)
		}

		return application.NewFleetService(repo, inspector)
	}))

	return container
}

func Fleet(container *di.Container) (*application.FleetService, error) {
	value, err := container.Resolve(FleetKey)
	if err != nil {
		return nil, err
	}

	park, ok := value.(*application.FleetService)
	if !ok {
		return nil, fmt.Errorf("неверный тип fleet: %T", value)
	}

	return park, nil
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
