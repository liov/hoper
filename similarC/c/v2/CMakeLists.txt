cmake_minimum_required(VERSION 3.14)
project(c C)

set(CMAKE_C_STANDARD 11)

add_library(c SHARED library.c library.h)

add_executable(d main.c)