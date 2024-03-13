/*
 * Copyright 2022 Xiongfa Li.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package processor

import (
	"github.com/ydx1011/gopher-core/bean"
	"github.com/ydx1011/yfig"
)

type ValueProcessor struct {
	conf      yfig.Properties
	tagPxName string
	tagName   string
}

type Opt func(processor *ValueProcessor)

func OptSetValueTag(tagPxName, tagName string) Opt {
	return func(processor *ValueProcessor) {
		if tagName != "" {
			if tagPxName == "" {
				tagPxName = yfig.TagPrefixName
			}
			processor.tagName = tagName
			processor.tagPxName = tagPxName
		}
	}
}

func NewValueProcessor(opts ...Opt) *ValueProcessor {
	ret := &ValueProcessor{}
	for _, opt := range opts {
		opt(ret)
	}
	return ret
}

func (p *ValueProcessor) Init(conf yfig.Properties, container bean.Container) error {
	p.conf = conf
	return nil
}

func (p *ValueProcessor) Classify(o interface{}) (bool, error) {
	if p.tagName == "" {
		return true, yfig.Fill(p.conf, o)
	} else {
		// 内部兼容tag 'fig'
		return true, yfig.FillExWithTagNames(p.conf, o, false,
			[]string{
				yfig.TagPrefixName,
				p.tagPxName,
			},
			[]string{
				yfig.TagName,
				p.tagName,
			})
	}
}

func (p *ValueProcessor) Process() error {
	return nil
}

func (p *ValueProcessor) BeanDestroy() error {
	return nil
}
