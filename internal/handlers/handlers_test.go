package handlers

import (
	"testing"
)

func TestParseIncompletePath(t *testing.T) {
	path := `/myorg/myproj`
	_, err := extractRepoPath(path)
	if err != errIncompletePath {
		t.Errorf("expecting incomplete path")
	}
}

func TestParseEmptyPath(t *testing.T) {
	path := `/`
	_, err := extractRepoPath(path)
	if err != errIncompletePath {
		t.Errorf("expecting incomplete path")
	}
}

func TestParsePath(t *testing.T) {
	tcs := []struct {
		path  string
		parts parts
	}{
		{"/myorg/myproj/myrepo", parts{org: "myorg", proj: "myproj", repo: "myrepo"}},
		{"/myorg/myproj/myrepo/sub1", parts{org: "myorg", proj: "myproj", repo: "myrepo", subs: "sub1"}},
		{"/myorg/myproj/myrepo/sub1/", parts{org: "myorg", proj: "myproj", repo: "myrepo", subs: "sub1/"}},
		{"/myorg/myproj/myrepo.git", parts{org: "myorg", proj: "myproj", repo: "myrepo.git"}},
		{"/myorg/myproj/myrepo.git/sub1", parts{org: "myorg", proj: "myproj", repo: "myrepo.git", subs: "sub1"}},
		{"/myorg/myproj/myrepo/sub1/sub2", parts{org: "myorg", proj: "myproj", repo: "myrepo", subs: "sub1/sub2"}},
	}
	for _, tc := range tcs {
		t.Run(tc.path, func(t *testing.T) {
			ps, err := extractRepoPath(tc.path)
			if err != nil {
				t.Errorf("error splitting path %s", err)
			}
			if ps != tc.parts {
				t.Errorf("got %+v, want %+v,", ps, tc.parts)
			}
		})
	}
}
