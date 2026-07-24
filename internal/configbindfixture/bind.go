package configbindfixture

import "github.com/shibukawa/tinybind-go/configbind"

// Register returns Bind handles for Load tests (discovered by tinybind-gen).
func Register() *WebServerConfig {
	return configbind.Bind[WebServerConfig]("webserver")
}

// RegisterMigrate returns migrate options only when that CLI branch is selected.
func RegisterMigrate() *MigrateOptions {
	return configbind.SubCommand[MigrateOptions]("migrate", "run database migrations")
}
