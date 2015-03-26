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
		val := s.Env.Get("TEST")
		Expect(val.(string)).To(Equal("secrets, maybe"))
	})
	It("can get strings", func() {
		Expect(s.Env.GetString("TEST")).To(Equal("secrets, maybe"))
	})
	It("Can get uints", func() {
		os.Setenv("NUMBERS", "1234")
		Expect(s.Env.GetUint("NUMBERS")).To(Equal(uint64(1234)))
		Expect(func() { s.Env.GetUint("TEST") }).To(Panic())
	})
	It("Can get ints", func() {
		os.Setenv("NUMBERS", "5678")
		Expect(s.Env.GetInt("NUMBERS")).To(Equal(int64(5678)))
		Expect(func() { s.Env.GetInt("TEST") }).To(Panic())
	})
	It("can get bools", func() {
		Expect(s.Env.GetBool("NOTHERE")).To(Equal(false))
		os.Setenv("HAPPENING", "true")
		Expect(s.Env.GetBool("HAPPENING")).To(Equal(true))
		os.Setenv("NOTHAPPENING", "false")
		Expect(s.Env.GetBool("NOTHAPPENING")).To(Equal(false))
	})
	It("can compose a tcp service address", func() {
		os.Setenv("SERVICE_PORT_1337_TCP_ADDR", "testerator")
		os.Setenv("SERVICE_PORT_1337_TCP_PORT", "7")
		Expect(s.Env.GetTCPServiceAddress("service", 1337)).To(Equal("testerator:7"))
	})
	It("can compose a udp service address", func() {
		os.Setenv("SERVICE_PORT_1337_UDP_ADDR", "testerator")
		os.Setenv("SERVICE_PORT_1337_UDP_PORT", "7")
		Expect(s.Env.GetUDPServiceAddress("service", 1337)).To(Equal("testerator:7"))
	})
	It("returns empty string if environment is not set up for service", func() {
		Expect(s.Env.GetUDPServiceAddress("nullservice", 1337)).To(Equal(""))
	})
	It("panics if you use the must version", func() {
		Expect(func() { s.Env.MustGetTCPServiceAddress("nullservice", 1337) }).To(Panic())
	})
	It("says we're in prod if the environment does", func() {
		Expect(s.Env.InProd()).To(BeFalse())
		os.Setenv("ENV", "prod")
		Expect(s.Env.InProd()).To(BeTrue())
	})
})

func TestEnv(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "OsEnv")
}
