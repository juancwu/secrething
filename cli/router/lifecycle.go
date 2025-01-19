package router

import tea "github.com/charmbracelet/bubbletea"

// LifecycleHooks defines the interface for page lifecycle events
type LifecycleHooks interface {
	// OnEnter is called when navigating to this page
	// Returns a command to be executed after entering
	// Params contain navigation data passed during route change
	OnEnter(params map[string]interface{}) tea.Cmd

	// OnExit is called when leaving this page (navigating away or going back)
	// Returns a cleanup command to be executed before leaving
	OnExit() tea.Cmd

	// BeforeNavigateAway is called before navigation occurs
	// Return true to allow navigation, false to prevent it
	// Can be used to show "unsaved changes" prompts
	BeforeNavigateAway() bool

	// AfterNavigateAway is called after navigation is confirmed but before OnExit
	// Use this to save state or clean up resources
	AfterNavigateAway() tea.Cmd
}
