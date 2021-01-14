#version 430 core
#extension GL_ARB_compute_shader: enable
#extension GL_ARB_shader_storage_buffer_object: enable

layout(std430, binding=4) buffer Buf1 {
	uint foo[];
};

layout(
	local_size_x = 4,
	local_size_y = 1,
	local_size_z = 1) in;

void main() {
	foo[gl_GlobalInvocationID.x] *= 2;
}
