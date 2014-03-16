// Copyright © 2014 Lawrence E. Bakst. All rights reserved.
// binding to allow AntTweakBar use from Go using GLFW only (for now)
// #define GLFW_CDECL
package atb
/*
#cgo CFLAGS:
#cgo LDFLAGS: -framework GLUT -framework OpenGL -framework Cocoa -framework IOKit -framework QuartzCore -L/Volumes/std/usr/leb/gotest/src/leb/atb/atb/lib -lAntTweakBar -lglfw3
#include <GLFW/glfw3.h>
#include "atb/include/AntTweakBar.h"
#include <stdio.h>
#include <stdlib.h>
void _SetAllTWCallBacks(GLFWwindow *window);
*/
import "C"
import "unsafe"
import "fmt"
import glfw "github.com/go-gl/glfw3"

/*
type CBHolder struct {
	MouseButtonHolder func(w *Window, button MouseButton, action Action, mod ModifierKey)
	CursorPosHolder   func(w *Window, xpos float64, ypos float64)
	CursorEnterHolder func(w *Window, entered bool)
	ScrollHolder      func(w *Window, xoff float64, yoff float64)
	KeyHolder         func(w *Window, key Key, scancode int, action Action, mods ModifierKey)
	CharHolder        func(w *Window, char uint)
}
*/

// Constants
const (
	TW_TYPE_UNDEF			= C.TW_TYPE_UNDEF
	TW_MOUSE_PRESSED		= C.TW_MOUSE_PRESSED
	TW_MOUSE_RELEASED		= C.TW_MOUSE_RELEASED
	TW_MOUSE_LEFT			= C.TW_MOUSE_LEFT
	TW_MOUSE_MIDDLE			= C.TW_MOUSE_MIDDLE
	TW_MOUSE_RIGHT			= C.TW_MOUSE_RIGHT
)


type TwBar C.TwBar

func TwInit(gapi, device string) int {
	var api C.TwGraphAPI

	devp := C.CString(device)
	switch gapi {
	case "opengl-compatibility", "opengl":
		api = C.TW_OPENGL
	case "D3D9":
	    api = C.TW_DIRECT3D9
	case "D3D10":
	    api = C.TW_DIRECT3D10
	case "D3D11":
	    api = C.TW_DIRECT3D11
	case "opengl-core":
		api = C.TW_OPENGL_CORE
	}
	ret := C.TwInit(api, nil)
	C.free(unsafe.Pointer(devp))
	return int(ret)
}

func TwNewBar(name string) *TwBar {
	namec := C.CString(name)
	defer C.free(unsafe.Pointer(namec))
	return (*TwBar)(C.TwNewBar(namec))
}

func TwDefine(avar string) int {
	avarc := C.CString(avar)
	defer C.free(unsafe.Pointer(avarc))
	return int(C.TwDefine(avarc))
}

