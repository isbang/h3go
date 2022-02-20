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

// Direction is H3 digit representing ijk+ axes direction.
// Values will be within the lowest 3 bits of an integer.
type Direction uint

const (
	// H3 digit in center
	CENTER_DIGIT Direction = 0

	// H3 digit in k-axes direction
	K_AXES_DIGIT Direction = 1

	// H3 digit in j-axes direction
	J_AXES_DIGIT Direction = 2

	// H3 digit in j == k direction
	JK_AXES_DIGIT Direction = J_AXES_DIGIT | K_AXES_DIGIT /* 3 */

	// H3 digit in i-axes direction
	I_AXES_DIGIT Direction = 4

	// H3 digit in i == k direction
	IK_AXES_DIGIT Direction = I_AXES_DIGIT | K_AXES_DIGIT /* 5 */

	// H3 digit in i == j direction
	IJ_AXES_DIGIT Direction = I_AXES_DIGIT | J_AXES_DIGIT /* 6 */

	// H3 digit in the invalid direction
	INVALID_DIGIT Direction = 7
)

// Valid digits will be less than this value. Same value as INVALID_DIGIT.
const NUM_DIGITS = int(INVALID_DIGIT)

// rotate60ccw Rotates indexing digit 60 degrees counter-clockwise.
func (digit Direction) rotate60ccw() Direction {
	switch digit {
	case K_AXES_DIGIT:
		return IK_AXES_DIGIT
	case IK_AXES_DIGIT:
		return I_AXES_DIGIT
	case I_AXES_DIGIT:
		return IJ_AXES_DIGIT
	case IJ_AXES_DIGIT:
		return J_AXES_DIGIT
	case J_AXES_DIGIT:
		return JK_AXES_DIGIT
	case JK_AXES_DIGIT:
		return K_AXES_DIGIT
	default:
		return digit
	}
}

// rotate60cw rotates indexing digit 60 degrees clockwise.
func (digit Direction) rotate60cw() Direction {
	switch digit {
	case K_AXES_DIGIT:
		return JK_AXES_DIGIT
	case JK_AXES_DIGIT:
		return J_AXES_DIGIT
	case J_AXES_DIGIT:
		return IJ_AXES_DIGIT
	case IJ_AXES_DIGIT:
		return I_AXES_DIGIT
	case I_AXES_DIGIT:
		return IK_AXES_DIGIT
	case IK_AXES_DIGIT:
		return K_AXES_DIGIT
	default:
		return digit
	}
}

// _rotate60ccw Rotates indexing digit 60 degrees counter-clockwise.
//
// Deprecated: Use (Direction).rotate60ccw instead
func _rotate60ccw(digit Direction) Direction {
	return digit.rotate60ccw()
}

// _rotate60cw rotates indexing digit 60 degrees clockwise.
//
// Deprecated: Use (Direction).rotate60cw instead
func _rotate60cw(digit Direction) Direction {
	return digit.rotate60cw()
}
