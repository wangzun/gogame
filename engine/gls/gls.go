// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gls

import (
	"bytes"
	"encoding/gob"
	"math"
	"unsafe"

	"golang.org/x/mobile/gl"
)

// GLS encapsulates the state of an OpenGL context and contains
// methods to call OpenGL functions.
type GLS struct {
	context             gl.Context
	stats               Stats             // statistics
	prog                *Program          // current active shader program
	programs            map[*Program]bool // shader programs cache
	checkErrors         bool              // check openGL API errors flag
	activeTexture       uint32            // cached last set active texture unit
	viewportX           int32             // cached last set viewport x
	viewportY           int32             // cached last set viewport y
	viewportWidth       int32             // cached last set viewport width
	viewportHeight      int32             // cached last set viewport height
	lineWidth           float32           // cached last set line width
	sideView            int64             // cached last set triangle side view mode
	frontFace           uint32            // cached last set glFrontFace value
	depthFunc           uint32            // cached last set depth function
	depthMask           int64             // cached last set depth mask
	capabilities        map[int]int       // cached capabilities (Enable/Disable)
	blendEquation       uint32            // cached last set blend equation value
	blendSrc            uint32            // cached last set blend src value
	blendDst            uint32            // cached last set blend equation destination value
	blendEquationRGB    uint32            // cached last set blend equation rgb value
	blendEquationAlpha  uint32            // cached last set blend equation alpha value
	blendSrcRGB         uint32            // cached last set blend src rgb
	blendSrcAlpha       uint32            // cached last set blend src alpha value
	blendDstRGB         uint32            // cached last set blend destination rgb value
	blendDstAlpha       uint32            // cached last set blend destination alpha value
	polygonModeFace     uint32            // cached last set polygon mode face
	polygonModeMode     uint32            // cached last set polygon mode mode
	polygonOffsetFactor float32           // cached last set polygon offset factor
	polygonOffsetUnits  float32           // cached last set polygon offset units
	gobuf               []byte            // conversion buffer with GO memory
	cbuf                []byte            // conversion buffer with C memory
}

// Stats contains counters of OpenGL resources being used as well
// the cumulative numbers of some OpenGL calls for performance evaluation.
type Stats struct {
	Shaders    int    // Current number of shader programs
	Vaos       int    // Number of Vertex Array Objects
	Buffers    int    // Number of Buffer Objects
	Textures   int    // Number of Textures
	Caphits    uint64 // Cumulative number of hits for Enable/Disable
	UnilocHits uint64 // Cumulative number of uniform location cache hits
	UnilocMiss uint64 // Cumulative number of uniform location cache misses
	Unisets    uint64 // Cumulative number of uniform sets
	Drawcalls  uint64 // Cumulative number of draw calls
}

// Polygon side view.
const (
	FrontSide = iota + 1
	BackSide
	DoubleSide
)

const (
	capUndef    = 0
	capDisabled = 1
	capEnabled  = 2
	uintUndef   = math.MaxUint32
	intFalse    = 0
	intTrue     = 1
)

const (
	FloatSize = int32(unsafe.Sizeof(float32(0)))
)

// New creates and returns a new instance of a GLS object,
// which encapsulates the state of an OpenGL context.
// This should be called only after an active OpenGL context
// is established, such as by creating a new window.
func New(context gl.Context) (*GLS, error) {

	gs := new(GLS)
	gs.reset()
	gs.SetContext(context)
	gs.setDefaultState()
	gs.checkErrors = true

	// Preallocate conversion buffers
	size := 1 * 1024
	gs.gobuf = make([]byte, size)
	return gs, nil
}

// SetCheckErrors enables/disables checking for errors after the
// call of any OpenGL function. It is enabled by default but
// could be disabled after an application is stable to improve the performance.
func (gs *GLS) SetCheckErrors(enable bool) {
	gs.checkErrors = enable
}

// CheckErrors returns if error checking is enabled or not.
func (gs *GLS) CheckErrors() bool {
	return gs.checkErrors
}

func (gs *GLS) SetContext(context gl.Context) {
	gs.context = context
}

