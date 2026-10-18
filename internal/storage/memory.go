package storage

import (
	"errors"
	"fmt"

	"vyshkat/internal/domain"
)

var ErrDuplicate = errors.New("инвентарный номер уже существует")

type MemoryRepository struct {
	records []domain.InventoryItem
	ids     map[string]struct{}
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		ids: make(map[string]struct{}),
	}
}

func (repo *MemoryRepository) Add(item domain.InventoryItem) error {
	if item == nil {
		return errors.New("нельзя сохранить пустую запись")
	}

	if _, found := repo.ids[item.InventoryNumber()]; found {
		return fmt.Errorf("%w: %s", ErrDuplicate, item.InventoryNumber())
	}

	repo.ids[item.InventoryNumber()] = struct{}{}
	repo.records = append(repo.records, item)

	return nil
}

func (repo *MemoryRepository) All() []domain.InventoryItem {
	items := make([]domain.InventoryItem, len(repo.records))
	copy(items, repo.records)

	return items
}
