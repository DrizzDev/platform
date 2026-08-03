package build

import "runtime/debug"

func Read() Info {
	info := Info{name: "drizz", version: unknown, revision: unknown}
	value, found := debug.ReadBuildInfo()
	if !found {
		return info
	}
	if value.Main.Version != "" && value.Main.Version != "(devel)" {
		info.version = value.Main.Version
	}
	for _, setting := range value.Settings {
		switch setting.Key {
		case "vcs.revision":
			info.revision = setting.Value
		case "vcs.modified":
			info.modified = setting.Value == "true"
		}
	}
	return info
}
