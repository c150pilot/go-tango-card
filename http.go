package tango

import (
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"
)

func ensureSuccessStatus(resp *resty.Response, operation string) error {
	if resp.StatusCode() >= 200 && resp.StatusCode() < 300 {
		return nil
	}

	return fmt.Errorf("%s failed with status %d (%s): %s", operation, resp.StatusCode(), resp.Status(), strings.TrimSpace(string(resp.Body())))
}
