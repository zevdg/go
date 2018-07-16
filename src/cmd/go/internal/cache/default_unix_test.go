// +build !windows,!darwin,!plan9

package cache

import (
	"strings"
	"testing"
)

func TestDefaultDir(t *testing.T) {
	goCacheDir := "/tmp/test-go-cache"
	xdgCacheDir := "/tmp/test-xdg-cache"
	homeDir := "/tmp/test-home"

	defer func() {
		env = nil
	}()

	env = envShim{
		"GOCACHE":        goCacheDir,
		"XDG_CACHE_HOME": xdgCacheDir,
		"HOME":           homeDir,
	}

	dir, showWarnings := defaultDir()
	if dir != goCacheDir {
		t.Errorf("Cache DefaultDir %q should be $GOCACHE %q", dir, goCacheDir)
	}
	if !showWarnings {
		t.Error("Warnings should be shown when $GOCACHE is set")
	}

	delete(env, "GOCACHE")
	dir, showWarnings = defaultDir()
	if !strings.HasPrefix(dir, xdgCacheDir+"/") {
		t.Errorf("Cache DefaultDir %q should be under $XDG_CACHE_HOME %q when $GOCACHE is unset", dir, xdgCacheDir)
	}
	if !showWarnings {
		t.Error("Warnings should be shown when $XDG_CACHE_HOME is set")
	}

	delete(env, "XDG_CACHE_HOME")
	dir, showWarnings = defaultDir()
	if !strings.HasPrefix(dir, homeDir+"/.cache/") {
		t.Errorf("Cache DefaultDir %q should be under $HOME/.cache %q when $GOCACHE and $XDG_CACHE_HOME are unset", dir, homeDir+"/.cache")
	}
	if !showWarnings {
		t.Error("Warnings should be shown when $HOME is not /")
	}

	delete(env, "HOME")
	if dir, _ := defaultDir(); dir != "off" {
		t.Error("Cache not disabled when $GOCACHE, $XDG_CACHE_HOME, and $HOME are unset")
	}

	env["HOME"] = "/"
	if _, showWarnings := defaultDir(); showWarnings {
		// https://golang.org/issue/26280
		t.Error("Cache initalization warnings should be squelched when $GOCACHE and $XDG_CACHE_HOME are unset and $HOME is /")
	}
}
