package spg

type (
	Page struct {
		Headers map[string]string
		Body    string
	}
	Node struct {
		page     *Page
		children map[string]*Node
	}
)
