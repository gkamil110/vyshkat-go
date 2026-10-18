package demo

import (
	"bytes"
	"strings"
	"testing"
)

func TestDemoShowsAllCoreScenarios(t *testing.T) {
	var buffer bytes.Buffer

	err := Run(&buffer)
	if err != nil {
		t.Fatal(err)
	}

	result := buffer.String()

	expectedText := []string{
		"Принято:",
		"Отклонено:",
		"Добавлено в парк:",
		"Не получилось добавить повторно:",
		"транспорт 2, вещи 3, всего 5, энергия 4.00",
		"Транспорт для новичков",
		"Вся техника и вещи в парке:",
	}

	for _, expected := range expectedText {
		if !strings.Contains(result, expected) {
			t.Errorf("в демо нет строки %q", expected)
		}
	}
}
