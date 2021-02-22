package main

// Source: https://github.com/sstephenson/launch_socket_server/blob/master/src/launch/socket.go

/*
#include <stdlib.h>
int launch_activate_socket(const char *name, int **fds, size_t *cnt);
*/
import "C"

import (
	"errors"
	"unsafe"
)

func activateSocket(name string) ([]int, error) {
	c_name := C.CString(name)
	var c_fds *C.int
	c_cnt := C.size_t(0)

	err := C.launch_activate_socket(c_name, &c_fds, &c_cnt)
	if err != 0 {
		return nil, errors.New("couldn't activate launchd socket " + name)
	}

	length := int(c_cnt)
	pointer := unsafe.Pointer(c_fds)
	fds := (*[1 << 30]C.int)(pointer)
	result := make([]int, length)

	for i := 0; i < length; i++ {
		result[i] = int(fds[i])
	}

	C.free(pointer)
	return result, nil
}
