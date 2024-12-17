// Code generated by gen.go using 'go generate'. DO NOT EDIT.

// Copyright 2024 The Ebitengine Authors
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

// This file is intended for precompiled shaders that will be introduced in the future.
// All constant names are underscores and not actually used,
// so they do not affect the binary file size.

package builtinshader

//ebitengine:shader
const _ = "//kage:unit pixels\n\npackage main\n\n\n\n\n\nfunc Fragment(dstPos vec4, srcPos vec2, color vec4) vec4 {\n\n\n\tclr := imageSrc0UnsafeAt(srcPos)\n\n\n\n\n\t// Apply the color scale.\n\tclr *= color\n\n\n\treturn clr\n}\n\n"

//ebitengine:shader
const _ = "//kage:unit pixels\n\npackage main\n\n\nvar ColorMBody mat4\nvar ColorMTranslation vec4\n\n\n\n\nfunc Fragment(dstPos vec4, srcPos vec2, color vec4) vec4 {\n\n\n\tclr := imageSrc0UnsafeAt(srcPos)\n\n\n\n\n\t// Un-premultiply alpha.\n\t// When the alpha is 0, 1-sign(alpha) is 1.0, which means division does nothing.\n\tclr.rgb /= clr.a + (1-sign(clr.a))\n\t// Apply the clr matrix.\n\tclr = (ColorMBody * clr) + ColorMTranslation\n\t// Premultiply alpha\n\tclr.rgb *= clr.a\n\t// Apply the color scale.\n\tclr *= color\n\t// Clamp the output.\n\tclr.rgb = min(clr.rgb, clr.a)\n\n\n\treturn clr\n}\n\n"

//ebitengine:shader
const _ = "//kage:unit pixels\n\npackage main\n\n\n\n\n\nfunc Fragment(dstPos vec4, srcPos vec2, color vec4) vec4 {\n\n\n\tclr := imageSrc0At(srcPos)\n\n\n\n\n\t// Apply the color scale.\n\tclr *= color\n\n\n\treturn clr\n}\n\n"

//ebitengine:shader
const _ = "//kage:unit pixels\n\npackage main\n\n\nvar ColorMBody mat4\nvar ColorMTranslation vec4\n\n\n\n\nfunc Fragment(dstPos vec4, srcPos vec2, color vec4) vec4 {\n\n\n\tclr := imageSrc0At(srcPos)\n\n\n\n\n\t// Un-premultiply alpha.\n\t// When the alpha is 0, 1-sign(alpha) is 1.0, which means division does nothing.\n\tclr.rgb /= clr.a + (1-sign(clr.a))\n\t// Apply the clr matrix.\n\tclr = (ColorMBody * clr) + ColorMTranslation\n\t// Premultiply alpha\n\tclr.rgb *= clr.a\n\t// Apply the color scale.\n\tclr *= color\n\t// Clamp the output.\n\tclr.rgb = min(clr.rgb, clr.a)\n\n\n\treturn clr\n}\n\n"

//ebitengine:shader
const _ = "//kage:unit pixels\n\npackage main\n\n\n\n\nfunc adjustSrcPosForAddressRepeat(p vec2) vec2 {\n\torigin := imageSrc0Origin()\n\tsize := imageSrc0Size()\n\treturn mod(p - origin, size) + origin\n}\n\n\nfunc Fragment(dstPos vec4, srcPos vec2, color vec4) vec4 {\n\n\n\tclr := imageSrc0At(adjustSrcPosForAddressRepeat(srcPos))\n\n\n\n\n\t// Apply the color scale.\n\tclr *= color\n\n\n\treturn clr\n}\n\n"

