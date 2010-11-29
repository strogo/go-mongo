// Copyright 2010 Gary Burd
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package mongo

import (
	"reflect"
	"strings"
	"strconv"
	"sync"
	"crypto/rand"
	"time"
	"encoding/binary"
)

type DateTime int64

type Timestamp int64

type CodeWithScope struct {
	Code  string
	Scope map[string]interface{}
}

type Regexp struct {
	Pattern string
	Flags   string
}

type ObjectId [12]byte

func NewObjectId() ObjectId {
	// The following format is used for object ids:
	//
	// [0:4] Time since epoch in seconds. This is compatible 
	//       with other drivers.
	// 
	// [4:12] Incrementing counter intialized with crypto random
	//        number. This ensures that object ids are unique, but
	//        is simpler than the format used by other drivers.
	t := time.Seconds()
	c := nextOidCounter()
	return ObjectId{
		byte(t >> 24),
		byte(t >> 16),
		byte(t >> 8),
		byte(t),
		byte(c >> 56),
		byte(c >> 48),
		byte(c >> 40),
		byte(c >> 32),
		byte(c >> 24),
		byte(c >> 16),
		byte(c >> 8),
		byte(c)}
}

var (
	oidLock    sync.Mutex
	oidCounter uint64
)

func nextOidCounter() uint64 {
	oidLock.Lock()
	defer oidLock.Unlock()
	if oidCounter == 0 {
		if err := binary.Read(rand.Reader, binary.BigEndian, &oidCounter); err != nil {
			panic(err)
		}
	}
	oidCounter += 1
	return oidCounter
}

type Symbol string

type Code string

type Doc []struct {
	Key   string
	Value interface{}
}

type Key int

const (
	MaxKey Key = 1
	MinKey Key = -1
)

const (
	kindFloat         = 0x1
	kindString        = 0x2
	kindDocument      = 0x3
	kindArray         = 0x4
	kindBinary        = 0x5
	kindObjectId      = 0x7
	kindBool          = 0x8
	kindDateTime      = 0x9
	kindNull          = 0xA
	kindRegexp        = 0xB
	kindCode          = 0xD
	kindSymbol        = 0xE
	kindCodeWithScope = 0xF
	kindInt32         = 0x10
	kindTimestamp     = 0x11
	kindInt64         = 0x12
	kindMinKey        = 0xff
	kindMaxKey        = 0x7f
)

var kindNames = map[int]string{
	kindFloat:         "float",
	kindString:        "string",
	kindDocument:      "document",
	kindArray:         "array",
	kindBinary:        "binary",
	kindObjectId:      "objectId",
	kindBool:          "bool",
	kindDateTime:      "dateTime",
	kindNull:          "null",
	kindRegexp:        "regexp",
	kindCode:          "code",
	kindSymbol:        "symbol",
	kindCodeWithScope: "codeWithScope",
	kindInt32:         "int32",
	kindTimestamp:     "timestamp",
	kindInt64:         "int64",
	kindMinKey:        "minKey",
	kindMaxKey:        "maxKey",
}

func kindName(kind int) string {
	name, ok := kindNames[kind]
	if !ok {
		name = strconv.Itoa(kind)
	}
	return name
}

func fieldName(f reflect.StructField) string {
	if f.Tag != "" {
		return f.Tag
	}
	if f.PkgPath == "" {
		return strings.ToLower(f.Name)
	}
	return ""
}