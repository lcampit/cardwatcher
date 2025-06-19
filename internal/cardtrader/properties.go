package cardtrader

type EditableProperties struct {
	Name           string `json:"name"`
	Type           string `json:"type"`
	DefaultValue   string `json:"default_value"`
	PossibleValues []any  `json:"possible_values"`
}
