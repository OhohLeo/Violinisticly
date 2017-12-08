cmake_minimum_required(VERSION 3.6)

project(
  JackTest
  VERSION 0.0.1
  LANGUAGES CXX)

file(GLOB src *.cpp)

add_executable(jacktest ${src})

set(CMAKE_MODULE_PATH ${PROJECT_SOURCE_DIR}/../../externals/cmake-modules)
find_package(Jack)
