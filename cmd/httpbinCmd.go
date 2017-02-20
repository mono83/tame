package cmd

import (
	"errors"
	"fmt"
	"github.com/mono83/tame/user"
	"github.com/spf13/cobra"
)

var httpbinCmd = &cobra.Command{
	Use:   "httpbin",
	Short: "Runs self tests using https://httpbin.org",
	RunE: func(cmd *cobra.Command, args []string) error {
		for _, a := range asserts {
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
}
