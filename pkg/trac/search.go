package trac

import "fmt"

// Search trac.
type Search struct {
	client *Client
}

// SearchFilters retrieve a list of search filters with each element in the
// form (name, description).
// Not implemented.
func (s *Search) SearchFilters() ([]string, error) {
	return nil, fmt.Errorf("Not implemented")
}

// Search using the given filters. Defaults to all if not provided. Results are
// returned as a list of tuples in the form (href, title, date, author,
// excerpt).
// Not implemented.
func (s *Search) Search(query string, filters []string) ([]string, error) {
	return nil, fmt.Errorf("Not implemented")
}
