# go-supertest

  HTTP assertions made easy for Go via [goreq](https://github.com/franela/goreq)

  We love [nodejs supertest](https://github.com/visionmedia/supertest) module and we wanted something like that for Go.

## Why?

  Because we want to keep HTTP testing in Go plain simple.
  

## Example


```go
import (
        . "github.com/franela/goblin"
        . "github.com/franela/go-supertest"
)

func MyTest(t *testing.T) {
  g := Goblin(t)

  g.Describe("GET /", function() {
    g.It("Should respond 200", function(done Done) {
      ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(200)
      }))
      defer ts.Close()

      NewRequest(ts.URL).
        Get("/").
        Expect(200, done)
    })
  })
}
```

## API


### .expect(status[, func])

  Assert response `status` code.

### .expect(status, body[, fn])

  Assert response `status` code and `body`.

## Notes

  Inspired by [supertest](https://github.com/visionmedia/supertest).
