package artifact

// ceiling is the default maximum bytes for one stored artifact, set well above a screenshot (bounded at 32 MiB) so no
// legitimate artifact is rejected. Options.Ceiling overrides it in one place.
const ceiling = 128 << 20
