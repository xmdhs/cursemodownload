package curseapi

import (
	"time"

	"github.com/maypok86/otter/v2"
)

var acache = otter.Must(&otter.Options[string, []byte]{
	MaximumSize:      1000,
	ExpiryCalculator: otter.ExpiryAccessing[string, []byte](12 * time.Hour),
})
