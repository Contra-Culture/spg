package spg

type (
	// Page - represents a generated page.
	Page struct {
		Headers map[string]string
		Body    string
	}
	// Node - represents a tree of rendered pages (website itself).
	Node struct {
		page     *Page
		children map[string]*Node
	}
)
