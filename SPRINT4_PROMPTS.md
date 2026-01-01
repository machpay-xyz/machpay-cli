# Sprint 4: Desktop App - Google OAuth Fix

**Sprint Goal:** Fix Google OAuth in Tauri desktop app using deep linking  
**Estimated Time:** 3-4 hours  
**Critical:** This has failed multiple times. Follow each step carefully and TEST before moving to the next.

---

## 🔴 Why Previous Attempts Failed

| Attempt | Problem |
|---------|---------|
| Clipboard relay | Browser blocks automatic clipboard writes for security |
| LocalStorage relay | Different origins (`console.machpay.xyz` vs `tauri://localhost`) don't share localStorage |
| Backend polling | Browser (HTTPS) cannot call localhost (HTTP) - Mixed Content blocked |
| Local HTTP server | Complex setup, Google doesn't allow localhost redirects easily |

## ✅ Why Deep Linking Works

**The Solution:** Custom URL protocol (`machpay://`)

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           OAUTH FLOW WITH DEEP LINKING                      │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  1. Desktop App                    2. System Browser                        │
│  ┌──────────────┐                  ┌──────────────────────────┐            │
│  │ Click Google │ ──opens browser──▶│ accounts.google.com     │            │
│  │    Login     │                  │ User completes OAuth     │            │
│  └──────────────┘                  └──────────┬───────────────┘            │
│                                               │                             │
│                                               ▼ redirect                    │
│                                    ┌──────────────────────────┐            │
│                                    │ console.machpay.xyz      │            │
│  3. Web Callback                   │ /auth/desktop-callback   │            │
│     (SERVER-SIDE)                  │                          │            │
│  ┌──────────────────────────────┐  │ • Receives ?code=xxx     │            │
│  │ Backend exchanges code       │◀─│ • Calls backend API      │            │
│  │ for token (no CORS issues)   │  │ • Gets token back        │            │
│  └──────────────────────────────┘  └──────────┬───────────────┘            │
│                                               │                             │
│                                               ▼ redirect to                 │
│                                    ┌──────────────────────────┐            │
│                                    │ machpay://auth/callback  │            │
│  4. Desktop App Receives           │ ?token=xxx               │            │
│  ┌──────────────┐                  │ &refresh_token=yyy       │            │
│  │ Tauri handles│◀─────────────────│                          │            │
│  │ machpay://   │                  └──────────────────────────┘            │
│  │ Logs in user │                                                          │
│  └──────────────┘                                                          │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

**Key Insight:** The token exchange happens SERVER-SIDE (backend → Google API), not in the browser. This bypasses ALL cross-origin restrictions.

---

## Prompts

### Prompt 4.1: Register Custom URL Scheme in Tauri

**Objective:** Make desktop app handle `machpay://` URLs

**Why First:** Without this, nothing else will work. Must verify deep links work before proceeding.

```
Register the machpay:// custom URL scheme in the Tauri desktop app.

Location: machpay-console/src-tauri/

Steps:

1. Update tauri.conf.json to register the URL scheme:
   
   Add to tauri.conf.json under "tauri":
   {
     "tauri": {
       "bundle": {
         "identifier": "xyz.machpay.console",
         "macOS": {
           "CFBundleURLTypes": [
             {
               "CFBundleURLSchemes": ["machpay"]
             }
           ]
         }
       }
     }
   }

2. Install deep-link plugin:
   cd src-tauri
   cargo add tauri-plugin-deep-link

3. Update main.rs to handle deep links:

   use tauri_plugin_deep_link;

   fn main() {
       tauri::Builder::default()
           .plugin(tauri_plugin_deep_link::init())
           .setup(|app| {
               // Register machpay:// scheme
               #[cfg(any(target_os = "macos", target_os = "windows"))]
               tauri_plugin_deep_link::register("machpay", |request| {
                   // This will be called when machpay:// URL is opened
                   dbg!(&request);
               })?;
               Ok(())
           })
           .run(tauri::generate_context!())
           .expect("error while running tauri application");
   }

4. Test deep linking manually:
   - Build the app: npm run tauri build
   - Install the app
   - Open terminal: open "machpay://test?foo=bar"
   - Check if app opens and receives the URL

DO NOT PROCEED until you can verify the app receives machpay:// URLs.
```

**Verification:**
```bash
# After building and installing:
open "machpay://test?hello=world"
# App should launch and receive the URL
```

---

### Prompt 4.2: Emit Deep Link Events to Frontend

**Objective:** Pass received URLs from Rust to React frontend

**Why:** The Rust backend receives the URL, but React needs to handle the login.

