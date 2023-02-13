package server

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"testing"
)

func TestApi(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "API test suite")
}

//var _ = Describe("/login", func() {
//	Context("When called with an invalid ID token", func() {
//		It("Returns a 401 error", func() {
//			e := echo.New()
//			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"auth_token":"TOKEN"}`))
//			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
//			rec := httptest.NewRecorder()
//			c := e.NewContext(req, rec)
//
//			err := loginUser(c)
//			Expect(err).To(BeNil())
//			Expect(rec.Code).To(Equal(http.StatusUnauthorized))
//		})
//	})
//})
