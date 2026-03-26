package tests

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"
)
import . "github.com/edmfilho/go-cielo-conecta"

func TestAuthorization(t *testing.T) {
	merchant := Merchant{
		ID:     os.Getenv("MERCHANT_ID"),
		Secret: os.Getenv("MERCHANT_SECRET"),
	}

	client, err := NewClient(merchant, HomologationEnvironment)
	if err != nil {
		t.Fatal(err)
	}

	defer client.Close()

	cc := CreditCard{
		InputMode:      Emv,
		ExpirationDate: "12/2030",
		TrackTwoData:   "A3CE8E13D475CE75BE954B53C804265E6F07202158933F5B",
		EncryptedCardData: EncryptedCardData{
			EncryptionType:       Dukpt3DesCBC,
			TrackTwoDataKSN:      "FFFFF999950000000001",
			InitializationVector: "0000000000000000",
			IsDataInTLVFormat:    false,
		},
		EmvData:              "9F02060000000001009F1A020076950542C00000005F2A0209869A0321100482025C009F360200019F100706011A039000009F260894A26A9922DE03EC9F2701005F3401018407A0000000031010",
		IssuerId:             1010,
		SaveCard:             false,
		PanSequenceNumber:    1,
		IsFallback:           false,
		AuthenticationMethod: OnlineAuthentication,
		BrandId:              1,
		PinBlock: PinBlock{
			EncryptedPinBlock: "ED9822DCD929BD14",
			EncryptionType:    Dukpt3Des,
			KsnIdentification: "FFFFF99999C19FC00072",
		},
	}

	sale := Sale{
		MerchantOrderId: "1234345676",
		Payment: &Payment{
			SubordinatedMerchantId: os.Getenv("MERCHANT_ID"),
			Installments:           1,
			Type:                   "PhysicalCreditCard",
			Interest:               ByMerchant,
			Capture:                true,
			SoftDescriptor:         "TestConecta",
			CreditCard:             &cc,
			PaymentDateTime:        time.Now().Format(time.RFC3339),
			Amount:                 100,
			ProductId:              1,
			PinPadInformation: &PinPadInformation{
				PhysicalCharacteristics: WithoutPinPad,
				ReturnDataInfo:          "00",
				SerialNumber:            "205fae17775b1955",
				TerminalID:              "00000001",
			},
		},
	}

	salePayed, err := client.Authorization(&sale)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println()

	var jsonData []byte

	jsonData, _ = json.MarshalIndent(salePayed, "", "    ")

	fmt.Println(string(jsonData))
}
