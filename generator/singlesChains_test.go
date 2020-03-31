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

func TestSinglesChains_twiceInAUnit(t *testing.T) {
	g := decode("023k74151409041s0g040w082874283k01020g011s023k04740w083k74042g022g0w0o0101080g3k040w1s02741s020w750o750o3k04080g01041s3k02740w801s028o017s0o043k80043k8o1k02010o1s")
	assert.True(t, g.singlesChains())
	assert.Equal(t, "023k74151409041s0g040w082874283k01020g011s023k04740w083k74042g022g0w0o0101080g3k040w1s02741s020w750o75083k04080g01041s3k02740w801s0288017c0o043k80043k8o0w0201081s", g.encode())
}

func TestSinglesChains_twoColorsElsewhere(t *testing.T) {
	g := decode("01023k080g1w1w740w1w080w1w7401023k0g741w0g3k020w08011w682o040g68020108746874012s081w4g0g02080g02014g744g1w1w1010081u010g741u3k023k1s7404080g0w010g01742q2o3k1w1y08")
	assert.True(t, g.singlesChains())
	assert.Equal(t, "01023k080g1w1w740w1w080w1w7401023k0g741w0g3k020w08011w682o040g4g020108744g74012s081w4g0g02080g02014g744g1w1w1010081u010g741u3k023k1s7404080g0w010g01740y2o3k1w1y08", g.encode())
}
