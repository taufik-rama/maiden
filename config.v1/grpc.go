package config

// ServiceGRPCList ...
type ServiceGRPCList map[string]ServiceGRPC

func (s ServiceGRPCList) from(cfg serviceWrapper) ServiceGRPCList {
	s = make(ServiceGRPCList)
	for key, val := range cfg.Services.GRPC {
		s[key] = ServiceGRPC{
			Port:       val.Port,
			Definition: val.Definition,
			Methods:    (ServiceGRPCMethods)(val.Methods),
			Conditions: (ServiceGRPCConditions)(val.Conditions),
		}
	}
	return s
}

func (s ServiceGRPCList) replace(other ServiceGRPCList) {
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

// ServiceGRPC ...
type ServiceGRPC struct {
	Port       uint16
	Definition string
	Methods    ServiceGRPCMethods
	Conditions ServiceGRPCConditions
}

func (s ServiceGRPC) replace(other *ServiceGRPC) {
	if s.Port != 0 {
		other.Port = s.Port
	}
	if !emptyString(s.Definition) {
		other.Definition = s.Definition
	}
	if other.Methods == nil {
		other.Methods = s.Methods
	} else if s.Methods != nil {
		s.Methods.replace(other.Methods)
	}
	if other.Conditions == nil {
		other.Conditions = s.Conditions
	} else if s.Conditions != nil {
		s.Conditions.replace(other.Conditions)
	}
}

// ServiceGRPCMethods ...
type ServiceGRPCMethods map[string]struct {
	Request  string
	Response string
}

func (s ServiceGRPCMethods) replace(other ServiceGRPCMethods) {
	for key := range s {
		if _, ok := other[key]; ok {
			value := other[key]
			if !emptyString(s[key].Request) {
				value.Request = s[key].Request
			}
			if !emptyString(s[key].Response) {
				value.Response = s[key].Response
			}
			other[key] = value
		} else {
			other[key] = s[key]
		}
	}
}

// ServiceGRPCConditions ...
type ServiceGRPCConditions map[string][]struct {
	Request  interface{}
	Response interface{}
}

func (s ServiceGRPCConditions) replace(other ServiceGRPCConditions) {
	for key := range s {
		if _, ok := other[key]; ok {
			other[key] = append(other[key], s[key]...)
		} else {
			other[key] = s[key]
		}
	}
}
