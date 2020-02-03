package config

// ServiceHTTPList ...
type ServiceHTTPList map[string]ServiceHTTP

func (s ServiceHTTPList) from(cfg serviceWrapper) ServiceHTTPList {
	for key, val := range cfg.Services.HTTP {
		s[key] = ServiceHTTP{
			Port:      val.Port,
			Endpoints: (ServiceHTTPEndpoints)(val.Endpoints),
		}
	}
	return s
}

func (s ServiceHTTPList) replace(other ServiceHTTPList) {
	for key := range s {
		if _, ok := other[key]; ok {
			value := other[key]
			s[key].replace(&value)
			other[key] = value
		} else {
			other[key] = s[key]
		}
	}
}

// ServiceHTTP ...
type ServiceHTTP struct {
	Port      uint16
	Endpoints ServiceHTTPEndpoints
}

func (s ServiceHTTP) replace(other *ServiceHTTP) {
	if s.Port != 0 {
		other.Port = s.Port
	}
	if other.Endpoints == nil {
		other.Endpoints = s.Endpoints
	} else if s.Endpoints != nil {
		s.Endpoints.replace(other.Endpoints)
	}
}

// ServiceHTTPEndpoints ...
type ServiceHTTPEndpoints map[string][]struct {
	Method   interface{}
	Request  interface{}
	Response interface{}
}

func (s ServiceHTTPEndpoints) replace(other ServiceHTTPEndpoints) {
	for key := range s {
		if _, ok := other[key]; ok {
			value := other[key]
			value = append(value, s[key]...)
			other[key] = value
		} else {
			other[key] = s[key]
		}
	}
}
