package recipes

import "fmt"

// recipeError is special error object, used by recipes
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

// Error builds and returns recipe error with text
func Error(recipe, message string) error {
	return recipeError{recipeName: recipe, message: message}
}

// Errorf builds and returns recipe error with text in sprintf format
func Errorf(recipe string, format string, args ...interface{}) error {
	return recipeError{recipeName: recipe, message: fmt.Sprintf(format, args...)}
}

// ErrorCaused builds and returns error, built over other error
func ErrorCaused(recipe string, cause error) error {
	return recipeError{recipeName: recipe, cause: cause}
}
