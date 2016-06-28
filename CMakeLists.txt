project(ACCELEROMETER)

cmake_minimum_required(VERSION 2.6)

include(${CMAKE_SOURCE_DIR}/scripts/arduino/arduino_uno.cmake)

set(ARDUINO_SRC
    ${CMAKE_SOURCE_DIR}/src/arduino/main.cpp
    ${ARDUINO_SOURCE_FILES}
)

add_executable(firmware ${ARDUINO_SRC})

add_custom_target(deploy  make
  COMMAND avr-objcopy -O ihex -R .eeprom firmware firmware.hex
  COMMAND sudo avrdude -v -c arduino -p ATMEGA328P -P /dev/ttyACM0 -b 115200 -U flash:w:firmware.hex)
