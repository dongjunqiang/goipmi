/*
Copyright (c) 2014 EOITek, Inc. All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package ipmi

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEscapedDSN(t *testing.T) {

	s0 := ""
	seg0 := _escapedSplit(s0, ":")
	assert.Equal(t, 1, len(seg0))
	assert.Equal(t, "", seg0[0])

	s0 = "abc"
	seg0 = _escapedSplit(s0, ":")
	assert.Equal(t, 1, len(seg0))
	assert.Equal(t, "abc", seg0[0])

	s0 = ":"
	seg0 = _escapedSplit(s0, ":")
	assert.Equal(t, 2, len(seg0))
	assert.Equal(t, "", seg0[0])
	assert.Equal(t, "", seg0[1])

	s1 := "user:password:host:port"
	seg1 := _escapedSplit(s1, ":")
	assert.Equal(t, 4, len(seg1))
	assert.Equal(t, "user", seg1[0])
	assert.Equal(t, "password", seg1[1])
	assert.Equal(t, "host", seg1[2])
	assert.Equal(t, "port", seg1[3])

	s2 := ":user:password"
	seg2 := _escapedSplit(s2, ":")
	assert.Equal(t, 3, len(seg2))
	assert.Equal(t, "", seg2[0])
	assert.Equal(t, "user", seg2[1])
	assert.Equal(t, "password", seg2[2])

	s3 := "user:password:"
	seg3 := _escapedSplit(s3, ":")
	assert.Equal(t, 3, len(seg3))
	assert.Equal(t, "user", seg3[0])
	assert.Equal(t, "password", seg3[1])
	assert.Equal(t, "", seg3[2])

	//escaped
	s4 := `user:password\:include\::host:port`
	seg4 := _escapedSplit(s4, ":")
	assert.Equal(t, 4, len(seg4))
	assert.Equal(t, "user", seg4[0])
	assert.Equal(t, `password:include:`, seg4[1])
	assert.Equal(t, "host", seg4[2])
	assert.Equal(t, "port", seg4[3])

	s5 := `:user:password\:p:host:port`
	seg5 := _escapedSplit(s5, ":")
	assert.Equal(t, 5, len(seg5))
	assert.Equal(t, "", seg5[0])
	assert.Equal(t, `user`, seg5[1])
	assert.Equal(t, "password:p", seg5[2])
	assert.Equal(t, "host", seg5[3])
	assert.Equal(t, "port", seg5[4])

	s6 := `user:password\:p:host:port:`
	seg6 := _escapedSplit(s6, ":")
	assert.Equal(t, 5, len(seg6))
	assert.Equal(t, `user`, seg6[0])
	assert.Equal(t, "password:p", seg6[1])
	assert.Equal(t, "host", seg6[2])
	assert.Equal(t, "port", seg6[3])
	assert.Equal(t, "", seg6[4])
}

func TestParseDSN(t *testing.T) {
	var err error
	s0 := "abc"
	_, _, _, _, err = ParseDSN(s0)
	assert.NotNil(t, err)

	s1 := "abc@def@g"
	_, _, _, _, err = ParseDSN(s1)
	assert.NotNil(t, err)

	s2 := "abc@def"
	_, _, _, _, err = ParseDSN(s2)
	assert.NotNil(t, err)

	s3 := "user:pwd@def"
	_, _, _, _, err = ParseDSN(s3)
	assert.NotNil(t, err)

	s4 := "user:pwd@def:not port"
	_, _, _, _, err = ParseDSN(s4)
	assert.NotNil(t, err)

	s5 := "user:pwd@host1:624"
	user, pwd, host, port, err := ParseDSN(s5)
	assert.Nil(t, err)
	assert.Equal(t, "user", user)
	assert.Equal(t, "pwd", pwd)
	assert.Equal(t, "host1", host)
	assert.Equal(t, uint16(624), port)
}
