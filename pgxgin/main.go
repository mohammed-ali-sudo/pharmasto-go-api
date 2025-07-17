package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	// Your full Supabase DSN
	DBURL = "postgresql://postgres:1010204080@" +
		"db.jqysfttmqmlkjlishljc.supabase.co:5432/postgres?sslmode=require"
	SUPABASE_HOST = "db.jqysfttmqmlkjlishljc.supabase.co"
	SUPABASE_PORT = "5432"
	DNS_SERVER    = "8.8.8.8:53" // Public Google DNS
)

func main() {
	// Parse the DSN
	cfg, err := pgxpool.ParseConfig(DBURL)
	if err != nil {
		log.Fatalf("ParseConfig failed: %v", err)
	}

	// Create a Go-resolver that uses Google DNS
	resolver := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{Timeout: time.Second}
			return d.DialContext(ctx, "udp", DNS_SERVER)
		},
	}

	// Override the DialFunc to do IPv4-only lookups via our resolver
	cfg.ConnConfig.DialFunc = func(ctx context.Context, _, _ string) (net.Conn, error) {
		// Lookup A records using our custom resolver
		ips, err := resolver.LookupIP(ctx, "ip4", SUPABASE_HOST)
		if err != nil {
			return nil, fmt.Errorf("IPv4 lookup via %s failed: %w", DNS_SERVER, err)
		}
		var lastErr error
		for _, ip := range ips {
			addr := net.JoinHostPort(ip.String(), SUPABASE_PORT)
			conn, err := (&net.Dialer{}).DialContext(ctx, "tcp4", addr)
			if err == nil {
				return conn, nil
			}
			lastErr = err
		}
		return nil, fmt.Errorf("could not dial any IPv4 for %q: %w", SUPABASE_HOST, lastErr)
	}

	// Pool settings
	cfg.MaxConns = 5
	cfg.MaxConnIdleTime = 30 * time.Second

	// Connect
	pool, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		log.Fatalf("Unable to connect: %v", err)
	}
	defer pool.Close()

	// Ping
	if err := pool.Ping(context.Background()); err != nil {
		log.Fatalf("Ping failed: %v", err)
	}
	log.Println("✅ Connected to Supabase Postgres over IPv4 with custom DNS!")
}
