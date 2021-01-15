#version 430 core
#extension GL_ARB_compute_shader: enable
#extension GL_ARB_shader_storage_buffer_object: enable

layout(std430, binding=0) buffer Buf1 {
  float foo[];
};

layout(
       local_size_x = 256,
       local_size_y = 1,
       local_size_z = 1) in;

void main() {
  uint i = gl_GlobalInvocationID.x;
  foo[i] = sqrt(foo[i]);
}
