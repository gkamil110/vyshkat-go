package bootstrap

import (
	"testing"

	"vyshkat/internal/application"
	"vyshkat/internal/di"
	"vyshkat/internal/domain"
)

type fakeCenter struct {
	allow     bool
	inspected int
}

func (f *fakeCenter) Inspect(domain.Transport) application.InspectionResult {
	f.inspected++

	return application.InspectionResult{
		Approved: f.allow,
		Reason:   "решение тестовой подмены",
	}
}

func TestCanReplaceInspectionThroughDI(t *testing.T) {
	for _, approved := range []bool{false, true} {
		container := NewContainer()
		center := &fakeCenter{allow: approved}

		err := container.Replace(InspectorKey, func(*di.Container) (any, error) {
			return application.Inspector(center), nil
		})
		if err != nil {
			t.Fatal(err)
		}

		park, err := Fleet(container)
		if err != nil {
			t.Fatal(err)
		}

		scooter, err := domain.NewScooter(
			"TEST-1", "Тестовый", 6, 1.0,
			domain.SafetyStatus{
				BrakesOK:     false,
				SteeringOK:   true,
				ElectricalOK: true,
			},
		)
		if err != nil {
			t.Fatal(err)
		}

		check, err := park.AdmitTransport(scooter)
		if err != nil || check.Approved != approved || center.inspected != 1 {
			t.Fatalf("неверный результат техосмотра: %+v, %v", check, err)
		}

		count := park.FleetSummary().TransportCount
		if (count == 1) != approved {
			t.Fatalf("неверное количество транспорта: %d", count)
		}

		err = container.Replace(InspectorKey, func(*di.Container) (any, error) {
			return center, nil
		})
		if err == nil {
			t.Fatal("нельзя менять зависимость после её использования")
		}
	}
}

func TestRealContainerReusesFleet(t *testing.T) {
	container := NewContainer()

	first, err := Fleet(container)
	if err != nil {
		t.Fatal(err)
	}

	second, err := Fleet(container)
	if err != nil || first != second {
		t.Fatal("контейнер создал два разных сервиса")
	}
}
