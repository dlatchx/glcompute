# GLcompute

GLcompute (aka glc) is a Go library  for GPGPU using OpenGL 4.3 compute shaders as the backend.

`glc/` contains the actual library
`demo/` contains a short example to get started with glc

`glc/buffer.go` is the most tricky part ; it uses `unsafe` and `reflect` to interface Go's safe slices with OpenGL's unsafe `mmap`-based GPU memory access.

The rest is basically wrappers for OpenGL objects/functions written with the help of [gotk3](https://github.com/gotk3/gotk3), plus some helper functions.
