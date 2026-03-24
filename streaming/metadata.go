package streaming

import (
	"encoding/json"
	"net/http"

	"gitlab.com/olaris/olaris-server/ffmpeg"
)

type metadataResponse struct {
	CheckCodecs []string `json:"checkCodecs"`
}

// serveMetadata generates a list of possible codecs that we could possibly serve and returns
// it to the client so that it can check them in advance of actually making a request for the
// manifest.
//
// We don't just return them all in the manifest because we want tight control over
// what the player will show and we want that control on the server. Otherwise,
// a) would have to pass all information that we use here (e.g.
// which stream is transmuxed) and b) we would need to implement this logic across all clients.
func serveMetadata(w http.ResponseWriter, r *http.Request) {
	fileLocator, statusErr := getFileLocatorOrFail(r)
	if statusErr != nil {
		http.Error(w, statusErr.Error(), statusErr.Status())
		return
	}

	streams, err := ffmpeg.GetStreams(fileLocator)
	if err != nil {
		http.Error(w, "Failed to get streams: "+err.Error(), http.StatusInternalServerError)
		return
	}

	checkCodecs := []string{}

	vStream, err := streams.GetVideoStream()
	if err != nil {
		http.Error(w, "No video stream found: "+err.Error(), http.StatusBadRequest)
		return
	}

	transmuxedVideo := ffmpeg.GetTransmuxedRepresentation(vStream)
	transcodedVideo := ffmpeg.GetSimilarTranscodedRepresentation(vStream)

	checkCodecs = append(checkCodecs,
		transmuxedVideo.Representation.Codecs,
		transcodedVideo.Representation.Codecs)

	lowQualityRepresentations := ffmpeg.GetStandardPresetVideoRepresentations(vStream)
	for _, r := range lowQualityRepresentations {
		checkCodecs = append(checkCodecs, r.Representation.Codecs)
	}

	for _, s := range streams.AudioStreams {
		transmuxedAudio := ffmpeg.GetTransmuxedRepresentation(s)
		transcodedAudio := ffmpeg.GetSimilarTranscodedRepresentation(s)
		lowQualityAudio, _ := ffmpeg.StreamRepresentationFromRepresentationId(
			s, "preset:128k-audio")

		checkCodecs = append(checkCodecs,
			transmuxedAudio.Representation.Codecs,
			transcodedAudio.Representation.Codecs,
			lowQualityAudio.Representation.Codecs)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(metadataResponse{CheckCodecs: checkCodecs})
}
