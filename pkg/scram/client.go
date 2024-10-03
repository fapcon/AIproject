package scram

import (
	"crypto/sha256"
	"crypto/sha512"

	"github.com/IBM/sarama"
	xscram "github.com/xdg/scram"
)

func NewClient(hashFunc xscram.HashGeneratorFcn) sarama.SCRAMClient {
	return &XDGSCRAMClient{HashGeneratorFcn: hashFunc} //nolint:exhaustruct // ok to fill not all fields
}

func NewSHA512Client() sarama.SCRAMClient {
	return NewClient(sha512.New)
}

func NewSHA256Client() sarama.SCRAMClient {
	return NewClient(sha256.New)
}

type XDGSCRAMClient struct {
	*xscram.Client
	*xscram.ClientConversation
	xscram.HashGeneratorFcn
}

func (x *XDGSCRAMClient) Begin(userName, password, authzID string) error {
	var err error
	x.Client, err = x.HashGeneratorFcn.NewClient(userName, password, authzID)
	if err != nil {
		return err
	}
	x.ClientConversation = x.Client.NewConversation()
	return nil
}

func (x *XDGSCRAMClient) Step(challenge string) (string, error) {
	return x.ClientConversation.Step(challenge)
}

func (x *XDGSCRAMClient) Done() bool {
	return x.ClientConversation.Done()
}
