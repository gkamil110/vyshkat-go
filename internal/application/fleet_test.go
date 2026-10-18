package application

import (
	"errors"
	"math"
	"testing"

	"vyshkat/internal/domain"
	"vyshkat/internal/storage"
)

type fakeInspector struct {
	result InspectionResult
	calls  int
}

func (f *fakeInspector) Inspect(domain.Transport) InspectionResult {
	f.calls++
	return f.result
}

func newFleet(t *testing.T, inspector Inspector) *FleetService {
	t.Helper()

	park, err := NewFleetService(storage.NewMemoryRepository(), inspector)
	if err != nil {
		t.Fatal(err)
	}

	return park
}

func TestAdmissionRequiresInspection(t *testing.T) {
	safety := domain.SafetyStatus{BrakesOK: true, SteeringOK: true, ElectricalOK: true}
	scooter, _ := domain.NewScooter("S-1", "Самокат", 10, 1.1, safety)

	inspector := &fakeInspector{
		result: InspectionResult{Approved: false, Reason: "не прошёл техосмотр"},
	}
	park := newFleet(t, inspector)

	check, err := park.AdmitTransport(scooter)
	if err != nil || check.Approved || inspector.calls != 1 || len(park.Inventory()) != 0 {
		t.Fatalf("самокат не должен был попасть в парк: %+v, %v", check, err)
	}

	if err := park.RegisterItem(scooter); !errors.Is(err, ErrTransportRequiresInspection) {
		t.Fatalf("получилось добавить транспорт без проверки: %v", err)
	}

	inspector.result = InspectionResult{Approved: true, Reason: "одобрено"}

	check, err = park.AdmitTransport(scooter)
	if err != nil || !check.Approved || inspector.calls != 2 || len(park.Inventory()) != 1 {
		t.Fatalf("самокат должен был пройти проверку: %+v, %v", check, err)
	}

	_, err = park.AdmitTransport(scooter)
	if !errors.Is(err, storage.ErrDuplicate) {
		t.Fatalf("повторное добавление должно выдать ошибку: %v", err)
	}
}

func TestSummaryAndBeginnerThreshold(t *testing.T) {
	safety := domain.SafetyStatus{BrakesOK: true, SteeringOK: true, ElectricalOK: true}

	inspector := &fakeInspector{
		result: InspectionResult{Approved: true, Reason: "test"},
	}
	park := newFleet(t, inspector)

	scooter, _ := domain.NewScooter("S", "Самокат", 6, 1.4, safety)
	bike, _ := domain.NewBicycle("B", "Велосипед", 5, safety)
	ebike, _ := domain.NewEBike("E", "Электровелосипед", 9, 1.1, safety)

	for _, transport := range []domain.Transport{scooter, bike, ebike} {
		if _, err := park.AdmitTransport(transport); err != nil {
			t.Fatal(err)
		}
	}

	helmet, _ := domain.NewHelmet("H", "Шлем")
	cabinet, _ := domain.NewChargingCabinet("C", "Шкаф", 2.5)
	dock, _ := domain.NewDockStation("D", "Док")

	for _, thing := range []domain.InventoryItem{helmet, cabinet, dock} {
		if err := park.RegisterItem(thing); err != nil {
			t.Fatal(err)
		}
	}

	stats := park.FleetSummary()

	if stats.TransportCount != 3 || stats.OtherItemCount != 3 ||
		stats.TotalCount != 6 || math.Abs(stats.TotalDailyEnergyKWh-5.0) > 1e-9 {
		t.Fatalf("неверная статистика парка: %+v", stats)
	}

	beginners := park.BeginnerTransport()

	if len(beginners) != 2 || beginners[0].InventoryNumber() != "S" ||
		beginners[1].InventoryNumber() != "E" {
		t.Fatalf("неверный список транспорта для новичков: %v", beginners)
	}
}

func TestNilInputs(t *testing.T) {
	if _, err := NewFleetService(nil, &fakeInspector{}); err == nil {
		t.Fatal("не проверяется пустой репозиторий")
	}

	if _, err := NewFleetService(storage.NewMemoryRepository(), nil); err == nil {
		t.Fatal("не проверяется отсутствие техосмотра")
	}

	park := newFleet(t, &fakeInspector{})

	if _, err := park.AdmitTransport(nil); err == nil {
		t.Fatal("можно передать пустой транспорт")
	}

	if err := park.RegisterItem(nil); err == nil {
		t.Fatal("можно добавить пустую вещь")
	}
}
