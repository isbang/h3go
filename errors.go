// Copyright 2022  Il Sub Bang
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package h3go

import "errors"

var (
	ErrCompactDuplicate     = errors.New("compact duplicated")
	ErrCompactLoopExceeded  = errors.New("compact loop exceeded")
	ErrUncompactResExceeded = errors.New("uncompact resolution exceeded")
)
