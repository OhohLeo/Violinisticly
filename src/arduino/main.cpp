#include <Arduino.h>
#include <stdio.h>

// ADXL335EB : accelerometer (x, y, z)
// COM  => GND
// Z    => A0
// Y    => A1
// X    => A2
// VSS  => 3.3V
// ST   => not linked

// STATUS : replaced by MPU-6050 6-axis accelerometer/gyroscope

bool has_changed(int previous_value, int new_value);

#define MARGIN    5

// these constants describe the pins. They won't change:
const int groundpin = 18;             // analog input pin 4 -- ground
const int powerpin = 19;              // analog input pin 5 -- voltage
const int xpin = A2;                  // x-axis of the accelerometer
const int ypin = A1;                  // y-axis
const int zpin = A0;                  // z-axis (only on 3-axis models)

int previous_x, previous_y, previous_z;

void loop() {

  char tmp[30];

  int x, y, z;

  x = analogRead(xpin);
  y = analogRead(ypin);
  z = analogRead(zpin);

  if (has_changed(previous_x, x)
      || has_changed(previous_y, y)
      || has_changed(previous_z, z))
  {
      // print the sensor values:
      sprintf(tmp, "%d-%d-%d", x, y, z);
	  Serial.println(tmp);

      previous_x = x;
      previous_y = y;
      previous_z = z;
  }

  // delay before next reading:
  delay(100);
}

bool has_changed(int previous_value, int new_value)
{
    if (new_value >= previous_value + MARGIN
        || new_value <= previous_value - MARGIN)
        return true;

    return false;
}

int main(void)
{
    // Mandatory init
    init();

    Serial.begin(115200);

    // Pin 13 has an LED connected on most Arduino boards
    pinMode(13, OUTPUT);

    while (true)
        loop();

    return 0;
}
