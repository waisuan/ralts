package chat

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"ralts/internal/db"
	testHelper "ralts/internal/testing"
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
		nowFn := func() time.Time {
			t, _ := time.Parse("200601021504", "202210201600")
			return t
		}

		BeforeEach(func() {
			_, _ = chat.SaveMessage("user", "I am a message", nowFn)
		})

		Context("There are messages from 3 hours ago", func() {
			It("Returns those messages", func() {
				msg, _ := chat.LoadAllMessages(nowFn(), nowFn)

				Expect(len(msg)).To(Equal(1))
				Expect(msg[0].Text).To(Equal("I am a message"))
			})
		})

		Context("There aren't any messages from 3 hours ago", func() {
			It("Returns an empty response", func() {
				t, _ := time.Parse("200601021504", "202310201600")
				msg, _ := chat.LoadAllMessages(t, nowFn)

				Expect(msg).To(BeEmpty())
			})
		})
	})
})
