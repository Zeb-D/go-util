//go2cpp.cpp
#include <iostream>

extern "C" {
    #include "go2cpp.h"
}

int SayHelloV3() {
  std::cout<<"Hello World v3";
    return 0;
}