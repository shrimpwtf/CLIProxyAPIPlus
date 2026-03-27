package util

import (
	"testing"

	"github.com/router-for-me/CLIProxyAPI/v6/internal/registry"
)

func TestGetProviderName_PrefersNativeStaticProvidersOverGenericAliases(t *testing.T) {
	reg := registry.GetGlobalRegistry()
	reg.RegisterClient("provider-test-acme", "acme", []*registry.ModelInfo{{ID: "gpt-5-codex"}})
	reg.RegisterClient("provider-test-codex", "codex", []*registry.ModelInfo{{ID: "gpt-5-codex"}})
	t.Cleanup(func() {
		reg.UnregisterClient("provider-test-acme")
		reg.UnregisterClient("provider-test-codex")
	})

	got := GetProviderName("gpt-5-codex")
	want := []string{"codex", "acme"}
	if len(got) != len(want) {
		t.Fatalf("provider count = %d, want %d (%v)", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("providers[%d] = %q, want %q (all=%v)", i, got[i], want[i], got)
		}
	}
}

func TestGetProviderName_PrefersAllNativeProvidersBeforeGenericAliases(t *testing.T) {
	reg := registry.GetGlobalRegistry()
	reg.RegisterClient("provider-test-acme-2", "acme", []*registry.ModelInfo{{ID: "gpt-5-codex"}})
	reg.RegisterClient("provider-test-codex-2", "codex", []*registry.ModelInfo{{ID: "gpt-5-codex"}})
	reg.RegisterClient("provider-test-ghcp-2", "github-copilot", []*registry.ModelInfo{{ID: "gpt-5-codex"}})
	t.Cleanup(func() {
		reg.UnregisterClient("provider-test-acme-2")
		reg.UnregisterClient("provider-test-codex-2")
		reg.UnregisterClient("provider-test-ghcp-2")
	})

	got := GetProviderName("gpt-5-codex")
	want := []string{"codex", "github-copilot", "acme"}
	if len(got) != len(want) {
		t.Fatalf("provider count = %d, want %d (%v)", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("providers[%d] = %q, want %q (all=%v)", i, got[i], want[i], got)
		}
	}
}

func TestGetProviderName_PreservesOriginalOrderWhenNoNativeProviderMatches(t *testing.T) {
	reg := registry.GetGlobalRegistry()
	reg.RegisterClient("provider-test-acme-3", "acme", []*registry.ModelInfo{{ID: "custom-model"}})
	reg.RegisterClient("provider-test-zeta-3", "zeta", []*registry.ModelInfo{{ID: "custom-model"}})
	t.Cleanup(func() {
		reg.UnregisterClient("provider-test-acme-3")
		reg.UnregisterClient("provider-test-zeta-3")
	})

	got := GetProviderName("custom-model")
	want := []string{"acme", "zeta"}
	if len(got) != len(want) {
		t.Fatalf("provider count = %d, want %d (%v)", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("providers[%d] = %q, want %q (all=%v)", i, got[i], want[i], got)
		}
	}
}
