cmake_minimum_required(VERSION 3.6)

set(CMAKE_TOOLCHAIN_FILE ${CMAKE_SOURCE_DIR}/scripts/arduino/ArduinoToolchain.cmake)

project(violin)

set(ARDUINO_SRC ${CMAKE_SOURCE_DIR}/src/arduino/main.cpp)

generate_arduino_firmware(test
    SRCS ${ARDUINO_SRC}
    PORT  /dev/ttyACM0
    SERIAL picocom @SERIAL_PORT@
    BOARD uno)
