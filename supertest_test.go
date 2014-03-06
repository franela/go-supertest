package supertest

import(
  . "github.com/franela/goblin"
  . "github.com/onsi/gomega"
  "testing"
  "net/http/httptest"
  "net/http"
  "io/ioutil"
  "encoding/json"
  "bytes"
)

func TestSuperTest(t *testing.T) {
  g := Goblin(t)

  RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

  g.Describe("Supertest", func() {
    g.Describe("Request(url)", func() {
      g.It("should be supported", func(done Done) {
        ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
          w.WriteHeader(201)
          done()
        }))
        defer ts.Close()

        NewRequest(ts.URL).
          Get("/").Expect(200)
      })

      g.Describe(".Send(interface{})", func() {
        g.It("should be supported", func(done Done) {
          var o map[string]string
          payload := map[string]string { "foo": "bar" }

          ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            if body, err := ioutil.ReadAll(r.Body); err != nil {
              g.Fail(err)
              return
            } else if err := json.Unmarshal(body, &o); err != nil {
              g.Fail(err)
              return
            }
            g.Assert(o).Equal(payload)
            done()
          }))
          defer ts.Close()

          NewRequest(ts.URL).
            Post("/").
            Send(payload).
            Expect(200)
        })
      })

      g.Describe(".Send(string)", func() {
        g.It("should be supported", func(done Done) {
          payload := "foo"

          ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            body, err := ioutil.ReadAll(r.Body)
            if err != nil {
              g.Fail(err)
              return
            }
            g.Assert(string(body)).Equal(payload)
            done()
          }))
          defer ts.Close()

          NewRequest(ts.URL).
            Post("/").
            Send(payload).
            Expect(200)
        })
      })

      g.Describe(".Set(string, string)", func() {
        g.It("should be supported", func(done Done) {
          ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            g.Assert(r.Header.Get("foo")).Equal("bar")
            g.Assert(r.Header.Get("a")).Equal("b")
            done()
          }))
          defer ts.Close()

          NewRequest(ts.URL).
            Get("/").
            Set("foo", "bar").
            Set("a", "b").
            Expect(200)
        })
      })

      g.Describe(".Query(string, string)", func() {
        g.It("should be supported", func(done Done) {
          ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            values := r.URL.Query()

            g.Assert(values.Get("foo")).Equal("bar")
            g.Assert(values.Get("a")).Equal("b")
            done()
          }))
          defer ts.Close()

          NewRequest(ts.URL).
            Get("/").
            Query("foo", "bar").
            Query("a", "b").
            Expect(200)
        })
      })
    })

    g.Describe(".Expect(status)", func() {
      g.It("should be supported", func(done Done) {
        ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
          w.WriteHeader(204)
        }))
        defer ts.Close()

        NewRequest(ts.URL).
          Get("/").
          Expect(204, done)
      })
    })

    g.Describe(".Expect(status, body)", func() {
      g.It("should assert the response body as json", func(done Done) {
        payload := map[string]string { "a": "b" }

        ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
          b, _ := json.Marshal(payload)
          w.WriteHeader(200)
          w.Write(b)
        }))
        defer ts.Close()

        NewRequest(ts.URL).
          Get("/").
          Expect(200, payload, done)
      })

      g.It("should assert the response body as json using structs", func(done Done) {
        type Test struct {
          Foo string
          Bar string
        }

        payload := Test{ Foo: "foo", Bar: "bar" }

        ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
          b, _ := json.Marshal(payload)
          w.WriteHeader(200)
          w.Write(b)
        }))
        defer ts.Close()

        NewRequest(ts.URL).
          Get("/").
          Expect(200, payload, done)
      })

      g.It("should fail when body is different", func() {
        type Test struct {
          Foo string
          Bar string
        }

        payload := Test{ Foo: "foo", Bar: "bar" }

        ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
          w.WriteHeader(200)
          w.Write(bytes.NewBufferString("[]").Bytes())
        }))
        defer ts.Close()

        err := NewRequest(ts.URL).
          Get("/").
          Expect(200, payload)

        g.Assert(err != nil).IsTrue()
      })

      g.It("should assert the response body as string", func(done Done) {
        ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
          w.WriteHeader(200)
          w.Write(bytes.NewBufferString("foo").Bytes())
        }))
        defer ts.Close()

        NewRequest(ts.URL).
          Get("/").
          Expect(200, "foo", done)
      })
    })
  })
}
