package console

import (
	"bytes"
	"strings"
	"testing"

	"vyshkat/internal/bootstrap"
)

func runScript(t *testing.T, input string) string {
	t.Helper()

	park, err := bootstrap.Fleet(bootstrap.NewContainer())
	if err != nil {
		t.Fatal(err)
	}

	var buffer bytes.Buffer
	app := New(strings.NewReader(input), &buffer, park)

	if err := app.Run(); err != nil {
		t.Fatal(err)
	}

	return buffer.String()
}

func TestInteractiveHappyPathAndReports(t *testing.T) {
	commands := []string{
		"1", "1", "s-001", "Самокат", "6", "1,5", "да", "да", "да",
		"2", "3", "c-001", "Шкаф", "2.0",
		"2", "1", "h-001", "Шлем",
		"3", "4", "5", "0",
	}

	result := runScript(t, strings.Join(commands, "\n")+"\n")

	expectedText := []string{
		"ПРИНЯТО: Самокат [S-001]",
		"НА БАЛАНСЕ: Шкаф [C-001]",
		"транспорт — 1, вещи — 2, всего — 3",
		"3.50 кВт·ч",
		"ТРАНСПОРТ ДЛЯ НОВИЧКОВ",
		"Самокат | № S-001",
		"ИНВЕНТАРНАЯ ВЕДОМОСТЬ:",
		"Шлем | Шлем | № H-001",
		"До встречи!",
	}

	for _, expected := range expectedText {
		if !strings.Contains(result, expected) {
			t.Errorf("в выводе нет строки %q:\n%s", expected, result)
		}
	}
}

func TestInteractiveRejectedAndErrorsAndEOF(t *testing.T) {
	commands := []string{
		"not-an-option",
		"1", "2", "b-001", "Велосипед", "5", "нет", "да",
		"3", "4", "5",
		"2", "1", "h-001", "Шлем",
		"2", "1", "H-001", "Дубликат",
		"2", "3", "c-001", "Шкаф", "bad-number",
		"0",
	}

	result := runScript(t, strings.Join(commands, "\n")+"\n")

	expectedText := []string{
		"Неизвестный пункт меню",
		"ОТКЛОНЕНО:",
		"транспорт — 0, вещи — 0, всего — 0",
		"Подходящего транспорта нет",
		"На балансе пока ничего нет",
		"инвентарный номер уже существует",
		"неверное число",
	}

	for _, expected := range expectedText {
		if !strings.Contains(result, expected) {
			t.Errorf("в выводе нет строки %q:\n%s", expected, result)
		}
	}

	unfinishedInput := runScript(t, "1\n1\n")
	if !strings.Contains(unfinishedInput, "Ввод завершён") {
		t.Fatalf("программа не завершилась после конца ввода: %s", unfinishedInput)
	}
}

func TestInvalidRatingsAndBoolean(t *testing.T) {
	commands := []string{
		"1", "2", "b-1", "Велосипед", "11", "да", "да",
		"1", "2", "b-2", "Велосипед", "4", "не знаю",
		"0",
	}

	result := runScript(t, strings.Join(commands, "\n")+"\n")

	if !strings.Contains(result, "простота использования должна быть от 1 до 10") {
		t.Errorf("не проверяется оценка простоты использования:\n%s", result)
	}

	if !strings.Contains(result, "введите да или нет") {
		t.Errorf("не проверяется ответ да/нет:\n%s", result)
	}
}
