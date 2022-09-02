package location

import (
	"errors"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"some-api/utils/db"
	"testing"
)

func TestLocation(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Location service test suite")
}

var _ = Describe("Get", func() {
	Context("When table item exists", func() {
		It("Returns the table item", func() {
			mockDB := &db.MockDatabaseClient{
				GetByPkFn: func(pk string) (map[string]types.AttributeValue, error) {
					return map[string]types.AttributeValue{
						"pk":        &types.AttributeValueMemberS{Value: "1234"},
						"personID":  &types.AttributeValueMemberS{Value: "1234"},
						"latitude":  &types.AttributeValueMemberS{Value: "1.12"},
						"longitude": &types.AttributeValueMemberS{Value: "0.9999"},
					}, nil
				},
			}
			location, err := Get(mockDB, "1234")
			Expect(err).To(BeNil())
			Expect(location.PersonID).To(Equal("1234"))
			Expect(location.Latitude).To(Equal("1.12"))
			Expect(location.Longitude).To(Equal("0.9999"))
		})
	})

	Context("When table item does not exist", func() {
		It("Returns an empty response", func() {
			mockDB := &db.MockDatabaseClient{
				GetByPkFn: func(pk string) (map[string]types.AttributeValue, error) {
					return nil, nil
				},
			}
			location, err := Get(mockDB, "1234")
			Expect(err).To(BeNil())
			Expect(location).To(BeNil())
		})
	})

	Context("When there is an error", func() {
		It("Returns the error", func() {
			mockDB := &db.MockDatabaseClient{
				GetByPkFn: func(pk string) (map[string]types.AttributeValue, error) {
					return nil, errors.New("bad stuff")
				},
			}
			location, err := Get(mockDB, "1234")
			Expect(err).ToNot(BeNil())
			Expect(location).To(BeNil())
		})
	})
})

var _ = Describe("Save", func() {
	Context("When there are no errors", func() {
		It("Returns an empty response", func() {
			mockDB := &db.MockDatabaseClient{
				UpdateItemFn: func(pk string, attributes map[string]string) error {
					Expect(pk).ToNot(BeEmpty())
					Expect(attributes).ToNot(BeEmpty())
					return nil
				},
			}
			err := Save(mockDB, "1234", "1.1", "0.00003")
			Expect(err).To(BeNil())
		})
	})

	Context("When there are errors", func() {
		It("Returns the error", func() {
			mockDB := &db.MockDatabaseClient{
				UpdateItemFn: func(pk string, attributes map[string]string) error {
					Expect(pk).ToNot(BeEmpty())
					Expect(attributes).ToNot(BeEmpty())
					return errors.New("bad stuff")
				},
			}
			err := Save(mockDB, "1234", "1.1", "0.00003")
			Expect(err).ToNot(BeNil())
		})
	})
})
