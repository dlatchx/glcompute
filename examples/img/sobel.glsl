#version 430 core
#extension GL_ARB_compute_shader: enable
#extension GL_ARB_shader_storage_buffer_object: enable

#define __at(pos)             ((pos).y * gl_NumWorkGroups.x + (pos).x)
#define get_at(buf, pos)      buf[__at(pos)]
#define set_at(buf, pos, col) buf[__at(pos)] = col;

layout(std430, binding=0) buffer Buf1 {
  vec4 imgIn[];
};

layout(std430, binding=1) buffer Buf2 {
  vec4 imgOut[];
};

layout(
       local_size_x = 1,
       local_size_y = 1,
       local_size_z = 1) in;

void main() {
  ivec2 pos = ivec2(gl_GlobalInvocationID.xy);
  ivec2 size = ivec2(gl_NumWorkGroups.xy);

  if (pos.x >= 1 && pos.x < size.x - 1 && pos.y > 1 && pos.y < size.y - 1) {
    vec4 gx =
        get_at(imgIn, pos + ivec2(-1, -1)) *  1.0
      + get_at(imgIn, pos + ivec2( 1, -1)) * -1.0
      + get_at(imgIn, pos + ivec2(-1,  0)) *  2.0
      + get_at(imgIn, pos + ivec2( 1,  0)) * -2.0
      + get_at(imgIn, pos + ivec2(-1,  1)) *  1.0
      + get_at(imgIn, pos + ivec2( 1,  1)) * -1.0;

    vec4 gy =
        get_at(imgIn, pos + ivec2(-1, -1)) *  1.0
      + get_at(imgIn, pos + ivec2( 0, -1)) *  2.0
      + get_at(imgIn, pos + ivec2( 1, -1)) *  1.0
      + get_at(imgIn, pos + ivec2(-1 , 1)) * -1.0
      + get_at(imgIn, pos + ivec2( 0,  1)) * -2.0
      + get_at(imgIn, pos + ivec2( 1,  1)) * -1.0;

    vec4 g = sqrt(gx*gx + gy*gy);

    set_at(imgOut, pos, g);
  } else
    set_at(imgOut, pos, vec4(0));
}