// reset resets the internal state kept of the OpenGL
func (gs *GLS) reset() {

	gs.lineWidth = 0.0
	gs.sideView = uintUndef
	gs.frontFace = 0
	gs.depthFunc = 0
	gs.depthMask = uintUndef
	gs.capabilities = make(map[int]int)
	gs.programs = make(map[*Program]bool)
	gs.prog = nil

	gs.activeTexture = uintUndef
	gs.blendEquation = uintUndef
	gs.blendSrc = uintUndef
	gs.blendDst = uintUndef
	gs.blendEquationRGB = 0
	gs.blendEquationAlpha = 0
	gs.blendSrcRGB = uintUndef
	gs.blendSrcAlpha = uintUndef
	gs.blendDstRGB = uintUndef
	gs.blendDstAlpha = uintUndef
	gs.polygonModeFace = 0
	gs.polygonModeMode = 0
	gs.polygonOffsetFactor = -1
	gs.polygonOffsetUnits = -1
}

// setDefaultState is used internally to set the initial state of OpenGL
// for this context.
func (gs *GLS) setDefaultState() {

	gs.context.ClearColor(0, 0, 0, 1)
	gs.context.ClearDepthf(1)
	gs.context.ClearStencil(0)

	gs.Enable(DEPTH_TEST)
	gs.DepthFunc(LEQUAL)
	gs.FrontFace(CCW)
	gs.CullFace(BACK)
	gs.Enable(CULL_FACE)
	gs.Enable(BLEND)
	gs.BlendEquation(FUNC_ADD)
	gs.BlendFunc(SRC_ALPHA, ONE_MINUS_SRC_ALPHA)
	gs.Enable(VERTEX_PROGRAM_POINT_SIZE)
	gs.Enable(PROGRAM_POINT_SIZE)
	gs.Enable(MULTISAMPLE)
	gs.Enable(POLYGON_OFFSET_FILL)
	gs.Enable(POLYGON_OFFSET_LINE)
	gs.Enable(POLYGON_OFFSET_POINT)
}

// Stats copy the current values of the internal statistics structure
// to the specified pointer.
func (gs *GLS) Stats(s *Stats) {

	*s = gs.stats
	s.Shaders = len(gs.programs)
}

// ActiveTexture selects which texture unit subsequent texture state calls
// will affect. The number of texture units an implementation supports is
// implementation dependent, but must be at least 48 in GL 3.3.
func (gs *GLS) ActiveTexture(texture uint32) {

	if gs.activeTexture == texture {
		return
	}
	gs.context.ActiveTexture(gl.Enum(texture))
	gs.activeTexture = texture
}

// AttachShader attaches the specified shader object to the specified program object.
// func (gs *GLS) AttachShader(program , shader uint32) {

// 	 C.glAttachShader(C.GLuint(program), C.GLuint(shader))
// }

func (gs *GLS) AttachShader(program gl.Program, shader gl.Shader) {

	gs.context.AttachShader(program, shader)
}

// BindBuffer binds a buffer object to the specified buffer binding point.
// func (gs *GLS) BindBuffer(target int, vbo uint32) {

// 	C.glBindBuffer(C.GLenum(target), C.GLuint(vbo))
// }

func (gs *GLS) BindBuffer(target gl.Enum, vbo gl.Buffer) {
	gs.context.BindBuffer(target, vbo)
}

// BindTexture lets you create or use a named texture.
func (gs *GLS) BindTexture(target gl.Enum, tex gl.Texture) {
	gs.context.BindTexture(target, tex)
}

// BindVertexArray binds the vertex array object.
func (gs *GLS) BindVertexArray(vao gl.VertexArray) {
	gs.context.BindVertexArray(vao)
}

// BlendEquation sets the blend equations for all draw buffers.
func (gs *GLS) BlendEquation(mode uint32) {

	if gs.blendEquation == mode {
		return
	}
	gs.context.BlendEquation(gl.Enum(mode))
	gs.blendEquation = mode
}

// BlendEquationSeparate sets the blend equations for all draw buffers
// allowing different equations for the RGB and alpha components.
func (gs *GLS) BlendEquationSeparate(modeRGB uint32, modeAlpha uint32) {

	if gs.blendEquationRGB == modeRGB && gs.blendEquationAlpha == modeAlpha {
		return
	}
	gs.context.BlendEquationSeparate(gl.Enum(modeRGB), gl.Enum(modeAlpha))
	gs.blendEquationRGB = modeRGB
	gs.blendEquationAlpha = modeAlpha
}

// BlendFunc defines the operation of blending for
// all draw buffers when blending is enabled.
func (gs *GLS) BlendFunc(sfactor, dfactor uint32) {

	if gs.blendSrc == sfactor && gs.blendDst == dfactor {
		return
	}
	gs.context.BlendFunc(gl.Enum(sfactor), gl.Enum(dfactor))
	gs.blendSrc = sfactor
	gs.blendDst = dfactor
}

