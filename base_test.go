
package statsd

import (
        "testing"
        "fmt"
)

func Test_NewClient(t *testing.T) {
        client := NewClient(&Config{Host:"localhost:8125"})

        if client == nil {
                fmt.Printf("Client nil!\n")
        }
}

