package inspection

import (
	"vyshkat/internal/application"
	"vyshkat/internal/domain"
)

type ServiceCenter struct{}

func (ServiceCenter) Inspect(transport domain.Transport) application.InspectionResult {
	safety := transport.InspectionState()

	if !safety.BrakesOK {
		return application.InspectionResult{
			Reason: "неисправна тормозная система",
		}
	}

	if !safety.SteeringOK {
		return application.InspectionResult{
			Reason: "неисправно рулевое управление",
		}
	}

	if _, usesEnergy := transport.(domain.EnergyConsumer); usesEnergy && !safety.ElectricalOK {
		return application.InspectionResult{
			Reason: "неисправна электрическая часть",
		}
	}

	return application.InspectionResult{
		Approved: true,
		Reason:   "техосмотр пройден",
	}
}
