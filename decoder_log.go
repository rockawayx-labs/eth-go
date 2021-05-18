// Copyright 2021 dfuse Platform Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package eth

import (
	"fmt"
	"io"
)

type LogDecoder struct {
	logEvent     *Log
	topicDecoder *Decoder
	DataDecoder  *Decoder

	topicIndex int
}

func NewLogDecoder(logEvent *Log) *LogDecoder {
	decoder := &LogDecoder{
		logEvent:     logEvent,
		topicDecoder: NewDecoder(nil),
	}

	if len(logEvent.Data) > 0 {
		decoder.DataDecoder = NewDecoder(logEvent.Data)
	}

	return decoder
}

func (d *LogDecoder) ReadTopic() ([]byte, error) {
	if d.topicIndex >= len(d.logEvent.Topics) {
		return nil, io.EOF
	}

	topic := d.logEvent.Topics[d.topicIndex]
	d.topicIndex++

	return topic, nil
}

func (d *LogDecoder) ReadTypedTopic(typeName string) (out interface{}, err error) {
	topic, err := d.ReadTopic()
	if err != nil {
		return nil, fmt.Errorf("read topic: %w", err)
	}

	return d.topicDecoder.SetBytes(topic).Read(typeName)
}

func (d *LogDecoder) ReadData(typeName string) (out interface{}, err error) {
	return d.DataDecoder.Read(typeName)
}
