package gotik

import (
	"bytes"
	"fmt"

	"github.com/jjcinaz/gotik/proto"
)

// Reply has all the sentences from a reply.
type Reply struct {
	Re   []*proto.Sentence
	Done *proto.Sentence
}

func (r *Reply) String() string {
	b := &bytes.Buffer{}
	for _, re := range r.Re {
		_, _ = fmt.Fprintf(b, "%s\n", re)
	}
	_, _ = fmt.Fprintf(b, "%s", r.Done)
	return b.String()
}

// readReply reads one reply synchronously. It returns the reply.
func (c *Client) readReply() (*Reply, error) {
	var lastErr error
	r := &Reply{}
	for {
		var (
			err  error
			sen  *proto.Sentence
			done bool
		)
		sen, err = c.r.ReadSentence()
		if err != nil {
			return nil, err
		}
		done, err = r.processSentence(sen)
		if err != nil {
			lastErr = err
		}
		if done {
			if lastErr != nil {
				return nil, lastErr
			}
			return r, nil
		}
	}
}

func (r *Reply) processSentence(sen *proto.Sentence) (bool, error) {
	switch sen.Word {
	case "!re":
		r.Re = append(r.Re, sen)
	case "!done":
		r.Done = sen
		return true, nil
	case "!empty":
		// New word added to ROS 7.18; make sure we return an empty reply slice
		// Note that "!done" will follow "!empty"
		r.Re = make([]*proto.Sentence, 0)
	case "!trap":
		return false, &DeviceError{sen}
	case "!fatal":
		return true, &DeviceError{sen}
	case "":
		// API docs say that empty sentences should be ignored
	default:
		return true, &UnknownReplyError{sen}
	}
	return false, nil
}
