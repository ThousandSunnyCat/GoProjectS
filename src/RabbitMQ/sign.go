package RabbitMQ

import (
	"crypto/sha256"
	"crypto/sha1"
	"fmt"
	"crypto/md5"
	
)

func (m *Message) Byte(secret string) []byte {
	return []byte(fmt.Sprintf("%s%s%d%s", secret, m.MessageId, m.CreateTime, m.MessageBody))
}

func (m *Message) MD5(secret string) string {
	return fmt.Sprintf("%x", md5.Sum(m.Byte(secret)))
}

func (m *Message) SHA1(secret string) string {
	return fmt.Sprintf("%x", sha1.Sum(m.Byte(secret)))
}

func (m *Message) SHA256(secret string) string {
	return fmt.Sprintf("%x", sha256.Sum256(m.Byte(secret)))
}

func (m *Message) Sign(secret, method string) string {
	switch method {
	case "MD5":
		return m.MD5(secret)
	case "SHA1":
		return m.SHA1(secret)
	case "SHA256":
		return m.SHA256(secret)
	default:
		return ""
	}
}