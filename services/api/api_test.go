package api

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
	"some-api/utils/db"
	"strings"
	"testing"
)

func TestApi(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "API test suite")
}

var _ = Describe("getLocation", func() {
	Context("Happy", func() {
		It("Runs", func() {
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
			a := NewApi(mockDB)

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := a.echo.NewContext(req, rec)
			c.SetPath("/location/:id")
			c.SetParamNames("id")
			c.SetParamValues("Evan")

			err := a.getLocation(c)
			Expect(err).To(BeNil())
			Expect(rec.Code).To(Equal(http.StatusOK))
			Expect(strings.TrimSuffix(rec.Body.String(), "\n")).To(
				Equal(`{"Pk":"1234","PersonID":"1234","Latitude":"1.12","Longitude":"0.9999"}`),
			)
		})
	})
})