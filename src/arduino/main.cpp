#include "I2Cdev.h"
#include "MPU6050.h"
#include "crc16.h"

void store_uh(uint8_t *p_dst, uint16_t uh_src);
uint8_t * frame_get(uint8_t *data, uint8_t len);

// class default I2C address is 0x68
// specific I2C addresses may be passed as a parameter here
// AD0 low = 0x68 (default for InvenSense evaluation board)
// AD0 high = 0x69
MPU6050 accelgyro;
//MPU6050 accelgyro(0x69); // <-- use for AD0 high

int16_t ax, ay, az;
int16_t gx, gy, gz;

#define HEADER_SIZE 2
#define DATA_SIZE   6 * sizeof(uint16_t)
#define CRC_SIZE    sizeof(uint16_t)
#define END_SIZE    1
#define FULL_SIZE   HEADER_SIZE + DATA_SIZE + CRC_SIZE + END_SIZE


// uncomment "OUTPUT_READABLE_ACCELGYRO" if you want to see a tab-separated
// list of the accel X/Y/Z and then gyro X/Y/Z values in decimal. Easy to read,
// not so easy to parse, and slow(er) over UART.
//#define OUTPUT_READABLE_ACCELGYRO

// uncomment "OUTPUT_BINARY_ACCELGYRO" to send all 6 axes of data as 16-bit
// binary, one right after the other. This is very fast (as fast as possible
// without compression or data loss), and easy to parse, but impossible to read
// for a human.
#define OUTPUT_BINARY_ACCELGYRO

#define LED_PIN 13
bool blinkState = false;

void setup() {
    // join I2C bus (I2Cdev library doesn't do this automatically)
    #if I2CDEV_IMPLEMENTATION == I2CDEV_ARDUINO_WIRE
        Wire.begin();
    #elif I2CDEV_IMPLEMENTATION == I2CDEV_BUILTIN_FASTWIRE
        Fastwire::setup(400, true);
    #endif

    // initialize serial communication
    Serial.begin(38400);

    // initialize device
    accelgyro.initialize();
}

void loop() {
    // read raw accel/gyro measurements from device
    accelgyro.getMotion6(&ax, &ay, &az, &gx, &gy, &gz);

    #ifdef OUTPUT_READABLE_ACCELGYRO
	char str[80];
	sprintf(str, "%d\t%d\t%d\t%d\t%d\t%d", ax, ay, az, gx, gy, gz);
	Serial.println(str);
    #endif

    #ifdef OUTPUT_BINARY_ACCELGYRO
	uint8_t rui_data[DATA_SIZE];
    uint8_t* buf;
	store_uh(rui_data, ax);
	store_uh(&rui_data[2], ay);
	store_uh(&rui_data[4], az);
	store_uh(&rui_data[6], gx);
	store_uh(&rui_data[8], gy);
	store_uh(&rui_data[10], gz);

    buf = frame_get(rui_data, DATA_SIZE);
    Serial.write(buf, FULL_SIZE);
    free(buf);
    #endif
}

void store_uh(uint8_t *p_dst, const uint16_t uh_src) {
  p_dst[0] = (uint8_t)(uh_src >> 8);
  p_dst[1] = (uint8_t)(uh_src & 0xFF);
}

// Ajoute l'entête, de la taille du buffer et le CRC16 du buffer et
// retourne le nouveau buffer prêt à l'envo.i
uint8_t * frame_get(uint8_t *data, uint8_t len)
{
  uint8_t *buf;
  int i;

  // allocation de la trame à émettre
  buf = (uint8_t *) malloc(HEADER_SIZE + len + CRC_SIZE);

  // ajout de l'entête
  buf[0] = ':';

  // ajout de la taille du buffer
  buf[1] = len;

  // recopie du buffer
  for (i = 0; i < len; i++) {
	buf[HEADER_SIZE+i] = data[i];
  }

  // ajout du CRC-16
  uint16_t crc = crc16(data, len);
  len += HEADER_SIZE;
  buf[len] = crc >> 8;
  buf[len+1] = crc & 0xFF;

  len += CRC_SIZE;
  buf[len] = '\n';

  return buf;
}

// Affiche le contenu de la trame
void frame_display(uint8_t *data, uint8_t len)
{
  int i;

  for (i = 0; i < len; i++) {
	printf(" %02X", data[i]);
  }

  printf("\n");
}
