package webrtcHelper

import (
	"github.com/pion/webrtc/v3"
)

var api *webrtc.API

func init() {
	// Optional: Customize media engine if necessário (H264, VP8, etc)
	m := &webrtc.MediaEngine{}
	m.RegisterCodec(webrtc.RTPCodecParameters{
		RTPCodecCapability: webrtc.RTPCodecCapability{
			MimeType:    webrtc.MimeTypeH264,
			ClockRate:   90000,
			Channels:    0,
			SDPFmtpLine: "level-asymmetry-allowed=1;packetization-mode=1;profile-level-id=42e01f",
		},
		PayloadType: 102,
	}, webrtc.RTPCodecTypeVideo)

	api = webrtc.NewAPI(webrtc.WithMediaEngine(m))
}

func CreatePeerConnection() (*webrtc.PeerConnection, *webrtc.TrackLocalStaticRTP, error) {
	peerConnection, err := api.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		return nil, nil, err
	}

	videoTrack, err := webrtc.NewTrackLocalStaticRTP(
		webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeH264},
		"video", "pion",
	)
	if err != nil {
		return nil, nil, err
	}

	_, err = peerConnection.AddTrack(videoTrack)
	if err != nil {
		return nil, nil, err
	}

	return peerConnection, videoTrack, nil
}
