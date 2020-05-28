package queryParamsCtrl

import "bitbucket.org/aldoft/printer-timeline-backend/errors"

func ExtractPrinterInfo(queryParameters map[string]string) (productNumber string, serialNumber string, err error) {
	productNumber = queryParameters["pn"]
	serialNumber = queryParameters["sn"]
	if productNumber == "" && serialNumber != "" {
		err = errors.QueryStringPnSn
		return
	}
	return productNumber, serialNumber, nil
}