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
	"errors"
	"fmt"
	"strconv"
)

// parse dsn which is used to describe a connection to an ipmi server
// format:
//   user:password@ip:port
//   use `\` to escape : @  \
func ParseDSN(dsn string) (user, password, host string, port uint16, _err error) {
	segs, err := escapedSplit(dsn, "@")
	if err != nil {
		_err = err
		return
	}

	if len(segs) != 2 {
		_err = errors.New("Invalid segments count using `@` as seperator")
		return
	}

	up, err := escapedSplit(segs[0], ":")
	if len(up) != 2 {
		_err = errors.New("Invalid user and password segments using `:` as seperator")
		return
	}

	user = up[0]
	password = up[1]
	hp, err := escapedSplit(segs[1], ":")
	if len(hp) != 2 {
		_err = errors.New("Invalid host and port segments using `:` as seperator")
		return
	}

	host = hp[0]
	_port, err := strconv.ParseUint(hp[1], 10, 16)
	if err != nil {
		_err = err
	}

	port = uint16(_port)
	return
}

// split string with seperator, only support single charactor seperator
// support escape using `\`
func escapedSplit(s, sep string) ([]string, error) {
	if len(sep) != 1 {
		return nil, errors.New(fmt.Sprintf("sperator must has length of 1, but %s supplied", sep))
	}

	return _escapedSplit(s, sep), nil
}

func _escapedSplit(s, sep string) []string {
	if s == "" {
		return []string{""}
	}

	r := make([]string, 0)
	split1 := func(s, sep string) (sleft, sright string, eos bool) {
		escaping := false
		last := 0
		tmp := ""
		for i := 0; i < len(s); i++ {
			if string(s[i]) == sep && !escaping {
				tmp = tmp + s[last:i]
				if i+1 >= len(s) {
					return tmp, "", false
				} else {
					return tmp, s[i+1:], false
				}
			}

			if escaping {
				tmp = tmp + s[last:i-1]
				last = i
				escaping = false
			}

			if string(s[i]) == "\\" && !escaping {
				escaping = true
			}
		}

		// not found sep
		return s, "", true
	}

	var left string
	var eos bool
	right := s

	for {
		left, right, eos = split1(right, sep)
		r = append(r, left)
		if eos {
			break
		}
	}

	return r
}
