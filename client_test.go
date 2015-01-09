package hut_test

import (
	"github.com/JustinTulloss/hut"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("clients", func() {
	var (
		s *hut.Service
	)
	BeforeEach(func() {
		s = hut.NewService(nil)
		s.Env = hut.MapEnv{
			"FIREBASE_URL":  "https://firebase-not-real",
			"FIREBASE_AUTH": "fake-auth",
		}
	})
	Describe("NewFirebaseClient", func() {
		It("should instantiate a firebase client when given creds", func() {
			client, err := s.NewFirebaseClient()
			Expect(err).To(BeNil())
			Expect(client).ToNot(BeNil())
		})
		It("Should not instantiate a firebase client if url is not defined", func() {
			s.Env.(hut.MapEnv)["FIREBASE_URL"] = ""
			client, err := s.NewFirebaseClient()
			Expect(err).ToNot(BeNil())
			Expect(client).To(BeNil())
		})
	})
})
