//
// Fragment Shader template
//

#ifdef GL_ES
precision lowp float;
#endif



varying vec3 Color;
// out vec4 FragColor;

void main() {

    gl_FragColor = vec4(Color, 1.0);
}

