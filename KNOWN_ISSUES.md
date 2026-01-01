# Known Issues

## 🐛 Active Issues

### Desktop App: Google OAuth Not Working
**Priority:** High  
**Status:** Parked  
**Created:** 2025-01-01

**Description:**
Google OAuth login in the desktop app gets stuck in "Connecting..." state after browser completes authentication.

**Attempted Solutions:**
1. ❌ Clipboard relay (browser → desktop app via clipboard)
   - Issue: Browser security blocks automatic clipboard writes
2. ❌ LocalStorage relay (browser → desktop app via localStorage)
   - Issue: Browser and desktop app don't share localStorage (different origins)
3. ❌ Local HTTP callback server (localhost:3737)
   - Issue: Google OAuth redirect URL changes, added complexity
4. ❌ Backend polling with in-memory cache
   - Issue: Browser (HTTPS) → localhost (HTTP) blocked by Mixed Content policy

**Root Cause:**
Cross-origin security restrictions prevent the browser callback page (running on `https://console.machpay.xyz`) from communicating with the desktop app or calling the local backend (`http://localhost:8081`).

**Technical Details:**
- Browser: Opens `https://console.machpay.xyz/auth/callback` after OAuth
- Desktop App: Embedded webview running on `tauri://localhost`
- Local Backend: Running on `http://localhost:8081`
- Problem: No reliable communication channel between these three contexts

**Proposed Solutions (For Future):**
1. **Deep Linking with Custom Protocol** (Recommended)
   - Register `machpay://` protocol handler on desktop
   - Modify Google OAuth redirect to `machpay://auth/callback?code=...`
   - Requires updating Google Cloud Console configuration
   - Reference: How VS Code, Discord, and Slack handle OAuth

2. **Dedicated Desktop Callback Endpoint**
   - Deploy separate HTTPS endpoint just for desktop callbacks
   - Use WebSocket or Server-Sent Events for real-time notification
   - More complex infrastructure

3. **QR Code Flow**
   - Desktop app shows QR code
   - User scans with phone, completes OAuth
   - Phone sends token to backend with session ID
   - Desktop app polls for token
   - Works but poor UX

**Workaround:**
Users can still login via:
- Email OTP (already working)
- Wallet signature (already working)
- Web app at https://console.machpay.xyz (full Google OAuth support)

**Related Files:**
- `/Users/abhishektomar/Desktop/git/machpay-console/src/pages/Auth.jsx`
- `/Users/abhishektomar/Desktop/git/machpay-console/src/pages/AuthCallback.jsx`
- `/Users/abhishektomar/Desktop/git/machpay-console/src/pages/DesktopAuthSuccess.jsx`
- `/Users/abhishektomar/Desktop/git/machpay-backend/internal/handlers/sql/auth_handler.go`
- `/Users/abhishektomar/Desktop/git/machpay-console/src-tauri/src/main.rs`

**Next Steps:**
- [ ] Research `machpay://` custom protocol registration with Tauri
- [ ] Test deep linking with Google OAuth redirect
- [ ] Update Google Cloud Console with custom protocol redirect URI
- [ ] Implement fallback UI in app to guide users to alternate login methods

---

## 📝 Notes
- Desktop app is fully functional for wallet and email login
- Google OAuth works perfectly in web app
- Issue only affects desktop app + Google OAuth combination

