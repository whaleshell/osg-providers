package provider

import (
	"testing"

	"github.com/whaleshell/whaleshell-core/policy"
)

func TestParseGitHubProfile(t *testing.T) {
	const yaml = `
id: github
category: source_control
binaries: [/usr/bin/gh, /usr/bin/git]
credentials:
  - name: api_token
    env_vars: [GITHUB_TOKEN, GH_TOKEN]
    required: true
endpoints:
  - id: api
    host: api.github.com
    port: 443
    protocol: rest
    tls: terminate
    access: read-only
  - id: git
    host: github.com
    port: 443
    protocol: rest
    tls: terminate
    rules:
      - allow: { method: GET, path: "**" }
      - allow: { method: POST, path: "/**/git-upload-pack" }
`
	p, err := ParseYAML([]byte(yaml))
	if err != nil {
		t.Fatal(err)
	}
	if p.ID != "github" || len(p.Endpoints) != 2 {
		t.Fatalf("%+v", p)
	}
	keys := p.EnvKeys()
	if len(keys) != 2 || keys[0] != "GITHUB_TOKEN" || keys[1] != "GH_TOKEN" {
		t.Fatalf("keys=%v", keys)
	}
}

func TestDiscoverEnvVars(t *testing.T) {
	p := Profile{
		ID: "github",
		Credentials: []Credential{{
			Name: "api_token", EnvVars: []string{"GITHUB_TOKEN", "GH_TOKEN"}, Required: true,
		}},
	}
	t.Setenv("GITHUB_TOKEN", "")
	t.Setenv("GH_TOKEN", "tok")
	keys, err := p.DiscoverEnvVars()
	if err != nil {
		t.Fatal(err)
	}
	if len(keys) != 1 || keys[0] != "GH_TOKEN" {
		t.Fatalf("%v", keys)
	}
	t.Setenv("GH_TOKEN", "")
	if _, err := p.DiscoverEnvVars(); err == nil {
		t.Fatal("expected missing credential error")
	}
}

func TestEffectivePolicy(t *testing.T) {
	base := policy.Document{Version: 1}
	p := Profile{
		ID: "nvidia",
		Endpoints: []policy.AllowRule{{
			Host: "integrate.api.nvidia.com", Port: 443,
			Protocol: "rest", TLS: "terminate", Access: "read-write",
		}},
		Binaries: []string{"/usr/bin/curl"},
		Credentials: []Credential{{
			Name: "api_key", EnvVars: []string{"NVIDIA_API_KEY"},
		}},
	}
	if err := p.Validate(); err != nil {
		t.Fatal(err)
	}
	out := EffectivePolicy(base, []Layer{{InstanceName: "nv", Profile: p}}, false)
	allows := out.NetworkAllows()
	if len(allows) != 1 {
		t.Fatalf("allow=%d", len(allows))
	}
	if allows[0].ID != "provider.nv.0" {
		t.Fatalf("id=%s", allows[0].ID)
	}
	if len(allows[0].Binaries) != 1 {
		t.Fatalf("binaries=%v", allows[0].Binaries)
	}
	if out.Credentials == nil || len(out.Credentials.EnvAllow) != 1 || out.Credentials.EnvAllow[0] != "NVIDIA_API_KEY" {
		t.Fatalf("env=%v", out.Credentials)
	}
	if len(base.NetworkAllows()) != 0 {
		t.Fatalf("base mutated: %d", len(base.NetworkAllows()))
	}
	fresh := policy.Document{Version: 1}
	suppressed := EffectivePolicy(fresh, []Layer{{InstanceName: "nv", Profile: p}}, true)
	if len(suppressed.NetworkAllows()) != 0 {
		t.Fatalf("expected suppress, got %d", len(suppressed.NetworkAllows()))
	}
}

func TestOmitProviderComposed(t *testing.T) {
	doc := policy.Document{Version: 1}
	doc.SetNetworkAllows([]policy.AllowRule{
		{ID: "github-push", Host: "github.com", Port: 443},
		{ID: "provider.gh.api", Host: "api.github.com", Port: 443},
		{ID: "provider.gh.git", Host: "github.com", Port: 443},
	})
	out, n := OmitProviderComposed(doc)
	if n != 2 {
		t.Fatalf("stripped=%d", n)
	}
	allows := out.NetworkAllows()
	if len(allows) != 1 || allows[0].ID != "github-push" {
		t.Fatalf("%+v", allows)
	}
	if len(doc.NetworkAllows()) != 3 {
		t.Fatalf("input mutated")
	}
	p := Profile{
		ID: "github",
		Endpoints: []policy.AllowRule{
			{ID: "api", Host: "api.github.com", Port: 443, Protocol: "rest", TLS: "terminate", Access: "read-only"},
		},
		Credentials: []Credential{{Name: "t", EnvVars: []string{"GITHUB_TOKEN"}}},
	}
	eff := EffectivePolicy(out, []Layer{{InstanceName: "gh", Profile: p}}, false)
	if len(eff.NetworkAllows()) != 2 {
		t.Fatalf("effective allow=%d %+v", len(eff.NetworkAllows()), eff.NetworkAllows())
	}
}

func TestAuditMatchHTTP(t *testing.T) {
	rule := policy.AllowRule{
		Host: "api.example.com", Port: 443, Protocol: "rest", TLS: "terminate",
		Access: "read-only", Enforcement: "audit",
	}
	ok, _ := rule.MatchHTTP("POST", "/v1")
	if ok {
		t.Fatal("POST should fail MatchHTTP on read-only")
	}
	if !rule.IsAudit() {
		t.Fatal("expected audit")
	}
}
