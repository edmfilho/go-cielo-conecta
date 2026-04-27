package go_cielo_conecta

import (
	"encoding/json"
	"fmt"
)

func (e *EncryptionType) String() string {
	return [...]string{"DukptDes", "MasterKey", "Dukpt3Des", "Dukpt3DesCBC"}[*e-1]
}

func ParseEncryptionType(s string) (*EncryptionType, error) {
	var e EncryptionType

	switch s {
	case "DukptDes":
		e = DukptDes
	case "MasterKey":
		e = MasterKey
	case "Dukpt3Des":
		e = Dukpt3Des
	case "Dukpt3DesCBC":
		e = Dukpt3DesCBC
	default:
		return nil, fmt.Errorf("invalid EncryptionType: %s", s)
	}

	return &e, nil
}

func (e *EncryptionType) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.String())
}

func (e *EncryptionType) UnmarshalJSON(data []byte) error {
	var asInt int
	if err := json.Unmarshal(data, &asInt); err == nil {
		*e = EncryptionType(asInt)
		return nil
	}

	var asString string
	if err := json.Unmarshal(data, &asString); err != nil {
		return err
	}

	et, err := ParseEncryptionType(asString)
	if err != nil {
		return err
	}

	*e = *et

	return nil
}
