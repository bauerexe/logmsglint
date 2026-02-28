package analysis

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	t.Parallel()
	t.Run("missing file uses defaults", func(t *testing.T) {
		t.Parallel()
		cfg, err := LoadConfig(filepath.Join(t.TempDir(), "missing.yml"))
		require.NoError(t, err)
		assert.True(t, cfg.Rules.Lowercase)
		assert.True(t, cfg.Rules.English)
		assert.True(t, cfg.Rules.NoSpecial)
		assert.True(t, cfg.Rules.Sensitive)
	})

	t.Run("yaml config overrides rules", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		path := filepath.Join(dir, "logmsglint.yml")
		err := os.WriteFile(path, []byte("rules:\n  lowercase: false\n  english: true\n  nospecial: false\n  sensitive: true\nsensitive:\n  patterns:\n    - '(?i)card\\s*\\d{4}'\n"), 0o600)
		require.NoError(t, err)

		cfg, err := LoadConfig(path)
		require.NoError(t, err)
		assert.False(t, cfg.Rules.Lowercase)
		assert.False(t, cfg.Rules.NoSpecial)
		assert.Equal(t, []string{"(?i)card\\s*\\d{4}"}, cfg.Sensitive.Patterns)
	})
}