func (t *TwBar) AddVarRW(key string, avar interface{}, doc string) int {
	keyc := C.CString(key)
	docc := C.CString(doc)
	f := func() {
		C.free(unsafe.Pointer(keyc))
		C.free(unsafe.Pointer(docc))
	}
	defer f()
	switch avar.(type) {
	case *bool:
		b := avar.(*bool)
		ret := int(C.TwAddVarRW((*C.TwBar)(t), keyc, C.TW_TYPE_BOOL8, unsafe.Pointer(b), docc))
		fmt.Printf("AddVarRW: b=%v\n", b)
		return ret
	case *int8:
		i8 := avar.(*int8)
		return int(C.TwAddVarRW((*C.TwBar)(t), keyc, C.TW_TYPE_INT8, unsafe.Pointer(i8), docc))
	case uint8:
		ui8 := avar.(*uint8)
		return int(C.TwAddVarRW((*C.TwBar)(t), keyc, C.TW_TYPE_UINT8, unsafe.Pointer(ui8), docc))
	case int16:
		i16 := avar.(*int16)
		return int(C.TwAddVarRW((*C.TwBar)(t), keyc, C.TW_TYPE_INT16, unsafe.Pointer(i16), docc))
	case uint16:
		ui16 := avar.(*uint16)
		return int(C.TwAddVarRW((*C.TwBar)(t), keyc, C.TW_TYPE_UINT16, unsafe.Pointer(ui16), docc))
	case *int32:
		i32 := avar.(*int32)
		return int(C.TwAddVarRW((*C.TwBar)(t), keyc, C.TW_TYPE_INT32, unsafe.Pointer(i32), docc))
	case *uint32:
		// uint32 gets you a color, use uint or uint64 for a number
		ui32 := avar.(*uint32)
		return int(C.TwAddVarRW((*C.TwBar)(t), keyc, C.TW_TYPE_COLOR32, unsafe.Pointer(ui32), docc))
	case *int64:
		i64 := avar.(*int64)
		//i32 := int32(i64)
		return int(C.TwAddVarRW((*C.TwBar)(t), keyc, C.TW_TYPE_INT32, unsafe.Pointer(i64), docc))
	case *int:
		i64 := avar.(*int)
		//i32 := int32(i64)
		return int(C.TwAddVarRW((*C.TwBar)(t), keyc, C.TW_TYPE_INT32, unsafe.Pointer(i64), docc))
	case *uint64:
		ui64 := avar.(*uint64)
		//ui32 := int32(ui64)
		return int(C.TwAddVarRW((*C.TwBar)(t), keyc, C.TW_TYPE_UINT32, unsafe.Pointer(ui64), docc))
	case *uint:
		ui64 := avar.(*uint)
		//ui32 := int32(ui64)
		return int(C.TwAddVarRW((*C.TwBar)(t), keyc, C.TW_TYPE_UINT32, unsafe.Pointer(ui64), docc))
	case *float32:
		f32 := avar.(*float32)
		return int(C.TwAddVarRW((*C.TwBar)(t), keyc, C.TW_TYPE_FLOAT, unsafe.Pointer(f32), docc))
	case *float64:
		f64 := avar.(*float64)
		return int(C.TwAddVarRW((*C.TwBar)(t), keyc, C.TW_TYPE_DOUBLE, unsafe.Pointer(f64), docc))
	
	// an array (or slice) of 3 or 4 float32's gets you a color, if you want a TW_TYPE_DIR3F or TW_TYPE_QUAT4F use float64
	case [3]float32:
		s := avar.([3]float32)
		return int(C.TwAddVarRW((*C.TwBar)(t), keyc, C.TW_TYPE_COLOR3F, unsafe.Pointer(&s[0]), docc))
		//TW_TYPE_DIR3F
	case *[4]float32:
		s := avar.(*[4]float32)
		return int(C.TwAddVarRW((*C.TwBar)(t), keyc, C.TW_TYPE_COLOR4F, unsafe.Pointer(&s[0]), docc))
		//TW_TYPE_QUAT4F
	case [3]float64:
		s := avar.([3]float64)
		return int(C.TwAddVarRW((*C.TwBar)(t), keyc, C.TW_TYPE_DIR3D, unsafe.Pointer(&s[0]), docc))
	case [4]float64:
		s := avar.([4]float32)
		return int(C.TwAddVarRW((*C.TwBar)(t), keyc, C.TW_TYPE_QUAT4D, unsafe.Pointer(&s[0]), docc))
	case []float32:
		s := avar.([]float32)
		switch len(s) {
		case 3:
			return int(C.TwAddVarRW((*C.TwBar)(t), keyc, C.TW_TYPE_COLOR3F, unsafe.Pointer(&s[0]), docc))
		case 4:
			return int(C.TwAddVarRW((*C.TwBar)(t), keyc, C.TW_TYPE_COLOR4F, unsafe.Pointer(&s[0]), docc))
		default:
			panic("TwAddVarRW: bad slice length")
		}
	case []float64:
		s := avar.([]float64)
		switch len(s) {
		case 3:
			return int(C.TwAddVarRW((*C.TwBar)(t), keyc, C.TW_TYPE_DIR3D, unsafe.Pointer(&s[0]), docc))
		case 4:
			return int(C.TwAddVarRW((*C.TwBar)(t), keyc, C.TW_TYPE_QUAT4D, unsafe.Pointer(&s[0]), docc))
		default:
			panic("TwAddVarRW: bad slice length")
		}
	case *[4]uint8:
		s := avar.(*[4]uint8)
		return int(C.TwAddVarRW((*C.TwBar)(t), keyc, C.TW_TYPE_COLOR32, unsafe.Pointer(&s[0]), docc))
	default:
		fmt.Printf("type=%T\n", avar)
		panic("TwAddVarRW: unknown type")
	}
}

func TwDraw() int {
	return int(C.TwDraw())

}

func TwWindowSize(width, height int) int {
	return int(C.TwWindowSize((C.int)(width), (C.int)(height)))
}

func TwEventKeyGLFW(key, action int) int {
	return int(C.TwEventKeyGLFW((C.int)(key), (C.int)(action)))
}

func TwEventCharGLFW(char, action int) int {
	return int(C.TwEventCharGLFW((C.int)(char), (C.int)(action)))
}