//ebitengine:shader
const _ = "//kage:unit pixels\n\npackage main\n\n\nvar ColorMBody mat4\nvar ColorMTranslation vec4\n\n\n\nfunc adjustSrcPosForAddressRepeat(p vec2) vec2 {\n\torigin := imageSrc0Origin()\n\tsize := imageSrc0Size()\n\treturn mod(p - origin, size) + origin\n}\n\n\nfunc Fragment(dstPos vec4, srcPos vec2, color vec4) vec4 {\n\n\n\tclr := imageSrc0At(adjustSrcPosForAddressRepeat(srcPos))\n\n\n\n\n\t// Un-premultiply alpha.\n\t// When the alpha is 0, 1-sign(alpha) is 1.0, which means division does nothing.\n\tclr.rgb /= clr.a + (1-sign(clr.a))\n\t// Apply the clr matrix.\n\tclr = (ColorMBody * clr) + ColorMTranslation\n\t// Premultiply alpha\n\tclr.rgb *= clr.a\n\t// Apply the color scale.\n\tclr *= color\n\t// Clamp the output.\n\tclr.rgb = min(clr.rgb, clr.a)\n\n\n\treturn clr\n}\n\n"

//ebitengine:shader
const _ = "//kage:unit pixels\n\npackage main\n\n\n\n\n\nfunc Fragment(dstPos vec4, srcPos vec2, color vec4) vec4 {\n\n\tp0 := srcPos - 1/2.0\n\tp1 := srcPos + 1/2.0\n\n\n\n\n\tc0 := imageSrc0UnsafeAt(p0)\n\tc1 := imageSrc0UnsafeAt(vec2(p1.x, p0.y))\n\tc2 := imageSrc0UnsafeAt(vec2(p0.x, p1.y))\n\tc3 := imageSrc0UnsafeAt(p1)\n\n\n\trate := fract(p1)\n\tclr := mix(mix(c0, c1, rate.x), mix(c2, c3, rate.x), rate.y)\n\n\n\n\t// Apply the color scale.\n\tclr *= color\n\n\n\treturn clr\n}\n\n"

//ebitengine:shader
const _ = "//kage:unit pixels\n\npackage main\n\n\nvar ColorMBody mat4\nvar ColorMTranslation vec4\n\n\n\n\nfunc Fragment(dstPos vec4, srcPos vec2, color vec4) vec4 {\n\n\tp0 := srcPos - 1/2.0\n\tp1 := srcPos + 1/2.0\n\n\n\n\n\tc0 := imageSrc0UnsafeAt(p0)\n\tc1 := imageSrc0UnsafeAt(vec2(p1.x, p0.y))\n\tc2 := imageSrc0UnsafeAt(vec2(p0.x, p1.y))\n\tc3 := imageSrc0UnsafeAt(p1)\n\n\n\trate := fract(p1)\n\tclr := mix(mix(c0, c1, rate.x), mix(c2, c3, rate.x), rate.y)\n\n\n\n\t// Un-premultiply alpha.\n\t// When the alpha is 0, 1-sign(alpha) is 1.0, which means division does nothing.\n\tclr.rgb /= clr.a + (1-sign(clr.a))\n\t// Apply the clr matrix.\n\tclr = (ColorMBody * clr) + ColorMTranslation\n\t// Premultiply alpha\n\tclr.rgb *= clr.a\n\t// Apply the color scale.\n\tclr *= color\n\t// Clamp the output.\n\tclr.rgb = min(clr.rgb, clr.a)\n\n\n\treturn clr\n}\n\n"

//ebitengine:shader
const _ = "//kage:unit pixels\n\npackage main\n\n\n\n\n\nfunc Fragment(dstPos vec4, srcPos vec2, color vec4) vec4 {\n\n\tp0 := srcPos - 1/2.0\n\tp1 := srcPos + 1/2.0\n\n\n\n\n\tc0 := imageSrc0At(p0)\n\tc1 := imageSrc0At(vec2(p1.x, p0.y))\n\tc2 := imageSrc0At(vec2(p0.x, p1.y))\n\tc3 := imageSrc0At(p1)\n\n\n\trate := fract(p1)\n\tclr := mix(mix(c0, c1, rate.x), mix(c2, c3, rate.x), rate.y)\n\n\n\n\t// Apply the color scale.\n\tclr *= color\n\n\n\treturn clr\n}\n\n"

