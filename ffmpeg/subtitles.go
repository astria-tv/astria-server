package ffmpeg

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	log "github.com/sirupsen/logrus"
)

func NewSubtitleSession(
	stream StreamRepresentation,
	outputDirBase string) (*TranscodingSession, error) {

	outputDir, err := ioutil.TempDir(outputDirBase, "subtitle-session-")
	if err != nil {
		return nil, err
	}

	cmd := exec.Command("ffmpeg",
		"-probesize", "25M", // Especially since we want to not read too much data for remote locations such as Rclone mounts we want to prevent reading too many bytes. Might need adjusting if it misses a bunch of subtitles.
		"-analyzeduration", "200M",
		"-i", buildFfmpegUrlFromFileLocator(stream.Stream.FileLocator),
		"-map", fmt.Sprintf("0:%d", stream.Stream.StreamId),
		"-f", "webvtt",
		"stream0_0.m4s")
	cmd.Stderr, _ = os.Open(os.DevNull)
	cmd.Dir = outputDir

	log.WithFields(log.Fields{"args": cmd.Args, "path": cmd.Dir}).Println("ffmpeg initialized for subtitles")

	return &TranscodingSession{
		cmd:       cmd,
		Stream:    stream,
		OutputDir: outputDir,
	}, nil
}

func GetSubtitleStreamRepresentation(stream Stream) StreamRepresentation {
	return StreamRepresentation{
		Stream: stream,
		Representation: Representation{
			RepresentationId: "webvtt",
		},
	}
}

func GetSubtitleStreamRepresentations(streams []Stream) []StreamRepresentation {
	subtitleRepresentations := []StreamRepresentation{}
	for _, s := range streams {
		subtitleRepresentations = append(subtitleRepresentations,
			GetSubtitleStreamRepresentation(s))
	}
	return subtitleRepresentations
}
