package server

import (
	"encoding/json"
	"errors"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
	locationService "some-api/internal/location"
	"some-api/utils/db"
	"testing"
)

func TestApi(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "API test suite")
}

var _ = Describe("getLocation", func() {
	Context("When location exists for a given user", func() {
		It("Returns the location coordinates", func() {
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
			a := NewServer(mockDB)

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := a.Router.NewContext(req, rec)
			c.SetPath("/location/:id")
			c.SetParamNames("id")
			c.SetParamValues("Evan")

			err := a.getLocation(c)
			Expect(err).To(BeNil())
			Expect(rec.Code).To(Equal(http.StatusOK))

			var location locationService.Location
			_ = json.Unmarshal(rec.Body.Bytes(), &location)
			Expect(location.Latitude).To(Equal("1.12"))
			Expect(location.Longitude).To(Equal("0.9999"))
		})
	})

	Context("When location does not exist for a given user", func() {
		It("Returns an empty response", func() {
			mockDB := &db.MockDatabaseClient{
				GetByPkFn: func(pk string) (map[string]types.AttributeValue, error) {
					return nil, nil
				},
			}
			a := NewServer(mockDB)

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := a.Router.NewContext(req, rec)
			c.SetPath("/location/:id")
			c.SetParamNames("id")
			c.SetParamValues("Evan")

			err := a.getLocation(c)
			Expect(err).To(BeNil())
			Expect(rec.Code).To(Equal(http.StatusNotFound))
			Expect(rec.Body.String()).To(Equal("Not found"))
		})
	})

	Context("When there is an error with the request", func() {
		It("Returns an error response", func() {
			mockDB := &db.MockDatabaseClient{
				GetByPkFn: func(pk string) (map[string]types.AttributeValue, error) {
					return nil, errors.New("bad stuff")
				},
			}
			a := NewServer(mockDB)

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := a.Router.NewContext(req, rec)
			c.SetPath("/location/:id")
			c.SetParamNames("id")
			c.SetParamValues("Evan")

			err := a.getLocation(c)
			Expect(err).To(BeNil())
			Expect(rec.Code).To(Equal(http.StatusInternalServerError))
		})
	})
})
