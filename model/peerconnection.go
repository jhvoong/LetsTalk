package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/metaclips/LetsTalk/values"
	"github.com/pion/rtcp"
	"github.com/pion/webrtc/v2"
)

func init() {
	var m = webrtc.MediaEngine{}
	// Setup the codecs you want to use.
	m.RegisterDefaultCodecs()
	m.RegisterCodec(webrtc.NewRTPVP8Codec(webrtc.DefaultPayloadTypeVP8, 90000))
	m.RegisterCodec(webrtc.NewRTPOpusCodec(webrtc.DefaultPayloadTypeOpus, 48000))

	classSessions.api = webrtc.NewAPI(webrtc.WithMediaEngine(m))
}

var (
	// Create a MediaEngine object to configure the supported codec
	classSessions = classSessionPeerConnections{
		publisherVideoTracks:  make(map[string]*webrtc.Track),
		publisherTrackMutexes: &sync.Mutex{},

		audioTracks:       make(map[string][]*webrtc.Track),
		audioTrackMutexes: &sync.Mutex{},

		peerConnection:        make(map[string]*webrtc.PeerConnection),
		peerConnectionMutexes: &sync.Mutex{},

		connectedUsers:      make(map[string][]string),
		connectedUsersMutex: &sync.Mutex{},
	}
)

