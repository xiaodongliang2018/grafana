package sources

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/grafana/grafana/pkg/plugins"
	"github.com/grafana/grafana/pkg/setting"
)

func TestSources_List(t *testing.T) {
	t.Run("Plugin sources are populated by default and listed in specific order", func(t *testing.T) {
		testdata, err := filepath.Abs("../testdata")
		require.NoError(t, err)

		cfg := &setting.Cfg{
			StaticRootPath:     testdata,
			PluginsPath:        filepath.Join(testdata, "pluginRootWithDist"),
			BundledPluginsPath: filepath.Join(testdata, "unsigned-panel"),
			PluginSettings: setting.PluginSettings{
				"foo": map[string]string{
					"path": filepath.Join(testdata, "test-app"),
				},
				"bar": map[string]string{
					"url": "https://grafana.plugin",
				},
			},
		}

		s := ProvideService(cfg)
		srcs := s.List(context.Background())

		ctx := context.Background()

		require.Len(t, srcs, 6)

		require.Equal(t, srcs[0].PluginClass(ctx), plugins.ClassCore)
		require.Equal(t, srcs[0].PluginURIs(ctx), []string{
			filepath.Join(testdata, "app", "plugins", "datasource"),
			filepath.Join(testdata, "app", "plugins", "panel"),
		})
		sig, exists := srcs[0].DefaultSignature(ctx)
		require.True(t, exists)
		require.Equal(t, plugins.SignatureStatusInternal, sig.Status)
		require.Equal(t, plugins.SignatureType(""), sig.Type)
		require.Equal(t, "", sig.SigningOrg)

		require.Equal(t, srcs[1].PluginClass(ctx), plugins.ClassBundled)
		require.Equal(t, srcs[1].PluginURIs(ctx), []string{filepath.Join(testdata, "unsigned-panel")})
		sig, exists = srcs[1].DefaultSignature(ctx)
		require.False(t, exists)
		require.Equal(t, plugins.Signature{}, sig)

		require.Equal(t, srcs[2].PluginClass(ctx), plugins.ClassExternal)
		require.Equal(t, srcs[2].PluginURIs(ctx), []string{
			filepath.Join(testdata, "pluginRootWithDist", "datasource"),
		})
		sig, exists = srcs[2].DefaultSignature(ctx)
		require.False(t, exists)
		require.Equal(t, plugins.Signature{}, sig)

		require.Equal(t, srcs[3].PluginClass(ctx), plugins.ClassExternal)
		require.Equal(t, srcs[3].PluginURIs(ctx), []string{
			filepath.Join(testdata, "pluginRootWithDist", "dist"),
		})
		sig, exists = srcs[3].DefaultSignature(ctx)
		require.False(t, exists)
		require.Equal(t, plugins.Signature{}, sig)

		require.Equal(t, srcs[4].PluginClass(ctx), plugins.ClassExternal)
		require.Equal(t, srcs[4].PluginURIs(ctx), []string{
			filepath.Join(testdata, "pluginRootWithDist", "panel"),
		})
		sig, exists = srcs[4].DefaultSignature(ctx)
		require.False(t, exists)
		require.Equal(t, plugins.Signature{}, sig)
	})
}
