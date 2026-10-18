package storage

import (
	"errors"
	"testing"

	"vyshkat/internal/domain"
)

func TestMemoryRepository(t *testing.T) {
	repo := NewMemoryRepository()

	if err := repo.Add(nil); err == nil {
		t.Fatal("получилось добавить пустую запись")
	}

	helmet, err := domain.NewHelmet("a-1", "Первый")
	if err != nil {
		t.Fatal(err)
	}

	dock, err := domain.NewDockStation("A-1", "Второй")
	if err != nil {
		t.Fatal(err)
	}

	if err := repo.Add(helmet); err != nil {
		t.Fatal(err)
	}

	if err := repo.Add(dock); !errors.Is(err, ErrDuplicate) {
		t.Fatalf("ожидалась ошибка повторного номера: %v", err)
	}

	items := repo.All()

	if len(items) != 1 || items[0].DisplayName() != "Первый" {
		t.Fatalf("неверный список вещей: %v", items)
	}

	items[0] = nil

	if repo.All()[0] == nil {
		t.Fatal("изменение полученного списка затронуло репозиторий")
	}
}
