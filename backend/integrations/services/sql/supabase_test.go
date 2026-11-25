package sql

import (
	"net/url"
	"testing"

	"ylem_integrations/entities"
)

func TestCreateSQLIntegrationConnectionSupabase(t *testing.T) {
	conn, err := CreateSQLIntegrationConnection(entities.SQLIntegrationTypeSupabase, DefaultSQLIntegrationConnectionConfiguration{
		Host:       "db.example.supabase.co",
		Port:       5432,
		User:       "postgres",
		Password:   "secret",
		Database:   "postgres",
		SslEnabled: true,
	})
	if err != nil {
		t.Fatalf("create connection: %v", err)
	}

	pg, ok := conn.(*PostgreSqlIntegrationConnection)
	if !ok {
		t.Fatalf("connection type %T, want *PostgreSqlIntegrationConnection", conn)
	}
	if !pg.SSLEnabled {
		t.Fatal("ssl should stay enabled when the caller requests it")
	}
	if len(pg.HiddenSchemas) == 0 {
		t.Fatal("supabase connection should hide internal schemas")
	}
	for _, schema := range pg.HiddenSchemas {
		if schema == "public" {
			t.Fatal("public schema must stay visible")
		}
	}
}

func TestBuildPostgresDSNEncodesPassword(t *testing.T) {
	raw := buildPostgresDSN("postgres", "p@ss:word/?", "db.example.supabase.co", 5432, "postgres", "require", nil)

	parsed, err := url.Parse(raw)
	if err != nil {
		t.Fatalf("parse dsn: %v", err)
	}
	if parsed.Hostname() != "db.example.supabase.co" {
		t.Fatalf("host %q", parsed.Hostname())
	}
	if parsed.Port() != "5432" {
		t.Fatalf("port %q", parsed.Port())
	}
	if parsed.Path != "/postgres" {
		t.Fatalf("path %q", parsed.Path)
	}
	password, ok := parsed.User.Password()
	if !ok || password != "p@ss:word/?" {
		t.Fatalf("password %q", password)
	}
	if parsed.Query().Get("sslmode") != "require" {
		t.Fatalf("sslmode %q", parsed.Query().Get("sslmode"))
	}
}

func TestIsSQLIntegrationTypeSupportedIncludesSupabase(t *testing.T) {
	if !entities.IsSQLIntegrationTypeSupported(entities.SQLIntegrationTypeSupabase) {
		t.Fatal("supabase should be a supported SQL integration type")
	}
	if entities.SQLIntegrationTypeSupabase != "supabase" {
		t.Fatalf("type %q", entities.SQLIntegrationTypeSupabase)
	}
}
