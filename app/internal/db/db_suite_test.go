package db_test

import (
	"bitbucket.org/aldoft/printer-timeline-backend/app/internal/db"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestDb(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Db Suite")
}

var ccPrinterSubscriptionCollection *db.CCPrinterSubscriptionCollection

var _ = BeforeSuite(func() {

	var err error

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Config: aws.Config{
			Region: aws.String(endpoints.UsEast1RegionID),
		},
	}))

	ccPrinterSubscriptionCollection, err = db.NewCCPrinterSubscriptionCollectionWithSession(sess)
	Expect(err).To(BeNil())
})
