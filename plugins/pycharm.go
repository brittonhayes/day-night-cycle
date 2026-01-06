package plugins

import (
	"fmt"
	"os"
	"path/filepath"
)

func PyCharm(config PluginConfig) error {
	lafClass := "com.intellij.ide.ui.laf.darcula.DarculaLaf"
	themeID := "Darcula"

	if config.IsLight {
		lafClass = "com.intellij.ide.ui.laf.IntelliJLaf"
		themeID = "IntelliJ"
	}

	// Allow custom theme class names via config
	if config.IsLight && config.Day != "" {
		themeID = config.Day
	} else if !config.IsLight && config.Night != "" {
		themeID = config.Night
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	// Find PyCharm configuration directory
	jetbrainsDir := filepath.Join(home, "Library/Application Support/JetBrains")

	entries, err := os.ReadDir(jetbrainsDir)
	if err != nil {
		return fmt.Errorf("JetBrains directory not found: %w", err)
	}

	// Look for PyCharm directories (e.g., PyCharm2025.3, PyCharmCE2025.1)
	var pycharmDir string
	for _, entry := range entries {
		if entry.IsDir() {
			name := entry.Name()
			if len(name) >= 7 && (name[:7] == "PyCharm" || name[:9] == "PyCharmCE") {
				pycharmDir = filepath.Join(jetbrainsDir, name)
				break
			}
		}
	}

	if pycharmDir == "" {
		return fmt.Errorf("PyCharm configuration directory not found in %s", jetbrainsDir)
	}

	// Create options directory if it doesn't exist
	optionsDir := filepath.Join(pycharmDir, "options")
	if err := os.MkdirAll(optionsDir, 0755); err != nil {
		return err
	}

	lafPath := filepath.Join(optionsDir, "laf.xml")

	// Create the laf.xml content
	content := fmt.Sprintf(`<application>
  <component name="LafManager" autodetect="false">
    <laf class-name="%s" themeId="%s" />
  </component>
</application>
`, lafClass, themeID)

	return os.WriteFile(lafPath, []byte(content), 0644)
}
