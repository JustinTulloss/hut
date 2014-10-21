package hut_test

import (
	"os"
	"strings"
	"testing"

	"github.com/JustinTulloss/hut"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("OsEnv", func() {
	var oldEnv []string
	var s *hut.Service
	BeforeEach(func() {
		s = hut.NewService(nil)
		oldEnv = os.Environ()
		os.Clearenv()
		os.Setenv("TEST", "secrets, maybe")
	})
	AfterEach(func() {
		s = nil
		for _, env := range oldEnv {
			s := strings.SplitN(env, "=", 2)
			os.Setenv(s[0], s[1])
		}
	})
	It("Can get configuration", func() {
		val, err := s.Env.Get("TEST")
		Expect(err).To(BeNil())
		Expect(val.(string)).To(Equal("secrets, maybe"))
	})
	It("can get strings", func() {
		val, err := s.Env.GetString("TEST")
		Expect(err).To(BeNil())
		Expect(val).To(Equal("secrets, maybe"))
	})
	It("Can get uints", func() {
		os.Setenv("NUMBERS", "1234")
		val, err := s.Env.GetUint("NUMBERS")
		Expect(err).To(BeNil())
		Expect(val).To(Equal(uint64(1234)))
		val, err = s.Env.GetUint("TEST")
		Expect(err).NotTo(BeNil())
	})
	It("Can get ints", func() {
		os.Setenv("NUMBERS", "5678")
		val, err := s.Env.GetInt("NUMBERS")
		Expect(err).To(BeNil())
		Expect(val).To(Equal(int64(5678)))
		val, err = s.Env.GetInt("TEST")
		Expect(err).NotTo(BeNil())
	})
	It("says we're in prod if the environment does", func() {
		Expect(s.Env.InProd()).To(Equal(false))
		os.Setenv("ENV", "prod")
		Expect(s.Env.InProd()).To(Equal(true))
	})
})

func TestEnv(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "OsEnv")
}