func (s *classSessionPeerConnections) startClassSession(msg []byte) {
	sessionID := uuid.New().String()
	sdp := sdpConstruct{}

	if err := json.Unmarshal(msg, &sdp); err != nil {
		// Send back a CreateSessionError
		return
	}

	// A single user might login using multiple devices. We close recent peerconnection if there's one.
	s.peerConnectionMutexes.Lock()
	if s.peerConnection[sdp.UserID] != nil {
		// ToDo: Return user already in session error
		return
	}
	s.peerConnectionMutexes.Unlock()

	peerConnection, err := classSessions.api.NewPeerConnection(values.PeerConnectionConfig)
	if err != nil {
		// Send back a CreateSessionError
		return
	}

	s.peerConnectionMutexes.Lock()
	s.peerConnection[sdp.UserID] = peerConnection
	s.peerConnectionMutexes.Unlock()

	sdp.ClassSessionID = sessionID

	s.connectedUsersMutex.Lock()
	s.connectedUsers[sessionID] = []string{sdp.UserID}
	s.connectedUsersMutex.Unlock()

	_, err = peerConnection.AddTransceiverFromKind(webrtc.RTPCodecTypeVideo, webrtc.RtpTransceiverInit{Direction: webrtc.RTPTransceiverDirectionRecvonly})
	if err != nil {
		// Close connection and send back a class session error back to user
		return
	}

	_, err = peerConnection.AddTransceiverFromKind(webrtc.RTPCodecTypeAudio, webrtc.RtpTransceiverInit{Direction: webrtc.RTPTransceiverDirectionRecvonly})
	if err != nil {
		// Close connection and send back a class session error back to user
		return
	}

	peerConnection.OnConnectionStateChange(func(cc webrtc.PeerConnectionState) {
		fmt.Printf("PeerConnection State has changed %s \n", cc.String())
		if cc == webrtc.PeerConnectionStateFailed {
			// Since this is the publisher, all video and audio tracks related to the user
			// should be cleared and all peer connections closed.
			// Note: We should check if audio track is nil BEFORE creating a peerconnection
			// when joining session.
			// ToDo Next: Do not mash up locks. Looks too bad

			s.peerConnectionMutexes.Lock()
			s.connectedUsersMutex.Lock()

			for _, user := range s.connectedUsers[sessionID] {
				closePeerConnection(s.peerConnection[user])

				delete(s.peerConnection, user)
			}

			delete(s.connectedUsers, sessionID)
			s.connectedUsersMutex.Unlock()
			s.peerConnectionMutexes.Unlock()

			s.audioTrackMutexes.Lock()
			delete(s.audioTracks, sessionID)
			s.audioTrackMutexes.Unlock()

			s.publisherTrackMutexes.Lock()
			delete(s.publisherVideoTracks, sessionID)
			s.publisherTrackMutexes.Unlock()
		}
	})

	peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		fmt.Printf("Connection State has changed %s \n", connectionState.String())
	})

	peerConnection.OnSignalingStateChange(func(cc webrtc.SignalingState) {
		fmt.Println("Session singaling", cc.String())
	})

	peerConnection.OnTrack(func(remoteTrack *webrtc.Track, receiver *webrtc.RTPReceiver) {
		// Publisher <-> Server is to receive both audio and video packets.
		// Packets are to be broadcasted to other users on the session.
		// For video, resolution is in 480px.
		// ToDo: How do I confirm this???
		if remoteTrack.PayloadType() == webrtc.DefaultPayloadTypeVP8 || remoteTrack.PayloadType() == webrtc.DefaultPayloadTypeVP9 || remoteTrack.PayloadType() == webrtc.DefaultPayloadTypeH264 {

			videoTrack, err := peerConnection.NewTrack(remoteTrack.PayloadType(), remoteTrack.SSRC(), sessionID, sdp.UserID)
			if err != nil {
				log.Println("unable to generate start session video track", err)
				// Return back a class session creation error back to client.
				peerConnection.Close()
				return
			}

			s.publisherTrackMutexes.Lock()
			s.publisherVideoTracks[sessionID] = videoTrack
			s.publisherTrackMutexes.Unlock()

			// Send a PLI on an interval so that the publisher is pushing a keyframe every rtcpPLIInterval
			go func() {
				ticker := time.NewTicker(values.RtcpPLIInterval)
				for range ticker.C {
					err := peerConnection.WriteRTCP([]rtcp.Packet{&rtcp.PictureLossIndication{MediaSSRC: videoTrack.SSRC()}})
					if err != nil {
						break
					}
				}
			}()

			rtpBuf := make([]byte, 1400)
			for {
				i, err := remoteTrack.Read(rtpBuf)
				if err != nil {
					log.Println("bbb", err)
					break
				}

				_, err = videoTrack.Write(rtpBuf[:i])
				if err != nil && !errors.Is(err, io.ErrClosedPipe) {
					log.Println("publisher video packet writed break", err)
					break
				}
			}

			log.Println("Publisher video track exited")

		} else {
			audioTrack, err := peerConnection.NewTrack(remoteTrack.PayloadType(), remoteTrack.SSRC(), sessionID, sdp.UserID)
			if err != nil {
				log.Println("ccc", err)
				// Return back a class session creation error back to client.
				// ToDo: Convert err==nil and see how peerConnectionState reacts.
				// Also, users might decide to disable video/audio on start.
				peerConnection.Close()
				return
			}

			s.audioTrackMutexes.Lock()
			s.audioTracks[sessionID] = []*webrtc.Track{audioTrack}
			s.audioTrackMutexes.Unlock()

			rtpBuf := make([]byte, 1400)
			for {
				i, err := remoteTrack.Read(rtpBuf)
				if err != nil {
					log.Println("ksks", err)
					break
				}

				_, err = audioTrack.Write(rtpBuf[:i])
				if err != nil && !errors.Is(err, io.ErrClosedPipe) {
					log.Println("publisher video packet writed break", err)
					break
				}
			}

			log.Println("Publisher audio track exited")
		}
	})

	s.negotiate(sdp, peerConnection)

	sdp.MsgType = "ClassSession"
	sdp.AuthorName = values.MapEmailToName[sdp.UserID]

	jsonContent, err := json.Marshal(sdp)
	if err != nil {
		peerConnection.Close()
		return
	}

	// Send back answer SDP to client and also class notification to all users in room.
	roomUsers, err := Message{
		RoomID:   sdp.RoomID,
		Name:     sdp.AuthorName,
		UserID:   sdp.UserID,
		Type:     "classSession",
		FileHash: sdp.ClassSessionID,
	}.SaveMessageContent()

	if err != nil {
		peerConnection.Close()
		return
	}

	for _, user := range roomUsers {

		HubConstruct.sendMessage(jsonContent, user)
	}
}