//ebitengine:shader
const _ = "//kage:unit pixels\n\npackage main\n\n\nvar ColorMBody mat4\nvar ColorMTranslation vec4\n\n\n\n\nfunc Fragment(dstPos vec4, srcPos vec2, color vec4) vec4 {\n\n\tp0 := srcPos - 1/2.0\n\tp1 := srcPos + 1/2.0\n\n\n\n\n\tc0 := imageSrc0At(p0)\n\tc1 := imageSrc0At(vec2(p1.x, p0.y))\n\tc2 := imageSrc0At(vec2(p0.x, p1.y))\n\tc3 := imageSrc0At(p1)\n\n\n\trate := fract(p1)\n\tclr := mix(mix(c0, c1, rate.x), mix(c2, c3, rate.x), rate.y)\n\n\n\n\t// Un-premultiply alpha.\n\t// When the alpha is 0, 1-sign(alpha) is 1.0, which means division does nothing.\n\tclr.rgb /= clr.a + (1-sign(clr.a))\n\t// Apply the clr matrix.\n\tclr = (ColorMBody * clr) + ColorMTranslation\n\t// Premultiply alpha\n\tclr.rgb *= clr.a\n\t// Apply the color scale.\n\tclr *= color\n\t// Clamp the output.\n\tclr.rgb = min(clr.rgb, clr.a)\n\n\n\treturn clr\n}\n\n"

//ebitengine:shader
const _ = "//kage:unit pixels\n\npackage main\n\n\n\n\nfunc adjustSrcPosForAddressRepeat(p vec2) vec2 {\n\torigin := imageSrc0Origin()\n\tsize := imageSrc0Size()\n\treturn mod(p - origin, size) + origin\n}\n\n\nfunc Fragment(dstPos vec4, srcPos vec2, color vec4) vec4 {\n\n\tp0 := srcPos - 1/2.0\n\tp1 := srcPos + 1/2.0\n\n\n\tp0 = adjustSrcPosForAddressRepeat(p0)\n\tp1 = adjustSrcPosForAddressRepeat(p1)\n\n\n\n\tc0 := imageSrc0At(p0)\n\tc1 := imageSrc0At(vec2(p1.x, p0.y))\n\tc2 := imageSrc0At(vec2(p0.x, p1.y))\n\tc3 := imageSrc0At(p1)\n\n\n\trate := fract(p1)\n\tclr := mix(mix(c0, c1, rate.x), mix(c2, c3, rate.x), rate.y)\n\n\n\n\t// Apply the color scale.\n\tclr *= color\n\n\n\treturn clr\n}\n\n"

//ebitengine:shader
const _ = "//kage:unit pixels\n\npackage main\n\n\nvar ColorMBody mat4\nvar ColorMTranslation vec4\n\n\n\nfunc adjustSrcPosForAddressRepeat(p vec2) vec2 {\n\torigin := imageSrc0Origin()\n\tsize := imageSrc0Size()\n\treturn mod(p - origin, size) + origin\n}\n\n\nfunc Fragment(dstPos vec4, srcPos vec2, color vec4) vec4 {\n\n\tp0 := srcPos - 1/2.0\n\tp1 := srcPos + 1/2.0\n\n\n\tp0 = adjustSrcPosForAddressRepeat(p0)\n\tp1 = adjustSrcPosForAddressRepeat(p1)\n\n\n\n\tc0 := imageSrc0At(p0)\n\tc1 := imageSrc0At(vec2(p1.x, p0.y))\n\tc2 := imageSrc0At(vec2(p0.x, p1.y))\n\tc3 := imageSrc0At(p1)\n\n\n\trate := fract(p1)\n\tclr := mix(mix(c0, c1, rate.x), mix(c2, c3, rate.x), rate.y)\n\n\n\n\t// Un-premultiply alpha.\n\t// When the alpha is 0, 1-sign(alpha) is 1.0, which means division does nothing.\n\tclr.rgb /= clr.a + (1-sign(clr.a))\n\t// Apply the clr matrix.\n\tclr = (ColorMBody * clr) + ColorMTranslation\n\t// Premultiply alpha\n\tclr.rgb *= clr.a\n\t// Apply the color scale.\n\tclr *= color\n\t// Clamp the output.\n\tclr.rgb = min(clr.rgb, clr.a)\n\n\n\treturn clr\n}\n\n"
