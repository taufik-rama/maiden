package config

// HTTPServices ...
type HTTPServices map[string]HTTPService

func (h HTTPServices) replace(other HTTPServices) {
	for name := range other {
		if _, ok := h[name]; !ok {
			h[name] = other[name]
			continue
		}
		service := h[name]

		if Args.Service.PreferImports {
			service.replace(other[name])
		} else {
			service.replaceIfEmpty(other[name])
		}

		h[name] = service
	}
}

// HTTPService ...
type HTTPService struct {
	Port      uint16        `yaml:"port"`
	Endpoints HTTPEndpoints `yaml:"endpoints"`
}

// Replace the config values for `other`.
// The `imports` field is not replaced because we only need that field to resolve
// other configs and not for the actual values
func (h *HTTPService) replace(other HTTPService) {
	if other.Port != 0 {
		h.Port = other.Port
	}
	h.Endpoints.append(other.Endpoints)
}

// Replace the config values for `other` on empty values.
// The `imports` field is not replaced because we only need that field to resolve
// other configs and not for the actual values
func (h *HTTPService) replaceIfEmpty(other HTTPService) {
	if h.Port == 0 {
		h.Port = other.Port
	}
	h.Endpoints.append(other.Endpoints)
}

// HTTPEndpoints ...
type HTTPEndpoints map[string][]struct {

	// Method should be one of: string, array of string
	Method interface{} `yaml:"method"`

	// Request should be one of: string, array of string, object
	Request interface{} `yaml:"request"`

	// Response should be one of: string, object
	Response interface{} `yaml:"response"`
}

func (h HTTPEndpoints) append(other HTTPEndpoints) {
	for url := range other {
		if _, ok := h[url]; !ok {
			h[url] = other[url]
		}
		h[url] = append(h[url], other[url]...)
	}
}
