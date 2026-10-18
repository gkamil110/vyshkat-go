package di

import "fmt"

type Key string

type Factory func(*Container) (any, error)

type Container struct {
	factories map[Key]Factory
	instances map[Key]any
	resolving map[Key]bool
}

func New() *Container {
	return &Container{
		factories: make(map[Key]Factory),
		instances: make(map[Key]any),
		resolving: make(map[Key]bool),
	}
}

func (container *Container) Provide(key Key, factory Factory) error {
	if key == "" || factory == nil {
		return fmt.Errorf("неверная регистрация зависимости %q", key)
	}

	if _, found := container.factories[key]; found {
		return fmt.Errorf("зависимость %q уже зарегистрирована", key)
	}

	container.factories[key] = factory
	return nil
}

func (container *Container) Replace(key Key, factory Factory) error {
	if factory == nil {
		return fmt.Errorf("пустая фабрика для %q", key)
	}

	if _, found := container.factories[key]; !found {
		return fmt.Errorf("зависимость %q не зарегистрирована", key)
	}

	if len(container.instances) > 0 {
		return fmt.Errorf("нельзя менять зависимости после создания сервисов")
	}

	container.factories[key] = factory
	return nil
}

func (container *Container) Resolve(key Key) (any, error) {
	if saved, found := container.instances[key]; found {
		return saved, nil
	}

	create, found := container.factories[key]
	if !found {
		return nil, fmt.Errorf("зависимость %q не зарегистрирована", key)
	}

	if container.resolving[key] {
		return nil, fmt.Errorf("циклическая зависимость %q", key)
	}

	container.resolving[key] = true
	defer delete(container.resolving, key)

	value, err := create(container)
	if err != nil {
		return nil, fmt.Errorf("создание %q: %w", key, err)
	}

	if value == nil {
		return nil, fmt.Errorf("фабрика %q вернула nil", key)
	}

	container.instances[key] = value
	return value, nil
}
