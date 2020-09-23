// Package queryparams  handles the logic to decide if there are any errorTypes in the requests based in the query parameters.
// Is the controller of the query parameters.
package queryparams

// ExtractPrinterInfo extracts the printer information from the query parameters argument and
// returns the appropiate data and an error, if any.
func ExtractPrinterInfo(queryParameters map[string]string) (productNumber string, serialNumber string, err error) {
	productNumber = queryParameters["pn"]
	serialNumber = queryParameters["sn"]
	if productNumber == "" && serialNumber != "" {
		err = ErrorQueryStringPnSn
		return
	}
	return productNumber, serialNumber, nil
}
