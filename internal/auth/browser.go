// ============================================================
// Browser Authentication - OAuth Browser Redirect Flow
// ============================================================
//
// Flow:
// 1. Start local HTTP server on random port
// 2. Open browser to console.machpay.xyz/auth/cli?port=PORT
// 3. User logs in normally (Google, Wallet, Email)
// 4. Console redirects to localhost:PORT/callback?token=JWT
// 5. CLI receives token, saves to config
//
// ============================================================

package auth

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"
)

const (
	// CallbackTimeout is how long to wait for the user to complete login
	CallbackTimeout = 5 * time.Minute

	// SuccessHTML is shown in the browser after successful login
	successHTML = `<!DOCTYPE html>
<html>
<head>
    <title>MachPay CLI - Login Success</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body { 
            font-family: system-ui, -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
            display: flex;
            justify-content: center;
            align-items: center;
            height: 100vh;
            background: linear-gradient(135deg, #0a0a0a 0%, #1a1a1a 100%);
            color: #fff;
        }
        .container { 
            text-align: center;
            padding: 2rem;
        }
        .check { 
            font-size: 72px; 
            margin-bottom: 24px;
            animation: pop 0.5s ease-out;
        }
        @keyframes pop {
            0% { transform: scale(0); opacity: 0; }
            50% { transform: scale(1.2); }
            100% { transform: scale(1); opacity: 1; }
        }
        h1 { 
            font-size: 24px;
            font-weight: 600;
            color: #10b981;
            margin-bottom: 8px;
        }
        p { 
            color: #71717a;
            font-size: 14px;
        }
        .hint {
            margin-top: 24px;
            padding: 12px 16px;
            background: rgba(16, 185, 129, 0.1);
            border: 1px solid rgba(16, 185, 129, 0.2);
            border-radius: 8px;
            font-size: 13px;
            color: #10b981;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="check">✅</div>
        <h1>Login Successful!</h1>
        <p>You can close this window and return to your terminal.</p>
        <div class="hint">
            Your CLI is now authenticated with MachPay
        </div>
    </div>
</body>
</html>`
)

// CallbackResult contains the result of the browser callback
type CallbackResult struct {
	Token string
	Error error
}

// StartCallbackServer starts a local HTTP server to receive the auth callback
func StartCallbackServer(port int, resultChan chan<- CallbackResult) *http.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token")
		errMsg := r.URL.Query().Get("error")

		if errMsg != "" {
			resultChan <- CallbackResult{Error: fmt.Errorf("auth error: %s", errMsg)}
			http.Error(w, errMsg, http.StatusBadRequest)
			return
		}

		if token == "" {
			resultChan <- CallbackResult{Error: fmt.Errorf("no token received")}
			http.Error(w, "No token", http.StatusBadRequest)
			return
		}

		// Send success page to browser
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(successHTML))

		// Send token to caller
		resultChan <- CallbackResult{Token: token}
	})

	// Health check endpoint for debugging
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	server := &http.Server{
		Addr:    fmt.Sprintf("localhost:%d", port),
		Handler: mux,
	}

	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			resultChan <- CallbackResult{Error: fmt.Errorf("callback server error: %w", err)}
		}
	}()

	return server
}

// FindFreePort finds an available TCP port on localhost
func FindFreePort() (int, error) {
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return 0, fmt.Errorf("find free port: %w", err)
	}
	defer listener.Close()
	return listener.Addr().(*net.TCPAddr).Port, nil
}

// ShutdownServer gracefully shuts down the callback server
func ShutdownServer(server *http.Server) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return server.Shutdown(ctx)
}

