package configbind_test

import (
	"testing"

	"github.com/shibukawa/tinybind-go/cliparser"
	"github.com/shibukawa/tinybind-go/configbind"
)

type multiPrefixConfig struct {
	Value string
}

func TestRegisterSupportsOneTypeAtMultiplePrefixes(t *testing.T) {
	configbind.ResetDefinitions()
	t.Cleanup(configbind.ResetDefinitions)
	register := func(prefix, value string) {
		key := prefix + ".value"
		configbind.Register[multiPrefixConfig](configbind.Definition{
			TypeName:  "example.test.MultiPrefixConfig",
			Prefix:    prefix,
			KnownKeys: []string{key},
			Defaults:  map[string]string{key: value},
			FlagMetas: []cliparser.FieldMeta{{Prefix: prefix, Key: "value"}},
			Apply: func(dst any, overlay *configbind.Overlay) error {
				if current, ok := overlay.GetString(key); ok {
					dst.(*multiPrefixConfig).Value = current
				}
				return nil
			},
		})
	}
	register("primary", "one")
	register("secondary", "two")
	configbind.ResetTargets()
	primary := configbind.Bind[multiPrefixConfig]("primary")
	secondary := configbind.Bind[multiPrefixConfig]("secondary")
	if _, err := configbind.Load(configbind.LoadOptions{
		Vendor:  "tinybind-test",
		Tool:    "multi-prefix",
		Args:    []string{},
		Environ: []string{},
	}); err != nil {
		t.Fatal(err)
	}
	if primary.Value != "one" || secondary.Value != "two" {
		t.Fatalf("binding metadata collided: primary=%q secondary=%q", primary.Value, secondary.Value)
	}
}
