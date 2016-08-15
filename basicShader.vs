#version 120

varying vec2 textureCoord0;

attribute vec3 position;
attribute vec2 textureCoord;

uniform mat4 transform;

void main()
{
	gl_Position = transform * vec4(position, 1.0);
	textureCoord0 = textureCoord;
}