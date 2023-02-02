package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/MrEhbr/media-normalizer/mkvmerge"
	mediasort "github.com/jpillora/media-sort/sort"
	"github.com/jpillora/sizestr"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
)

// These values are private which ensures they can only be set with the build flags.
var (
	version = "unknown"
	commit  = "unknown"
	date    = "unknown"
	builtBy = "unknown"
)

func main() {
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Printf("%s - version %s\n", c.App.Name, version)
		fmt.Printf("  commit: \t%s\n", commit)
		fmt.Printf("  build date: \t%s\n", date)
		fmt.Printf("  build user: \t%s\n", builtBy)
		fmt.Printf("  go version: \t%s\n", runtime.Version())
	}
	flags := []cli.Flag{
		&cli.StringFlag{
			Name:     "config",
			Category: "general",
			Aliases:  []string{"c"},
			EnvVars:  []string{"MN_CONFIG"},
			Usage:    "config file path",
			Required: false,
		},
		altsrc.NewBoolFlag(&cli.BoolFlag{
			Name:     "merge.dry_run",
			Usage:    "only print files info that would be processed",
			Category: "mkvmerge",
		}),
		altsrc.NewStringSliceFlag(&cli.StringSliceFlag{
			Name:        "merge.video_extensions",
			Category:    "mkvmerge",
			DefaultText: "",
			Usage:       "video extensions that should be processed",
			Value:       cli.NewStringSlice("mp4", "avi", "mkv", "mpeg", "mov", "webm"),
		}),
		altsrc.NewStringSliceFlag(&cli.StringSliceFlag{
			Name:        "merge.audio_extensions",
			Category:    "mkvmerge",
			DefaultText: "",
			Usage:       "audio extensions that should be processed",
			Value:       cli.NewStringSlice("ogg", "mka", "wav"),
		}),
		altsrc.NewStringSliceFlag(&cli.StringSliceFlag{
			Name:        "merge.subtitles_extensions",
			Category:    "mkvmerge",
			DefaultText: "",
			Usage:       "subtitles extensions that should be processed",
			Value:       cli.NewStringSlice("ass"),
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:     "sort.tv_template",
			Category: "media sort",
			Usage:    "tv series path template",
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:     "sort.movie_template",
			Category: "media sort",
			Usage:    "movie path template",
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:        "sort.tv_dir",
			Category:    "media sort",
			Usage:       "tv series base directory",
			DefaultText: "defaults to current directory",
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:        "sort.movie_dir",
			Category:    "media sort",
			Usage:       "movie base directory",
			DefaultText: "defaults to current directory",
		}),
		altsrc.NewStringSliceFlag(&cli.StringSliceFlag{
			Name:     "sort.extensions",
			Category: "media sort",
			Usage:    "types of files that should be sorted",
			Value:    cli.NewStringSlice("mp4", "m4v", "avi", "mkv", "mpeg", "mpg", "mov", "webm"),
		}),
		altsrc.NewIntFlag(&cli.IntFlag{
			Name:     "sort.file_limit",
			Category: "media sort",
			Usage:    "maximum number of files to search",
			Value:    1000,
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:     "sort.min_file_size",
			Category: "media sort",
			Usage:    "minimum file size(25MB, 1GB)",
			Value:    "25MB",
			Action: func(_ *cli.Context, s string) error {
				if _, err := sizestr.Parse(s); err != nil {
					return fmt.Errorf("--sort.min_file_size: %q invalid format", s)
				}

				return nil
			},
		}),
		altsrc.NewIntFlag(&cli.IntFlag{
			Name:     "sort.accuracy_threshold",
			Category: "media sort",
			Usage:    "filename match accuracy threshold(0-100)",
			Value:    95,
			Action: func(_ *cli.Context, v int) error {
				if v > 100 {
					return errors.New("--sort.accuracy_threshold must be in range 0-100")
				}
				return nil
			},
		}),
		altsrc.NewIntFlag(&cli.IntFlag{
			Name:     "sort.concurrency",
			Category: "media sort",
			Usage:    "search concurrency [warning] setting this too high can cause rate-limiting errors",
			Value:    6,
		}),
		altsrc.NewBoolFlag(&cli.BoolFlag{
			Name:     "sort.recursive",
			Category: "media sort",
			Usage:    "also search through subdirectories",
		}),
		altsrc.NewBoolFlag(&cli.BoolFlag{
			Name:     "sort.skip_hidden",
			Category: "media sort",
			Usage:    "skip dot files",
		}),
		altsrc.NewBoolFlag(&cli.BoolFlag{
			Name:     "sort.skip_subs",
			Category: "media sort",
			Usage:    "skip subtitles (srt files)",
		}),
		altsrc.NewBoolFlag(&cli.BoolFlag{
			Name:     "sort.dry_run",
			Category: "media sort",
			Usage:    "perform sort but don't actually move any files",
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:     "sort.action",
			Category: "media sort",
			Usage:    "filesystem action used to sort files (copy|link|move)",
			Value:    "move",
		}),
	}
	app := &cli.App{
		Name:    "media-normalizer",
		Usage:   "media-normalizer <folder with media>",
		Version: version,
		Before:  altsrc.InitInputSourceWithContext(flags, altsrc.NewYamlSourceFromFlagFunc("config")),
		Action:  normalize,
		Flags:   flags,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func normalize(c *cli.Context) error {
	target := c.Args().First()
	if c.NArg() == 0 {
		target = "."
	}

	mkvmergeCfg := mkvmerge.Config{
		Target:              target,
		VideoExtensions:     c.StringSlice("merge.video_extensions"),
		AudioExtensions:     c.StringSlice("merge.audio_extensions"),
		SubtitlesExtensions: c.StringSlice("merge.subtitles_extensions"),
		DryRun:              c.Bool("merge.dry_run"),
	}

	log.Println("merging external audio/subtitles")
	if err := mkvmerge.Merge(mkvmergeCfg); err != nil {
		return err
	}

	mediasortCfg := mediasort.Config{
		Targets: []string{target},
		PathConfig: mediasort.PathConfig{
			TVTemplate:    c.String("sort.tv_template"),
			MovieTemplate: c.String("sort.movie_template"),
		},
		TVDir:             c.String("sort.tv_dir"),
		MovieDir:          c.String("sort.movie_dir"),
		Extensions:        strings.Join(c.StringSlice("sort.extensions"), ","),
		FileLimit:         c.Int("sort.file_limit"),
		AccuracyThreshold: c.Int("sort.accuracy_threshold"),
		Recursive:         c.Bool("sort.recursive"),
		SkipHidden:        c.Bool("sort.skip_hidden"),
		SkipSubs:          c.Bool("sort.skip_subs"),
		Concurrency:       c.Int("sort.concurrency"),
		MinFileSize:       sizestr.Bytes((sizestr.MustParse(c.String("sort.min_file_size")))),
		DryRun:            c.Bool("sort.dry_run"),
		Action:            mediasort.Action(c.String("sort.action")),
	}

	log.Println("Sorting media")
	if err := mediasort.FileSystemSort(mediasortCfg); err != nil {
		return err
	}

	return nil
}
