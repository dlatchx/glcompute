package glc

import (
	"reflect"
	"runtime"
	"unsafe"

	"github.com/go-gl/gl/v4.3-core/gl"
)

// const (
// 	ARRAY_BUFFER              = gl.ARRAY_BUFFER              // Vertex attributes
// 	ATOMIC_COUNTER_BUFFER     = gl.ATOMIC_COUNTER_BUFFER     // Atomic counter storage
// 	COPY_READ_BUFFER          = gl.COPY_READ_BUFFER          // Buffer copy source
// 	COPY_WRITE_BUFFER         = gl.COPY_WRITE_BUFFER         // Buffer copy destination
// 	DISPATCH_INDIRECT_BUFFER  = gl.DISPATCH_INDIRECT_BUFFER  // Indirect compute dispatch commands
// 	DRAW_INDIRECT_BUFFER      = gl.DRAW_INDIRECT_BUFFER      // Indirect command arguments
// 	ELEMENT_ARRAY_BUFFER      = gl.ELEMENT_ARRAY_BUFFER      // Vertex array indices
// 	PIXEL_PACK_BUFFER         = gl.PIXEL_PACK_BUFFER         // Pixel read target
// 	PIXEL_UNPACK_BUFFER       = gl.PIXEL_UNPACK_BUFFER       // Texture data source
// 	QUERY_BUFFER              = gl.QUERY_BUFFER              // Query result buffer
// 	SHADER_STORAGE_BUFFER     = gl.SHADER_STORAGE_BUFFER     // Read-write storage for shaders
// 	TEXTURE_BUFFER            = gl.TEXTURE_BUFFER            // Texture data buffer
// 	TRANSFORM_FEEDBACK_BUFFER = gl.TRANSFORM_FEEDBACK_BUFFER // Transform feedback buffer
// 	UNIFORM_BUFFER            = gl.UNIFORM_BUFFER            // Uniform block storage
// )

// usage flags
const (
	BUF_STREAM_DRAW  = gl.STREAM_DRAW
	BUF_STREAM_READ  = gl.STREAM_READ
	BUF_STREAM_COPY  = gl.STREAM_COPY
	BUF_STATIC_DRAW  = gl.STATIC_DRAW
	BUF_STATIC_READ  = gl.STATIC_READ
	BUF_STATIC_COPY  = gl.STATIC_COPY
	BUF_DYNAMIC_DRAW = gl.DYNAMIC_DRAW
	BUF_DYNAMIC_READ = gl.DYNAMIC_READ
	BUF_DYNAMIC_COPY = gl.DYNAMIC_COPY
)

type Buffer struct {
	id uint32

	target uint32

	nbElem   int
	elemSize int
}

func (b *Buffer) Ready() bool {
	return b.id != 0
}

func (b *Buffer) Delete() {
	if b.Ready() {
		gl.DeleteBuffers(0, &b.id)
		b.id = 0
	}
}

func newBuffer(target uint32, length int, elemSize uintptr, usageHint uint32) *Buffer {
	nbElem := length
	size := nbElem * int(elemSize)

	var id uint32
	gl.GenBuffers(1, &id)

	gl.BindBuffer(target, id)
	gl.BufferData(target, size, nil, usageHint)

	buf := &Buffer{
		id: id,

		target: target,

		nbElem:   nbElem,
		elemSize: int(elemSize),
	}

	runtime.SetFinalizer(buf, func(b *Buffer) {
		b.Delete()
	})

	return buf
}

func NewBufferAtomic(length int, elemSize uintptr, usageHint uint32) *Buffer {
	return newBuffer(gl.ATOMIC_COUNTER_BUFFER, length, elemSize, usageHint)
}

func NewBufferStorage(length int, elemSize uintptr, usageHint uint32) *Buffer {
	return newBuffer(gl.SHADER_STORAGE_BUFFER, length, elemSize, usageHint)
}

func NewBufferUniform(length int, elemSize uintptr, usageHint uint32) *Buffer {
	return newBuffer(gl.UNIFORM_BUFFER, length, elemSize, usageHint)
}

func loadBuffer(target uint32, sourceSlicePtr interface{}, usageHint uint32) *Buffer {
	v := reflect.Indirect(reflect.ValueOf(sourceSlicePtr))
	length := v.Len()
	elemSize := v.Type().Elem().Size()

	buf := newBuffer(target, length, elemSize, usageHint)

	buf.Upload(sourceSlicePtr)

	return buf
}

func LoadBufferAtomic(sourceSlicePtr interface{}, usageHint uint32) *Buffer {
	return loadBuffer(gl.ATOMIC_COUNTER_BUFFER, sourceSlicePtr, usageHint)
}