// BlendFuncSeparate defines the operation of blending for all draw buffers when blending
// is enabled, allowing different operations for the RGB and alpha components.
func (gs *GLS) BlendFuncSeparate(srcRGB uint32, dstRGB uint32, srcAlpha uint32, dstAlpha uint32) {

	if gs.blendSrcRGB == srcRGB && gs.blendDstRGB == dstRGB &&
		gs.blendSrcAlpha == srcAlpha && gs.blendDstAlpha == dstAlpha {
		return
	}
	gs.context.BlendFuncSeparate(gl.Enum(srcRGB), gl.Enum(dstRGB), gl.Enum(srcAlpha), gl.Enum(dstAlpha))
	gs.blendSrcRGB = srcRGB
	gs.blendDstRGB = dstRGB
	gs.blendSrcAlpha = srcAlpha
	gs.blendDstAlpha = dstAlpha
}

// BufferData creates a new data store for the buffer object currently
// bound to target, deleting any pre-existing data store.

// func (gs *GLS) BufferData(target uint32, size int, data interface{}, usage uint32) {

// 	C.glBufferData(C.GLenum(target), C.GLsizeiptr(size), ptr(data), C.GLenum(usage))
// }

func (gs *GLS) BufferData(target uint32, size int, data []byte, usage uint32) {
	gs.context.BufferData(gl.Enum(target), data, gl.Enum(usage))
}

// ClearColor specifies the red, green, blue, and alpha values
// used by glClear to clear the color buffers.
func (gs *GLS) ClearColor(r, g, b, a float32) {
	gs.context.ClearColor(r, g, b, a)
}

// Clear sets the bitplane area of the window to values previously
// selected by ClearColor, ClearDepth, and ClearStencil.
func (gs *GLS) Clear(mask uint) {

	gs.context.Clear(gl.Enum(mask))

}

// CompileShader compiles the source code strings that
// have been stored in the specified shader object.
func (gs *GLS) CompileShader(shader gl.Shader) {

	gs.context.CompileShader(shader)
}

// CreateProgram creates an empty program object and returns
// a non-zero value by which it can be referenced.
func (gs *GLS) CreateProgram() gl.Program {

	p := gs.context.CreateProgram()

	return p
}

// CreateShader creates an empty shader object and returns
// a non-zero value by which it can be referenced.
func (gs *GLS) CreateShader(stype uint32) gl.Shader {

	h := gs.context.CreateShader(gl.Enum(stype))
	return h
}

// DeleteBuffers deletes n​buffer objects named
// by the elements of the provided array.
func (gs *GLS) DeleteBuffers(bufs ...gl.Buffer) {

	for _, buf := range bufs {
		gs.context.DeleteBuffer(buf)
	}
	gs.stats.Buffers -= len(bufs)
}

// DeleteShader frees the memory and invalidates the name
// associated with the specified shader object.
func (gs *GLS) DeleteShader(shader gl.Shader) {
	gs.context.DeleteShader(shader)
}

// DeleteProgram frees the memory and invalidates the name
// associated with the specified program object.
func (gs *GLS) DeleteProgram(program gl.Program) {
	gs.context.DeleteProgram(program)

}

// DeleteTextures deletes n​textures named
// by the elements of the provided array.
func (gs *GLS) DeleteTextures(tex ...gl.Texture) {
	for _, v := range tex {
		gs.context.DeleteTexture(v)
	}

	gs.stats.Textures -= len(tex)
}

// DeleteVertexArrays deletes n​vertex array objects named
// by the elements of the provided array.
func (gs *GLS) DeleteVertexArrays(vaos ...gl.VertexArray) {
	for _, vao := range vaos {
		gs.context.DeleteVertexArray(vao)
	}
	gs.stats.Vaos -= len(vaos)
}

// DepthFunc specifies the function used to compare each incoming pixel
// depth value with the depth value present in the depth buffer.
func (gs *GLS) DepthFunc(mode uint32) {

	if gs.depthFunc == mode {
		return
	}
	gs.context.DepthFunc(gl.Enum(mode))
	gs.depthFunc = mode
}

