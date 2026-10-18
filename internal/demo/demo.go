package demo

import (
	"fmt"
	"io"

	"vyshkat/internal/bootstrap"
	"vyshkat/internal/domain"
)

func Run(out io.Writer) error {
	park, err := bootstrap.Fleet(bootstrap.NewContainer())
	if err != nil {
		return err
	}

	fmt.Fprintln(out, "=== Демонстрация работы программы ===")

	working := domain.SafetyStatus{
		BrakesOK:     true,
		SteeringOK:   true,
		ElectricalOK: true,
	}

	brokenBrakes := domain.SafetyStatus{
		BrakesOK:     false,
		SteeringOK:   true,
		ElectricalOK: true,
	}

	scooter, err := domain.NewScooter("V-001", "Самокат у корпуса", 8, 1.4, working)
	if err != nil {
		return err
	}

	bike, err := domain.NewBicycle("V-002", "Городской велосипед", 5, working)
	if err != nil {
		return err
	}

	brokenBike, err := domain.NewEBike("V-003", "Электровелосипед с поломкой", 9, 1.2, brokenBrakes)
	if err != nil {
		return err
	}

	for _, transport := range []domain.Transport{scooter, bike, brokenBike} {
		check, err := park.AdmitTransport(transport)
		if err != nil {
			return err
		}

		if check.Approved {
			fmt.Fprintf(out, "Принято: %s [%s]\n", transport.DisplayName(), transport.InventoryNumber())
		} else {
			fmt.Fprintf(out, "Отклонено: %s [%s] — %s\n", transport.DisplayName(), transport.InventoryNumber(), check.Reason)
		}
	}

	helmet, err := domain.NewHelmet("A-001", "Защитный шлем")
	if err != nil {
		return err
	}

	dock, err := domain.NewDockStation("A-002", "Док-станция")
	if err != nil {
		return err
	}

	cabinet, err := domain.NewChargingCabinet("A-003", "Зарядный шкаф", 2.6)
	if err != nil {
		return err
	}

	for _, thing := range []domain.InventoryItem{helmet, dock, cabinet} {
		if err := park.RegisterItem(thing); err != nil {
			return err
		}

		fmt.Fprintf(out, "Добавлено в парк: %s [%s]\n", thing.DisplayName(), thing.InventoryNumber())
	}

	duplicateHelmet, err := domain.NewHelmet("V-001", "Дубликат номера")
	if err != nil {
		return err
	}

	if err := park.RegisterItem(duplicateHelmet); err != nil {
		fmt.Fprintf(out, "Не получилось добавить повторно: %v\n", err)
	} else {
		return fmt.Errorf("демонстрация: ожидалось отклонение дубликата")
	}

	stats := park.FleetSummary()

	fmt.Fprintf(out, "Итого: транспорт %d, вещи %d, всего %d, энергия %.2f кВт·ч/сутки\n",
		stats.TransportCount, stats.OtherItemCount, stats.TotalCount, stats.TotalDailyEnergyKWh)

	fmt.Fprintln(out, "Транспорт для новичков (оценка от 6):")

	for _, transport := range park.BeginnerTransport() {
		fmt.Fprintf(out, "  %s [%s], простота %d\n",
			transport.DisplayName(), transport.InventoryNumber(), transport.EaseOfUse())
	}

	fmt.Fprintln(out, "Вся техника и вещи в парке:")

	for _, thing := range park.Inventory() {
		fmt.Fprintf(out, "  %s | %s | %s\n",
			thing.Category(), thing.DisplayName(), thing.InventoryNumber())
	}

	fmt.Fprintln(out, "=== Демо закончено. Теперь можно добавлять свою технику ===")

	return nil
}