func (s *classSessionPeerConnections) joinClassSession(msg []byte) {
	sdp := sdpConstruct{}
	if err := json.Unmarshal(msg, &sdp); err != nil {
		// Send back a CreateSessionError
		return
	}

	// A single user might login using multiple devices. We close recent peerconnection if there's one.
	s.peerConnectionMutexes.Lock()
	if s.peerConnection[sdp.UserID] != nil {
		// ToDo: Return user already in session error
		return
	}
	s.peerConnectionMutexes.Unlock()

	peerConnection, err := classSessions.api.NewPeerConnection(values.PeerConnectionConfig)
	if err != nil {
		// Send back a CreateSessionError
		peerConnection.Close()
		return
	}

	// ToDo: Since this is a single video, multiple audio. We need to call a send only.
	// Pion currently does not support sendOnly so there would only be support for chrome.
	_, err = peerConnection.AddTransceiverFromKind(webrtc.RTPCodecTypeAudio, webrtc.RtpTransceiverInit{Direction: webrtc.RTPTransceiverDirectionRecvonly})
	if err != nil {
		// Close connection and send back a class session error back to user
		peerConnection.Close()
		return
	}

	peerConnection.OnConnectionStateChange(func(cc webrtc.PeerConnectionState) {
		fmt.Printf("PeerConnection State has changed for joined class %s \n", cc.String())

		if cc == webrtc.PeerConnectionStateFailed {
			s.peerConnectionMutexes.Lock()
			closePeerConnection(peerConnection)
			delete(s.peerConnection, sdp.UserID)
			s.peerConnectionMutexes.Unlock()

			// Remove user from connected users list.
			s.connectedUsersMutex.Lock()
			connectedUsers := s.connectedUsers[sdp.ClassSessionID]
			for i := range connectedUsers {
				if connectedUsers[i] == sdp.UserID {
					if len(connectedUsers) > i+1 {
						connectedUsers = append(connectedUsers[:i], connectedUsers[i+1:]...)
					} else {
						connectedUsers = connectedUsers[:i]
					}
					break
				}
			}
			s.connectedUsers[sdp.ClassSessionID] = connectedUsers
			s.connectedUsersMutex.Unlock()

			s.audioTrackMutexes.Lock()
			audioTracks := s.audioTracks[sdp.ClassSessionID]

			for i, audioTrack := range audioTracks {
				if audioTrack != nil && audioTrack.Label() == sdp.UserID {
					if len(audioTracks) > i+1 {
						audioTracks = append(audioTracks[:i], audioTracks[i+1:]...)
					} else {
						audioTracks = audioTracks[:i]
					}
					break
				}
			}

			s.audioTracks[sdp.ClassSessionID] = audioTracks
			s.audioTrackMutexes.Unlock()
		}
	})

	peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		fmt.Printf("Join class session Connection State has changed %s \n", connectionState.String())
	})

	peerConnection.OnSignalingStateChange(func(cc webrtc.SignalingState) {
		fmt.Println("join class session singaling", cc.String())
	})

	peerConnection.OnTrack(func(remoteTrack *webrtc.Track, receiver *webrtc.RTPReceiver) {
		audioTrack, err := peerConnection.NewTrack(remoteTrack.PayloadType(), remoteTrack.SSRC(), sdp.ClassSessionID, sdp.UserID)
		if err != nil {
			// ToDo: Close peerconnection and send error message to user.
			peerConnection.Close()
			return
		}

		// Confirm both video and audio track from publisher are both enabled and publisher is still up.
		s.peerConnectionMutexes.Lock()

		publisherPeerConnection := s.peerConnection[sdp.Author]
		if publisherPeerConnection == nil {
			// Send back a JoinSessionError. Class session is closed.
			closePeerConnection(peerConnection)

			s.peerConnectionMutexes.Unlock()
			return
		}

		s.peerConnectionMutexes.Unlock()

		s.publisherTrackMutexes.Lock()
		if s.publisherVideoTracks[sdp.ClassSessionID] == nil {
			closePeerConnection(peerConnection)

			s.publisherTrackMutexes.Unlock()
			return
		}

		_, err = peerConnection.AddTrack(s.publisherVideoTracks[sdp.ClassSessionID])
		if err != nil {
			closePeerConnection(peerConnection)

			s.publisherTrackMutexes.Unlock()
			return
		}
		s.publisherTrackMutexes.Unlock()

		// Send audio track to other session users.
		s.peerConnectionMutexes.Lock()
		s.connectedUsersMutex.Lock()

		for _, user := range s.connectedUsers[sdp.ClassSessionID] {
			pc := s.peerConnection[user]
			if pc != nil {
				if _, err := pc.AddTransceiverFromTrack(audioTrack); err != nil {
					log.Println("Could not add track to other users", err)
				}
			}
		}

		s.connectedUsers[sdp.ClassSessionID] = append(s.connectedUsers[sdp.ClassSessionID], sdp.UserID)
		s.peerConnection[sdp.UserID] = peerConnection

		s.connectedUsersMutex.Unlock()
		s.peerConnectionMutexes.Unlock()

		// Add other users audio tracks.
		s.audioTrackMutexes.Lock()

		for _, track := range s.audioTracks[sdp.ClassSessionID] {
			if track != nil {
				_, err := peerConnection.AddTransceiverFromTrack(track)
				if err != nil {
					log.Println("error adding audio track in join session", err)
				}

			} else {
				log.Println("failed track here")
			}
		}

		s.audioTracks[sdp.ClassSessionID] = append(s.audioTracks[sdp.ClassSessionID], audioTrack)

		s.audioTrackMutexes.Unlock()

		if err = s.sendRenegotiateOffer(sdp.ClassSessionID); err != nil {
			log.Println("Failed to send renegotiation offer, closing now", err)
			closePeerConnection(peerConnection)
		}

		rtpBuf := make([]byte, 1400)
		for {
			i, err := remoteTrack.Read(rtpBuf)
			if err != nil {
				log.Println("error reading remote track at join session", err)
				break
			}

			_, err = audioTrack.Write(rtpBuf[:i])
			if err != nil && !errors.Is(err, io.ErrClosedPipe) {
				log.Println("publisher video packet writed break", err)
				break
			}
		}

		log.Println("subscriber audio track exited")
	})

	// Negotiate SDP.
	if err = s.negotiate(sdp, peerConnection); err != nil {
		// send error message to user
		log.Println("Unable to negotiate on join class session", err)
	}
}

