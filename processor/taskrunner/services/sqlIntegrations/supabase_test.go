package sqlIntegrations

import (
	"net/url"
	"testing"
)

func TestCreateSQLIntegrationConnectionSupabase(t *testing.T) {
	conn, err := CreateSQLIntegrationConnection(sqlIntegrationTypeSupabase, DefaultSQLIntegrationConnectionConfiguration{
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

	pg, ok := conn.(*PostgreSqlSQLIntegrationConnection)
	if !ok {
		t.Fatalf("connection type %T, want *PostgreSqlSQLIntegrationConnection", conn)
	}
	if pg.Host != "db.example.supabase.co" || pg.Port != 5432 || !pg.SSLEnabled {
		t.Fatalf("connection %+v", pg)
	}
}

func TestBuildPostgresDSNEncodesPassword(t *testing.T) {
	raw := buildPostgresDSN("postgres", "p@ss:word/?", "db.example.supabase.co", 5432, "postgres", "require", map[string]string{
		"sshHost": "bastion.internal",
	})

	parsed, err := url.Parse(raw)
	if err != nil {
		t.Fatalf("parse dsn: %v", err)
	}
	password, ok := parsed.User.Password()
	if !ok || password != "p@ss:word/?" {
		t.Fatalf("password %q", password)
	}
	if parsed.Query().Get("sslmode") != "require" {
		t.Fatalf("sslmode %q", parsed.Query().Get("sslmode"))
	}
	if parsed.Query().Get("sshHost") != "bastion.internal" {
		t.Fatalf("sshHost %q", parsed.Query().Get("sshHost"))
	}
}
