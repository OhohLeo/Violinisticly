#define POLYNOM 0x1081

// CRC-CCITT (Kermit)
unsigned short crc16(const unsigned char* buf, int len){

  unsigned short crc = 0;

  while (len-- > 0) {
	unsigned char byte = *buf++;
	unsigned short q = (crc ^ byte) & 0x0f;
	crc = (crc >> 4) ^ (q * POLYNOM);
	q = (crc ^ (byte >> 4)) & 0xf;
	crc = (crc >> 4) ^ (q * POLYNOM);
  }

  return (crc >> 8) ^ (crc << 8);
}