func (s *classSessionPeerConnections) negotiate(sdp sdpConstruct, peerConnection *webrtc.PeerConnection) error {
	if peerConnection == nil {
		s.peerConnectionMutexes.Lock()

		if peerConnection = s.peerConnection[sdp.UserID]; peerConnection == nil {
			return values.ErrPeerConnectionNotFound
		}

		s.peerConnectionMutexes.Unlock()
	}

	err := peerConnection.SetRemoteDescription(
		webrtc.SessionDescription{
			SDP:  sdp.SDP,
			Type: webrtc.SDPTypeOffer,
		})

	if err != nil {
		// ToDo: Do we need to do anything further??
		return err
	}

	answer, err := peerConnection.CreateAnswer(nil)
	if err != nil {
		// ToDo: Do we need to do anything further??
		return err
	}

	if err := peerConnection.SetLocalDescription(answer); err != nil {
		return err
	}

	user := sdp.UserID

	sdp = sdpConstruct{}
	sdp.SDP = answer.SDP
	sdp.MsgType = "Negotiate"

	jsonContent, err := json.Marshal(sdp)
	if err != nil {
		return err
	}

	HubConstruct.sendMessage(jsonContent, user)

	return nil
}

func (s *classSessionPeerConnections) sendRenegotiateOffer(session string) error {
	data := struct {
		MsgType   string `json:"msgType"`
		SessionID string `json:"sessionID"`
		// UserID    string // ToDo: UserID was removed so as to parse once. Revisit.
	}{
		"Renegotiate",
		session,
	}

	jsonContent, err := json.Marshal(data)
	if err != nil {
		return err
	}

	s.connectedUsersMutex.Lock()
	for _, user := range s.connectedUsers[session] {
		HubConstruct.sendMessage(jsonContent, user)
	}
	s.connectedUsersMutex.Unlock()

	return nil
}

func closePeerConnection(pc *webrtc.PeerConnection) {
	if pc == nil || pc.ConnectionState() == webrtc.PeerConnectionStateClosed {
		return
	}

	pc.Close()
}
