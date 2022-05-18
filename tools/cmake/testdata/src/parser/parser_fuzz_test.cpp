#include <iostream>

#include "parser.h"
#include <cifuzz/cifuzz.h>

FUZZ_TEST(const uint8_t *data, size_t size) {
  parse(std::string(reinterpret_cast<const char*>(data), size));
}
