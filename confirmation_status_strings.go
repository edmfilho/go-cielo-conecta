package go_cielo_conecta

import (
	"encoding/json"
	"fmt"
	"strings"
)

func ParseConfirmationStatus(s string) (ConfirmationStatus, error) {
	var c ConfirmationStatus

	switch strings.ToLower(s) {
	case "pendente":
		c = Pendente
	case "confirmado":
		c = Confirmado
	case "desfeito":
		c = Desfeito
	default:
		return 0, fmt.Errorf("invalid ConfirmationStatus: %s", s)
	}

	return c, nil
}

func (c ConfirmationStatus) String() string {
	return [...]string{"Pendente", "Confirmado", "Desfeito"}[c]
}

func (c *ConfirmationStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String())
}

func (c *ConfirmationStatus) UnmarshalJSON(data []byte) error {
	var asInt uint
	if err := json.Unmarshal(data, &asInt); err == nil {
		*c = ConfirmationStatus(asInt)
		return nil
	}

	var asString string
	if err := json.Unmarshal(data, &asString); err != nil {
		return err
	}

	cs, err := ParseConfirmationStatus(asString)
	if err != nil {
		return err
	}

	*c = cs

	return nil
}
