cmake_minimum_required(VERSION 3.6)

set(CMAKE_TOOLCHAIN_FILE ${CMAKE_SOURCE_DIR}/scripts/arduino/ArduinoToolchain.cmake)

project(arduino)

set(ARDUINO_SRCS ${CMAKE_SOURCE_DIR}/src/arduino/main.cpp
                 ${CMAKE_SOURCE_DIR}/src/arduino/I2Cdev.cpp
				 ${CMAKE_SOURCE_DIR}/src/arduino/MPU6050.cpp)

set(ARDUINO_HDRS ${CMAKE_SOURCE_DIR}/src/arduino/I2Cdev.h
                 ${CMAKE_SOURCE_DIR}/src/arduino/MPU6050.h
                 ${CMAKE_SOURCE_DIR}/src/arduino/MPU6050_9Axis_MotionApps41.h
                 ${CMAKE_SOURCE_DIR}/src/arduino/MPU6050_6Axis_MotionApps20.h
                 ${CMAKE_SOURCE_DIR}/src/arduino/helper_3dmath.h)

generate_arduino_firmware(arduino
    SRCS ${ARDUINO_SRCS}
	HDRS ${ARDUINO_HDRS}
    PORT  /dev/ttyACM0
    SERIAL picocom @SERIAL_PORT@
    BOARD uno)
