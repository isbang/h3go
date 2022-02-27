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

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

// _ipow does integer exponentiation efficiently. Taken from StackOverflow.
//
// Return the exponentiated value
func _ipow(base, exp int) int {
	result := 1
	for exp > 0 {
		if exp&1 > 0 {
			result *= base
		}
		exp >>= 1
		base *= base
	}

	return result
}
