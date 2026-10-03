package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadDefaultUsesConfiguredPath(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.env")
	contents := "SERVERPORT=:9393\nDATABASEUSER=test-user\nDATABASEPASS=test-pass\nDBCONNECTION=localhost:5432\nDATABASE=billing\nNETPROTOCOL=tcp\n"
	if err := os.WriteFile(path, []byte(contents), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("BILLINGGO_CONFIG_FILE", path)

	got, err := LoadDefault()
	if err != nil {
		t.Fatal(err)
	}
	if got.ServerPort != ":9393" {
		t.Fatalf("ServerPort = %q, want %q", got.ServerPort, ":9393")
	}
}

func TestLoadRejectsMalformedServerPort(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.env")
	contents := "SERVERPORT=not-a-listen-address\nDATABASEUSER=test-user\nDATABASEPASS=test-pass\nDBCONNECTION=localhost:5432\nDATABASE=billing\nNETPROTOCOL=tcp\n"
	if err := os.WriteFile(path, []byte(contents), 0o600); err != nil {
		t.Fatal(err)
	}

	if _, err := Load(path); err == nil {
		t.Fatal("Load() succeeded with malformed SERVERPORT")
	}
}

func TestLoadReadsInvoiceConfiguration(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.env")
	contents := "SERVERPORT=:8383\nDATABASEUSER=test-user\nDATABASEPASS=test-pass\nDBCONNECTION=localhost:5432\nDATABASE=billing\nNETPROTOCOL=tcp\n"
	if err := os.WriteFile(path, []byte(contents), 0o600); err != nil {
		t.Fatal(err)
	}

	got, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}

	if got.ServerPort != ":8383" {
		t.Fatalf("ServerPort = %q, want %q", got.ServerPort, ":8383")
	}
	if got.DatabaseUser != "test-user" {
		t.Fatalf("DatabaseUser = %q, want %q", got.DatabaseUser, "test-user")
	}
	if got.DatabasePassword != "test-pass" {
		t.Fatalf("DatabasePassword = %q, want %q", got.DatabasePassword, "test-pass")
	}
	if got.DatabaseHost != "localhost:5432" {
		t.Fatalf("DatabaseHost = %q, want %q", got.DatabaseHost, "localhost:5432")
	}
	if got.DatabaseName != "billing" {
		t.Fatalf("DatabaseName = %q, want %q", got.DatabaseName, "billing")
	}
	if got.NetworkProtocol != "tcp" {
		t.Fatalf("NetworkProtocol = %q, want %q", got.NetworkProtocol, "tcp")
	}
}
