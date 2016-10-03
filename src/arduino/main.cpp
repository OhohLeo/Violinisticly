#include "I2Cdev.h"
//#include "MPU6050.h"
#include "MPU6050_6Axis_MotionApps20.h"
#include "crc16.h"
#include "PinChangeInt.h"

uint8_t store_float32(uint8_t *pua_buf, float value);
void send_status(uint8_t ua_type, uint8_t ua_status);
uint8_t *frame_get(uint8_t *data, uint8_t len);

// MPU6050    -     ARDUINO
// VCC => black  => 3.3V
// GND => white  => GND
// SCL => grey   => A5
// SDA => purple => A4
// INT => orange => DIGITAL-2

// class default I2C address is 0x68
// specific I2C addresses may be passed as a parameter here
// AD0 low = 0x68 (default for InvenSense evaluation board)
// AD0 high = 0x69
MPU6050 mpu;
//MPU6050 accelgyro(0x69); // <-- use for AD0 high

int16_t ax, ay, az;
int16_t gx, gy, gz;

#define HEADER_SIZE 2
#define STATUS_SIZE 2 * sizeof(uint8_t)
#define DATA_SIZE   6 * sizeof(uint16_t)
#define CRC_SIZE    sizeof(uint16_t)
#define END_SIZE    1
#define FULL_SIZE   HEADER_SIZE + CRC_SIZE + END_SIZE

// uncomment "OUTPUT_QUATERNION" if you want to see the actual
// quaternion components in a [w, x, y, z] format (not best for parsing
// on a remote host such as Processing or something though)
#define OUTPUT_QUATERNION   0x01

// uncomment "OUTPUT_EULER" if you want to see Euler angles
// (in degrees) calculated from the quaternions coming from the FIFO.
// Note that Euler angles suffer from gimbal lock (for more info, see
// http://en.wikipedia.org/wiki/Gimbal_lock)
// #define OUTPUT_EULER        0x02

// uncomment "OUTPUT_YAWPITCHROLL" if you want to see the yaw/
// pitch/roll angles (in degrees) calculated from the quaternions coming
// from the FIFO. Note this also requires gravity vector calculations.
// Also note that yaw/pitch/roll angles suffer from gimbal lock (for
// more info, see: http://en.wikipedia.org/wiki/Gimbal_lock)
// #define OUTPUT_YAWPITCHROLL 0x04

// uncomment "OUTPUT_REALACCEL" if you want to see acceleration
// components with gravity removed. This acceleration reference frame is
// not compensated for orientation, so +X is always +X according to the
// sensor, just without the effects of gravity. If you want acceleration
// compensated for orientation, us OUTPUT_WORLDACCEL instead.
// #define OUTPUT_REALACCEL    0x08

// uncomment "OUTPUT_WORLDACCEL" if you want to see acceleration
// components with gravity removed and adjusted for the world frame of
// reference (yaw is relative to initial orientation, since no magnetometer
// is present in this case). Could be quite handy in some cases.
// #define OUTPUT_WORLDACCEL   0x10

//#define OUTPUT_BUFFER       0x20
//#define BUFFER_SIZE         8

#define MPU_INITIALIZE 0
#define MPU_CONNECTION 1
#define DMP_INITIALIZE 2
#define DMP_INTERRUPT  3
#define FIFO_OVERFLOW  4

#define STATUS_OK     0x00
#define STATUS_FAIL   0xff

#define INTERRUPT_PIN 2
#define LED_PIN 13
bool b_blink_state = false;

// MPU control/status vars
bool b_dmp_ready = false;        // set true if DMP init was successful
uint8_t ua_mpu_interrupt_status; // holds actual interrupt status byte from MPU
uint8_t ua_dev_status;           // return status after each device operation (0 = success, !0 = error)
uint16_t uh_packet_size;         // expected DMP packet size (default is 42 bytes)
uint16_t uh_fifo_count;          // count of all bytes currently in FIFO
uint8_t ua_fifo_buffer[64];      // FIFO storage buffer

// orientation/motion vars
Quaternion s_quaternion;             // [w, x, y, z]         quaternion container
VectorInt16 s_acceleration;          // [x, y, z]            accel sensor measurements
VectorInt16 s_acceleration_real;     // [x, y, z]            gravity-free accel sensor measurements
VectorInt16 s_acceleration_world;    // [x, y, z]            world-frame accel sensor measurements
VectorFloat s_gravity;               // [x, y, z]            gravity vector
float rf_euler[3];                   // [psi, theta, phi]    Euler angle container
float rf_ypr[3];                     // [yaw, pitch, roll]   yaw/pitch/roll container and gravity vector

volatile bool mpuInterrupt = false;     // indicates whether MPU interrupt pin has gone high
void dmpDataReady()
{
    mpuInterrupt = true;
}

