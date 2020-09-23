package api

import (
	"bitbucket.org/aldoft/printer-timeline-backend/app/internal/configs"
	"bitbucket.org/aldoft/printer-timeline-backend/app/internal/db"
	"context"
)

// GetPrinterSubscriptions is the responsible of retrieving subscriptions based in the queryParameters.
func GetPrinterSubscriptions(queryParameters map[string]string, printerSubscriptionFetcher db.PrinterSubscriptionFetcher) (status int, result []*db.CCPrinterSubscriptionModel, err error) {

	printerId := queryParameters[configs.ProductNumberQueryParam] + db.PrinterIdSeparator + queryParameters[configs.SerialNumberQueryParam]
	subs, err := printerSubscriptionFetcher.GetPrinterSubscriptions(context.Background(), printerId)
	if err != nil {
		status = SelectHTTPStatus(err)
		return
	}

	status = SelectHTTPStatus(err)
	return status, subs, err
}
