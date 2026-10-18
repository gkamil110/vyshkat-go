package inspection

import (
	"strings"
	"testing"

	"vyshkat/internal/domain"
)

func TestServiceCenter(t *testing.T) {
	tests := []struct {
		name       string
		safety     domain.SafetyStatus
		isElectric bool
		approved   bool
		reason     string
	}{
		{"исправный самокат", domain.SafetyStatus{BrakesOK: true, SteeringOK: true, ElectricalOK: true}, true, true, "техосмотр пройден"},
		{"сломаны тормоза", domain.SafetyStatus{BrakesOK: false, SteeringOK: true, ElectricalOK: true}, true, false, "тормоз"},
		{"неисправный руль", domain.SafetyStatus{BrakesOK: true, SteeringOK: false, ElectricalOK: true}, true, false, "рулев"},
		{"проблема с электрикой", domain.SafetyStatus{BrakesOK: true, SteeringOK: true, ElectricalOK: false}, true, false, "электр"},
		{"обычный велосипед", domain.SafetyStatus{BrakesOK: true, SteeringOK: true, ElectricalOK: false}, false, true, "техосмотр пройден"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var transport domain.Transport
			var err error

			if test.isElectric {
				transport, err = domain.NewScooter("TEST-1", "Самокат", 6, 1, test.safety)
			} else {
				transport, err = domain.NewBicycle("TEST-1", "Велосипед", 6, test.safety)
			}

			if err != nil {
				t.Fatal(err)
			}

			check := (ServiceCenter{}).Inspect(transport)

			if check.Approved != test.approved {
				t.Fatalf("неверный результат техосмотра: %+v", check)
			}

			if !strings.Contains(check.Reason, test.reason) {
				t.Fatalf("неверная причина отказа: %s", check.Reason)
			}
		})
	}
}
