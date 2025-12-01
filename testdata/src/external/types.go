// Package external simulates a third-party package
// where we cannot add "// required" comments.
package external

type User struct {
	ID    string
	Name  string
	Email string
}

type Config struct {
	APIKey  string
	Timeout int
}
