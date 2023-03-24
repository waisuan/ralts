package chat

import "testing"

func TestChat_LoadAllMessages(t *testing.T) {

}

//func TestChat(t *testing.T) {
//	RegisterFailHandler(Fail)
//	RunSpecs(t, "Chat test suite")
//}
//
//var _ = Describe("LoadAllMessages", func() {
//	var dbClient *db.RaltsDatabase
//	var chat *Chat
//
//	BeforeEach(func() {
//		dbClient = db.NewRaltsDatabase("postgres://postgres:password@localhost:5433/ralts_test?sslmode=disable")
//		chat = NewChat(dbClient)
//		testHelper.ClearDB() // Just in case...
//		testHelper.InitDB()
//	})
//
//	AfterEach(func() {
//		dbClient.Close()
//		testHelper.ClearDB()
//	})
//
//	Context("When there is only one message available", func() {
//		BeforeEach(func() {
//			_, _ = chat.SaveMessage("evansia", "Lorem Ipsum", time.Now)
//		})
//
//		It("Returns that one message", func() {
//			msgs, _ := chat.LoadAllMessages()
//
//			Expect(len(msgs)).To(Equal(1))
//			msg := msgs[0]
//			Expect(msg.ChatId).To(Equal(int64(1)))
//			Expect(msg.Username).To(Equal("evansia"))
//			Expect(msg.Message).To(Equal("Lorem Ipsum"))
//			Expect(msg.CreatedAt).ToNot(BeNil())
//		})
//	})
//
//	Context("When there are multiple messages available", func() {
//		BeforeEach(func() {
//			for i := 0; i < 10; i++ {
//				_, _ = chat.SaveMessage("evansia", "Lorem Ipsum", time.Now)
//			}
//		})
//
//		It("Returns all of the messages", func() {
//			msgs, _ := chat.LoadAllMessages()
//
//			Expect(len(msgs)).To(Equal(10))
//			for i, m := range msgs {
//				Expect(m.ChatId).To(Equal(int64(i + 1)))
//			}
//		})
//	})
//})
//
//var _ = Describe("SaveMessage", func() {
//	var dbClient *db.RaltsDatabase
//	var chat *Chat
//
//	BeforeEach(func() {
//		dbClient = db.NewRaltsDatabase("postgres://postgres:password@localhost:5433/ralts_test?sslmode=disable")
//		chat = NewChat(dbClient)
//		testHelper.ClearDB() // Just in case...
//		testHelper.InitDB()
//	})
//
//	AfterEach(func() {
//		dbClient.Close()
//		testHelper.ClearDB()
//	})
//
//	Context("When called", func() {
//		It("Returns the saved message record", func() {
//			msg, err := chat.SaveMessage("esia", "ABCD == EFGH", time.Now)
//			Expect(err).To(BeNil())
//			Expect(msg.ChatId).To(Equal(int64(1)))
//			Expect(msg.Username).To(Equal("esia"))
//			Expect(msg.Message).To(Equal("ABCD == EFGH"))
//
//			formattedCreatedAt := msg.CreatedAt.Format("2006-01-02 15:04")
//			formattedNow := time.Now().Format("2006-01-02 15:04")
//			Expect(formattedCreatedAt).To(Equal(formattedNow))
//		})
//	})
//})