func TwEventMouseButtonGLFW(button glfw.MouseButton, action glfw.Action) int {
	// ret := int(C.TwEventMouseButtonGLFW((C.int)(button), (C.int)(action)))
    var twa C.TwMouseAction
    var ret int
    if action == glfw.Press {
    	twa = TW_MOUSE_PRESSED
    } else {
    	twa = TW_MOUSE_RELEASED
    }
    switch button {
    case glfw.MouseButtonLeft:
    	ret = int(C.TwMouseButton(twa, TW_MOUSE_LEFT))
	case glfw.MouseButtonMiddle:
		ret = int(C.TwMouseButton(twa, TW_MOUSE_MIDDLE))
	case glfw.MouseButtonRight:
		ret = int(C.TwMouseButton(twa, TW_MOUSE_RIGHT))
    }
	//ret := int(C.TwMouseButton((C.TwMouseAction)(action), (C.TwMouseButtonID)(button)))
	//fmt.Printf("TwMouseButton: button=%d, action=%d, ret=%d\n", button, action, ret)
	return ret
}

func TwEventMousePosGLFW(x, y int) int {
	// ret := int(C.TwEventMouseButtonGLFW((C.int)(x), (C.int)(y)))
	ret := int(C.TwMouseMotion((C.int)(x), (C.int)(y)))
	//fmt.Printf("TwEventMousePosGLFW: x=%d, y=%d, ret=%d\n", x, y, ret)
	return ret
}

func Terminate() int {
	return int(C.TwTerminate())
}

/*
func Char(w *glfw.Window, c uint) {
	fmt.Printf("c=%d\n", c)
}

func Mouse(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
    fmt.Printf("mouse: button=%v, action=%v, mod=%v\n", button, action, mod)
}

func Cursor(w *glfw.Window, xpos float64, ypos float64) {
    //fmt.Printf("xpos=%.3f, ypos=%.3f\n", xpos, ypos)

}

func Enter(w *glfw.Window, entered bool) {
    fmt.Printf("entered=%v\n", entered)
}

func Scroll(w *glfw.Window, xoff float64, yoff float64) {
    fmt.Printf("xoff=%.2f, yoff=%.2f\n", xoff, yoff)
}


	MouseButtonHolder func(w *Window, button MouseButton, action Action, mod ModifierKey)
	CursorPosHolder   func(w *Window, xpos float64, ypos float64)
	CursorEnterHolder func(w *Window, entered bool)
	ScrollHolder      func(w *Window, xoff float64, yoff float64)
	KeyHolder         func(w *Window, key Key, scancode int, action Action, mods ModifierKey)
	CharHolder        func(w *Window, char uint)


func SetAllTWCallBacks(w *glfw.Window) { // unsafe.Pointer *glfw.Window
	w.SetCharacterCallback(Char)
	//C._SetAllTWCallBacks((*C.GLFWwindow)(w))
}
*/


/*
//export goMouseButtonCB
func goMouseButtonCB(window unsafe.Pointer, button, action, mods C.int) {
        w := windows.get((*C.GLFWwindow)(window))
        w.fMouseButtonHolder(w, MouseButton(button), Action(action), ModifierKey(mods))
}

//export goCursorPosCB
func goCursorPosCB(window unsafe.Pointer, xpos, ypos C.float) {
        w := windows.get((*C.GLFWwindow)(window))
        w.fCursorPosHolder(w, float64(xpos), float64(ypos))
}

//export goCursorEnterCB
func goCursorEnterCB(window unsafe.Pointer, entered C.int) {
        w := windows.get((*C.GLFWwindow)(window))
        hasEntered := glfwbool(entered)
        w.fCursorEnterHolder(w, hasEntered)
}

//export goScrollCB
func goScrollCB(window unsafe.Pointer, xoff, yoff C.float) {
        w := windows.get((*C.GLFWwindow)(window))
        w.fScrollHolder(w, float64(xoff), float64(yoff))
}

//export goKeyCB
func goKeyCB(window unsafe.Pointer, key, scancode, action, mods C.int) {
        w := windows.get((*C.GLFWwindow)(window))
        w.fKeyHolder(w, Key(key), int(scancode), Action(action), ModifierKey(mods))
}

//export goCharCB
func goCharCB(window unsafe.Pointer, character C.uint) {
        w := windows.get((*C.GLFWwindow)(window))
        w.fCharHolder(w, uint(character))
}
*/

//void _SetAllTWCallBacks(GLFWwindow *window);
//void glfwSetKeyCallbackCB(GLFWwindow *window);
//void glfwSetCharCallbackCB(GLFWwindow *window);
//void glfwSetMouseButtonCallbackCB(GLFWwindow *window);
//void glfwSetCursorPosCallbackCB(GLFWwindow *window);
//void glfwSetCursorEnterCallbackCB(GLFWwindow *window);
//void glfwSetScrollCallbackCB(GLFWwindow *window);
//float GetAxisAtIndex(float *axis, int i);
//unsigned char GetButtonsAtIndex(unsigned char *buttons, int i);

