# Parsing the Request Line

Now it's time for the best part of any project, the string parsing! By building on top of TCP, we already have code that handles plain-text data, now we just need to take that plain text and turn it into structured data, ensuring that it follows the HTTP protocol.

For example, given:

```http
POST /coffee HTTP/1.1
Host: localhost:42069
User-Agent: curl/7.81.0
Accept: */*
Content-Length: 21

{"flavor":"dark mode"}
```

We want our HTTP parser to return a struct that looks like this:

```go
type Request struct {
    RequestLine RequestLine
    Headers     map[string]string
    Body        []byte
}
```

That way our _application logic_ (the server code) has nicely structured HTTP data to work with.

Our goal is to take the server we created in the last section and have it parse out the [start-line](https://datatracker.ietf.org/doc/html/rfc9112#name-message-format) according to the RFC [message parsing section](https://datatracker.ietf.org/doc/html/rfc9112#name-message-parsing).

## The Request-Line

Remember how HTTP messages start with a `start-line`? Well, if it's a _request_ (not a response), then the `start-line` is called the `request-line` and has a specific format.

```http
HTTP-version  = HTTP-name "/" DIGIT "." DIGIT
HTTP-name     = %s"HTTP"
request-line  = method SP request-target SP HTTP-version
```

Which is a bit hairy to read (because its accounting for all the possible options), but for `HTTP/1.1` it's pretty simple, an example request-line looks like this:

```http
GET /coffee HTTP/1.1
```

## Assignment

For now, we're only going to parse the `request-line`.

Here we are, using test-driven-development, Uncle Bob would be proud.

On a real note, when it comes to blackbox RFCs like this, **good tests help us move (blazingly) faster** and avoid bugs.

1.  Remove your useless `assert` from the `request_test.go` file.
2.  Add the following structs to `request.go`:
    
    ```go
    type Request struct {
        RequestLine RequestLine
    }
    
    type RequestLine struct {
        HttpVersion   string
        RequestTarget string
        Method        string
    }
    ```
    
    Don't forget to [capitalize your exports](https://www.youtube.com/shorts/6FPVsJXD39o) smh...
    
    Hey, it's still better than `public static void main`...
    
3.  Stub out a new function with this signature: `func RequestFromReader(reader io.Reader) (*Request, error)`. The `Request` struct represents a full parsed HTTP request, but for now, we'll just be working on the `RequestLine` part.
4.  Add the following tests to `request_test.go`:
    
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
    
    // Test: Invalid number of parts in request line
    _, err = RequestFromReader(strings.NewReader("/coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
    require.Error(t, err)
    ```
    
5.  Implement the `RequestFromReader` function to parse the `request-line` from the reader. Here are some things to keep in mind:
    
    -   For now, you can slurp the entire request into memory using [`io.ReadAll`](https://pkg.go.dev/io#ReadAll) and work with the entire thing as a string.
    -   Create a `parseRequestLine` function to do the _parsing_.
    -   Remember that newlines in HTTP are `\r\n`, not just `\n`.
    -   You can discard everything that comes after the `request-line` for now.
    -   There are always just 3 parts to the request line: [`strings.Split`](https://pkg.go.dev/strings#Split) is your friend here.
    -   Verify that the "method" part only contains capital alphabetic characters.
    -   Verify that the http version part is `1.1`, extracted from the literal `HTTP/1.1` format, as we only support `HTTP/1.1` for now.
    -   Here are some additional references that might help (but don't overthink it, this step should be pretty straightforward):
        -   [2.3](https://datatracker.ietf.org/doc/html/rfc9112#section-2.3)
        -   [3.1](https://datatracker.ietf.org/doc/html/rfc9112#name-method)
        -   [3.2](https://datatracker.ietf.org/doc/html/rfc9112#name-request-target)
6.  Add more test cases to `request_test.go` to cover any edge cases you can think of. Here are the names of all the tests I wrote:
    
    -   Good Request line
    -   Good Request line with path
    -   Good POST Request with path
    -   Invalid number of parts in request line
    -   Invalid method (out of order) Request line
    -   Invalid version in Request line

When you're done with each lesson, look at the tests in the solution files to see if I did anything you didn't think of.

Don't worry about connecting the `request` package to any application (`cmd`) code yet - let's just get our tests passing first.

**Run and submit** the CLI tests.