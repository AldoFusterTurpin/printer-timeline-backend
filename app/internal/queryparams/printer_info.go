// Package queryparams  extracts the query parameters from the http requests.
// It handles the logic to decide if there are any errors in the requests based in the query parameters.
package queryparams

import "bitbucket.org/aldoft/printer-timeline-backend/app/internal/errors"

// ExtractPrinterInfo extracts the printer information from the query parameters argument and
// returns the appropiate data and an error, if any.
func ExtractPrinterInfo(queryParameters map[string]string) (productNumber string, serialNumber string, err error) {
	productNumber = queryParameters["pn"]
	serialNumber = queryParameters["sn"]
	if productNumber == "" && serialNumber != "" {
		err = errors.QueryStringPnSn
		return
	}
	return productNumber, serialNumber, nil
}