func LoadBufferStorage(sourceSlicePtr interface{}, usageHint uint32) *Buffer {
	return loadBuffer(gl.SHADER_STORAGE_BUFFER, sourceSlicePtr, usageHint)
}

func LoadBufferUniform(sourceSlicePtr interface{}, usageHint uint32) *Buffer {
	return loadBuffer(gl.UNIFORM_BUFFER, sourceSlicePtr, usageHint)
}

// access flags
const (
	MAP_READ       = gl.MAP_READ_BIT
	MAP_WRITE      = gl.MAP_WRITE_BIT
	MAP_PERSISTENT = gl.MAP_PERSISTENT_BIT
	MAP_COHERENT   = gl.MAP_COHERENT_BIT

	MAP_INVALIDATE_RANGE  = gl.MAP_INVALIDATE_RANGE_BIT
	MAP_INVALIDATE_BUFFER = gl.MAP_INVALIDATE_BUFFER_BIT
	MAP_FLUSH_EXPLICIT    = gl.MAP_FLUSH_EXPLICIT_BIT
	MAP_UNSYNCHRONIZED    = gl.MAP_UNSYNCHRONIZED_BIT
)

func (b *Buffer) Map(slicePtr interface{}, accessFlags uint32) {
	b.MapRange(slicePtr, 0, b.nbElem, accessFlags)
}

// use black magic to get slice descriptor
func getSliceHeader(slicePtr interface{}) *reflect.SliceHeader {
	return (*reflect.SliceHeader)(unsafe.Pointer(reflect.ValueOf(slicePtr).Pointer()))
}

func (b *Buffer) MapRange(slicePtr interface{}, offset, length int, accessFlags uint32) {
	// get mapped memory address
	gl.BindBuffer(b.target, b.id)
	baseAddr := gl.MapBufferRange(b.target, offset*b.elemSize, length*b.elemSize, accessFlags)

	sliceHeader := getSliceHeader(slicePtr)
	if baseAddr == nil {
		sliceHeader.Data = 0
		sliceHeader.Len = 0
		sliceHeader.Cap = 0
	} else {
		sliceHeader.Data = uintptr(baseAddr)
		sliceHeader.Len = length
		sliceHeader.Cap = length
	}
}

func (b *Buffer) mapByte(slicePtr interface{}, accessFlags uint32) {
	// get mapped memory address
	gl.BindBuffer(b.target, b.id)
	baseAddr := gl.MapBufferRange(b.target, 0, b.nbElem*b.elemSize, accessFlags)

	sliceHeader := getSliceHeader(slicePtr)
	if baseAddr == nil {
		sliceHeader.Data = 0
		sliceHeader.Len = 0
		sliceHeader.Cap = 0
	} else {
		sliceHeader.Data = uintptr(baseAddr)
		sliceHeader.Len = b.nbElem * b.elemSize
		sliceHeader.Cap = b.nbElem * b.elemSize
	}
}

func (b *Buffer) Upload(slicePtr interface{}) {
	var mmapBindingPoint []byte
	b.mapByte(&mmapBindingPoint, MAP_WRITE|MAP_INVALIDATE_BUFFER)

	// get a []byte version of input slice to avoid element size problems
	var byteSlice []byte
	byteSliceHeader := getSliceHeader(&byteSlice)
	inputSliceHeader := getSliceHeader(slicePtr)
	byteSliceHeader.Data = inputSliceHeader.Data
	byteSliceHeader.Len = inputSliceHeader.Len * b.elemSize
	byteSliceHeader.Cap = inputSliceHeader.Cap * b.elemSize

	copy(mmapBindingPoint, byteSlice)

	ok := b.Unmap()
	if !ok {
		panic("unmap error")
	}
}

func (b *Buffer) Download(slicePtr interface{}) {
	var mmapBindingPoint []byte
	b.mapByte(&mmapBindingPoint, MAP_READ)

	// get a []byte version of input slice to avoid element size problems
	var byteSlice []byte
	byteSliceHeader := getSliceHeader(&byteSlice)
	outputSliceHeader := getSliceHeader(slicePtr)
	byteSliceHeader.Data = outputSliceHeader.Data
	byteSliceHeader.Len = outputSliceHeader.Len * b.elemSize
	byteSliceHeader.Cap = outputSliceHeader.Cap * b.elemSize

	copy(byteSlice, mmapBindingPoint)

	ok := b.Unmap()
	if !ok {
		panic("unmap error")
	}
}

func (b *Buffer) Unmap() bool {
	gl.BindBuffer(b.target, b.id)
	return gl.UnmapBuffer(b.target)
}

func (b *Buffer) Bind(slot uint32) {
	gl.BindBufferBase(b.target, slot, b.id)
}
