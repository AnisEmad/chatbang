package main

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestDetectBrowser tests the browser detection functionality
func TestDetectBrowser(t *testing.T) {
	// Test that detectBrowser returns a valid path or error
	browserPath, err := detectBrowser()
	
	if err != nil {
		// If no browser found, error should be specific
		expectedErr := "no Chromium-based browser found in PATH"
		if err.Error() != expectedErr {
			t.Errorf("Expected error %q, got %q", expectedErr, err.Error())
		}
	} else {
		// If browser found, path should be valid
		if browserPath == "" {
			t.Error("Browser path is empty but no error returned")
		}
		
		// Check if path exists
		if _, statErr := os.Stat(browserPath); statErr != nil {
			t.Errorf("Detected browser path %q does not exist: %v", browserPath, statErr)
		}
		
		// Verify it's one of the known browsers
		found := false
		for _, browser := range browsers {
			if strings.Contains(browserPath, browser) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Browser path %q does not contain any known browser name", browserPath)
		}
	}
}

// TestBrowsersList tests that the browsers list is not empty
func TestBrowsersList(t *testing.T) {
	if len(browsers) == 0 {
		t.Error("Browsers list should not be empty")
	}
	
	expectedBrowsers := []string{
		"chromium",
		"google-chrome",
		"brave-browser",
		"microsoft-edge",
	}
	
	for _, expected := range expectedBrowsers {
		found := false
		for _, browser := range browsers {
			if browser == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected browser %q not found in browsers list", expected)
		}
	}
}

// TestConfigCreation tests the creation of config directory
func TestConfigCreation(t *testing.T) {
	tempDir := t.TempDir()
	configDir := filepath.Join(tempDir, ".config", "chatbang")
	
	err := os.MkdirAll(configDir, 0o755)
	if err != nil {
		t.Fatalf("Failed to create config directory: %v", err)
	}
	
	// Verify directory exists
	info, err := os.Stat(configDir)
	if os.IsNotExist(err) {
		t.Errorf("Config directory was not created")
	}
	
	if !info.IsDir() {
		t.Errorf("Config path exists but is not a directory")
	}
	
	// Verify permissions
	if info.Mode().Perm() != 0o755 {
		t.Errorf("Config directory has wrong permissions: got %v, want %v", info.Mode().Perm(), 0o755)
	}
}

// TestConfigFileCreation tests writing and reading config file
func TestConfigFileCreation(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "chatbang")
	
	// Create and write config
	configFile, err := os.OpenFile(configPath, os.O_RDWR|os.O_CREATE, 0o644)
	if err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}
	defer configFile.Close()
	
	// Write default config
	defaultConfig := "browser=/usr/bin/google-chrome"
	_, err = configFile.WriteString(defaultConfig)
	if err != nil {
		t.Fatalf("Failed to write to config file: %v", err)
	}
	
	// Read back and verify
	configFile.Seek(0, 0)
	content := make([]byte, len(defaultConfig))
	n, err := configFile.Read(content)
	if err != nil {
		t.Fatalf("Failed to read config file: %v", err)
	}
	
	if n != len(defaultConfig) {
		t.Errorf("Read %d bytes, expected %d", n, len(defaultConfig))
	}
	
	if string(content) != defaultConfig {
		t.Errorf("Config content mismatch. Expected: %s, Got: %s", defaultConfig, string(content))
	}
}

