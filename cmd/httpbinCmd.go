package cmd

import (
	"errors"
	"fmt"
	"github.com/mono83/tame/user"
	"github.com/spf13/cobra"
)

var httpbinCmd = &cobra.Command{
	Use:   "httpbin [name]",
	Short: "Runs self tests using https://httpbin.org",
	RunE: func(cmd *cobra.Command, args []string) error {
		toUse := asserts
		if len(args) == 1 {
			toUse = []assertion{}
			for _, a := range asserts {
				if a.name == args[0] {
					toUse = append(toUse, a)
				}
			}
		}

		for _, a := range toUse {
			u := user.New()
			err := a.assert(u)
			if err != nil {
				fmt.Println("[x]", a.name, err.Error())
				return err
			}
			fmt.Println(" + ", a.name)
		}

		return nil
	},
}

type assertion struct {
	name   string
	assert func(*user.User) error
}

var asserts = []assertion{
	{
		name: "User agent",
		assert: func(u *user.User) error {
			page, err := u.Get("https://httpbin.org/user-agent")
			if err != nil {
				return err
			}

			var t struct {
				UserAgent string `json:"user-agent"`
			}
			err = page.AsJSON(&t)
			if err != nil {
				return err
			}
			if u.UserAgent != t.UserAgent {
				return fmt.Errorf(
					"User agents not match\n%s\n%s\n",
					u.UserAgent,
					t.UserAgent,
				)
			}
			return nil
		},
	},
	{
		name: "Headers",
		assert: func(u *user.User) error {
			page, err := u.Get("https://httpbin.org/headers")
			if err != nil {
				return err
			}

			var t struct {
				Headers map[string]string `json:"headers"`
			}
			err = page.AsJSON(&t)
			if err != nil {
				return err
			}
			for name, value := range u.Header {
				v, ok := t.Headers[name]
				if !ok {
					return fmt.Errorf("Missing header %s", name)
				}
				if value != v {
					return fmt.Errorf(
						"Header value mismatch\n%s\n%s\n",
						value,
						v,
					)
				}
			}
			return nil
		},
	},
	{
		name: "Gzip",
		assert: func(u *user.User) error {
			page, err := u.Get("https://httpbin.org/gzip")
			if err != nil {
				return err
			}

			var t struct {
				Compress bool `json:"gzipped"`
			}
			err = page.AsJSON(&t)
			if err != nil {
				return err
			}
			if !t.Compress {
				return errors.New("Malformed response")
			}
			return nil
		},
	},
	{
		name: "Deflate",
		assert: func(u *user.User) error {
			page, err := u.Get("https://httpbin.org/deflate")
			if err != nil {
				return err
			}

			var t struct {
				Compress bool `json:"deflated"`
			}
			err = page.AsJSON(&t)
			if err != nil {
				return err
			}
			if !t.Compress {
				return errors.New("Malformed response")
			}
			return nil
		},
	},
	{
		name: "Status 404",
		assert: func(u *user.User) error {
			page, err := u.Get("https://httpbin.org/status/404")
			if err != nil {
				return err
			}

			if page.StatusCode != 404 {
				return fmt.Errorf("Expected 404, but got %d", page.StatusCode)
			}

			return nil
		},
	},
	{
		name: "Cookies",
		assert: func(u *user.User) error {
			page, err := u.Get("https://httpbin.org/cookies/set?foo=bar")
			if err != nil {
				return err
			}

			if v, ok := u.GetCookie(page.URL, "foo"); !ok || v != "bar" {
				return errors.New("Cookie wasnt set")
			}

			// More cookie check
			page, err = u.Get("https://httpbin.org/cookies")
			if err != nil {
				return err
			}

			if v, ok := u.GetCookie(page.URL, "foo"); !ok || v != "bar" {
				return errors.New("Cookie wasnt set")
			}

			// Removing cookies
			page, err = u.Get("https://httpbin.org/cookies/delete?foo")
			if err != nil {
				return err
			}

			cookieList := u.Cookies(page.URL)
			if len(cookieList) != 0 {
				return errors.New("Cookie wasnt deleted")
			}

			return nil
		},
	},
}
