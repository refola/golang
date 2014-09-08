package core

import (
	"testing"
)

func testTranslation(data map[string]string, fn func(string) string, t *testing.T) {
	for i, v := range data {
		go func(i, v string) {
			t.Logf("Trying to convert:\n\t\"%s\"\nInto:\n\t\"%s\"\n", i, v)
			s := fn(i)
			if v != s {
				t.Errorf("FAIL:\n\t\"%s\"\n\n", s)
			} else {
				t.Log("SUCCESS.\n")
			}
		}(i, v)
	}
}

func TestSanitize(t *testing.T) {
	data := map[string]string{
		"abcdefghijklmnopqrstuvwxyz/ABCDEFGHIJKLMNOPQRSTUVWXYZ//0123456789": "abcdefghijklmnopqrstuvwxyz/ABCDEFGHIJKLMNOPQRSTUVWXYZ//0123456789", // both should be the same
		"`~!@#$%^&*()_+|-={}[]:\" \\;'<>?,./":                               "%60%7E%21%40%23%24%25%5E%26%2A%28%29%5F%2B%7C%2D%3D%7B%7D%5B%5D%3A%22%20%5C%3B%27%3C%3E%3F%2C%2E/",
		"hello, world.":                                                     "hello%2C%20world%2E"}
	testTranslation(data, Sanitize, t)
}

func TestUnsanitize(t *testing.T) {
	data := map[string]string{"%2c%20%2e": ", .",
		"Hello, world":                               "Hello, world",
		"%22%48%65%6c%6c%6f%2c%20%77%6f%72%6c%64%22": "\"Hello, world\"",
		"%20": " "}
	testTranslation(data, Unsanitize, t)
}

func TestUncleanChars(t *testing.T) {
	data := map[string]string{"%20": " ",
		"%68%69":       "hi",
		"%25%66%6f%6F": "%foo"}
	f := func(s string) string {
		return string(uncleanChars([]byte(s)))
	}
	testTranslation(data, f, t)
}
