package mkvmerge

type Config struct {
	Target              string
	VideoExtensions     []string
	AudioExtensions     []string
	SubtitlesExtensions []string
	DryRun              bool
}
