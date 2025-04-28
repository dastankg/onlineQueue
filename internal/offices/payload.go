package offices

type OfficeCreateRequest struct {
	Name         string `json:"name"`
	Address      string `json:"address"`
	WorkingHours string `json:"working_hours"`
}
