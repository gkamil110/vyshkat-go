package domain

import (
	"math"
	"strings"
	"testing"
)

var safe = SafetyStatus{
	BrakesOK:     true,
	SteeringOK:   true,
	ElectricalOK: true,
}

func TestConstructorsAndInterfaces(t *testing.T) {
	scooter, err := NewScooter("  a-01 ", "  Самокат ", 6, 1.25, safe)
	if err != nil {
		t.Fatal(err)
	}

	if scooter.InventoryNumber() != "A-01" ||
		scooter.DisplayName() != "Самокат" ||
		scooter.EaseOfUse() != 6 {
		t.Fatalf("неверные данные самоката: %#v", scooter)
	}

	var _ InventoryItem = scooter
	var _ Transport = scooter
	var _ EnergyConsumer = scooter

	bike, err := NewBicycle("B-01", "Велосипед", 2, safe)
	if err != nil {
		t.Fatal(err)
	}

	if _, usesEnergy := any(bike).(EnergyConsumer); usesEnergy {
		t.Fatal("обычный велосипед не должен тратить электричество")
	}

	ebike, err := NewEBike("E-01", "Электровелосипед", 7, 1.5, safe)
	if err != nil {
		t.Fatal(err)
	}

	if ebike.DailyEnergyKWh() != 1.5 {
		t.Fatalf("неверный расход энергии электровелосипеда: %v", ebike.DailyEnergyKWh())
	}

	helmet, err := NewHelmet("H-01", "Шлем")
	if err != nil {
		t.Fatal(err)
	}

	if _, usesEnergy := any(helmet).(EnergyConsumer); usesEnergy {
		t.Fatal("шлем не должен потреблять энергию")
	}

	dock, err := NewDockStation("D-01", "Док")
	if err != nil {
		t.Fatal(err)
	}

	if dock.Category() != "Док-станция" {
		t.Fatalf("неверная категория док-станции: %s", dock.Category())
	}

	cabinet, err := NewChargingCabinet("C-01", "Зарядный шкаф", 0.4)
	if err != nil {
		t.Fatal(err)
	}

	if cabinet.DailyEnergyKWh() != 0.4 {
		t.Fatalf("неверный расход энергии шкафа: %v", cabinet.DailyEnergyKWh())
	}
}

func TestInvalidDomainData(t *testing.T) {
	tests := []struct {
		name string
		make func() error
	}{
		{"пустой номер", func() error {
			_, err := NewHelmet(" ", "Шлем")
			return err
		}},
		{"пустое название", func() error {
			_, err := NewHelmet("ID", " ")
			return err
		}},
		{"длинный номер", func() error {
			_, err := NewHelmet(strings.Repeat("x", 41), "Шлем")
			return err
		}},
		{"длинное название", func() error {
			_, err := NewHelmet("ID", strings.Repeat("x", 101))
			return err
		}},
		{"простота 0", func() error {
			_, err := NewBicycle("B", "Bike", 0, safe)
			return err
		}},
		{"простота 11", func() error {
			_, err := NewBicycle("B", "Bike", 11, safe)
			return err
		}},
		{"нулевая энергия самоката", func() error {
			_, err := NewScooter("S", "Scooter", 6, 0, safe)
			return err
		}},
		{"отрицательная энергия", func() error {
			_, err := NewEBike("E", "EBike", 6, -1, safe)
			return err
		}},
		{"энергия NaN", func() error {
			_, err := NewChargingCabinet("C", "Cabinet", math.NaN())
			return err
		}},
		{"энергия Inf", func() error {
			_, err := NewChargingCabinet("C", "Cabinet", math.Inf(1))
			return err
		}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := test.make(); err == nil {
				t.Fatal("неверные данные не вызвали ошибку")
			}
		})
	}
}
