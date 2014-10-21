package hut_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/JustinTulloss/hut"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const jsonContentType = "application/json; charset=utf-8"

func unmarshalBody(body *bytes.Buffer) map[string]interface{} {
	rep := map[string]interface{}{}
	err := json.Unmarshal(body.Bytes(), &rep)
	if err != nil {
		Fail(err.Error())
	}
	return rep
}

var _ = Describe("Replies", func() {
	var (
		w *httptest.ResponseRecorder
		s *hut.Service
	)
	BeforeEach(func() {
		w = httptest.NewRecorder()
		s = hut.NewService(nil)
	})
	It("should send a payload as a json encoded string", func() {
		payload := map[string]interface{}{
			"good morning": "vietnam",
		}
		s.Reply(payload, w)
		Expect(w.Code).To(Equal(http.StatusOK))
		ct := w.Header().Get("Content-Type")
		Expect(ct).To(Equal(jsonContentType))
		rep := unmarshalBody(w.Body)
		Expect(rep).To(Equal(payload))
	})
	It("should send an error if the payload is not jsonifiable", func() {
		payload := map[string]interface{}{
			"good morning": s.Reply,
		}
		s.Reply(payload, w)
		Expect(w.Code).To(Equal(http.StatusInternalServerError))
	})
	Describe("Success", func() {
		It("Sends back a standard message without metadata", func() {
			s.SuccessReply(nil, w)
			Expect(w.Code).To(Equal(http.StatusOK))
			rep := unmarshalBody(w.Body)
			Expect(len(rep)).To(Equal(1))
			Expect(rep["status"]).To(Equal("OK"))
		})
		It("Sends back standard message with provided metadata", func() {
			payload := map[string]interface{}{
				"id": float64(1234), // More convenient since JSON is all floats
				"tx": "12njcn28",
			}
			s.SuccessReply(payload, w)
			Expect(w.Code).To(Equal(http.StatusOK))
			rep := unmarshalBody(w.Body)
			Expect(len(rep)).To(Equal(3))
			Expect(rep["status"]).To(Equal("OK"))
			for k, v := range payload {
				Expect(rep[k]).To(Equal(v))
			}
		})
	})
	Describe("Error", func() {
		It("Sends back appropriate code and message", func() {
			myErr := errors.New("A problem occurred!")
			s.ErrorReply(myErr, w)
			Expect(w.Code).To(Equal(http.StatusInternalServerError))
			rep := unmarshalBody(w.Body)
			Expect(int(rep["code"].(float64))).To(Equal(http.StatusInternalServerError))
			Expect(rep["message"]).ToNot(BeNil())
		})
	})
	Describe("HttpError", func() {
		It("Sends JSON encoded reply", func() {
			s.HttpErrorReply(w, "things are bad", http.StatusGone)
			Expect(w.Code).To(Equal(http.StatusGone))
			Expect(w.Header().Get("Content-Type")).To(Equal(jsonContentType))
			rep := unmarshalBody(w.Body)
			Expect(rep["message"]).To(Equal("things are bad"))
			Expect(rep["code"]).To(Equal(float64(http.StatusGone)))
		})
	})
	Describe("NotFound", func() {
		It("Returns a JSON NotFound message", func() {
			s.NotFound(w)
			Expect(w.Code).To(Equal(http.StatusNotFound))
			rep := unmarshalBody(w.Body)
			Expect(rep["code"]).To(Equal(float64(http.StatusNotFound)))
		})
	})
})

func TestReplies(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Replies")
}
