// SPDX-License-Identifier: AGPL-3.0-only
// Copyright (c) 2026 Arturo (AutoHost)

package node

import (
	"autohost-cli/internal/platform/config"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

// UpService handles the interactive node enrollment flow (Tailscale-style).
type UpService struct{}

type upCallbackResult struct {
	token  string
	apiURL string
	err    error
}

func (s *UpService) Up(cloudURL, apiURL, name string) error {
	cfg := config.Load()
	if cloudURL == "" {
		cloudURL = cfg.CloudURL
	}
	if apiURL == "" {
		apiURL = cfg.APIURL
	}
	cloudURL = strings.TrimRight(cloudURL, "/")
	apiURL = strings.TrimRight(apiURL, "/")

	nd := (&EnrollService{}).GatherNodeData()
	if name != "" {
		nd.HostName = name
	}

	stateBytes := make([]byte, 16)
	if _, err := rand.Read(stateBytes); err != nil {
		return fmt.Errorf("generating state nonce: %w", err)
	}
	state := hex.EncodeToString(stateBytes)

	listener, err := net.Listen("tcp", "0.0.0.0:0")
	if err != nil {
		return fmt.Errorf("starting local callback server: %w", err)
	}
	localPort := listener.Addr().(*net.TCPAddr).Port

	callbackIP := DetectLocalIP()
	if callbackIP == "unknown" {
		callbackIP = "127.0.0.1"
	}
	callbackBase := fmt.Sprintf("http://%s:%d", callbackIP, localPort)

	resultCh := make(chan upCallbackResult, 1)
	var once sync.Once

	mux := http.NewServeMux()
	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		receivedState := q.Get("state")
		token := q.Get("token")
		apiURL := q.Get("api_url")

		if receivedState != state {
			http.Error(w, "state mismatch", http.StatusBadRequest)
			once.Do(func() {
				resultCh <- upCallbackResult{err: fmt.Errorf("state mismatch: possible CSRF attack")}
			})
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, successHTML)
		once.Do(func() {
			resultCh <- upCallbackResult{token: token, apiURL: apiURL}
		})
	})

	srv := &http.Server{
		Handler:     mux,
		ReadTimeout: 10 * time.Second,
	}
	go func() { _ = srv.Serve(listener) }()
	defer srv.Close()

	authURL := fmt.Sprintf("%s/cli-auth?state=%s&callback=%s&node=%s&api_url=%s",
		cloudURL,
		url.QueryEscape(state),
		url.QueryEscape(callbackBase),
		url.QueryEscape(nd.HostName),
		url.QueryEscape(apiURL),
	)

	fmt.Println()
	fmt.Println("🌐 Opening the browser to authorize this node...")
	fmt.Printf("   Node: %s\n", nd.HostName)
	fmt.Println()
	fmt.Println("   URL:", authURL)
	fmt.Println()
	fmt.Println("💡 If the browser does not open automatically, copy and paste the URL.")
	fmt.Println()

	fmt.Println("⏳ Waiting for authorization in the browser (5 min)...")

	select {
	case res := <-resultCh:
		if res.err != nil {
			return fmt.Errorf("authorization failed: %w", res.err)
		}
		fmt.Println()
		fmt.Println("✅ Authorization received. Registering node...")
		finalAPI := res.apiURL
		if finalAPI == "" {
			finalAPI = apiURL
		}
		return (&EnrollService{}).Link(finalAPI, res.token, nd.HostName)
	case <-time.After(5 * time.Minute):
		return fmt.Errorf("timed out after 5 min. Run `autohost up` again")
	}
}

const successHTML = `<!DOCTYPE html>
<html lang="en">
<head><meta charset="UTF-8"><title>AutoHost CLI — Authorized</title>
<style>
  body{font-family:system-ui,sans-serif;display:flex;align-items:center;justify-content:center;min-height:100vh;margin:0;background:#0a0a0a;color:#e5e5e5}
  .card{text-align:center;padding:2.5rem 3rem;border:1px solid #27272a;border-radius:1rem;background:#18181b;max-width:440px}
  h2{font-size:1.5rem;margin:0 0 .75rem}
  p{color:#a1a1aa;margin:0}
</style>
</head>
<body>
  <div class="card">
    <h2>✅ Node authorized</h2>
    <p>You can close this window and return to the terminal.</p>
  </div>
</body>
</html>`

