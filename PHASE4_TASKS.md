# Phase 4: Desktop Bundling - Detailed Task Breakdown

**Status:** 50% Complete  
**Priority:** Medium  
**Estimated Time:** 2-3 days

---

## 📊 Current Status Analysis

### ✅ **COMPLETED** (CLI Side)
- `open` command fully implemented (`internal/cmd/open.go`)
- Desktop app path detection for macOS, Windows, Linux
- Browser fallback mechanism
- Route support for deep linking
- `--web` flag for forcing browser mode

### ✅ **COMPLETED** (Console Side - Partial)
- Tauri configuration exists (`src-tauri/tauri.conf.json`)
- First launch modal component exists (`src/components/FirstLaunchModal.jsx`)
- Tauri dependencies installed (`@tauri-apps/api`, `@tauri-apps/cli`)
- Build scripts configured (`tauri:dev`, `tauri:build`)
- Icons configured (`src-tauri/icons/icon.svg`)
- External binaries configured (CLI and Gateway bundled)

### ⚠️ **PENDING VERIFICATION**
- Desktop app actually builds successfully
- Desktop app can be installed and launched
- First launch modal triggers correctly
- CLI-to-desktop communication works
- Deep linking routes work in desktop app
- Bundled CLI/Gateway binaries work

---

## 🎯 Pending Tasks

### **Task 1: Verify Tauri Desktop App Builds**

#### 1.1 Build Desktop App for Development
**Prompt for AI/Human:**
```
Navigate to machpay-console and verify the Tauri desktop app builds successfully:

1. Check that all dependencies are installed
2. Run the development build
3. Verify the app launches
4. Test basic functionality (navigation, auth, etc.)

Commands to run:
```bash
cd /Users/abhishektomar/Desktop/git/machpay-console
npm install
npm run tauri:dev
```

Expected outcome:
- App opens in a native window
- No console errors
- Can navigate between pages
- Looks identical to web version
```

**Success Criteria:**
- [ ] Desktop app window opens successfully
- [ ] No TypeScript/build errors
- [ ] App UI renders correctly
- [ ] Navigation works
- [ ] No console errors in dev tools

---

#### 1.2 Build Production Desktop App
**Prompt for AI/Human:**
```
Build the production version of the desktop app for your platform:

Commands:
```bash
cd /Users/abhishektomar/Desktop/git/machpay-console
npm run tauri:build
```

This will create:
- macOS: .dmg and .app in src-tauri/target/release/bundle/dmg/
- Windows: .msi in src-tauri/target/release/bundle/msi/
- Linux: .AppImage and .deb in src-tauri/target/release/bundle/

Test the installation:
1. Install the generated package
2. Launch MachPay.app (or equivalent)
3. Verify it works standalone
```

**Success Criteria:**
- [ ] Build completes without errors
- [ ] Installation package is created
- [ ] App installs successfully
- [ ] App launches from Applications folder/Start Menu
- [ ] App runs independently (not from terminal)

---

### **Task 2: Verify First Launch Modal Integration**

#### 2.1 Test First Launch Detection
**Prompt for AI/Human:**
```
Verify the first launch modal appears on initial app startup:

File to check: machpay-console/src/App.jsx

Ensure:
1. useFirstLaunch() hook is imported and used
2. FirstLaunchModal is rendered when isFirstLaunch === true
3. Modal shows welcome screen
4. User can choose to install CLI or skip

Test:
1. Delete any existing app preferences/config:
   - macOS: ~/Library/Application Support/xyz.machpay.console/
   - Windows: %APPDATA%\xyz.machpay.console\
   - Linux: ~/.config/xyz.machpay.console/
2. Launch the app
3. Verify first launch modal appears
4. Complete the wizard
5. Relaunch app - modal should NOT appear again
```

**Success Criteria:**
- [ ] Modal appears on first launch only
- [ ] Welcome screen displays correctly
- [ ] "Install CLI" option works (or shows appropriate message)
- [ ] "Skip" option works
- [ ] Modal doesn't appear on subsequent launches
- [ ] Preference is saved correctly

