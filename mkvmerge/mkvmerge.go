package mkvmerge

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/elboletaire/remuxing/models"
)

type Videos []Video

type Video struct {
	Name       string
	Path       string
	OutputPath string
	Audio      []string
	Subtitles  []string
}

func (v Video) NeedMerge() bool {
	return len(v.Audio) > 0 || len(v.Subtitles) > 0
}

func (v Video) ReplaceOrigin() error {
	return os.Rename(v.OutputPath, v.Path)
}

func (v Video) CmdArgs() []string {
	tracks := models.BuildTracks(append([]string{v.Path}, append(v.Audio, v.Subtitles...)...))
	command := []string{"-o", v.OutputPath, "--abort-on-warnings"}
	command = videoArgs(tracks.GetBestVideo(), command)
	command = audiosArgs(tracks.GetBestAudios(nil), command)
	command = subtitlesArgs(tracks.GetBestSubtitles(nil), command)
	return command
}

func Merge(c Config) error {
	if len(c.VideoExtensions) == 0 {
		c.VideoExtensions = []string{"mkv"}
	}

	if len(c.AudioExtensions) == 0 {
		c.AudioExtensions = []string{"mka"}
	}

	if len(c.SubtitlesExtensions) == 0 {
		c.AudioExtensions = []string{"ass"}
	}

	files, err := collectFiles(c)
	if err != nil {
		return err
	}

	for i, v := range files {
		if !v.NeedMerge() {
			log.Printf("[#%d/%d] %q => nothing found to be merged", i+1, len(files), v.Path)
			continue
		}

		log.Printf("[#%d/%d] %q => audio count: %d, subtitles count: %d", i+1, len(files), v.Path, len(v.Audio), len(v.Subtitles))
		if c.DryRun {
			continue
		}
		out, err := exec.Command("mkvmerge", v.CmdArgs()...).CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to run mkvmerge: %w  out: %s", err, string(out))
		}

		if err := v.ReplaceOrigin(); err != nil {
			return fmt.Errorf("failed to rename file: %w", err)
		}

		log.Printf("file %s processed", v.Path)
		log.Println("removing merged audio/subtitle files")
		for _, vv := range append(append([]string(nil), v.Subtitles...), v.Audio...) {
			if err := os.Remove(vv); err != nil {
				log.Printf("failed to remove file: %s, err: %s", vv, err)
			}
		}
	}

	return nil
}

func collectFiles(c Config) ([]Video, error) {
	if !isDir(c.Target) {
		return nil, nil
	}

	dirEntities, err := os.ReadDir(c.Target)
	if err != nil {
		return nil, err
	}

	files := make(map[string]Video, 0)
	for _, f := range dirEntities {
		if f.IsDir() || !inStrings(c.VideoExtensions, ext(f.Name())) {
			continue
		}

		fName := strings.TrimSuffix(f.Name(), filepath.Ext(f.Name()))
		files[fName] = Video{
			Name:       fName,
			Path:       filepath.Join(c.Target, f.Name()),
			OutputPath: filepath.Join(c.Target, outputFile(f.Name())),
			Audio:      []string{},
			Subtitles:  []string{},
		}
	}

	if err := filepath.WalkDir(c.Target, func(path string, d fs.DirEntry, _ error) error {
		fName := strings.TrimSuffix(d.Name(), filepath.Ext(d.Name()))
		switch {
		case inStrings(c.AudioExtensions, ext(d.Name())):
			if v, ok := files[fName]; ok {
				v.Audio = append(v.Audio, path)
				files[fName] = v
			}
		case inStrings(c.SubtitlesExtensions, ext(d.Name())):
			if v, ok := files[fName]; ok {
				v.Subtitles = append(v.Subtitles, path)
				files[fName] = v
			}
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("failed to walk dir %w", err)
	}

	fileSlice := make([]Video, 0, len(files))
	for _, v := range files {
		fileSlice = append(fileSlice, v)
	}

	sort.Slice(fileSlice, func(i, j int) bool {
		return fileSlice[i].Name < fileSlice[j].Name
	})

	return fileSlice, nil
}

func videoArgs(video *models.TrackController, command []string) []string {
	command = append(
		command,
		"--title", "",
		"-A", "-T", "-S",
		"-d", video.Track.GetID(),
		video.Input.FileName,
	)

	return command
}

func audiosArgs(audios models.Tracks, command []string) []string {
	var defaultSet bool
	for _, audio := range audios {
		command = append(command, "-T")
		if !defaultSet && isExternal(audio) {
			command = append(command, "--default-track", audio.Track.GetID())
		}
		command = append(
			command,
			"--language", audio.Track.GetArgIDLabel(audio.Track.Properties.Language),
			"-a", audio.Track.GetID(),
			"--track-name", trackName(audio),
			"-D", "-S",
			audio.Input.FileName,
		)
	}

	return command
}

func subtitlesArgs(subtitles models.Tracks, command []string) []string {
	for _, subtitle := range subtitles {
		command = append(
			command,
			"-T",
			"--default-track-flag", subtitle.Track.GetArgIDLabel("false"),
			"-s", strconv.Itoa(int(subtitle.Track.ID)),
			"--track-name", trackName(subtitle),
			"-D", "-A",
		)

		if subtitle.Track.Properties.Forced {
			command = append(
				command,
				"--forced-track", subtitle.Track.GetArgIDLabel("true"),
			)
		}

		command = append(command, subtitle.Input.FileName)
	}

	return command
}

func outputFile(path string) string {
	return fmt.Sprintf("%s_merged.mkv", trimmExt(path))
}

func isDir(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}

	return fileInfo.IsDir()
}
