package db

import (
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	testHelper "some-api/internal/testing"
	"testing"
)

func TestDb(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "DB service test suite")
}

type Test struct {
	Pk        string
	FirstName string
}

var _ = Describe("GetByPk", func() {
	var db *DatabaseClient

	BeforeEach(func() {
		db = New()
	})

	AfterEach(func() {
		testHelper.ClearDB(db)
	})

	Context("When table item does not exist", func() {
		It("Returns an empty response", func() {
			out, err := db.GetByPk("test")
			Expect(err).To(BeNil())
			Expect(out).To(BeNil())
		})
	})

	Context("When table item exists", func() {
		BeforeEach(func() {
			_ = db.UpdateItem("test", map[string]string{"firstName": "Evan"})
		})

		It("Returns a non-empty response", func() {
			out, err := db.GetByPk("test")
			Expect(err).To(BeNil())

			var t Test
			err = attributevalue.UnmarshalMap(out, &t)
			Expect(t.FirstName).To(Equal("Evan"))
		})
	})
})

var _ = Describe("UpdateItem", func() {
	var db *DatabaseClient

	BeforeEach(func() {
		db = New()
	})

	AfterEach(func() {
		testHelper.ClearDB(db)
	})

	Context("When table item does not exist", func() {
		It("Saves as a new item", func() {
			err := db.UpdateItem("test", map[string]string{"firstName": "Evan"})
			Expect(err).To(BeNil())

			out, _ := db.GetByPk("test")

			var t Test
			err = attributevalue.UnmarshalMap(out, &t)
			Expect(t.FirstName).To(Equal("Evan"))
		})
	})

	Context("When table item does exist", func() {
		It("Updates the item", func() {
			err := db.UpdateItem("test", map[string]string{"firstName": "Evan"})
			Expect(err).To(BeNil())

			err = db.UpdateItem("test", map[string]string{"firstName": "EvanSia"})
			Expect(err).To(BeNil())

			out, err := db.GetByPk("test")

			var t Test
			err = attributevalue.UnmarshalMap(out, &t)
			Expect(t.FirstName).To(Equal("EvanSia"))
		})
	})
})