```
Update the Tauri app to emit events when deep links are received.

Location: machpay-console/src-tauri/src/main.rs

1. Update the deep link handler to emit events:

   .setup(|app| {
       let handle = app.handle().clone();
       
       tauri_plugin_deep_link::register("machpay", move |request| {
           println!("[Deep Link] Received: {}", request);
           
           // Emit event to frontend
           if let Err(e) = handle.emit_all("deep-link", request.to_string()) {
               eprintln!("[Deep Link] Failed to emit: {}", e);
           }
       })?;
       
       Ok(())
   })

2. Create frontend listener in React (App.jsx or Auth.jsx):

   import { listen } from '@tauri-apps/api/event';

   useEffect(() => {
       // Only in Tauri
       if (!window.__TAURI__) return;

       const unlisten = listen('deep-link', (event) => {
           console.log('[Deep Link] Received:', event.payload);
           
           // Parse the URL
           const url = new URL(event.payload);
           
           if (url.pathname === '/auth/callback' || url.pathname === '//auth/callback') {
               const token = url.searchParams.get('token');
               const refreshToken = url.searchParams.get('refresh_token');
               
               if (token) {
                   console.log('[Deep Link] Got auth token!');
                   // Handle login with token
                   handleDeepLinkAuth(token, refreshToken);
               }
           }
       });

       return () => {
           unlisten.then(fn => fn());
       };
   }, []);

3. Test the event flow:
   - Run app in dev mode: npm run tauri dev
   - Open terminal: open "machpay://auth/callback?token=test123"
   - Check browser console for the log

DO NOT PROCEED until frontend receives the deep link event.
```

**Verification:**
```bash
# With app running in dev mode:
open "machpay://auth/callback?token=test123&refresh_token=test456"
# Console should show: [Deep Link] Received: machpay://auth/callback?token=test123...
```

---

### Prompt 4.3: Create Desktop OAuth Callback Endpoint (Backend)

**Objective:** Backend endpoint that exchanges code for token and returns redirect URL

**Why:** This is the SERVER-SIDE token exchange that avoids all CORS issues.

```
Create a new endpoint in the backend specifically for desktop OAuth callbacks.

Location: machpay-backend/internal/handlers/sql/auth_handler.go

1. Add new endpoint: POST /v1/auth/google/desktop-callback

   This endpoint:
   - Receives: { "code": "xxx", "redirect_uri": "https://console.machpay.xyz/auth/desktop-callback" }
   - Exchanges code with Google (server-side, no CORS)
   - Returns: { "redirect_url": "machpay://auth/callback?token=xxx&refresh_token=yyy" }

   func (h *AuthHandler) GoogleDesktopCallback(w http.ResponseWriter, r *http.Request) {
       var req struct {
           Code        string `json:"code"`
           RedirectURI string `json:"redirect_uri"`
       }
       
       if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
           http.Error(w, "Invalid request", http.StatusBadRequest)
           return
       }
       
       // Exchange code for tokens (same as existing Google OAuth)
       token, user, err := h.authService.ExchangeGoogleCode(r.Context(), req.Code, req.RedirectURI)
       if err != nil {
           http.Error(w, err.Error(), http.StatusUnauthorized)
           return
       }
       
       // Build deep link URL with tokens
       deepLink := fmt.Sprintf(
           "machpay://auth/callback?token=%s&refresh_token=%s&user_id=%s&email=%s",
           url.QueryEscape(token.AccessToken),
           url.QueryEscape(token.RefreshToken),
           url.QueryEscape(user.ID),
           url.QueryEscape(user.Email),
       )
       
       json.NewEncoder(w).Encode(map[string]string{
           "redirect_url": deepLink,
       })
   }

2. Register the route:
   router.HandleFunc("/v1/auth/google/desktop-callback", h.GoogleDesktopCallback).Methods("POST")

3. Test with curl:
   curl -X POST http://localhost:8081/v1/auth/google/desktop-callback \
     -H "Content-Type: application/json" \
     -d '{"code":"test","redirect_uri":"https://console.machpay.xyz/auth/desktop-callback"}'

   Should return: {"redirect_url":"machpay://auth/callback?token=..."}
```

**Verification:**
- Endpoint exists and responds
- Returns properly formatted deep link URL

---

### Prompt 4.4: Create Desktop Callback Page (Frontend)

**Objective:** Web page that receives OAuth code and triggers deep link redirect

**Why:** This page runs on `console.machpay.xyz`, receives the Google redirect, calls the backend, and redirects to `machpay://`

