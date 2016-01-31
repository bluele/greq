# greq

Yet another HTTP client library for go(golang).

# Example

```go
// examples/retry.go
package main

import (
	"errors"
	"fmt"
	"github.com/bluele/greq"
	"github.com/cenkalti/backoff"
	"io/ioutil"
	"net/http"
)

func main() {
	res, err := greq.Get("http://example.com/notfound.html").
		RequestHandler(greq.RetryBackoff(3, backoff.NewExponentialBackOff())).
		ResponseHandler(func(res *http.Response, err error) error {
		if res.StatusCode >= 400 && res.StatusCode < 500 {
			return errors.New("40X error")
		}
		return err
	}).Do()
	if err != nil {
		fmt.Println("error:", err.Error())
		return
	}
	body, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(body))
}
```

Output

```
$ go run examples/example.go
error: 40X error
```

# Author

**Jun Kimura**

* <http://github.com/bluele>
* <junkxdev@gmail.com>