---

#### 2.2 Implement Missing Tauri Backend Commands
**Prompt for AI/Human:**
```
The FirstLaunchModal component calls several Tauri commands that may not be implemented yet:
- check_cli_installed
- install_cli
- complete_first_launch
- is_first_launch

File to edit: machpay-console/src-tauri/src/main.rs

Implement these commands:

```rust
use tauri::command;
use std::process::Command;
use std::fs;
use std::path::PathBuf;

#[command]
fn is_first_launch() -> Result<bool, String> {
    let app_dir = dirs::config_dir()
        .ok_or("Could not get config directory")?
        .join("xyz.machpay.console");
    
    let flag_file = app_dir.join(".first_launch_complete");
    Ok(!flag_file.exists())
}

#[command]
fn complete_first_launch() -> Result<(), String> {
    let app_dir = dirs::config_dir()
        .ok_or("Could not get config directory")?
        .join("xyz.machpay.console");
    
    fs::create_dir_all(&app_dir)
        .map_err(|e| format!("Failed to create app dir: {}", e))?;
    
    let flag_file = app_dir.join(".first_launch_complete");
    fs::write(flag_file, "")
        .map_err(|e| format!("Failed to write flag file: {}", e))?;
    
    Ok(())
}

#[command]
fn check_cli_installed() -> Result<CliStatus, String> {
    let output = Command::new("which")
        .arg("machpay")
        .output();
    
    match output {
        Ok(output) => {
            let installed = output.status.success();
            let path = if installed {
                Some(String::from_utf8_lossy(&output.stdout).trim().to_string())
            } else {
                None
            };
            
            Ok(CliStatus {
                installed,
                path,
                version: None, // Could run machpay version to get this
            })
        }
        Err(_) => Ok(CliStatus {
            installed: false,
            path: None,
            version: None,
        }),
    }
}

#[command]
fn install_cli() -> Result<String, String> {
    // Attempt to install via Homebrew on macOS
    #[cfg(target_os = "macos")]
    {
        let output = Command::new("brew")
            .args(&["install", "machpay-xyz/tap/machpay"])
            .output()
            .map_err(|e| format!("Failed to run brew: {}", e))?;
        
        if output.status.success() {
            Ok("/opt/homebrew/bin/machpay".to_string())
        } else {
            Err(String::from_utf8_lossy(&output.stderr).to_string())
        }
    }
    
    #[cfg(not(target_os = "macos"))]
    {
        Err("CLI installation not yet supported on this platform".to_string())
    }
}

#[derive(serde::Serialize)]
struct CliStatus {
    installed: bool,
    path: Option<String>,
    version: Option<String>,
}

// Register commands in main():
fn main() {
    tauri::Builder::default()
        .invoke_handler(tauri::generate_handler![
            is_first_launch,
            complete_first_launch,
            check_cli_installed,
            install_cli,
        ])
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}
```

Add to Cargo.toml if needed:
```toml
[dependencies]
dirs = "5.0"
serde = { version = "1.0", features = ["derive"] }
```
```

**Success Criteria:**
- [ ] All Tauri commands implemented in main.rs
- [ ] Commands compile without errors
- [ ] FirstLaunchModal can call these commands
- [ ] CLI installation works (at least on macOS)
- [ ] First launch flag is persisted

---

### **Task 3: Test CLI-to-Desktop Integration**

#### 3.1 Test `machpay open` Command
**Prompt for AI/Human:**
```
Test that the CLI can successfully launch the desktop app:

Prerequisites:
1. Desktop app is installed in /Applications/MachPay.app (macOS)
2. CLI is built and available

Tests:
```bash
# Test 1: Open desktop app (no route)
./machpay open

# Test 2: Open with specific route
./machpay open marketplace
./machpay open funding
./machpay open settings

# Test 3: Force web mode
./machpay open --web

