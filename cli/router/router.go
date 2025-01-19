package router

import (
	"errors"

	tea "github.com/charmbracelet/bubbletea"
)

var (
	errInvalidPage error = errors.New("Invalid page")
)

// Page represents an ID for a specific tea.Model
type Page string

// Router represents a basic router to manage tea.Model as pages.
type Router struct {
	currentPage Page
	pages       map[Page]tea.Model
}

func New(currentPage Page, pages map[Page]tea.Model) Router {
	return Router{
		currentPage: currentPage,
		pages:       pages,
	}
}

// Page returns the current tea.Model or page is on
func (r *Router) Page() tea.Model {
	return r.pages[r.currentPage]
}

// Navigate sets the current page to the given one
func (r *Router) Navigate(to Page) error {
	if !r.IsValidPageID(to) {
		return errInvalidPage
	}
	r.currentPage = to
	return nil
}

// PageID gets the current page id
func (r *Router) PageID() Page {
	return r.currentPage
}

// IsValidPageID is a helper method to verify that the page id is valid within the router
func (r *Router) IsValidPageID(p Page) bool {
	_, ok := r.pages[p]
	return ok
}
