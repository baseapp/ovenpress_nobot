package lib

import (
	"bytes"
	"html/template"
	"strings"
	"testing"

	"git.gammaspectra.live/git/go-away/lib/settings"
)

func TestChallengeMinimalTemplate(t *testing.T) {
	tpl := templates["challenge-minimal.gohtml"]
	if tpl == nil {
		t.Fatal("challenge-minimal.gohtml template not registered")
	}

	t.Run("default loading state", func(t *testing.T) {
		var buf bytes.Buffer
		err := tpl.Execute(&buf, map[string]any{
			"Id":        "req-123",
			"Random":    "abc",
			"Path":      "/.well-known/go-away",
			"Title":     settings.DefaultStrings.Get("title_challenge"),
			"Challenge": "js-pow-sha256",
			"Strings":   settings.DefaultStrings,
			"EndTags": []template.HTML{
				`<script async type="module" src="/.well-known/go-away/js-pow-sha256/script.mjs"></script>`,
			},
		})
		if err != nil {
			t.Fatalf("Execute failed: %v", err)
		}
		out := buf.String()
		if !strings.Contains(out, `<p id="status">Loading...</p>`) {
			t.Errorf("expected Loading... status, got:\n%s", out)
		}
		if !strings.Contains(out, `id="title"`) {
			t.Errorf("expected hidden #title element for script.mjs compatibility")
		}
		if strings.Contains(out, "style.css") || strings.Contains(out, "logo.png") || strings.Contains(out, "Why am I seeing this?") {
			t.Errorf("minimal template should not include anubis assets or details, got:\n%s", out)
		}
		if !strings.Contains(out, "<noscript>") {
			t.Errorf("expected noscript warning when EndTags is present")
		}
	})

	t.Run("error state", func(t *testing.T) {
		var buf bytes.Buffer
		err := tpl.Execute(&buf, map[string]any{
			"Id":       "req-456",
			"Random":   "abc",
			"Path":     "/.well-known/go-away",
			"Title":    "Oh no! Bad Request",
			"Error":    "challenge expired or invalid",
			"Redirect": "/",
			"Strings":  settings.DefaultStrings,
		})
		if err != nil {
			t.Fatalf("Execute failed: %v", err)
		}
		out := buf.String()
		if !strings.Contains(out, `<p id="status">Error: challenge expired or invalid</p>`) {
			t.Errorf("expected error message, got:\n%s", out)
		}
		if !strings.Contains(out, `<a href="/">Refresh page</a>`) {
			t.Errorf("expected refresh link, got:\n%s", out)
		}
	})

	t.Run("captcha state", func(t *testing.T) {
		var buf bytes.Buffer
		err := tpl.Execute(&buf, map[string]any{
			"Id":         "req-789",
			"Random":     "abc",
			"Path":       "/.well-known/go-away",
			"Title":      settings.DefaultStrings.Get("title_captcha"),
			"CaptchaSVG": template.HTML(`<svg></svg>`),
			"VerifyURL":  "/verify",
			"Error":      "Incorrect code. Please try again.",
			"Strings":    settings.DefaultStrings,
		})
		if err != nil {
			t.Fatalf("Execute failed: %v", err)
		}
		out := buf.String()
		if !strings.Contains(out, `<svg></svg>`) || !strings.Contains(out, `action="/verify"`) {
			t.Errorf("expected captcha svg and verify form, got:\n%s", out)
		}
		if !strings.Contains(out, "Incorrect code. Please try again.") {
			t.Errorf("expected captcha error message, got:\n%s", out)
		}
	})
}