# Test 4: Desktop app not installed (should fallback to browser)
# Temporarily rename /Applications/MachPay.app
mv /Applications/MachPay.app /Applications/MachPay.app.bak
./machpay open
mv /Applications/MachPay.app.bak /Applications/MachPay.app
```

Expected behavior:
- Desktop app launches (not browser)
- App navigates to correct route
- CLI shows success message
- --web flag forces browser
```

**Success Criteria:**
- [ ] `machpay open` launches desktop app
- [ ] Routes are passed correctly to desktop app
- [ ] Desktop app navigates to correct page
- [ ] Browser fallback works when app not installed
- [ ] `--web` flag forces browser mode

---

#### 3.2 Implement Deep Link Route Handling
**Prompt for AI/Human:**
```
The desktop app needs to handle --route argument passed by the CLI.

File to edit: machpay-console/src-tauri/src/main.rs

Add route handling:

```rust
use tauri::Manager;

fn main() {
    tauri::Builder::default()
        .setup(|app| {
            // Handle CLI arguments
            let args: Vec<String> = std::env::args().collect();
            let route = args.iter()
                .find(|arg| arg.starts_with("--route="))
                .and_then(|arg| arg.strip_prefix("--route="))
                .unwrap_or("/explorer");
            
            // Get the main window
            if let Some(window) = app.get_window("main") {
                // Emit route event to frontend
                window.emit("navigate", route).ok();
            }
            
            Ok(())
        })
        .invoke_handler(tauri::generate_handler![
            is_first_launch,
            complete_first_launch,
            check_cli_installed,
            install_cli,
        ])
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}
```

File to edit: machpay-console/src/App.jsx

Add event listener:

```jsx
import { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';

function App() {
  const navigate = useNavigate();
  
  useEffect(() => {
    // Listen for route changes from Tauri
    if (window.__TAURI__) {
      const { listen } = window.__TAURI__.event;
      
      const unlisten = listen('navigate', (event) => {
        navigate(event.payload);
      });
      
      return () => {
        unlisten.then(f => f());
      };
    }
  }, [navigate]);
  
  // ... rest of app
}
```
```

**Success Criteria:**
- [ ] Desktop app receives --route argument
- [ ] App navigates to correct route on launch
- [ ] Routes work for all supported paths
- [ ] Navigation is smooth (no flicker)

---

### **Task 4: App Icons and Branding**

#### 4.1 Verify App Icons
**Prompt for AI/Human:**
```
Check that the desktop app has proper icons at all sizes.

Current icon: machpay-console/src-tauri/icons/icon.svg

Tauri should auto-generate all required sizes, but verify:

1. Check src-tauri/icons/ directory contains:
   - 32x32.png
   - 128x128.png
   - 128x128@2x.png
   - icon.icns (macOS)
   - icon.ico (Windows)
   - icon.png (Linux)

2. If missing, regenerate icons:
```bash
cd machpay-console/src-tauri
npx @tauri-apps/cli icon /path/to/your-icon.png
```

3. Verify installed app shows correct icon:
   - In dock/taskbar when running
   - In Applications folder/Start Menu
   - In window title bar
```

**Success Criteria:**
- [ ] All icon sizes generated
- [ ] App shows correct icon in Applications folder
- [ ] App shows correct icon when running
- [ ] Icon is clear and recognizable

---

#### 4.2 Update App Metadata
**Prompt for AI/Human:**
```
Verify app metadata in tauri.conf.json:

File: machpay-console/src-tauri/tauri.conf.json

Check/update:
```json
{
  "package": {
    "productName": "MachPay",
    "version": "1.0.0"  // Should match package.json
  },
  "tauri": {
    "bundle": {
      "identifier": "xyz.machpay.console",
      "publisher": "MachPay",
      "copyright": "Copyright © 2025 MachPay",
      "category": "Finance",
      "shortDescription": "Pay-per-use API monetization platform",
      "longDescription": "MachPay enables instant API monetization with blockchain-powered micropayments. Vendors can monetize APIs instantly, agents access services at fraction of traditional costs.",
      "macOS": {
        "minimumSystemVersion": "10.15"
      }
    }
  }
}
```
```

**Success Criteria:**
- [ ] Product name is correct
- [ ] Version matches package.json
- [ ] Bundle identifier is unique
- [ ] Description is accurate
- [ ] Copyright is current

---

### **Task 5: Bundle CLI and Gateway Binaries**

#### 5.1 Verify External Binaries are Included
**Prompt for AI/Human:**
```
The desktop app should bundle CLI and Gateway binaries for offline use.

Current config (tauri.conf.json):
```json
"externalBin": [
  "bin/machpay",
  "bin/machpay-gateway"
]
```

Verify binaries exist:
```bash
cd machpay-console/src-tauri
ls -lh bin/
# Should show:
# machpay
# machpay-gateway
# Platform-specific variants (-aarch64-apple-darwin, etc.)
```

If missing, copy from releases:
```bash
# Copy CLI binary
cp /Users/abhishektomar/Desktop/git/machpay-cli/machpay src-tauri/bin/

# Copy Gateway binary (from gateway repo)
cp /Users/abhishektomar/Desktop/git/machpay-gateway/machpay src-tauri/bin/machpay-gateway

# Rename with platform suffix for Tauri
mv src-tauri/bin/machpay src-tauri/bin/machpay-aarch64-apple-darwin
mv src-tauri/bin/machpay-gateway src-tauri/bin/machpay-gateway-aarch64-apple-darwin
```

Test bundled binaries work:
```bash
# After building desktop app
/Applications/MachPay.app/Contents/MacOS/machpay version
/Applications/MachPay.app/Contents/MacOS/machpay-gateway --version
```
```

**Success Criteria:**
- [ ] Binaries exist in src-tauri/bin/
- [ ] Binaries are included in desktop app bundle
- [ ] Binaries are executable
- [ ] App can invoke bundled binaries
- [ ] Binaries work without external dependencies

---

#### 5.2 Test Bundled Binary Invocation
**Prompt for AI/Human:**
```
Add functionality to use bundled binaries from the desktop app.

File to create: machpay-console/src-tauri/src/commands.rs

```rust
use tauri::api::process::{Command, CommandEvent};
use tauri::Manager;

#[tauri::command]
async fn run_cli_command(args: Vec<String>) -> Result<String, String> {
    let (mut rx, _child) = Command::new_sidecar("machpay")
        .expect("failed to create `machpay` sidecar")
        .args(args)
        .spawn()
        .map_err(|e| format!("Failed to spawn CLI: {}", e))?;

    let mut output = String::new();
    while let Some(event) = rx.recv().await {
        match event {
            CommandEvent::Stdout(line) => output.push_str(&line),
            CommandEvent::Stderr(line) => output.push_str(&line),
            _ => {}
        }
    }

    Ok(output)
}

#[tauri::command]
async fn run_gateway(config_path: String) -> Result<(), String> {
    Command::new_sidecar("machpay-gateway")
        .expect("failed to create `machpay-gateway` sidecar")
        .args(["--config", &config_path])
        .spawn()
        .map_err(|e| format!("Failed to start gateway: {}", e))?;

    Ok(())
}
```

Update main.rs to include these commands.
```

**Success Criteria:**
- [ ] Desktop app can invoke bundled CLI
- [ ] Desktop app can start bundled gateway
- [ ] Binaries run with correct permissions
- [ ] Output is captured correctly

---

### **Task 6: Cross-Platform Testing**

#### 6.1 Test on macOS
**Prompt for AI/Human:**
```
Complete end-to-end testing on macOS:

1. Build desktop app
2. Install from .dmg
3. Launch app from Applications
4. Complete first launch wizard
5. Test CLI integration (machpay open)
6. Test bundled binaries work
7. Test app updates (if implemented)
8. Verify app signature (for distribution)

Document any issues or platform-specific bugs.
```

**Success Criteria:**
- [ ] App installs from DMG
- [ ] App passes macOS Gatekeeper
- [ ] All features work on macOS
- [ ] No crashes or errors

---

#### 6.2 Test on Windows (if applicable)
**Prompt for AI/Human:**
```
Complete end-to-end testing on Windows:

1. Build desktop app (.msi)
2. Install from installer
3. Launch app from Start Menu
4. Complete first launch wizard
5. Test CLI integration (machpay open)
6. Test bundled binaries work
7. Verify app shows in installed programs

Note: May need Windows-specific adjustments in CLI open command.
```

**Success Criteria:**
- [ ] App installs from MSI
- [ ] App works on Windows 10/11
- [ ] CLI can launch desktop app
- [ ] No Windows-specific bugs

---

#### 6.3 Test on Linux (if applicable)
**Prompt for AI/Human:**
```
Complete end-to-end testing on Linux:

1. Build desktop app (.AppImage or .deb)
2. Install/run application
3. Test all functionality
4. Verify CLI integration
5. Test bundled binaries

Note: May need to adjust binary paths for Linux.
```

**Success Criteria:**
- [ ] App runs on Ubuntu/Debian
- [ ] CLI can launch desktop app
- [ ] Binaries work on Linux

---

## 📋 Quick Checklist

### CLI Side (Already Complete ✅)
- [x] `machpay open` command implemented
- [x] Desktop app path detection
- [x] Browser fallback
- [x] Route support
- [x] Web flag

### Desktop App Side
- [ ] App builds successfully
- [ ] Production bundle works
- [ ] First launch modal triggers
- [ ] Tauri backend commands implemented
- [ ] CLI can launch desktop app
- [ ] Deep link routing works
- [ ] Icons are correct
- [ ] Metadata is updated
- [ ] Binaries are bundled
- [ ] Cross-platform tested

---

## 🚀 Recommended Execution Order

1. **Day 1: Build and Basic Testing**
   - Task 1.1: Build desktop app for development
   - Task 1.2: Build production bundle
   - Task 4.1: Verify icons
   - Task 4.2: Update metadata

2. **Day 2: First Launch and Integration**
   - Task 2.2: Implement Tauri backend commands
   - Task 2.1: Test first launch modal
   - Task 3.1: Test CLI-to-desktop integration
   - Task 3.2: Implement deep link routing

3. **Day 3: Bundled Binaries and Polish**
   - Task 5.1: Verify bundled binaries
   - Task 5.2: Test binary invocation
   - Task 6.1/6.2/6.3: Cross-platform testing

---

## 🎯 Success Metrics

**Phase 4 will be 100% complete when:**
- ✅ Desktop app builds and installs on target platforms
- ✅ First launch modal appears and works correctly
- ✅ `machpay open` successfully launches desktop app
- ✅ Deep linking routes work correctly
- ✅ Bundled CLI/Gateway binaries are functional
- ✅ App passes testing on all target platforms
- ✅ App is ready for distribution (signed, notarized if needed)

---

## 📝 Notes

- **Priority:** Medium (not blocking v0.1.0 CLI release)
- **Can ship CLI v0.1.0 without desktop app**
- Desktop app is value-add, not critical path
- Focus on macOS first, then Windows, then Linux
- Desktop app can be v1.1.0 or later release

---

## 🔗 Related Files

**CLI:**
- `/Users/abhishektomar/Desktop/git/machpay-cli/internal/cmd/open.go`

**Console:**
- `/Users/abhishektomar/Desktop/git/machpay-console/src-tauri/tauri.conf.json`
- `/Users/abhishektomar/Desktop/git/machpay-console/src-tauri/src/main.rs`
- `/Users/abhishektomar/Desktop/git/machpay-console/src/components/FirstLaunchModal.jsx`
- `/Users/abhishektomar/Desktop/git/machpay-console/package.json`

---

**Last Updated:** 2025-01-01

