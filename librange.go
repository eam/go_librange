package librange

/*
#cgo CFLAGS: -DPNG_DEBUG=1 -I/home/eam/git/libcrange/source/src/  -I/usr/include/apr-1 -I/Users/evan/git/libcrange/source/src/
#cgo amd64 386 CFLAGS: -DX86=1 -I/home/eam/git/libcrange/source/src/  -I/usr/include/apr-1
#cgo LDFLAGS: -L/home/eam/prefix/usr/lib  -L/usr/local/lib -L/usr/lib64/perl5/CORE -lpcre -lperl -lpython2.6 -lcrange -L/Users/evan/prefix/usr/lib  -L/usr/local/lib  -L/System/Library/Perl/5.16/darwin-thread-multi-2level/CORE -L/usr/local/Cellar/python/2.7.6/Frameworks/Python.framework/Versions/2.7/lib/python2.7/config  -L/usr/local/Cellar/pcre/8.34/lib 
#include <range.h>
#include <stdio.h>
#include <stdlib.h>

// These three functions cribbed from https://github.com/jgallagher/go-libpq

static void setArrayString(char **a, char *s, int n) {
	a[n] = s;
}

static void freeArrayElements(int n, char **a) {
	int i;
	for (i = 0; i < n; i++) {
		free(a[i]);
		a[i] = NULL;
	}
}

char **makeCharArray(int num_ptrs) {
  return calloc(sizeof(char *), num_ptrs);
}

static void print_array(char **a) {
  char **i = a;
  while (*i) {
    printf("addr: %p, ", *i);
    printf("content: %s\n", *i);
    i++;
  }
}

*/
import "C"

import "fmt"
import "unsafe"


type RangeLib struct {
        Lr *_Ctype_easy_lr
	want_warnings bool
}

func NewRangeLib(config_file string) *RangeLib {
	c_config_file := C.CString(config_file)
	range_obj := C.range_easy_create(c_config_file)
	C.free(unsafe.Pointer(c_config_file))
        //fmt.Printf("%v\n", range_obj)
	rl := &RangeLib{range_obj, true}
	return rl
}

func (rl *RangeLib) ExpandRange(range_expr string) []string {
	range_string := range_expr
	c_range_string := C.CString(range_string)
        ret := C.range_easy_expand(rl.Lr, c_range_string)
	C.free(unsafe.Pointer(c_range_string))
        var strings []string
        q := ret
	for {
            p := (**C.char)(q)
            if *p == nil {
                break
            }
            strings = append(strings, C.GoString(*p))
            q = (**_Ctype_char)(unsafe.Pointer(unsafe.Sizeof(q) + uintptr(unsafe.Pointer(q))))
        }
	//fmt.Printf("%v\n", strings)
	return strings
}

func (rl *RangeLib) CompressRange(nodes []string) string {
	fmt.Printf("%v\n", len(nodes))
	node_count := len(nodes)
	var q **_Ctype_char
        q = C.makeCharArray(C.int(node_count + 1))
        fmt.Printf("%v\n", q)
	for index, str := range nodes {
          fmt.Printf("index: %v\n", index)
          C.setArrayString(q, C.CString(str), C.int(index))
        }
	fmt.Printf("setting %v\n", C.int(node_count))
        C.setArrayString(q, (*_Ctype_char)(unsafe.Pointer(uintptr(0))), C.int(node_count))
	fmt.Printf("rl.Lr: %v\n", rl.Lr)
	fmt.Printf("q: %v\n", q)
        C.print_array(q)
	compressed := C.range_easy_compress(rl.Lr, q)
	C.freeArrayElements(C.int(node_count + 1), q)
	return C.GoString(compressed)
}
