package console

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"vyshkat/internal/application"
	"vyshkat/internal/domain"
)

type Console struct {
	service *application.FleetService
	input   *bufio.Reader
	output  io.Writer
}

func New(in io.Reader, out io.Writer, service *application.FleetService) *Console {
	return &Console{service: service, input: bufio.NewReader(in), output: out}
}

func (ui *Console) Run() error {
	fmt.Fprintln(ui.output, "\nДобро пожаловать в «ВышКат»!")

	for {
		fmt.Fprintln(ui.output, "\n1. Добавить транспорт (с техосмотром)")
		fmt.Fprintln(ui.output, "2. Добавить вещь")
		fmt.Fprintln(ui.output, "3. Статистика парка")
		fmt.Fprintln(ui.output, "4. Транспорт для новичков")
		fmt.Fprintln(ui.output, "5. Список имущества")
		fmt.Fprintln(ui.output, "0. Выход")

		choice, err := ui.ask("Ваш выбор: ")
		if errors.Is(err, io.EOF) {
			fmt.Fprintln(ui.output, "\nВвод завершён. До встречи!")
			return nil
		}
		if err != nil {
			return err
		}

		switch choice {
		case "0":
			fmt.Fprintln(ui.output, "До встречи!")
			return nil
		case "1":
			err = ui.addTransport()
		case "2":
			err = ui.addEquipment()
		case "3":
			ui.showSummary()
		case "4":
			ui.showBeginners()
		case "5":
			ui.showInventory()
		default:
			fmt.Fprintln(ui.output, "Неизвестный пункт меню. Попробуйте ещё раз.")
		}

		if errors.Is(err, io.EOF) {
			fmt.Fprintln(ui.output, "\nВвод завершён. До встречи!")
			return nil
		}
		if err != nil {
			return err
		}
	}
}

func (ui *Console) ask(question string) (string, error) {
	fmt.Fprint(ui.output, question)
	text, err := ui.input.ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return "", err
	}

	text = strings.TrimSpace(text)
	if errors.Is(err, io.EOF) && text == "" {
		return "", io.EOF
	}
	return text, nil
}

func (ui *Console) askInt(question string) (int, error) {
	text, err := ui.ask(question)
	if err != nil {
		return 0, err
	}

	number, err := strconv.Atoi(text)
	if err != nil {
		return 0, fmt.Errorf("нужно ввести целое число: %q", text)
	}
	return number, nil
}

func (ui *Console) askEnergy() (float64, error) {
	text, err := ui.ask("Расход энергии за сутки (кВт·ч, например 1.5): ")
	if err != nil {
		return 0, err
	}

	energy, err := strconv.ParseFloat(strings.ReplaceAll(text, ",", "."), 64)
	if err != nil {
		return 0, fmt.Errorf("неверное число: %q", text)
	}
	return energy, nil
}

func (ui *Console) askBool(question string) (bool, error) {
	answer, err := ui.ask(question + " (да/нет): ")
	if err != nil {
		return false, err
	}

	switch strings.ToLower(answer) {
	case "да", "д", "yes", "y", "+", "1":
		return true, nil
	case "нет", "н", "no", "n", "-", "0":
		return false, nil
	default:
		return false, fmt.Errorf("введите да или нет, получено %q", answer)
	}
}

func (ui *Console) showInputError(err error) error {
	if errors.Is(err, io.EOF) {
		return err
	}
	fmt.Fprintf(ui.output, "Ошибка: %v\n", err)
	return nil
}