// DepthMask enables or disables writing into the depth buffer.
func (gs *GLS) DepthMask(flag bool) {

	if gs.depthMask == intTrue && flag {
		return
	}
	if gs.depthMask == intFalse && !flag {
		return
	}
	gs.context.DepthMask(flag)
	if flag {
		gs.depthMask = intTrue
	} else {
		gs.depthMask = intFalse
	}
}

// DrawArrays renders primitives from array data.
func (gs *GLS) DrawArrays(mode uint32, first int32, count int32) {

	gs.context.DrawArrays(gl.Enum(mode), int(first), int(count))

	gs.stats.Drawcalls++
}

// DrawBuffer specifies which color buffers are to be drawn into.
// func (gs *GLS) DrawBuffer(mode uint32) {

// 	C.glDrawBuffer(C.GLenum(mode))
// }

// DrawElements renders primitives from array data.
func (gs *GLS) DrawElements(mode uint32, count int32, itype uint32, start uint32) {
	gs.context.DrawElements(gl.Enum(mode), int(count), gl.Enum(itype), int(start))

	gs.stats.Drawcalls++
}

// Enable enables the specified capability.
func (gs *GLS) Enable(cap int) {

	if gs.capabilities[cap] == capEnabled {
		gs.stats.Caphits++
		return
	}
	gs.context.Enable(gl.Enum(cap))
	gs.capabilities[cap] = capEnabled
}

// Disable disables the specified capability.
func (gs *GLS) Disable(cap int) {

	if gs.capabilities[cap] == capDisabled {
		gs.stats.Caphits++
		return
	}
	gs.context.Disable(gl.Enum(cap))
	gs.capabilities[cap] = capDisabled
}

// EnableVertexAttribArray enables a generic vertex attribute array.
func (gs *GLS) EnableVertexAttribArray(attr gl.Attrib) {

	gs.context.EnableVertexAttribArray(attr)

}

// CullFace specifies whether front- or back-facing facets can be culled.
func (gs *GLS) CullFace(mode uint32) {

	gs.context.CullFace(gl.Enum(mode))

}

// FrontFace defines front- and back-facing polygons.
func (gs *GLS) FrontFace(mode uint32) {

	if gs.frontFace == mode {
		return
	}
	gs.context.FrontFace(gl.Enum(mode))
	gs.frontFace = mode
}

// GenBuffer generates a​buffer object name.
func (gs *GLS) GenBuffer() gl.Buffer {

	buf := gs.context.CreateBuffer()
	gs.stats.Buffers++
	return buf
}

// GenerateMipmap generates mipmaps for the specified texture target.
func (gs *GLS) GenerateMipmap(target uint32) {
	gs.context.GenerateMipmap(gl.Enum(target))

}

// GenTexture generates a texture object name.
func (gs *GLS) GenTexture() gl.Texture {

	tex := gs.context.CreateTexture()

	gs.stats.Textures++
	return tex
}

// GenVertexArray generates a vertex array object name.
func (gs *GLS) GenVertexArray() gl.VertexArray {

	vao := gs.context.CreateVertexArray()

	gs.stats.Vaos++
	return vao
}

// GetAttribLocation returns the location of the specified attribute variable.
func (gs *GLS) GetAttribLocation(program gl.Program, name string) gl.Attrib {

	loc := gs.context.GetAttribLocation(program, name)
	return loc
}

// GetProgramiv returns the specified parameter from the specified program object.
func (gs *GLS) GetProgramiv(program gl.Program, pname uint32, params *int32) {

	p := gs.context.GetProgrami(program, gl.Enum(pname))

	*params = int32(p)

	// C.glGetProgramiv(C.GLuint(program), C.GLenum(pname), (*C.GLint)(params))
}

// GetProgramInfoLog returns the information log for the specified program object.
func (gs *GLS) GetProgramInfoLog(program gl.Program) string {

	return gs.context.GetProgramInfoLog(program)
	// C.glGetProgramInfoLog(C.GLuint(program), C.GLsizei(length), nil, gs.gobufSize(uint32(length)))
	// return string(gs.gobuf[:length])
}

// GetShaderInfoLog returns the information log for the specified shader object.
func (gs *GLS) GetShaderInfoLog(shader gl.Shader) string {

	var length int32
	gs.GetShaderiv(shader, INFO_LOG_LENGTH, &length)
	if length == 0 {
		return ""
	}

	return gs.context.GetShaderInfoLog(shader)

	// C.glGetShaderInfoLog(C.GLuint(shader), C.GLsizei(length), nil, gs.gobufSize(uint32(length)))
	// return string(gs.gobuf[:length])
}

