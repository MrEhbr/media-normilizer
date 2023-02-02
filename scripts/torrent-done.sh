#!/bin/sh
set -xe
dir=$TR_TORRENT_DIR
target=$TR_TORRENT_NAME
if [ -d "$dir/$TR_TORRENT_NAME" ]
then
    dir=$TR_TORRENT_DIR/$TR_TORRENT_NAME
    target="."
fi
cd $dir
/bin/media-normalizer $target
