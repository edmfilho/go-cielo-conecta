package go_cielo_conecta

const (
	ByMerchant Interest = "ByMerchant"
	ByIssuer   Interest = "ByIssuer"

	CPF  IdentityType = "CPF"
	CNPJ IdentityType = "CNPJ"

	Typed          InputMode = "Typed"
	MagStripe      InputMode = "MagStripe"
	Emv            InputMode = "Emv"
	ContactlessEmv InputMode = "ContactlessEmv"

	NoPassword            AuthenticationMethod = "NoPassword"
	OnlineAuthentication  AuthenticationMethod = "OnlineAuthentication"
	OfflineAuthentication AuthenticationMethod = "OfflineAuthentication"

	WithoutPinPad                                   PhysicalCharacteristics = "WithoutPinPad"
	PinPadWithoutChipReader                         PhysicalCharacteristics = "PinPadWithoutChipReader"
	PinPadWithChipReaderWithoutSamModule            PhysicalCharacteristics = "PinPadWithChipReaderWithoutSamModule"
	PinPadWithChipReaderWithSamModule               PhysicalCharacteristics = "PinPadWithChipReaderWithSamModule"
	NotCertifiedPinPad                              PhysicalCharacteristics = "NotCertifiedPinPad"
	PinPadWithChipReaderWithoutSamAndContactless    PhysicalCharacteristics = "PinPadWithChipReaderWithoutSamAndContactless"
	PinPadWithChipReaderWithSamModuleAndContactless PhysicalCharacteristics = "PinPadWithChipReaderWithSamAndContactless"

	Collected   SecurityCodeStatus = "Collected"
	Unreadable  SecurityCodeStatus = "Unreadable"
	Nonexistent SecurityCodeStatus = "Nonexistent"

	CurrencyBRL = currency("BRL")
	CurrencyUSD = currency("USD")
)

const (
	DukptDes EncryptionType = iota + 1
	MasterKey
	Dukpt3Des
	Dukpt3DesCBC
)

const (
	Pending StatusPayment = iota + 1
	Confirmed
	Cancelled
	Reversed
	Processing
	Denied
	Unreachable
	WaitingValidation
	WaitingCapture
	RefundedDevolution
	Refunded
	Approved
)

var (
	SandboxEnvironment = Environment{
		OAuthURL:    "https://authsandbox.cieloecommerce.cielo.com.br/oauth2/token",
		ParamsURL:   "https://parametersdownloadsandbox.cieloecommerce.cielo.com.br/api/v0.1/initialization/{SubordinatedMerchantId}/{TerminalId}",
		APIUrl:      "https://apisandbox.cieloecommerce.cielo.com.br",
		APIQueryUrl: "https://apiquerysandbox.cieloecommerce.cielo.com.br",
	}

	HomologationEnvironment = Environment{
		OAuthURL:     "https://authsandbox.cieloecommerce.cielo.com.br/oauth2/token",
		ParamsURL:    "https://parametersdownloadsandbox.cieloecommerce.cielo.com.br/api/v0.1/initialization/{SubordinatedMerchantId}/{TerminalId}",
		APIUrl:       "https://apisandbox.cieloecommerce.cielo.com.br",
		APIQueryUrl:  "https://apiquerysandbox.cieloecommerce.cielo.com.br",
		Homologation: true,
	}

	ProductionEnvironment = Environment{
		OAuthURL:    "https://auth.cieloecommerce.cielo.com.br/oauth2/token",
		ParamsURL:   "https://parametersdownload.cieloecommerce.cielo.com.br/api/v0.1/initialization/{SubordinatedMerchantId}/{TerminalId}",
		APIUrl:      "https://api.cieloecommerce.cielo.com.br",
		APIQueryUrl: "https://apiquery.cieloecommerce.cielo.com.br/",
	}
)