// GetString returns a string describing the specified aspect of the current GL connection.
func (gs *GLS) GetString(name uint32) string {

	return gs.context.GetString(gl.Enum(name))

	// cs := C.glGetString(C.GLenum(name))
	// return C.GoString((*C.char)(unsafe.Pointer(cs)))
}

// GetUniformLocation returns the location of a uniform variable for the specified program.
func (gs *GLS) GetUniformLocation(program gl.Program, name string) gl.Uniform {

	loc := gs.context.GetUniformLocation(program, name)

	return loc
}

// GetViewport returns the current viewport information.
func (gs *GLS) GetViewport() (x, y, width, height int32) {

	return gs.viewportX, gs.viewportY, gs.viewportWidth, gs.viewportHeight
}

// LineWidth specifies the rasterized width of both aliased and antialiased lines.
func (gs *GLS) LineWidth(width float32) {

	if gs.lineWidth == width {
		return
	}

	gs.context.LineWidth(width)
	gs.lineWidth = width
}

// LinkProgram links the specified program object.
func (gs *GLS) LinkProgram(program gl.Program) {

	gs.context.LinkProgram(program)

}

// GetShaderiv returns the specified parameter from the specified shader object.
func (gs *GLS) GetShaderiv(shader gl.Shader, pname uint32, params *int32) {

	p := gs.context.GetShaderi(shader, gl.Enum(pname))

	*params = int32(p)
}

// Scissor defines the scissor box rectangle in window coordinates.
func (gs *GLS) Scissor(x, y int32, width, height uint32) {

	gs.context.Scissor(x, y, int32(width), int32(height))

}

// ShaderSource sets the source code for the specified shader object.
func (gs *GLS) ShaderSource(shader gl.Shader, src string) {

	gs.context.ShaderSource(shader, src)
}

// TexImage2D specifies a two-dimensional texture image.

// "encoding/gob"
// "bytes"

func (gs *GLS) TexImage2D(target uint32, level int32, iformat int32, width int32, height int32, border int32, format uint32, itype uint32, tex interface{}) {

	data, _ := GetBytes(tex)
	gs.context.TexImage2D(gl.Enum(target), int(level), int(iformat), int(width), int(height), gl.Enum(format), gl.Enum(itype), data)

}

// func (gs *GLS) TexImage2D(target uint32, level int32, iformat int32, width int32, height int32, border int32, format uint32, itype uint32, data []byte) {

// 	gs.context.TexImage2D(gl.Enum(target), int(level), int(iformat), int(width), int(height), gl.Enum(format), gl.Enum(itype), data)

// }

// TexParameteri sets the specified texture parameter on the specified texture.
func (gs *GLS) TexParameteri(target uint32, pname uint32, param int32) {

	gs.context.TexParameteri(gl.Enum(target), gl.Enum(pname), int(param))

}

// PolygonMode controls the interpretation of polygons for rasterization.
// func (gs *GLS) PolygonMode(face, mode uint32) {

// 	if gs.polygonModeFace == face && gs.polygonModeMode == mode {
// 		return
// 	}

// 	// gs.PolygonMode(face, mode)
// 	gs.polygonModeFace = face
// 	gs.polygonModeMode = mode
// }

// PolygonOffset sets the scale and units used to calculate depth values.
func (gs *GLS) PolygonOffset(factor float32, units float32) {

	if gs.polygonOffsetFactor == factor && gs.polygonOffsetUnits == units {
		return
	}
	gs.context.PolygonOffset(factor, units)
	gs.polygonOffsetFactor = factor
	gs.polygonOffsetUnits = units
}

// Uniform1i sets the value of an int uniform variable for the current program object.
func (gs *GLS) Uniform1i(location gl.Uniform, v0 int32) {

	gs.context.Uniform1i(location, int(v0))

	gs.stats.Unisets++
}

// Uniform1f sets the value of a float uniform variable for the current program object.
func (gs *GLS) Uniform1f(location gl.Uniform, v0 float32) {

	gs.context.Uniform1f(location, v0)

	gs.stats.Unisets++
}

// Uniform2f sets the value of a vec2 uniform variable for the current program object.
func (gs *GLS) Uniform2f(location gl.Uniform, v0, v1 float32) {

	gs.context.Uniform2f(location, v0, v1)

	gs.stats.Unisets++
}

