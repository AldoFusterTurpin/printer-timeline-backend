package db_test

import (
	"context"
	"strings"
	"time"

	. "bitbucket.org/aldoft/printer-timeline-backend/app/internal/db"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	printerId1    = "K0Q45A!SGCCTEAM125"
	printerId2    = "CZ056A!HP333777AX"
	printerIDFake = "prin!Item"
	accountId1    = "pr12345678"
	accountId2    = "se00000000"
	accountId3    = "fibo24"
	accFake       = "anAccID"
	fakeServId    = "fakeServId"
)

var _ = Describe("CCPrinterSubscription integration tests", func() {
	var (
		subscriptionsArr []*CCPrinterSubscriptionModel
		ctx              context.Context
		err              error
	)

	BeforeEach(func() {
		ctx = context.Background()

		subscriptionsArr = []*CCPrinterSubscriptionModel{
			createSubscription(printerId1, accountId1, ServicePrintOS),
			createSubscription(printerId2, accountId2, ServiceSeals),
			createSubscription(printerId2, accountId3, ServiceFibo24),
		}

		for _, subscription := range subscriptionsArr {
			err = ccPrinterSubscriptionCollection.Put(ctx, subscription)
			Expect(err).To(BeNil())
		}

	})

	AfterEach(func() {
		for _, subscription := range subscriptionsArr {
			err = ccPrinterSubscriptionCollection.Delete(ctx, subscription.PrinterID, subscription.AccountID)
			Expect(err).To(BeNil())
		}
	})

	Describe("test cases", func() {
		It("Should get the expected printer subscription", func() {
			subscription, err := ccPrinterSubscriptionCollection.Get(ctx, printerId1, accountId1)

			Expect(err).To(BeNil())
			Expect(subscription.PrinterID).To(BeEquivalentTo(printerId1))
			Expect(subscription.AccountID).To(BeEquivalentTo(accountId1))
		})

		It("Should return error while attempted to retrieve non existing subscriptions", func() {
			_, err := ccPrinterSubscriptionCollection.Get(ctx, "fakePrinter", "fakeAccountId")

			Expect(err).To(BeEquivalentTo(NotFoundErr))
		})

		It("Should get all printer subscriptions", func() {
			subscriptions, err := ccPrinterSubscriptionCollection.GetPrinterSubscriptions(ctx, printerId2)

			Expect(err).To(BeNil())
			Expect(len(subscriptions)).To(BeEquivalentTo(2))
		})

		It("Should return error while attempted to retrieve non existing printer subscriptions", func() {
			_, err := ccPrinterSubscriptionCollection.GetPrinterSubscriptions(ctx, "fakePrinterId")

			Expect(err).To(BeEquivalentTo(NotFoundErr))
		})

		It("Should return error if adding twice", func() {

			tempSubscription := createSubscription(printerIDFake, accFake, fakeServId)

			err1 := ccPrinterSubscriptionCollection.Put(ctx, tempSubscription)
			Expect(err1).To(BeNil())

			err2 := ccPrinterSubscriptionCollection.Put(ctx, tempSubscription)
			Expect(err2).NotTo(BeNil())
			Expect(err2.Error()).To(BeEquivalentTo(ConditionalPutErr.Error()))

			subscription, err := ccPrinterSubscriptionCollection.Get(ctx, printerIDFake, accFake)

			Expect(err).To(BeNil())
			Expect(subscription.PrinterID).To(BeEquivalentTo(printerIDFake))
			Expect(subscription.AccountID).To(BeEquivalentTo(accFake))

			errs := ccPrinterSubscriptionCollection.Delete(ctx, tempSubscription.PrinterID, tempSubscription.AccountID)
			Expect(errs).To(BeNil())

		})

		It("Should return error retrieving unexistent item", func() {
			err := ccPrinterSubscriptionCollection.Delete(ctx, printerIDFake, accFake)
			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(BeEquivalentTo(ConditionalDelErr.Error()))
		})

	})

})

func createSubscription(printerId string, accountId string, serviceId string) *CCPrinterSubscriptionModel {
	productNumberSerialNumber := strings.Split(printerId, "!")
	return &CCPrinterSubscriptionModel{
		PrinterID:             printerId,
		AccountID:             accountId,
		SerialNumber:          productNumberSerialNumber[1],
		ProductNumber:         productNumberSerialNumber[0],
		ServiceID:             Service(serviceId),
		RegistrationTimeEpoch: time.Now().Unix(),
	}
}
