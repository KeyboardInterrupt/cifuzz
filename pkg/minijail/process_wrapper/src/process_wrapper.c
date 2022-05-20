#include <errno.h>
#include <stddef.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

// Includes argv[0] and the working directory.
#define NUM_FIXED_ARGS 2u

// Executes argv[2] with all following arguments after changing the
// working directory to argv[1].
int main(int argc, char **argv) {
  if (argc < 4) {
    fprintf(stderr,
            "Usage: %s <directory> <executable_path> <executable_arg1>...\n",
            argv[0]);
    return 1;
  }

  if (chdir(argv[1]) == -1) {
    fprintf(stderr, "chdir(%s) failed: %s\n", argv[1], strerror(errno));
    return 1;
  }

  char **executable_argv;
  if (strcmp("--", argv[2]) != 0) {
    executable_argv = &argv[3];
  } else {
    executable_argv = &argv[2];
  }

  if (execv(*executable_argv, executable_argv) == -1) {
    fprintf(stderr, "execv(%s) failed: %s\n", *executable_argv,
            strerror(errno));
    return 1;
  }
  // Not reached.
}