```
Create a new page specifically for desktop OAuth callbacks.

Location: machpay-console/src/pages/DesktopCallback.jsx

1. Create the new page:

   // This page is shown briefly in the browser after Google OAuth
   // It calls the backend to exchange code, then redirects to machpay://
   
   import { useEffect, useState } from 'react';
   import { useSearchParams } from 'react-router-dom';

   const API_URL = import.meta.env.VITE_API_URL || 'https://api.machpay.xyz';

   export default function DesktopCallback() {
       const [searchParams] = useSearchParams();
       const [status, setStatus] = useState('Processing...');
       const [error, setError] = useState(null);

       useEffect(() => {
           const code = searchParams.get('code');
           const state = searchParams.get('state');
           const errorParam = searchParams.get('error');

           if (errorParam) {
               setError(`OAuth error: ${errorParam}`);
               return;
           }

           if (!code) {
               setError('No authorization code received');
               return;
           }

           // Call backend to exchange code and get deep link URL
           async function exchangeCode() {
               try {
                   setStatus('Exchanging authorization code...');
                   
                   const response = await fetch(`${API_URL}/v1/auth/google/desktop-callback`, {
                       method: 'POST',
                       headers: { 'Content-Type': 'application/json' },
                       body: JSON.stringify({
                           code: code,
                           redirect_uri: 'https://console.machpay.xyz/auth/desktop-callback'
                       })
                   });

                   if (!response.ok) {
                       throw new Error('Failed to exchange code');
                   }

                   const { redirect_url } = await response.json();
                   
                   setStatus('Opening MachPay app...');
                   
                   // Redirect to deep link
                   window.location.href = redirect_url;
                   
                   // Show fallback after delay
                   setTimeout(() => {
                       setStatus('If the app did not open, click the button below:');
                   }, 2000);
                   
               } catch (err) {
                   console.error('[DesktopCallback] Error:', err);
                   setError(err.message);
               }
           }

           exchangeCode();
       }, [searchParams]);

       return (
           <div className="min-h-screen flex items-center justify-center bg-gray-900">
               <div className="text-center p-8 bg-gray-800 rounded-lg max-w-md">
                   <h1 className="text-2xl font-bold text-white mb-4">
                       {error ? '❌ Login Failed' : '🔐 Completing Login...'}
                   </h1>
                   
                   {error ? (
                       <p className="text-red-400 mb-4">{error}</p>
                   ) : (
                       <p className="text-gray-300 mb-4">{status}</p>
                   )}
                   
                   {/* Fallback button */}
                   <button 
                       onClick={() => window.location.href = 'machpay://'}
                       className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
                   >
                       Open MachPay App
                   </button>
                   
                   <p className="text-gray-500 text-sm mt-4">
                       You can close this window after the app opens.
                   </p>
               </div>
           </div>
       );
   }

2. Add route in App.jsx or router:
   <Route path="/auth/desktop-callback" element={<DesktopCallback />} />

3. Deploy this page to console.machpay.xyz
```

**Verification:**
- Page renders correctly
- Can be accessed at `/auth/desktop-callback`

---

### Prompt 4.5: Update Google OAuth Flow in Desktop App

**Objective:** Change desktop OAuth to use the new callback URL

**Why:** We need to send users to the desktop-specific callback

```
Update the Google OAuth flow in the desktop app to use the new callback.

Location: machpay-console/src/pages/Auth.jsx

1. Update the redirect URI for desktop:

   const isTauriApp = window.__TAURI__ !== undefined;
   const GOOGLE_REDIRECT_URI = isTauriApp 
       ? 'https://console.machpay.xyz/auth/desktop-callback'  // NEW: Desktop-specific
       : `${window.location.origin}/auth/callback`;           // Web stays the same

2. Simplify the desktop flow (remove polling):

   if (isTauri) {
       console.log('[Auth] Desktop app - opening browser for OAuth');
       
       const { open } = await import('@tauri-apps/api/shell');
       await open(googleAuthUrl.toString());
       
       toast.info('Complete login in your browser', {
           title: 'Browser Opened',
           duration: 5000,
       });
       
       // No more polling! Deep link will handle the callback
       // Just set a flag that we're waiting
       sessionStorage.setItem('awaiting_google_auth', 'true');
       
       // Don't setConnecting(null) yet - let deep link handler do it
       return;
   }

3. Add deep link handler (from Prompt 4.2) that completes the login:

   const handleDeepLinkAuth = async (token, refreshToken) => {
       try {
           // Store tokens
           localStorage.setItem('auth_token', token);
           if (refreshToken) {
               localStorage.setItem('refresh_token', refreshToken);
           }
           
           // Get user info with token
           const response = await fetch(`${API_URL}/v1/user/me`, {
               headers: { 'Authorization': `Bearer ${token}` }
           });
           const user = await response.json();
           
           // Complete login
           setConnecting(null);
           toast.success('Login Successful!', {
               title: 'Welcome ' + (user.displayName || user.email),
           });
           
           const route = user.role === 'architect' ? '/dashboard' : 
                         user.role === 'operator' ? '/studio' : '/explorer';
           onAuthenticated?.({ ...user, token, route });
           
       } catch (err) {
           console.error('[Deep Link] Auth failed:', err);
           setConnecting(null);
           toast.error('Login failed. Please try again.');
       }
   };
```

