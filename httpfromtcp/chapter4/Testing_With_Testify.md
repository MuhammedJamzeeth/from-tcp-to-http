# Testing With Testify

In most of our Go courses, we encourage the use of the [standard library's testing package](https://pkg.go.dev/testing), and tend to use [table-driven](https://go.dev/wiki/TableDrivenTests) [tests](https://dave.cheney.net/2019/05/07/prefer-table-driven-tests)... but ThePrimeagen (ex-Netflix btw) has a different preferred approach, so we'll roll with that for this course!

No complex logic in my tests, please!!! I'll explain my approach below.

1.  I prefer bigger tests (higher level) to smaller tests (tiny units) because it means I can refactor my codebase more aggressively without having to rewrite a bunch of tests.
2.  I don't write tests for every function. I write tests for functions that I'm _unlikely to get right the first time_. So complicated logic, like our parser, is a perfect candidate for lots of tests.
3.  **I keep my tests simple and declarative**.

Point `#3` is why I tend to use [`testify`](https://github.com/stretchr/testify) and avoid the loops required for table-driven tests. This is what the test suite looked like in _Lane's_ way of doing things:

```go
tests := []struct {
		name            string
		request         string
		wantMethod      string
		wantTarget      string
		wantVersion     string
		wantErr         bool
	}{
		// test cases here...
}
for _, tt := range tests {
	t.Run(tt.name, func(t *testing.T) {
		r, err := RequestFromReader(strings.NewReader(tt.request))
		if tt.wantErr && err == nil {
			t.Fatalf("expected an error, but got none")
			return
		}
		if !tt.wantErr && err != nil {
			t.Fatalf("expected no error, but got one")
			return
		}
		if err != nil {
			return
		}
		if r == nil {
			t.Fatalf("expected a request, but got none")
		}
		if r.RequestLine.Method != tt.wantMethod {
			t.Fatalf("got method = %v, want %v", r.RequestLine.Method, tt.wantMethod)
		}
		if r.RequestLine.RequestTarget != tt.wantTarget {
			t.Fatalf("got target = %v, want %v", r.RequestLine.RequestTarget, tt.wantTarget)
		}
		if r.RequestLine.HttpVersion != tt.wantVersion {
			t.Fatalf("got version = %v, want %v", r.RequestLine.HttpVersion, tt.wantVersion)
		}
	})
```

I prefer a simple procedural approach:

```go
// Test: Good GET Request line
r, err := RequestFromReader(strings.NewReader("GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
require.NoError(t, err)
require.NotNil(t, r)
assert.Equal(t, "GET", r.RequestLine.Method)
assert.Equal(t, "/", r.RequestLine.RequestTarget)
assert.Equal(t, "1.1", r.RequestLine.HttpVersion)

// Test: Good GET Request line with path
r, err = RequestFromReader(strings.NewReader("GET /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
require.NoError(t, err)
require.NotNil(t, r)
assert.Equal(t, "GET", r.RequestLine.Method)
assert.Equal(t, "/coffee", r.RequestLine.RequestTarget)
assert.Equal(t, "1.1", r.RequestLine.HttpVersion)
```

Fewer conditionals, fewer things for me to get _wrong_ in my tests.

## Assignment

We'll be working on the request line parser, so let's get our directory structure and test dependency set up.

1.  Create an [`internal` directory](https://dave.cheney.net/2019/10/06/use-internal-packages-to-reduce-your-public-api-surface), and a `request` directory inside of it.

    ```sh
    mkdir -p ./internal/request
    ```

2.  Create a `request.go` file. Declare that it's part of the `request` package.
3.  Create a `request_test.go` file, it's also part of the `request` package. Our tests will go here.
4.  Install the [`testify`](https://github.com/stretchr/testify) package as a dependency in your module.

    ```sh
    go get -u github.com/stretchr/testify/assert
    ```

5.  Add a simple test (that doesn't test any application logic yet, but ensures you have `testify` set up correctly) to `request_test.go`:

    ```go
    package request

    import (
        "testing"

        "github.com/stretchr/testify/assert"
    )

    func TestRequestLineParse(t *testing.T) {
        assert.Equal(t, "TheTestagen", "TheTestagen")
    }
    ```

6.  Run your tests to ensure they work. You can run `go test ./...` from the root of your module to run all tests for all packages.

**Run and submit** the CLI tests from the **root of your project**.