void setup()
{
    // join I2C bus (I2Cdev library doesn't do this automatically)
    #if I2CDEV_IMPLEMENTATION == I2CDEV_ARDUINO_WIRE
        Wire.begin();
        //Wire.setClock(400000); // 400kHz I2C clock. Comment this line if having compilation difficulties
    #elif I2CDEV_IMPLEMENTATION == I2CDEV_BUILTIN_FASTWIRE
        Fastwire::setup(400, true);
    #endif

    // initialize serial communication
    Serial.begin(38400);
    while (!Serial);

    mpu.initialize();
    pinMode(INTERRUPT_PIN, INPUT);

	send_status(MPU_INITIALIZE, STATUS_OK);

    // verify connection
	send_status(MPU_CONNECTION, mpu.testConnection() ? STATUS_OK : STATUS_FAIL);

    // load and configure the DMP
    // 0 = DMP OK
    // 1 = initial memory load failed
    // 2 = DMP configuration updates failed
    ua_dev_status = mpu.dmpInitialize();
	send_status(DMP_INITIALIZE, ua_dev_status);

    // supply your own gyro offsets here, scaled for min sensitivity
    mpu.setXGyroOffset(120);
    mpu.setYGyroOffset(76);
    mpu.setZGyroOffset(-185);
    mpu.setZAccelOffset(1688); // 1688 factory default for my test chip

    // make sure it worked (returns 0 if so)
    if (ua_dev_status == 0)
	{
        // turn on the DMP, now that it's ready
        mpu.setDMPEnabled(true);

        // enable Arduino interrupt detection
        attachPinChangeInterrupt(INTERRUPT_PIN, dmpDataReady, RISING);
        ua_mpu_interrupt_status = mpu.getIntStatus();
        send_status(DMP_INTERRUPT, ua_mpu_interrupt_status);

        b_dmp_ready = true;

        // get expected DMP packet size for later comparison
        uh_packet_size = mpu.dmpGetFIFOPacketSize();
    }

    // configure LED for output
    pinMode(LED_PIN, OUTPUT);
}

