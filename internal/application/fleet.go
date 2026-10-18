package application

import (
	"errors"
	"fmt"

	"vyshkat/internal/domain"
)

const MinimumBeginnerEase = 6

var ErrTransportRequiresInspection = errors.New("транспорт добавляется только через техосмотр")

type InspectionResult struct {
	Approved bool
	Reason   string
}

// Проверяет транспорт перед добавлением в парк
type Inspector interface {
	Inspect(domain.Transport) InspectionResult
}

// Хранит всю технику и остальные вещи
type Repository interface {
	Add(domain.InventoryItem) error
	All() []domain.InventoryItem
}

type Summary struct {
	TransportCount      int
	OtherItemCount      int
	TotalCount          int
	TotalDailyEnergyKWh float64
}

type FleetService struct {
	repo      Repository
	inspector Inspector
}

func NewFleetService(repo Repository, inspector Inspector) (*FleetService, error) {
	if repo == nil || inspector == nil {
		return nil, errors.New("репозиторий и техосмотр обязательны")
	}

	return &FleetService{
		repo:      repo,
		inspector: inspector,
	}, nil
}

func (fleet *FleetService) AdmitTransport(transport domain.Transport) (InspectionResult, error) {
	if transport == nil {
		return InspectionResult{}, errors.New("транспорт не задан")
	}

	check := fleet.inspector.Inspect(transport)

	if !check.Approved {
		return check, nil
	}

	err := fleet.repo.Add(transport)
	if err != nil {
		return InspectionResult{}, fmt.Errorf("приёмка: %w", err)
	}

	return check, nil
}

func (fleet *FleetService) RegisterItem(item domain.InventoryItem) error {
	if item == nil {
		return errors.New("вещь не задана")
	}

	if _, isTransport := item.(domain.Transport); isTransport {
		return ErrTransportRequiresInspection
	}

	return fleet.repo.Add(item)
}

func (fleet *FleetService) Inventory() []domain.InventoryItem {
	return fleet.repo.All()
}

func (fleet *FleetService) BeginnerTransport() []domain.Transport {
	var beginners []domain.Transport

	for _, item := range fleet.repo.All() {
		transport, ok := item.(domain.Transport)
		if !ok {
			continue
		}

		if transport.EaseOfUse() < MinimumBeginnerEase {
			continue
		}

		beginners = append(beginners, transport)
	}

	return beginners
}

func (fleet *FleetService) FleetSummary() Summary {
	var stats Summary

	for _, item := range fleet.repo.All() {
		stats.TotalCount++

		if _, isTransport := item.(domain.Transport); isTransport {
			stats.TransportCount++
		} else {
			stats.OtherItemCount++
		}

		if device, usesEnergy := item.(domain.EnergyConsumer); usesEnergy {
			stats.TotalDailyEnergyKWh += device.DailyEnergyKWh()
		}
	}

	return stats
}