**Verification:**
- Click Google login in desktop app
- Browser opens to Google
- After login, browser redirects to desktop-callback page
- Page calls backend, gets deep link
- Browser redirects to machpay://
- Desktop app receives token and logs in

---

### Prompt 4.6: Update Google Cloud Console

**Objective:** Add the new redirect URI to Google OAuth settings

**Why:** Google only allows redirects to registered URIs

```
Add the desktop callback URL to Google Cloud Console.

1. Go to: https://console.cloud.google.com/apis/credentials

2. Find your OAuth 2.0 Client ID (MachPay project)

3. Under "Authorized redirect URIs", ADD:
   https://console.machpay.xyz/auth/desktop-callback

4. Keep existing URIs (don't remove):
   - https://console.machpay.xyz/auth/callback
   - http://localhost:5173/auth/callback (for dev)

5. Save changes

6. Wait 5 minutes for propagation

7. Test the full flow:
   - Open desktop app
   - Click Google login
   - Complete OAuth in browser
   - Verify redirect to desktop-callback
   - Verify deep link redirect
   - Verify app receives token
```

---

### Prompt 4.7: Full Integration Test

**Objective:** Test the complete flow end-to-end

```
Perform a complete integration test of the Google OAuth flow.

Test Cases:

1. Happy Path:
   - [ ] Open desktop app (fresh install)
   - [ ] Click Google login
   - [ ] Browser opens to Google
   - [ ] Complete Google sign-in
   - [ ] Browser redirects to desktop-callback page
   - [ ] Page shows "Completing login..."
   - [ ] Desktop app comes to foreground
   - [ ] User is logged in
   - [ ] Correct user info displayed

2. App Not Running:
   - [ ] Close desktop app completely
   - [ ] In browser, go to: https://console.machpay.xyz/auth/desktop-callback?code=test
   - [ ] Should attempt deep link
   - [ ] Fallback button shown
   - [ ] Clicking button should try to open app

3. Error Cases:
   - [ ] Invalid code → Shows error message
   - [ ] Network error → Shows error message
   - [ ] Cancel OAuth → Shows error message

4. Security:
   - [ ] Tokens not exposed in browser URL bar
   - [ ] Tokens not logged to console (except debug mode)
   - [ ] State parameter validated (CSRF protection)

Debug Commands:
   # Check if deep link is registered (macOS)
   /System/Library/Frameworks/CoreServices.framework/Versions/A/Frameworks/LaunchServices.framework/Versions/A/Support/lsregister -dump | grep machpay

   # Test deep link
   open "machpay://auth/callback?token=test"

   # View Tauri logs
   # In terminal where you ran `npm run tauri dev`
```

---

## Checklist

### Phase 1: Deep Link Setup
- [ ] 4.1: Custom URL scheme registered in Tauri
- [ ] Test: `open "machpay://test"` launches app

### Phase 2: Event Bridge
- [ ] 4.2: Deep link events emitted to frontend
- [ ] Test: Frontend console shows received URLs

### Phase 3: Backend
- [ ] 4.3: Desktop callback endpoint created
- [ ] Test: Endpoint returns deep link URL

### Phase 4: Web Callback
- [ ] 4.4: Desktop callback page created
- [ ] 4.6: Google Console updated
- [ ] Test: Page exchanges code and redirects

### Phase 5: Integration
- [ ] 4.5: Desktop app uses new flow
- [ ] 4.7: Full end-to-end test passes

---

## Troubleshooting

### Deep link not working on macOS
```bash
# Re-register URL schemes
/System/Library/Frameworks/CoreServices.framework/Versions/A/Frameworks/LaunchServices.framework/Versions/A/Support/lsregister -u /Applications/MachPay.app

# Then re-install the app
```

### App doesn't receive URL
- Check `tauri.conf.json` has correct bundle identifier
- Check `Info.plist` (in built .app) has `CFBundleURLTypes`
- Restart the app after changes

### Backend returns error
- Check Google Client Secret is configured
- Check redirect URI matches exactly
- Check code hasn't expired (10 min lifetime)

### Token not working
- Verify token format (should be JWT)
- Check token hasn't expired
- Verify API endpoint accepts the token

---

## Files Modified

| File | Purpose |
|------|---------|
| `src-tauri/tauri.conf.json` | Register URL scheme |
| `src-tauri/Cargo.toml` | Add deep-link plugin |
| `src-tauri/src/main.rs` | Handle deep links |
| `src/pages/DesktopCallback.jsx` | NEW: Web callback page |
| `src/pages/Auth.jsx` | Update OAuth flow |
| `src/App.jsx` | Add deep link listener |
| `backend/.../auth_handler.go` | NEW: Desktop callback endpoint |

---

**Document Version:** 1.0  
**Created:** 2025-01-01  
**Critical Note:** Test EACH step before proceeding to the next!