void loop()
{
  // if programming failed, don't try to do anything
  if (!b_dmp_ready)
    return;

  // wait for MPU interrupt or extra packet(s) available
  while (!mpuInterrupt && uh_fifo_count < uh_packet_size)
  {
  }

  // reset interrupt flag and get INT_STATUS byte
  mpuInterrupt = false;
  ua_mpu_interrupt_status = mpu.getIntStatus();

  // get current FIFO count
  uh_fifo_count = mpu.getFIFOCount();

  // check for overflow (this should never happen unless our code is too inefficient)
  if ((ua_mpu_interrupt_status & 0x10) || uh_fifo_count == 1024)
  {
	// reset so we can continue cleanly
	mpu.resetFIFO();
	send_status(FIFO_OVERFLOW, ua_mpu_interrupt_status);

    // otherwise, check for DMP data ready interrupt (this should happen frequently)
  }
  else if (ua_mpu_interrupt_status & 0x02)
  {
	uint8_t ua_idx, ua_nb = 0;
	uint8_t ua_data_len = 1;
	uint8_t ua_types = 0;
    uint8_t *pua_buf, *pua_data_buf, *pua_data_start;

	// wait for correct available data length, should be a VERY short wait
	while (uh_fifo_count < uh_packet_size)
	{
	  uh_fifo_count = mpu.getFIFOCount();
	}

	// read a packet from FIFO
	mpu.getFIFOBytes(ua_fifo_buffer, uh_packet_size);

	// track FIFO count here in case there is > 1 packet available
	// (this lets us immediately read more without waiting for an interrupt)
	uh_fifo_count -= uh_packet_size;


#ifdef OUTPUT_BUFFER
	ua_data_len += BUFFER_SIZE;
	ua_types |= OUTPUT_BUFFER;
#else
	// display quaternion values in easy matrix form: w x y z
	mpu.dmpGetQuaternion(&s_quaternion, ua_fifo_buffer);
#endif

#ifdef OUTPUT_QUATERNION
	ua_data_len += 4 * sizeof(float);
	ua_types |= OUTPUT_QUATERNION;
#endif

#ifdef OUTPUT_EULER
	mpu.dmpGetEuler(rf_euler, &s_quaternion);
    ua_data_len += 3 * sizeof(float);
	ua_types |= OUTPUT_EULER;
#endif

#if defined(OUTPUT_YAWPITCHROLL) || defined(OUTPUT_REALACCEL) || defined(OUTPUT_WORLDACCEL)
	mpu.dmpGetGravity(&s_gravity, &s_quaternion);
#endif

#ifdef OUTPUT_YAWPITCHROLL
	mpu.dmpGetYawPitchRoll(rf_ypr, &s_quaternion, &s_gravity);
	ua_data_len += 3 * sizeof(float);
	ua_types |= OUTPUT_YAWPITCHROLL;
#endif

#if defined(OUTPUT_REALACCEL) || defined(OUTPUT_WORLDACCEL)
	// display real acceleration, adjusted to remove gravity
	mpu.dmpGetAccel(&s_acceleration, ua_fifo_buffer);
	mpu.dmpGetLinearAccel(&s_acceleration_real, &s_acceleration, &s_gravity);
#endif

#ifdef OUTPUT_REALACCEL
	ua_data_len += 3 * sizeof(float);
	ua_types |= OUTPUT_REALACCEL;
#endif

#ifdef OUTPUT_WORLDACCEL
	// display initial world-frame acceleration, adjusted to remove gravity
	mpu.dmpGetLinearAccelInWorld(&s_acceleration_world, &s_acceleration_real, &s_quaternion);
	ua_data_len += 3 * sizeof(float);
	ua_types |= OUTPUT_WORLDACCEL;
#endif

	// allocate the buffe to store the values
    pua_data_start = (uint8_t*)malloc(ua_data_len);

	// Store the start of the buffer
	pua_data_buf = pua_data_start;

	// Setup type of data expected
	*pua_data_buf++ = ua_types;

#ifdef OUTPUT_BUFFER
	*pua_data_buf++ = ua_fifo_buffer[0];
	*pua_data_buf++ = ua_fifo_buffer[1];
	*pua_data_buf++ = ua_fifo_buffer[4];
	*pua_data_buf++ = ua_fifo_buffer[5];
	*pua_data_buf++ = ua_fifo_buffer[8];
	*pua_data_buf++ = ua_fifo_buffer[9];
	*pua_data_buf++ = ua_fifo_buffer[12];
	*pua_data_buf++ = ua_fifo_buffer[13];
#endif

#ifdef OUTPUT_QUATERNION
	pua_data_buf += store_float32(pua_data_buf, s_quaternion.w);
	pua_data_buf += store_float32(pua_data_buf, s_quaternion.x);
	pua_data_buf += store_float32(pua_data_buf, s_quaternion.y);
	pua_data_buf += store_float32(pua_data_buf, s_quaternion.z);
#endif

#ifdef OUTPUT_EULER
	pua_data_buf += store_float32(pua_data_buf, rf_euler[0] * 180/M_PI);
	pua_data_buf += store_float32(pua_data_buf, rf_euler[1] * 180/M_PI);
	pua_data_buf += store_float32(pua_data_buf, rf_euler[2] * 180/M_PI);
#endif

#ifdef OUTPUT_YAWPITCHROLL
	pua_data_buf += store_float32(pua_data_buf, rf_ypr[0] * 180/M_PI);
	pua_data_buf += store_float32(pua_data_buf, rf_ypr[1] * 180/M_PI);
	pua_data_buf += store_float32(pua_data_buf, rf_ypr[2] * 180/M_PI);
#endif

#ifdef OUTPUT_REALACCEL
	pua_data_buf += store_float32(pua_data_buf, s_acceleration_real.x);
	pua_data_buf += store_float32(pua_data_buf, s_acceleration_real.y);
	pua_data_buf += store_float32(pua_data_buf, s_acceleration_real.z);
#endif

#ifdef OUTPUT_WORLDACCEL
	pua_data_buf += store_float32(pua_data_buf, s_acceleration_world.x);
	pua_data_buf += store_float32(pua_data_buf, s_acceleration_world.y);
	pua_data_buf += store_float32(pua_data_buf, s_acceleration_world.z);
#endif

	// send the result
	pua_buf = frame_get((uint8_t*)pua_data_start, ua_data_len);
	if (pua_buf != NULL)
	  {
		Serial.write(pua_buf, FULL_SIZE + ua_data_len);
		free(pua_buf);
	  }

	free(pua_data_start);


	// blink LED to indicate activity
	b_blink_state = !b_blink_state;
	digitalWrite(LED_PIN, b_blink_state);
  }
}

uint8_t store_float32(uint8_t *pua_buf, float value)
{
  memcpy(pua_buf, &value, sizeof(float));
  return sizeof(float);
}

void send_status(uint8_t ua_type, uint8_t ua_status)
{
  uint8_t rua_data[STATUS_SIZE] = { ua_type, ua_status };
  uint8_t *pua_buf;
  pua_buf = frame_get(rua_data, STATUS_SIZE);
  Serial.write(pua_buf, FULL_SIZE + STATUS_SIZE);
  free(pua_buf);
}

// Ajoute l'entête, de la taille du buffer et le CRC16 du buffer et
// retourne le nouveau buffer prêt à l'envo.i
uint8_t *frame_get(uint8_t *pua_data, uint8_t len)
{
  uint8_t *pua_buf, *pua_start;
  uint8_t *pua_data_start = pua_data;
  uint16_t crc;
  int i;

  // allocation de la trame à émettre
  pua_start = (uint8_t *)malloc(FULL_SIZE + len);

  // pointe au début du buffer
  pua_buf = pua_start;

  // ajout de l'entête
  *pua_buf++ = ':';

  // ajout de la taille du buffer
  *pua_buf++ = len;

  // recopie de données
  for (i = 0; i < len; i++) {
	*pua_buf++ = *pua_data++;
  }

  // ajout du CRC-16
  crc = crc16(pua_data_start, len);
  *pua_buf++ = (uint8_t)(crc >> 8);
  *pua_buf++ = (uint8_t)crc;

  // ajout de la fin de la trame
  *pua_buf++ = '\n';

  return pua_start;
}
