# media-normalizer

[![Go](https://github.com/MrEhbr/media-normalizer/actions/workflows/go.yml/badge.svg)](https://github.com/MrEhbr/media-normalizer/actions/workflows/go.yml)
[![License](https://img.shields.io/badge/license-Apache--2.0%20%2F%20MIT-%2397ca00.svg)](https://github.com/MrEhbr/media-normalizer/blob/master/COPYRIGHT)
[![GitHub release](https://img.shields.io/github/release/MrEhbr/media-normalizer.svg)](https://github.com/MrEhbr/media-normalizer/releases)
[![codecov](https://codecov.io/gh/MrEhbr/media-normalizer/branch/master/graph/badge.svg)](https://codecov.io/gh/MrEhbr/media-normalizer)
![Made by Aleksei Burmistrov](https://img.shields.io/badge/made%20by-Aleksei%20Burmistrov-blue.svg?style=flat)

This utility was made for serveral purposes:

- Rename files that [Jellyfish](http://Jellyfin.org/) and similar apps could recognize Movie/TV Show with [media-sort](github.com/jpillora/media-sort)
- Merge external audio and subtitles in one file with [remuxing](github.com/elboletaire/remuxing)
- Use as done scipt in transmission-daemon

## Usage

```console
NAME:
   media-normalizer - media-normalizer <folder with media>

USAGE:
   media-normalizer [global options] command [command options] [arguments...]

VERSION:
   v0.0.0

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   general

   --config value, -c value  config file path [$MN_CONFIG]

   media sort

   --sort.accuracy_threshold value                      filename match accuracy threshold(0-100) (default: 95)
   --sort.action value                                  filesystem action used to sort files (copy|link|move) (default: "move")
   --sort.concurrency value                             search concurrency [warning] setting this too high can cause rate-limiting errors (default: 6)
   --sort.dry_run                                       perform sort but don't actually move any files (default: false)
   --sort.extensions value [ --sort.extensions value ]  types of files that should be sorted (default: "mp4", "m4v", "avi", "mkv", "mpeg", "mpg", "mov", "webm")
   --sort.file_limit value                              maximum number of files to search (default: 1000)
   --sort.min_file_size value                           minimum file size(25MB, 1GB) (default: "25MB")
   --sort.movie_dir value                               movie base directory (default: defaults to current directory)
   --sort.movie_template value                          movie path template
   --sort.recursive                                     also search through subdirectories (default: false)
   --sort.skip_hidden                                   skip dot files (default: false)
   --sort.skip_subs                                     skip subtitles (srt files) (default: false)
   --sort.tv_dir value                                  tv series base directory (default: defaults to current directory)
   --sort.tv_template value                             tv series path template

   mkvmerge

   --merge.audio_extensions value [ --merge.audio_extensions value ]          audio extensions that should be processed (default: "ogg", "mka", "wav")
   --merge.dry_run                                                            only print files info that would be processed (default: false)
   --merge.subtitles_extensions value [ --merge.subtitles_extensions value ]  subtitles extensions that should be processed (default: "ass")
   --merge.video_extensions value [ --merge.video_extensions value ]          video extensions that should be processed (default: "mp4", "avi", "mkv", "mpeg", "mov", "webm")
```

### Config example

```yaml
sort:
  movie_dir: /media/movies
  tv_dir: /media/tv_shows
  accuracy_threshold: 80
  recursive: true
  concurrency: 1
  dry_run: true
merge:
  dry_run: true
  video_extensions: [mkv]
  audio_extensions: [mka]
  subtitles_extensions: [ass]
```

### Transmission config example

```json
 ...
 "script-torrent-done-enabled": true,
 "script-torrent-done-filename": "/usr/bin/torrent-done.sh",
 ...
```


## Install

### Using go

```console
go get -u github.com/MrEhbr/media-normalizer/cmd/media-normalizer
```

### Download releases

<https://github.com/MrEhbr/media-normalizer/releases>

## License

© 2022 [Aleksei Burmistrov]

Licensed under the [Apache License, Version 2.0](https://www.apache.org/licenses/LICENSE-2.0) ([`LICENSE`](LICENSE)). See the [`COPYRIGHT`](COPYRIGHT) file for more details.

`SPDX-License-Identifier: Apache-2.0`
