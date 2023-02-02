package mkvmerge

import (
	"path/filepath"
	"strings"

	"github.com/elboletaire/remuxing/models"
)

func trackName(track models.TrackController) string {
	if filepath.Ext(track.Input.FileName) == ".mkv" {
		return track.Track.GetArgID()
	}

	return track.Track.GetArgIDLabel(filepath.Base(filepath.Dir(track.Input.FileName)))
}

func isExternal(track models.TrackController) bool {
	return filepath.Ext(track.Input.FileName) != ".mkv"
}

func inStrings(s []string, v string) bool {
	for _, vv := range s {
		if vv == v {
			return true
		}
	}
	return false
}

func ext(path string) string {
	return strings.TrimPrefix(filepath.Ext(path), ".")
}

func trimmExt(path string) string {
	return strings.TrimSuffix(path, filepath.Ext(path))
}
