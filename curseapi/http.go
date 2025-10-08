package curseapi

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/maypok86/otter/v2"
)

var c = http.Client{Timeout: 10 * time.Second}

var key = os.Getenv("CURSE_API_KEY")

func httpget(ctx context.Context, url string) ([]byte, error) {
	reqs, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("httpget: %w", err)
	}
	reqs.Header.Set("Accept", "*/*")
	reqs.Header.Set("x-api-key", key)
	reqs.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.105 Safari/537.36")
	rep, err := c.Do(reqs)
	if err != nil {
		return nil, fmt.Errorf("httpget: %w", err)
	}
	defer rep.Body.Close()
	if rep.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("httpget: %w", ErrHttpCode{Code: rep.StatusCode})
	}
	b, err := io.ReadAll(rep.Body)
	if err != nil {
		return nil, fmt.Errorf("httpget: %w", err)
	}
	return b, err
}

type ErrHttpCode struct {
	Code int
}

func (e ErrHttpCode) Error() string {
	return fmt.Sprintf("ErrHttpCode: %d", e.Code)
}

func httpcache(ctx context.Context, url string, acache *otter.Cache[string, []byte]) ([]byte, error) {
	return acache.Get(ctx, url, otter.LoaderFunc[string, []byte](func(ctx context.Context, key string) ([]byte, error) {
		return httpget(ctx, url)
	}))

}
