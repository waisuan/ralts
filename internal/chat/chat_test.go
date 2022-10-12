package chat

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"some-api/internal/db"
	testHelper "some-api/internal/testing"
	"testing"
	"time"
)

func TestChat(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Location service test suite")
}

var _ = Describe("Load", func() {
	var dbClient *db.DatabaseClient
	var chat *Chat

	BeforeEach(func() {
		dbClient = db.New()
		chat = NewChat(dbClient)
	})

	AfterEach(func() {
		testHelper.ClearDB(dbClient)
	})

	Context("Given a day and the current time", func() {
		It("Returns messages ranging from 3 hours ago and on that day", func() {
			msg, _ := chat.Load(time.Now(), time.Now)
			Expect(msg).To(BeEmpty())
		})
	})
})