// Uniform3f sets the value of a vec3 uniform variable for the current program object.
func (gs *GLS) Uniform3f(location gl.Uniform, v0, v1, v2 float32) {

	gs.context.Uniform3f(location, v0, v1, v2)

	gs.stats.Unisets++
}

// Uniform4f sets the value of a vec4 uniform variable for the current program object.
func (gs *GLS) Uniform4f(location gl.Uniform, v0, v1, v2, v3 float32) {

	gs.context.Uniform4f(location, v0, v1, v2, v3)

	gs.stats.Unisets++
}

// UniformMatrix3fv sets the value of one or many 3x3 float matrices for the current program object.
func (gs *GLS) UniformMatrix3fv(location gl.Uniform, count int32, transpose bool, pm []float32) {

	gs.context.UniformMatrix3fv(location, pm)

	gs.stats.Unisets++
}

// UniformMatrix4fv sets the value of one or many 4x4 float matrices for the current program object.
func (gs *GLS) UniformMatrix4fv(location gl.Uniform, count int32, transpose bool, pm []float32) {

	gs.context.UniformMatrix4fv(location, pm)

	gs.stats.Unisets++
}

// Uniform1fv sets the value of one or many float uniform variables for the current program object.
func (gs *GLS) Uniform1fv(location gl.Uniform, count int32, v []float32) {
	gs.context.Uniform1fv(location, v)

	gs.stats.Unisets++
}

// Uniform2fv sets the value of one or many vec2 uniform variables for the current program object.
func (gs *GLS) Uniform2fv(location gl.Uniform, src []float32) {

	gs.context.Uniform2fv(location, src)

	gs.stats.Unisets++
}

// func (gs *GLS) Uniform2fvUP(location int32, count int32, v unsafe.Pointer) {

// 	gs.context.Uniform2fv(location)

// 	C.glUniform2fv(C.GLint(location), C.GLsizei(count), (*C.GLfloat)(v))
// 	gs.stats.Unisets++
// }

// Uniform3fv sets the value of one or many vec3 uniform variables for the current program object.
func (gs *GLS) Uniform3fv(location gl.Uniform, count int32, src []float32) {

	gs.context.Uniform3fv(location, src)

	gs.stats.Unisets++
}

// func (gs *GLS) Uniform3fvUP(location int32, count int32, v unsafe.Pointer) {

// 	C.glUniform3fv(C.GLint(location), C.GLsizei(count), (*C.GLfloat)(v))
// 	gs.stats.Unisets++
// }

// Uniform4fv sets the value of one or many vec4 uniform variables for the current program object.
func (gs *GLS) Uniform4fv(location gl.Uniform, count int32, v []float32) {

	gs.context.Uniform4fv(location, v)
	gs.stats.Unisets++
}

// func (gs *GLS) Uniform4fvUP(location int32, count int32, v unsafe.Pointer) {

// 	C.glUniform4fv(C.GLint(location), C.GLsizei(count), (*C.GLfloat)(v))
// 	gs.stats.Unisets++
// }

// VertexAttribPointer defines an array of generic vertex attribute data.
func (gs *GLS) VertexAttribPointer(index gl.Attrib, size int32, xtype uint32, normalized bool, stride int32, offset uint32) {

	// dst Attrib, size int, ty Enum, normalized bool, stride, offset int
	gs.context.VertexAttribPointer(index, int(size), gl.Enum(xtype), normalized, int(stride), int(offset))

	// C.glVertexAttribPointer(C.GLuint(index), C.GLint(size), C.GLenum(xtype), bool2c(normalized), C.GLsizei(stride), unsafe.Pointer(uintptr(offset)))
}

// Viewport sets the viewport.
func (gs *GLS) Viewport(x, y, width, height int32) {

	gs.context.Viewport(int(x), int(y), int(width), int(height))

	gs.viewportX = x
	gs.viewportY = y
	gs.viewportWidth = width
	gs.viewportHeight = height
}

// UseProgram sets the specified program as the current program.
func (gs *GLS) UseProgram(prog *Program) {

	if prog.handle.Value == 0 {
		panic("Invalid program")
	}

	//----- Todo

	gs.context.UseProgram(prog.handle)
	gs.prog = prog

	// Inserts program in cache if not already there.
	if !gs.programs[prog] {
		gs.programs[prog] = true
		log.Debug("New Program activated. Total: %d", len(gs.programs))
	}
}

func GetBytes(key interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
