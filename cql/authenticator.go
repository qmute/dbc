package cql

import (
	"fmt"

	"github.com/gocql/gocql"
)

const (
	approvedAliAuth = "org.apache.cassandra.auth.AliLocationAwarePwdAuthenticator"
)

// AliAuthenticator 阿里去cassandra验证
// gocql: unable to create session: control: unable to connect to initial hosts: unexpected authenticator
// "org.apache.cassandra.auth.AliLocationAwarePwdAuthenticator" 没在白名单中
// 自定义阿里去cassandra验证
type AliAuthenticator struct {
	Username string
	Password string
}

func approve(authenticator string) bool {
	return authenticator == approvedAliAuth
}

// Challenge  挑战
func (p AliAuthenticator) Challenge(req []byte) ([]byte, gocql.Authenticator, error) {
	if !approve(string(req)) {
		return nil, nil, fmt.Errorf("not ali %q", req)
	}

	resp := make([]byte, 2+len(p.Username)+len(p.Password))
	resp[0] = 0
	copy(resp[1:], p.Username)
	resp[len(p.Username)+1] = 0
	copy(resp[2+len(p.Username):], p.Password)
	return resp, nil, nil
}

// Success 成功
func (p AliAuthenticator) Success(data []byte) error {
	return nil
}
