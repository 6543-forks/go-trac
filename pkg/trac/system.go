package trac

import (
	"encoding/json"
	"fmt"
)

// APIVersion represents the remote version.
// First element is the epoch (0=Trac 0.10, 1=Trac 0.11 or higher). Second
// element is the major version number, third is the minor.
type APIVersion struct {
	Epoch int
	Major int
	Minor int
}

// UnmarshalJSON unmarshals returned array into an APIVersion
func (v *APIVersion) UnmarshalJSON(in []byte) error {
	ver := []interface{}{
		&v.Epoch,
		&v.Major,
		&v.Minor,
	}
	if err := json.Unmarshal(in, &ver); err != nil {
		return err
	}
	return nil
}

// System represents the core of the RPC system.
type System struct {
	client *Client
}

// APIVersion returns the version of the API.
func (s *System) APIVersion() (APIVersion, error) {
	var v = APIVersion{}
	r, err := s.client.Query("system.getAPIVersion")
	if err != nil {
		return v, err
	}
	if err := json.Unmarshal(r.Result, &v); err != nil {
		return v, err
	}
	return v, nil
}

// Methods  returns a list of strings, one for each (non-system) method
// supported by the RPC server.
func (s *System) Methods() ([]string, error) {
	return s.client.All("system.listMethods")
}

// MethodHelp method takes one parameter, the name of a method implemented by
// the RPC server. It returns a documentation string describing the use of that
// method. If no such string is available, an empty string is returned. The
// documentation string may contain HTML markup.
func (s *System) MethodHelp(method string) (string, error) {
	var m string
	r, err := s.client.Query("system.methodHelp", method)
	if err != nil {
		return m, err
	}
	if err := json.Unmarshal(r.Result, &m); err != nil {
		return m, err
	}
	return m, nil
}

// MethodSignature is not implemented.
func (s *System) MethodSignature(method string) ([]string, error) {
	return nil, fmt.Errorf("Not implemented")
}
