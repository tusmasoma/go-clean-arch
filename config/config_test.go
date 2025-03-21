package config

import (
	"context"
	"testing"
	"time"

	"github.com/sethvargo/go-envconfig"
	"github.com/stretchr/testify/require"
)

func Test_NewDBConfig(t *testing.T) {
	ctx := context.Background()

	patterns := []struct {
		name  string
		setup func(t *testing.T)
		want  *DBConfig
		err   error
	}{
		{
			name: "default",
			setup: func(t *testing.T) {
				t.Helper()
			},
			want: nil,
			err:  envconfig.ErrMissingRequired,
		},
		{
			name: "set env",
			setup: func(t *testing.T) {
				t.Helper()
				t.Setenv("MYSQL_USER", "root")
				t.Setenv("MYSQL_PASSWORD", "campfinder")
				t.Setenv("MYSQL_HOST", "mysql")
				t.Setenv("MYSQL_PORT", "3306")
				t.Setenv("MYSQL_DB_NAME", "campfinderdb")
			},
			want: &DBConfig{
				Host:     "mysql",
				Port:     "3306",
				User:     "root",
				Password: "campfinder",
				DBName:   "campfinderdb",
			},
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(t)

			got, err := NewDBConfig(ctx, "MYSQL_")
			if err != nil {
				require.ErrorIs(t, err, tt.err)
			}
			require.Equal(t, tt.want, got)
		})
	}
}

func Test_NewServerConfig(t *testing.T) {
	ctx := context.Background()

	patterns := []struct {
		name  string
		setup func(t *testing.T)
		want  *ServerConfig
		err   error
	}{
		{
			name: "default",
			setup: func(t *testing.T) {
				t.Helper()
			},
			want: &ServerConfig{
				ReadTimeout:               5 * time.Second,
				WriteTimeout:              10 * time.Second,
				IdleTimeout:               15 * time.Second,
				GracefulShutdownTimeout:   5 * time.Second,
				PreflightCacheDurationSec: 300,
			},
			err: nil,
		},
		{
			name: "set env",
			setup: func(t *testing.T) {
				t.Helper()
				t.Setenv("SERVER_READ_TIMEOUT", "2s")
				t.Setenv("SERVER_WRITE_TIMEOUT", "4s")
				t.Setenv("SERVER_IDLE_TIMEOUT", "10s")
				t.Setenv("SERVER_GRACEFUL_SHUTDOWN_TIMEOUT", "3s")
				t.Setenv("SERVER_PREFLIGHT_CACHE_DURATION_SEC", "150")
			},
			want: &ServerConfig{
				ReadTimeout:               2 * time.Second,
				WriteTimeout:              4 * time.Second,
				IdleTimeout:               10 * time.Second,
				GracefulShutdownTimeout:   3 * time.Second,
				PreflightCacheDurationSec: 150,
			},
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(t)

			got, err := NewServerConfig(ctx)
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
