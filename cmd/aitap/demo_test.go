package main

import "testing"

func TestDemoCallsCoverProvidersAndContent(t *testing.T) {
	calls := demoCalls()

	if len(calls) < 5 {
		t.Fatalf("expected at least 5 demo calls, got %d", len(calls))
	}

	foundOpenRouter := false
	foundStreaming := false

	for _, call := range calls {
		if call.Provider == "" {
			t.Fatal("expected demo call provider to be set")
		}
		if call.Model == "" {
			t.Fatal("expected demo call model to be set")
		}
		if call.StatusCode == 0 {
			t.Fatal("expected demo call status code to be set")
		}
		if call.Provider == "openrouter" {
			foundOpenRouter = true
		}
		if call.IsStreaming {
			foundStreaming = true
		}
	}

	if !foundOpenRouter {
		t.Fatal("expected demo feed to include an openrouter sample")
	}
	if !foundStreaming {
		t.Fatal("expected demo feed to include streaming samples")
	}
}
