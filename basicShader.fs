#version 120

varying vec2 textureCoord0;

uniform sampler2D diffuse;

void main()
{
	gl_FragColor = texture2D(diffuse, textureCoord0);
}
