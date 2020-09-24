// Package queryparams  handles the logic to decide if there are any errorTypes in the requests based in the query parameters.
// Is the controller of the query parameters.
package queryparams

import "bitbucket.org/aldoft/printer-timeline-backend/app/internal/configs"

// ExtractPrinterInfo extracts the printer information from the query parameters argument and
// returns the appropiate data and an error, if any.
func ExtractPrinterInfo(queryParameters map[string]string) (productNumber string, serialNumber string, err error) {
	productNumber = queryParameters[configs.ProductNumberQueryParam]
	serialNumber = queryParameters[configs.SerialNumberQueryParam]
	if productNumber == "" && serialNumber != "" {
		err = ErrorQueryStringPnSn
		return
	}
	return productNumber, serialNumber, nil
}
