package di

import (
	"errors"
	"strings"
	"testing"
)

func TestResolveCachingAndValidation(t *testing.T) {
	container := New()

	if err := container.Provide("", func(*Container) (any, error) {
		return 1, nil
	}); err == nil {
		t.Fatal("можно добавить зависимость без ключа")
	}

	if err := container.Provide("value", nil); err == nil {
		t.Fatal("можно добавить пустую фабрику")
	}

	calls := 0
	err := container.Provide("value", func(*Container) (any, error) {
		calls++
		return 42, nil
	})
	if err != nil {
		t.Fatal(err)
	}

	if err := container.Provide("value", func(*Container) (any, error) {
		return 0, nil
	}); err == nil {
		t.Fatal("получилось дважды зарегистрировать зависимость")
	}

	for i := 0; i < 2; i++ {
		value, err := container.Resolve("value")
		if err != nil || value != 42 {
			t.Fatalf("неверное значение: %v, ошибка: %v", value, err)
		}
	}

	if calls != 1 {
		t.Fatalf("фабрика вызвалась %d раз вместо одного", calls)
	}

	if _, err := container.Resolve("missing"); err == nil {
		t.Fatal("неизвестная зависимость не вызвала ошибку")
	}

	if err := container.Replace("value", func(*Container) (any, error) {
		return 9, nil
	}); err == nil {
		t.Fatal("получилось заменить уже созданную зависимость")
	}
}

func TestCycleAndFailure(t *testing.T) {
	container := New()

	if err := container.Replace("unknown", func(*Container) (any, error) {
		return 1, nil
	}); err == nil {
		t.Fatal("получилось заменить несуществующую зависимость")
	}

	err := container.Provide("a", func(deps *Container) (any, error) {
		return deps.Resolve("b")
	})
	if err != nil {
		t.Fatal(err)
	}

	err = container.Provide("b", func(deps *Container) (any, error) {
		return deps.Resolve("a")
	})
	if err != nil {
		t.Fatal(err)
	}

	if _, err := container.Resolve("a"); err == nil || !strings.Contains(err.Error(), "циклическая") {
		t.Fatalf("не обнаружена циклическая зависимость: %v", err)
	}

	err = container.Provide("error", func(*Container) (any, error) {
		return nil, errors.New("oops")
	})
	if err != nil {
		t.Fatal(err)
	}

	if _, err := container.Resolve("error"); err == nil {
		t.Fatal("ошибка фабрики не обработана")
	}

	err = container.Provide("nil", func(*Container) (any, error) {
		return nil, nil
	})
	if err != nil {
		t.Fatal(err)
	}

	if _, err := container.Resolve("nil"); err == nil {
		t.Fatal("фабрика вернула nil, но ошибки нет")
	}
}

func TestReplaceBeforeResolution(t *testing.T) {
	container := New()

	err := container.Provide("k", func(*Container) (any, error) {
		return "original", nil
	})
	if err != nil {
		t.Fatal(err)
	}

	if err := container.Replace("k", nil); err == nil {
		t.Fatal("можно заменить зависимость пустой фабрикой")
	}

	err = container.Replace("k", func(*Container) (any, error) {
		return "fake", nil
	})
	if err != nil {
		t.Fatal(err)
	}

	value, err := container.Resolve("k")
	if err != nil || value != "fake" {
		t.Fatalf("зависимость не заменилась: %v, %v", value, err)
	}
}