func (ui *Console) addTransport() error {
	fmt.Fprintln(ui.output, "Какой транспорт? 1 — самокат, 2 — велосипед, 3 — электровелосипед")
	kind, err := ui.ask("Номер варианта: ")
	if err != nil {
		return ui.showInputError(err)
	}
	if kind != "1" && kind != "2" && kind != "3" {
		fmt.Fprintln(ui.output, "Такого транспорта нет в списке")
		return nil
	}

	number, err := ui.ask("Инвентарный номер: ")
	if err != nil {
		return ui.showInputError(err)
	}
	name, err := ui.ask("Название: ")
	if err != nil {
		return ui.showInputError(err)
	}
	ease, err := ui.askInt("Простота использования (от 1 до 10): ")
	if err != nil {
		return ui.showInputError(err)
	}

	var dailyEnergy float64
	if kind != "2" {
		dailyEnergy, err = ui.askEnergy()
		if err != nil {
			return ui.showInputError(err)
		}
	}

	fmt.Fprintln(ui.output, "Проверим состояние транспорта:")
	brakesOK, err := ui.askBool("Тормоза работают?")
	if err != nil {
		return ui.showInputError(err)
	}
	steeringOK, err := ui.askBool("Руль исправен?")
	if err != nil {
		return ui.showInputError(err)
	}

	safety := domain.SafetyStatus{BrakesOK: brakesOK, SteeringOK: steeringOK}
	if kind != "2" {
		safety.ElectricalOK, err = ui.askBool("Электрика исправна?")
		if err != nil {
			return ui.showInputError(err)
		}
	}

	var transport domain.Transport
	switch kind {
	case "1":
		transport, err = domain.NewScooter(number, name, ease, dailyEnergy, safety)
	case "2":
		transport, err = domain.NewBicycle(number, name, ease, safety)
	case "3":
		transport, err = domain.NewEBike(number, name, ease, dailyEnergy, safety)
	}
	if err != nil {
		return ui.showInputError(err)
	}

	check, err := ui.service.AdmitTransport(transport)
	if err != nil {
		return ui.showInputError(err)
	}
	if check.Approved {
		fmt.Fprintf(ui.output, "ПРИНЯТО: %s [%s]. %s\n", transport.DisplayName(), transport.InventoryNumber(), check.Reason)
	} else {
		fmt.Fprintf(ui.output, "ОТКЛОНЕНО: %s [%s]. Причина: %s\n", transport.DisplayName(), transport.InventoryNumber(), check.Reason)
	}
	return nil
}

func (ui *Console) addEquipment() error {
	fmt.Fprintln(ui.output, "Что добавляем? 1 — шлем, 2 — док-станция, 3 — зарядный шкаф")
	kind, err := ui.ask("Номер варианта: ")
	if err != nil {
		return ui.showInputError(err)
	}
	if kind != "1" && kind != "2" && kind != "3" {
		fmt.Fprintln(ui.output, "Такого варианта нет в списке")
		return nil
	}

	number, err := ui.ask("Инвентарный номер: ")
	if err != nil {
		return ui.showInputError(err)
	}
	name, err := ui.ask("Название: ")
	if err != nil {
		return ui.showInputError(err)
	}

	var item domain.InventoryItem
	switch kind {
	case "1":
		item, err = domain.NewHelmet(number, name)
	case "2":
		item, err = domain.NewDockStation(number, name)
	case "3":
		var energy float64
		energy, err = ui.askEnergy()
		if err != nil {
			return ui.showInputError(err)
		}
		item, err = domain.NewChargingCabinet(number, name, energy)
	}
	if err != nil {
		return ui.showInputError(err)
	}
	if err := ui.service.RegisterItem(item); err != nil {
		return ui.showInputError(err)
	}

	fmt.Fprintf(ui.output, "НА БАЛАНСЕ: %s [%s]\n", item.DisplayName(), item.InventoryNumber())
	return nil
}

func (ui *Console) showSummary() {
	stats := ui.service.FleetSummary()
	fmt.Fprintf(ui.output, "\nСВОДКА: транспорт — %d, вещи — %d, всего — %d\n", stats.TransportCount, stats.OtherItemCount, stats.TotalCount)
	fmt.Fprintf(ui.output, "Суммарное суточное потребление: %.2f кВт·ч\n", stats.TotalDailyEnergyKWh)
}

func (ui *Console) showBeginners() {
	transportList := ui.service.BeginnerTransport()
	fmt.Fprintf(ui.output, "\nТРАНСПОРТ ДЛЯ НОВИЧКОВ (простота >= %d):\n", application.MinimumBeginnerEase)
	if len(transportList) == 0 {
		fmt.Fprintln(ui.output, "Подходящего транспорта нет")
		return
	}

	for _, transport := range transportList {
		fmt.Fprintf(ui.output, "- %s | %s | № %s | простота: %d/10\n", transport.Category(), transport.DisplayName(), transport.InventoryNumber(), transport.EaseOfUse())
	}
}

func (ui *Console) showInventory() {
	items := ui.service.Inventory()
	fmt.Fprintln(ui.output, "\nИНВЕНТАРНАЯ ВЕДОМОСТЬ:")
	if len(items) == 0 {
		fmt.Fprintln(ui.output, "На балансе пока ничего нет")
		return
	}

	for _, item := range items {
		fmt.Fprintf(ui.output, "- %s | %s | № %s\n", item.Category(), item.DisplayName(), item.InventoryNumber())
	}
}
