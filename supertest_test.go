package supertest

import(
  . "github.com/franela/goblin"
  . "github.com/onsi/gomega"
  "testing"
)

func TestSuperTest(t *testing.T) {
  g := Goblin(t)

  RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

  g.Describe("Request(url)", func() {
    g.It("should be supported")

    g.Describe(".Send(interface{})", func() {
      g.It("should be supported")
    })

    g.Describe(".Send(string)", func() {
      g.It("should be supported")
    })

    g.Describe(".Set(map[string]string)", func() {
      g.It("should be supported")
      g.It("should support calling multiple times")
    })

    g.Describe(".Set(string, string)", func() {
      g.It("should be supported")
      g.It("should support calling multiple times")
    })

    g.Describe(".Query(map[string]string)", func() {
      g.It("should be supported")
      g.It("should support calling multiple times")
    })

    g.Describe(".Query(string, string)", func() {
      g.It("should be supported")
      g.It("should support calling multiple times")
    })

    g.Describe(".Query(string)", func() {
      g.It("should be supported")
      g.It("should support calling multiple times")
    })
  })

  g.Describe(".Expect(status[, fn])", func() {
    g.It("should assert the response status")
  })

  g.Describe(".Expect(status, body[, fn])", func() {
    g.It("should assert the response body and status")
  })
}
