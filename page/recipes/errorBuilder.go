package recipes

import "fmt"

// recipeError is special error object, emitted
type recipeError struct {
	recipeName, message string
	cause               error
}

func (r recipeError) Error() string {
	if len(r.message) > 0 && r.cause != nil {
		return "Recipe [" + r.recipeName + "] error. " + r.message + " caused by " + r.cause.Error()
	} else if r.cause != nil {
		return "Recipe [" + r.recipeName + "] error. " + r.cause.Error()
	} else {
		return "Recipe [" + r.recipeName + "] error. " + r.message
	}
}

// ErrorBuilder is wrapper, used to generate errors
type ErrorBuilder struct {
	RecipeName string
}

// Error builds and returns recipe error with text
func (b ErrorBuilder) Error(message string) error {
	return recipeError{recipeName: b.RecipeName, message: message}
}

// Errorf builds and returns recipe error with text in sprintf format
func (b ErrorBuilder) Errorf(format string, args ...interface{}) error {
	return recipeError{recipeName: b.RecipeName, message: fmt.Sprintf(format, args...)}
}

// CausedBy builds and returns error, built over other error
func (b ErrorBuilder) CausedBy(cause error) error {
	return recipeError{recipeName: b.RecipeName, cause: cause}
}

// ErrorEmptyPage return error struct for cases, when nil supplied to recipe instead of Page
func (b ErrorBuilder) ErrorEmptyPage() error {
	return recipeError{recipeName: b.RecipeName, message: "Empty page struct supplied for recipe"}
}
