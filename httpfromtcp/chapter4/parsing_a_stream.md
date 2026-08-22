This lesson's interactive features are locked, please upgrade to keep using them

# Parsing a Stream

Unfortunately parsing code tends to be just one edge case after another. Remember how I said TCP **guarantees data to be in order**? That's true, but I never said it had to be _complete_. TCP (and by extension, HTTP) is a _streaming protocol_, which means we receive data in chunks and should be able to parse it as it comes in.

So, instead of a full HTTP request, we might just get the first few characters, like this:

```text
GE
```

We need to manage the _state_ of our parser to handle incomplete reads. For example, maybe in the first pass, our parser only gets:

```text
GE
```

It needs to be smart enough to know that it's not done yet and keep reading until it gets the full request line:

```text
GET /coffee HTTP/1.1
```

## Assignment

1.  Paste this code into your `request_test.go` file:
    
    ```go
    type chunkReader struct {
        data            string
        numBytesPerRead int
        pos             int
    }
    
    // Read reads up to len(p) or numBytesPerRead bytes from the string per call
    // its useful for simulating reading a variable number of bytes per chunk from a network connection
    func (cr *chunkReader) Read(p []byte) (n int, err error) {
        if cr.pos >= len(cr.data) {
            return 0, io.EOF
        }
        endIndex := cr.pos + cr.numBytesPerRead
        if endIndex > len(cr.data) {
            endIndex = len(cr.data)
        }
        n = copy(p, cr.data[cr.pos:endIndex])
        cr.pos += n
    
        return n, nil
    }
    ```
    
2.  Update your test suite to use the `chunkReader` type and test for different numbers of bytes read per chunk:
    
    ```go
    // Test: Good GET Request line
    reader := &chunkReader{
        data:            "GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
        numBytesPerRead: 3,
    }
    r, err := RequestFromReader(reader)
    require.NoError(t, err)
    require.NotNil(t, r)
    assert.Equal(t, "GET", r.RequestLine.Method)
    assert.Equal(t, "/", r.RequestLine.RequestTarget)
    assert.Equal(t, "1.1", r.RequestLine.HttpVersion)
    
    // Test: Good GET Request line with path
    reader = &chunkReader{
        data:            "GET /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
        numBytesPerRead: 1,
    }
    r, err = RequestFromReader(reader)
    require.NoError(t, err)
    require.NotNil(t, r)
    assert.Equal(t, "GET", r.RequestLine.Method)
    assert.Equal(t, "/coffee", r.RequestLine.RequestTarget)
    assert.Equal(t, "1.1", r.RequestLine.HttpVersion)
    ```
    
    Be sure to test values as low as 1 and as high as the length of the request string. _Our code should work under all conditions._
    
3.  Update your `parseRequestLine` to return the number of bytes it consumed. If it can't find an `\r\n` (this is important!) it should return `0` and _no error_. This just means that it needs more data before it can parse the request line.
4.  Add a new internal "enum" (I just used an `int`) to your `Request` struct to track the state of the parser. For now, you just need 2 states:
    
    -   "initialized"
    -   "done".
5.  Implement a new `func (r *Request) parse(data []byte) (int, error)` method.
    
    -   It accepts all currently unparsed bytes from the buffer
    -   It updates the "state" of the parser, and the parsed `RequestLine` field.
    -   It returns the number of bytes it consumed (meaning successfully parsed) and an error if it encountered one.
    
    If you want additional help, see the `Tips` section below.
    
6.  Update the `RequestFromReader` function.
    
    -   Instead of reading _all_ the bytes, and then parsing the request line, it should use a loop to continually read from the reader and parse new chunks using the `parse` method.
    -   The loop should continue until the parser is in the "done" state.
    -   You'll need to keep track of:
        -   A buffer to read data into (`[]byte`). I started with a size of `8` and grew it as needed. I also shifted data in and out of it so I don't need to keep storing already-parsed data.
        -   How many bytes you've _read_ from the reader
        -   How many bytes you've _parsed_ from the buffer
    -   The end result is the same (aside from the fact that it properly handles chunks _as they arrive_) in that it returns a parsed `Request` struct once the `reader` is exhausted.

**Run and submit** the CLI tests.

## Tips

**Implementation help for `func (r *Request) parse(data []byte) (int, error)`**:

-   If the state of the parser is "initialized", it should call `parseRequestLine`.
    -   If there is an error, it should just return the error.
    -   If zero bytes are parsed, but no error is returned, it should return `0` and `nil`: it needs more data.
    -   If bytes are consumed successfully, it should update the `.RequestLine` field and change the state to "done".
-   If the state of the parser is "done", it should return an error that says something like "error: trying to read data in a done state"
-   If the state is anything else, it should return an error that says something like "error: unknown state"

**Implementation help for `RequestFromReader`**:

-   It shouldn't call `io.ReadAll` anymore. Instead, it should create a new buffer: `buf := make([]byte, bufferSize, bufferSize)`. Set `bufferSize` as a constant at the top of the file, and for now, just a size of `8`. We want to test with small buffers to make sure our parser can handle it.
-   Create a new `readToIndex` variable and set it to `0`. This will keep track of how much data we've read from the `io.Reader` into the buffer.
-   Create a new `Request` struct and set the state to "initialized".
-   While the state of the parser is not "done":
    -   If the buffer is full (we've read data into the entire buffer), grow it. Create a new slice that's twice the size and [`copy`](https://pkg.go.dev/builtin#copy) the old data into the new slice.
    -   Read from the `io.Reader` into the buffer starting at `readToIndex`.
        -   If you hit the end of the reader ([`io.EOF`](https://pkg.go.dev/io#pkg-variables)) set the state to "done" and break out of the loop.
        -   Update `readToIndex` with the number of bytes you actually read
        -   Call `r.parse` passing the slice of the buffer that has data that you've actually read so far
        -   Remove the data that was parsed successfully from the buffer (this keeps our buffer small and memory efficient). I used the `copy` function and a new slice to do this.
        -   Decrement the `readToIndex` by the number of bytes that were parsed so that it matches the new length of the buffer.