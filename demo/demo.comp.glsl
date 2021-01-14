#version 430 core
#extension GL_ARB_compute_shader: enable
#extension GL_ARB_shader_storage_buffer_object: enable
//#extension GL_ARB_compute_variable_group_size: enable

layout(std430, binding=1) buffer Buf1 {
	uint foo[];
};

layout(
	local_size_x = 1,
	local_size_y = 1,
	local_size_z = 1) in;

void main() {
	int b = 0;
	vec4 c = vec4(1, 0.5, 2, sqrt(3));
	for (int i = 0; i < 30000; i++) {
		b++;
		c *= float(sqrt(b));
	}

	foo[gl_GlobalInvocationID.x] *= uint(sqrt(dot(c, c)));
}
