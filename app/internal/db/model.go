package db

const PrinterIdSeparator = "!"

type CCPrinterSubscriptionModel struct {
	PrinterID             string  `json:"PrinterID"`           // Primary Key
	AccountID             string  `json:"AccountID,omitempty"` // Sort Key
	SerialNumber          string  `json:"SN,omitempty"`
	ProductNumber         string  `json:"PN,omitempty"`
	ServiceID             Service `json:"ServiceID,omitempty"`
	RegistrationTimeEpoch int64   `json:"RegistrationTimeEpoch,omitempty"`
}

type Service string

const (
	ServiceUnknown Service = "Unknown"
	ServicePrintOS         = "PRINTOS"
	ServiceLatexGO         = "LATEX2GO"
	ServiceSeals           = "seals"
	ServiceFibo24          = "fibo24"
	ServicePPU             = "HP-PPU"
)