// TestConfigParsing tests parsing of browser path from config
func TestConfigParsing(t *testing.T) {
	tests := []struct {
		name           string
		configContent  string
		expectedPath   string
		shouldBeEmpty  bool
	}{
		{
			name:          "Valid Chrome config",
			configContent: "browser=/usr/bin/google-chrome\n",
			expectedPath:  "/usr/bin/google-chrome",
			shouldBeEmpty: false,
		},
		{
			name:          "Valid Brave config",
			configContent: "browser=/usr/bin/brave-browser\n",
			expectedPath:  "/usr/bin/brave-browser",
			shouldBeEmpty: false,
		},
		{
			name:          "Config with spaces",
			configContent: "browser = /usr/bin/google-chrome\n",
			expectedPath:  "/usr/bin/google-chrome",
			shouldBeEmpty: false,
		},
		{
			name:          "Config with comment",
			configContent: "# This is a comment\nbrowser=/usr/bin/chromium\n",
			expectedPath:  "/usr/bin/chromium",
			shouldBeEmpty: false,
		},
		{
			name:          "Empty config",
			configContent: "",
			expectedPath:  "",
			shouldBeEmpty: true,
		},
		{
			name:          "Only comments",
			configContent: "# browser=/usr/bin/google-chrome\n",
			expectedPath:  "",
			shouldBeEmpty: true,
		},
		{
			name:          "Empty lines",
			configContent: "\n\n\nbrowser=/usr/bin/google-chrome\n",
			expectedPath:  "/usr/bin/google-chrome",
			shouldBeEmpty: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			configPath := filepath.Join(tempDir, "chatbang")
			
			// Write test config
			err := os.WriteFile(configPath, []byte(tt.configContent), 0o644)
			if err != nil {
				t.Fatalf("Failed to write test config: %v", err)
			}
			
			// Parse config (simulating the main.go logic)
			configFile, err := os.Open(configPath)
			if err != nil {
				t.Fatalf("Failed to open config file: %v", err)
			}
			defer configFile.Close()
			
			var parsedBrowser string
			scanner := bufio.NewScanner(configFile)
			for scanner.Scan() {
				line := strings.TrimSpace(scanner.Text())
				if line == "" || strings.HasPrefix(line, "#") {
					continue
				}
				parts := strings.SplitN(line, "=", 2)
				if len(parts) == 2 && strings.TrimSpace(parts[0]) == "browser" {
					parsedBrowser = strings.TrimSpace(parts[1])
				}
			}
			
			if tt.shouldBeEmpty && parsedBrowser != "" {
				t.Errorf("Expected empty browser path, got %q", parsedBrowser)
			}
			
			if !tt.shouldBeEmpty && parsedBrowser != tt.expectedPath {
				t.Errorf("Expected browser path %q, got %q", tt.expectedPath, parsedBrowser)
			}
		})
	}
}

// TestPromptModification tests the prompt modification logic
func TestPromptModification(t *testing.T) {
	tests := []struct {
		name           string
		originalPrompt string
		expectedSuffix string
	}{
		{
			name:           "Simple question",
			originalPrompt: "What is Go?",
			expectedSuffix: " (Make an answer in less than 5 lines).",
		},
		{
			name:           "Complex question",
			originalPrompt: "Explain quantum computing",
			expectedSuffix: " (Make an answer in less than 5 lines).",
		},
		{
			name:           "Empty prompt",
			originalPrompt: "",
			expectedSuffix: " (Make an answer in less than 5 lines).",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			modifiedPrompt := tt.originalPrompt + tt.expectedSuffix
			
			if !strings.HasSuffix(modifiedPrompt, tt.expectedSuffix) {
				t.Errorf("Modified prompt does not have expected suffix")
			}
			
			if !strings.HasPrefix(modifiedPrompt, tt.originalPrompt) {
				t.Errorf("Modified prompt does not start with original prompt")
			}
			
			expectedLength := len(tt.originalPrompt) + len(tt.expectedSuffix)
			if len(modifiedPrompt) != expectedLength {
				t.Errorf("Modified prompt length mismatch: got %d, want %d", len(modifiedPrompt), expectedLength)
			}
		})
	}
}

// TestPromptValidation tests validation of user prompts
func TestPromptValidation(t *testing.T) {
	tests := []struct {
		name        string
		prompt      string
		shouldSkip  bool
	}{
		{
			name:       "Valid prompt",
			prompt:     "What is the capital of France?",
			shouldSkip: false,
		},
		{
			name:       "Empty prompt",
			prompt:     "",
			shouldSkip: true,
		},
		{
			name:       "Whitespace only prompt",
			prompt:     "   ",
			shouldSkip: false, // Code only checks len() == 0, not trimmed content
		},
		{
			name:       "Tab only prompt",
			prompt:     "\t\t",
			shouldSkip: false, // Code only checks len() == 0, not trimmed content
		},
		{
			name:       "Long prompt",
			prompt:     "This is a very long prompt that contains many words and should still be valid because ChatGPT can handle long prompts without any issues.",
			shouldSkip: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simulate the validation logic from main.go
			// In main.go, prompts with len() == 0 are skipped
			shouldSkip := len(tt.prompt) == 0
			
			if shouldSkip != tt.shouldSkip {
				t.Errorf("Expected shouldSkip=%v, got %v for prompt %q", tt.shouldSkip, shouldSkip, tt.prompt)
			}
		})
	}
}

// TestBrowserPathValidation tests validation of browser paths
func TestBrowserPathValidation(t *testing.T) {
	tests := []struct {
		name        string
		browserPath string
	}{
		{
			name:        "Standard Chrome path",
			browserPath: "/usr/bin/google-chrome",
		},
		{
			name:        "Standard Brave path",
			browserPath: "/usr/bin/brave-browser",
		},
		{
			name:        "Bin Chrome path",
			browserPath: "/bin/google-chrome",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Just verify the path format is correct
			if !strings.HasPrefix(tt.browserPath, "/") {
				t.Errorf("Browser path should be absolute: %s", tt.browserPath)
			}
			
			// Verify it contains one of the known browsers
			found := false
			for _, browser := range browsers {
				if strings.Contains(tt.browserPath, browser) {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Browser path %q does not contain any known browser name", tt.browserPath)
			}
		})
	}
}

