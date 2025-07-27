package testing

import (
	"sync"

	cm "github.com/cysp/terraform-provider-censusworkspace/internal/census-management-go"
)

type Handler struct {
	mu sync.Mutex

	Sources      map[string]*cm.SourceData
	sourceIDLast int64
}

var _ cm.Handler = (*Handler)(nil)
