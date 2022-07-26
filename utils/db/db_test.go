package db

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"testing"
)

func TestDb(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Location service test suite")
}

func clearDB(db *DatabaseClient) {
	_, err := db.client.DeleteTable(context.TODO(), &dynamodb.DeleteTableInput{TableName: aws.String(tableName)})
	if err != nil {
		panic(err)
	}
}

var _ = Describe("GetByPk", func() {
	var db *DatabaseClient

	BeforeEach(func() {
		db = New()
	})

	AfterEach(func() {
		clearDB(db)
	})

	Context("When table item does not exist", func() {
		It("Returns an empty response", func() {
			out, err := db.GetByPk("test")
			Expect(err).To(BeNil())
			Expect(out).To(BeNil())
		})
	})

	//Context("When table item exists", func() {
	//	BeforeAll(func() {
	//
	//	})
	//
	//	It("Returns a non-empty response", func() {
	//		out, err := db.GetByPk("test")
	//		Expect(err).To(BeNil())
	//		Expect(out).ToNot(BeNil())
	//	})
	//})
})

var _ = Describe("UpdateItem", func() {
	var db *DatabaseClient

	BeforeEach(func() {
		db = New()
	})

	AfterEach(func() {
		clearDB(db)
	})

	Context("When table item does not exist", func() {
		It("Saves as a new item", func() {
			err := db.UpdateItem("test", map[string]string{"first_name": "Evan"})
			Expect(err).To(BeNil())
		})
	})
})