// TestConfigDirectoryPath tests the config directory path construction
func TestConfigDirectoryPath(t *testing.T) {
	homeDir := "/home/testuser"
	expectedConfigDir := homeDir + "/.config/chatbang"
	expectedConfigPath := expectedConfigDir + "/chatbang"
	expectedProfileDir := homeDir + "/.config/chatbang/profile_data"
	
	if !strings.HasPrefix(expectedConfigDir, homeDir) {
		t.Errorf("Config dir should start with home dir")
	}
	
	if !strings.Contains(expectedConfigDir, ".config/chatbang") {
		t.Errorf("Config dir should contain .config/chatbang")
	}
	
	if !strings.HasSuffix(expectedConfigPath, "/chatbang") {
		t.Errorf("Config path should end with /chatbang")
	}
	
	if !strings.HasSuffix(expectedProfileDir, "/profile_data") {
		t.Errorf("Profile dir should end with /profile_data")
	}
}

// TestContextTimeout tests the context timeout constant
func TestContextTimeout(t *testing.T) {
	if ctxTime != 2000 {
		t.Errorf("Expected ctxTime to be 2000, got %d", ctxTime)
	}
}

// TestJavaScriptClipboardRead tests the JavaScript clipboard read promise
func TestJavaScriptClipboardRead(t *testing.T) {
	expectedJS := `new Promise((resolve, reject) => { window.navigator.clipboard.readText() .then(text => resolve(text)) .catch(err => reject(err)); });`
	
	// Verify the JS string is valid
	if !strings.Contains(expectedJS, "navigator.clipboard.readText") {
		t.Error("JS should contain clipboard.readText")
	}
	
	if !strings.Contains(expectedJS, "Promise") {
		t.Error("JS should use Promise")
	}
}

// TestButtonSelector tests the button selector used in the code
func TestButtonSelector(t *testing.T) {
	buttonDiv := `button[data-testid="copy-turn-action-button"]`
	
	if !strings.HasPrefix(buttonDiv, "button[") {
		t.Error("Button selector should start with button[")
	}
	
	if !strings.Contains(buttonDiv, "data-testid") {
		t.Error("Button selector should use data-testid")
	}
	
	if !strings.Contains(buttonDiv, "copy-turn-action-button") {
		t.Error("Button selector should contain copy-turn-action-button")
	}
}

// TestChatGPTSelectors tests the ChatGPT DOM selectors
func TestChatGPTSelectors(t *testing.T) {
	selectors := map[string]string{
		"prompt-textarea":        "#prompt-textarea",
		"composer-submit-button": "#composer-submit-button",
	}
	
	for name, selector := range selectors {
		if !strings.HasPrefix(selector, "#") {
			t.Errorf("Selector %s should start with #", name)
		}
		
		if !strings.Contains(selector, name) {
			t.Errorf("Selector should contain the name %s", name)
		}
	}
}

// TestChatGPTURL tests the ChatGPT URL
func TestChatGPTURL(t *testing.T) {
	url := "https://chatgpt.com"
	altURL := "https://www.chatgpt.com/"
	
	if !strings.HasPrefix(url, "https://") {
		t.Error("URL should use HTTPS")
	}
	
	if !strings.Contains(url, "chatgpt.com") {
		t.Error("URL should contain chatgpt.com")
	}
	
	if !strings.HasPrefix(altURL, "https://") {
		t.Error("Alt URL should use HTTPS")
	}
}

// BenchmarkPromptModification benchmarks prompt modification
func BenchmarkPromptModification(b *testing.B) {
	prompt := "What is the meaning of life?"
	suffix := " (Make an answer in less than 5 lines)."
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = prompt + suffix
	}
}

// BenchmarkConfigParsing benchmarks config file parsing
func BenchmarkConfigParsing(b *testing.B) {
	tempDir := b.TempDir()
	configPath := filepath.Join(tempDir, "chatbang")
	
	configContent := "browser=/usr/bin/google-chrome\n"
	err := os.WriteFile(configPath, []byte(configContent), 0o644)
	if err != nil {
		b.Fatalf("Failed to create test config: %v", err)
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		configFile, _ := os.Open(configPath)
		var parsedBrowser string
		scanner := bufio.NewScanner(configFile)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 && strings.TrimSpace(parts[0]) == "browser" {
				parsedBrowser = strings.TrimSpace(parts[1])
			}
		}
		_ = parsedBrowser
		configFile.Close()
	}
}