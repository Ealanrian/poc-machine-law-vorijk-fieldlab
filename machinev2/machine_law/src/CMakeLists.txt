# The Flutter tooling requires that developers have CMake 3.10 or later
# installed. You should not increase this version, as doing so will cause
# the plugin to fail to compile for some customers of the plugin.
#cmake_minimum_required(VERSION 3.10)
#
#project(machine_law_library VERSION 0.0.1 LANGUAGES C)
#
#add_library(machine_law SHARED
#  "machine_law.c"
#)

#set_target_properties(machine_law PROPERTIES
#  PUBLIC_HEADER machine_law.h
#  OUTPUT_NAME "machine_law"
#)

#target_compile_definitions(machine_law PUBLIC DART_SHARED_LIB)#
#
#if (ANDROID)
  # Support Android 15 16k page size
#  target_link_options(machine_law PRIVATE "-Wl,-z,max-page-size=16384")
#endif()
