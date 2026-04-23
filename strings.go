package go_cielo_conecta

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Custom Marshal & Unmarshal for EncryptionType and StatusPayment to handle both integer and string representations in JSON.
func (e *EncryptionType) String() string {
	return [...]string{"DukptDes", "MasterKey", "Dukpt3Des", "Dukpt3DesCBC"}[*e-1]
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

	switch asString {
	case "DukptDes":
		*e = DukptDes
	case "MasterKey":
		*e = MasterKey
	case "Dukpt3Des":
		*e = Dukpt3Des
	case "Dukpt3DesCBC":
		*e = Dukpt3DesCBC
	default:
		return fmt.Errorf("invalid EncryptionTypeEnum=%s", asString) // ou retornar um erro indicando valor desconhecido
	}

	return nil
}

func (s *StatusPayment) String() string {
	return [...]string{"Pending", "Confirmed", "Cancelled", "Reversed", "Processing", "Denied", "Unreachable", "WaitingValidation", "WaitingCapture", "RefundedDevolution", "Refunded", "Approved"}[*s-1]
}

func (s *StatusPayment) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (s *StatusPayment) UnmarshalJSON(data []byte) error {
	var asInt int
	if err := json.Unmarshal(data, &asInt); err == nil {
		*s = StatusPayment(asInt)
		return nil
	}

	var asString string
	if err := json.Unmarshal(data, &asString); err != nil {
		return err
	}

	switch strings.ToLower(asString) {
	case "pending":
		*s = Pending
	case "confirmed":
		*s = Confirmed
	case "cancelled":
		*s = Cancelled
	case "reversed":
		*s = Reversed
	case "processing":
		*s = Processing
	case "denied":
		*s = Denied
	case "unreachable":
		*s = Unreachable
	case "waitingvalidation":
		*s = WaitingValidation
	case "waitingcapture":
		*s = WaitingCapture
	case "refundeddevolution":
		*s = RefundedDevolution
	case "refunded":
		*s = Refunded
	case "approved":
		*s = Approved
	default:
		return fmt.Errorf("invalid StatusPaymentEnum=%s", asString) // ou retornar um erro indicando valor desconhecido
	}

	return nil
}

func (c *ConfirmationStatus) String() string {
	return [...]string{"Pendente", "Confirmado", "Desfeito"}[*c]
}

func (c *ConfirmationStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String())
}

func (c *ConfirmationStatus) UnmarshalJSON(data []byte) error {
	var asInt int
	if err := json.Unmarshal(data, &asInt); err == nil {
		*c = ConfirmationStatus(asInt)
		return nil
	}

	var asString string
	if err := json.Unmarshal(data, &asString); err != nil {
		return err
	}

	switch strings.ToLower(asString) {
	case "pendente":
		*c = Pendente
	case "confirmado":
		*c = Confirmado
	case "desfeito":
		*c = Desfeito
	default:
		return fmt.Errorf("invalid ConfirmationStatus=%s", asString)
	}

	return nil
}
