/*
 * Copyright © 2020, G.Ralph Kuntz, MD.
 *
 * Licensed under the Apache License, Version 2.0(the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIC
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package generator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNakedSingleBox(t *testing.T) {
	g := decode("02e7e7e71se7e7043ke7e7e7e7e70we71se704e7e7e708e70we7e7e7e73ke702e71se7e701e7e7e7e7e7e7e70we7e71se704e708e7e7e7e708e73ke7e7e774e70we708e7e7e7e7e77401e7e70we7e7e702")
	assert.True(t, g.nakedSingleGroup(&box))
	assert.Equal(t, "02e1e1bb1sbb7v043ke1e1e1bbbb0w7v1s7v04e1e1bb08bb0w7v7v8u8u3ke102e11sbbbb018u8ue1e1e1bbbb0w8u8u1se104e108bbbb5y5y089j3k9j7171745y0w5y089j9j71717174015y9j0w9j717102", g.encode())
}

func TestNakedSingleCol(t *testing.T) {
	g := decode("02e1e1bb1sbb7v043ke1e1e1bbbb0w7v1s7v04e1e1bb08bb0w7v7v8u8u3ke102e11sbbbb018u8ue1e1e1bbbb0w8u8u1se104e108bbbb5y5y089j3k9j7171745y0w5y089j9j71717174015y9j0w9j717102")
	assert.True(t, g.nakedSingleGroup(&col))
	assert.Equal(t, "02d48hbb1sbb7n043k6wd48hbb7l0w7n1s0p04d48hbb08bb0w7v0p1k7y3kdt02d51sb70l017y8mdt7ld5bbb70w1k7y1sdt04d508b70l5s5y089j3k9j4555745s0w0m087l9j45552l74010m9j0w9j455502", g.encode())
}

func TestNakedSingleRow(t *testing.T) {
	g := decode("02d48hbb1sbb7n043k6wd48hbb7l0w7n1s0p04d48hbb08bb0w7v0p1k7y3kdt02d51sb70l017y8mdt7ld5bbb70w1k7y1sdt04d508b70l5s5y089j3k9j4555745s0w0m087l9j45552l74010m9j0w9j455502")
	assert.True(t, g.nakedSingleGroup(&row))
	assert.Equal(t, "027s8h7l1s7l7l043k48bc7lbb7l0w7n1s0p04cw7lb708b70w7n0h1k7w3k8h027t1s7l0l017y7qcw7kd4bab60w1c7m1sc104b508b70h282e082f3k2f0l1d745s0w0m087l9j45412d74010k2c0w2c444802", g.encode())
}
