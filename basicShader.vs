#version 120

varying vec2 textureCoord0;

attribute vec3 position;
attribute vec2 textureCoord;

void main()
{
	gl_Position = vec4(position, 1.0);
	textureCoord0 = textureCoord;
}