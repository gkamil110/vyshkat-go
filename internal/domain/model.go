package domain

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"unicode/utf8"
)

type InventoryItem interface {
	InventoryNumber() string
	DisplayName() string
	Category() string
}

type EnergyConsumer interface {
	DailyEnergyKWh() float64
}

type Transport interface {
	InventoryItem
	EaseOfUse() int
	InspectionState() SafetyStatus
}

type SafetyStatus struct {
	BrakesOK     bool
	SteeringOK   bool
	ElectricalOK bool
}

type baseItem struct {
	id   string
	name string
}

func (item baseItem) InventoryNumber() string {
	return item.id
}

func (item baseItem) DisplayName() string {
	return item.name
}

func newBaseItem(number, name string) (baseItem, error) {
	number = strings.ToUpper(strings.TrimSpace(number))
	name = strings.TrimSpace(name)

	if number == "" {
		return baseItem{}, errors.New("инвентарный номер не может быть пустым")
	}

	if name == "" {
		return baseItem{}, errors.New("наименование не может быть пустым")
	}

	if utf8.RuneCountInString(number) > 40 || utf8.RuneCountInString(name) > 100 {
		return baseItem{}, errors.New("инвентарный номер или наименование слишком длинные")
	}

	return baseItem{id: number, name: name}, nil
}

type transportBase struct {
	baseItem
	ease   int
	safety SafetyStatus
}

func (transport transportBase) EaseOfUse() int {
	return transport.ease
}

func (transport transportBase) InspectionState() SafetyStatus {
	return transport.safety
}

func newTransportBase(number, name string, ease int, safety SafetyStatus) (transportBase, error) {
	item, err := newBaseItem(number, name)
	if err != nil {
		return transportBase{}, err
	}

	if ease < 1 || ease > 10 {
		return transportBase{}, fmt.Errorf("простота использования должна быть от 1 до 10, получено %d", ease)
	}

	return transportBase{
		baseItem: item,
		ease:     ease,
		safety:   safety,
	}, nil
}

func validateEnergy(energy float64) error {
	if math.IsNaN(energy) || math.IsInf(energy, 0) || energy <= 0 {
		return errors.New("суточная энергия должна быть конечным положительным числом")
	}

	return nil
}

type Scooter struct {
	transportBase
	energy float64
}

func NewScooter(number, name string, ease int, energy float64, safety SafetyStatus) (*Scooter, error) {
	transport, err := newTransportBase(number, name, ease, safety)
	if err != nil {
		return nil, err
	}

	if err := validateEnergy(energy); err != nil {
		return nil, err
	}

	return &Scooter{
		transportBase: transport,
		energy:        energy,
	}, nil
}

func (scooter *Scooter) Category() string {
	return "Электросамокат"
}

func (scooter *Scooter) DailyEnergyKWh() float64 {
	return scooter.energy
}

type Bicycle struct {
	transportBase
}

func NewBicycle(number, name string, ease int, safety SafetyStatus) (*Bicycle, error) {
	transport, err := newTransportBase(number, name, ease, safety)
	if err != nil {
		return nil, err
	}

	return &Bicycle{
		transportBase: transport,
	}, nil
}

func (bike *Bicycle) Category() string {
	return "Велосипед"
}

type EBike struct {
	transportBase
	energy float64
}

func NewEBike(number, name string, ease int, energy float64, safety SafetyStatus) (*EBike, error) {
	transport, err := newTransportBase(number, name, ease, safety)
	if err != nil {
		return nil, err
	}

	if err := validateEnergy(energy); err != nil {
		return nil, err
	}

	return &EBike{
		transportBase: transport,
		energy:        energy,
	}, nil
}

func (bike *EBike) Category() string {
	return "Электровелосипед"
}

func (bike *EBike) DailyEnergyKWh() float64 {
	return bike.energy
}

type Helmet struct {
	baseItem
}

func NewHelmet(number, name string) (*Helmet, error) {
	item, err := newBaseItem(number, name)
	if err != nil {
		return nil, err
	}

	return &Helmet{baseItem: item}, nil
}

func (helmet *Helmet) Category() string {
	return "Шлем"
}

type DockStation struct {
	baseItem
}

func NewDockStation(number, name string) (*DockStation, error) {
	item, err := newBaseItem(number, name)
	if err != nil {
		return nil, err
	}

	return &DockStation{baseItem: item}, nil
}

func (dock *DockStation) Category() string {
	return "Док-станция"
}

type ChargingCabinet struct {
	baseItem
	energy float64
}

func NewChargingCabinet(number, name string, energy float64) (*ChargingCabinet, error) {
	item, err := newBaseItem(number, name)
	if err != nil {
		return nil, err
	}

	if err := validateEnergy(energy); err != nil {
		return nil, err
	}

	return &ChargingCabinet{
		baseItem: item,
		energy:   energy,
	}, nil
}

func (cabinet *ChargingCabinet) Category() string {
	return "Зарядный шкаф"
}

func (cabinet *ChargingCabinet) DailyEnergyKWh() float64 {
	return cabinet.energy
}
