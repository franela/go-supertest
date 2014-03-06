# Go-supertest

  [Goblin](https://github.com/franela/goblin) HTTP assertions made easy for Go 

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
      
      NewRequest(ts.URL).Get("/").Expect(200, done)
    })
  })
}

```

## API


### .expect(status[, Done])

  Assert response `status` code.

### .expect(status, payload[, Done])

  Assert response `status` code and `payload`.



## Notes

  Inspired by [supertest](https://github.com/visionmedia/supertest).
