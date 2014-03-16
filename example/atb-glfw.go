//  @file       TwSimpleGLFW.c
//  @brief      A simple example that uses AntTweakBar with 
//              OpenGL and the GLFW windowing system.
//  @author     Philippe Decaudin
//  @date       2006/05/20

package main

import (
    "fmt"
    "os"
    "github.com/go-gl/gl"
    "github.com/go-gl/glu"
    glfw "github.com/go-gl/glfw3"
    //"github.com/tildeleb/glam"
    "leb/atb"
)

// This example program draws a possibly transparent cube 
func DrawModel(wireframe bool) {
    var numPass int

    //fmt.Printf("DrawModel: wireframe=%v\n", wireframe)
    // Enable OpenGL transparency and light (could have been done once at init)
    gl.Enable(gl.BLEND)
    gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
    gl.Enable(gl.DEPTH_TEST)
    gl.Enable(gl.LIGHT0)   // use default light diffuse and position
    gl.LightModeli(gl.LIGHT_MODEL_TWO_SIDE, 1)
    gl.Enable(gl.COLOR_MATERIAL)
    gl.ColorMaterial(gl.FRONT_AND_BACK, gl.DIFFUSE)
    gl.Enable(gl.LINE_SMOOTH)
    gl.LineWidth(3.0)
    
    if wireframe {
        gl.Disable(gl.CULL_FACE)  
        gl.Disable(gl.LIGHTING)
        gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
        numPass = 1
    } else {
        gl.Enable(gl.CULL_FACE)
        gl.Enable(gl.LIGHTING)
        gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
        numPass = 2
    }

    for pass := 0; pass < numPass; pass++ {
        var face gl.GLenum
        // Since the material could be transparent, we draw the convex model in 2 passes:
        // first its back faces, and second its front faces.
        if pass == 0 {
            face = gl.FRONT
        } else {
            face = gl.BACK  
        }
        gl.CullFace(face)

        // Draw the model (a cube)
        gl.Begin(gl.QUADS)
            gl.Normal3f(0,0,-1); gl.Vertex3f(0,0,0); gl.Vertex3f(0,1,0); gl.Vertex3f(1,1,0); gl.Vertex3f(1,0,0); // front face
            gl.Normal3f(0,0,+1); gl.Vertex3f(0,0,1); gl.Vertex3f(1,0,1); gl.Vertex3f(1,1,1); gl.Vertex3f(0,1,1); // back face
            gl.Normal3f(-1,0,0); gl.Vertex3f(0,0,0); gl.Vertex3f(0,0,1); gl.Vertex3f(0,1,1); gl.Vertex3f(0,1,0); // left face
            gl.Normal3f(+1,0,0); gl.Vertex3f(1,0,0); gl.Vertex3f(1,1,0); gl.Vertex3f(1,1,1); gl.Vertex3f(1,0,1); // right face
            gl.Normal3f(0,-1,0); gl.Vertex3f(0,0,0); gl.Vertex3f(1,0,0); gl.Vertex3f(1,0,1); gl.Vertex3f(0,0,1); // bottom face
            gl.Normal3f(0,+1,0); gl.Vertex3f(0,1,0); gl.Vertex3f(0,1,1); gl.Vertex3f(1,1,1); gl.Vertex3f(1,1,0); // top face
        gl.End()
    }
}

// Callback function called by GLFW when window size changes

func WindowSize(window *glfw.Window, width, height int) {
    fmt.Printf("WindowSizeCB: width=%d, height=%d\n", width, height)
    // Set OpenGL viewport and camera
    gl.Viewport(0, 0, width, height)
    gl.MatrixMode(gl.PROJECTION)
    gl.LoadIdentity()
    ar := float64(width)/float64(height)
    glu.Perspective(40, ar, 1, 10) // ???
    glu.LookAt(-1,0,3, 0,0,0, 0,1,0)
    
    // Send the new window size to AntTweakBar
    atb.TwWindowSize(width, height)
}

/*
    TW_API int TW_CDECL_CALL TwEventMouseButtonGLFWcdecl(int glfwButton, int glfwAction);
    TW_API int TW_CDECL_CALL TwEventKeyGLFWcdecl(int glfwKey, int glfwAction);
    TW_API int TW_CDECL_CALL TwEventCharGLFWcdecl(int glfwChar, int glfwAction);
    TW_API int TW_CDECL_CALL TwEventMousePosGLFWcdecl(int mouseX, int mouseY);
    TW_API int TW_CDECL_CALL TwEventMouseWheelGLFWcdecl(int wheelPos);
*/

func Char(w *glfw.Window, char uint) {
    //fmt.Printf("Char: char=%d\n", char)
    atb.TwEventCharGLFW(int(char), int(0))
}

func Key(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mod glfw.ModifierKey) {
    //fmt.Printf("Key: key=%v, scancode=%d, action=%v, mod=%v\n", key, scancode, action, mod)
    atb.TwEventKeyGLFW(int(key), int(action))
}

func MouseButton(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
    //fmt.Printf("MouseButton: button=%v, action=%v, mod=%v\n", button, action, mod)
    atb.TwEventMouseButtonGLFW(button, action)
}

func Cursor(w *glfw.Window, xpos float64, ypos float64) {
    //fmt.Printf("Cursor: xpos=%.3f, ypos=%.3f\n", xpos, ypos)
    atb.TwEventMousePosGLFW(int(xpos), int(ypos))
}

func Enter(w *glfw.Window, entered bool) {
    //fmt.Printf("Enter: entered=%v\n", entered)
}

func Scroll(w *glfw.Window, xoff float64, yoff float64) {
    fmt.Printf("Scroll: xoff=%.2f, yoff=%.2f\n", xoff, yoff)
    //atb.TwMouseWheelGLFW(xoff, yoff)
}


// Main
func main() {
    // GLFWvidmode mode;        // GLFW video mode
    var bar *atb.TwBar          // Pointer to a tweak bar
    var time, dt float64        // Current time and enlapsed time
    var turn float64            // Model turn counter
    var speed float64 = 0.3     // Model rotation speed
    var wire bool               // Draw model in wireframe?
    var bgColor = [4]float32{0.1, 0.2, 0.4}          // Background color 
    var cubeColor = [4]uint8{255, 0, 0, 128}         // Model color (32bits RGBA)
    var tw64 int64

    // Intialize GLFW   
    if !glfw.Init() {
        // An error occured
        fmt.Fprintf(os.Stderr, "GLFW initialization failed\n")
        return
    }

    //glfw.SetErrorCallback(errorCallback)

    if !glfw.Init() {
        panic("Can't init glfw!")
    }
    defer glfw.Terminate()

    // Create a window
    window, err := glfw.CreateWindow(640, 480, "AntTweakBar simple example using GLFW Go version", nil, nil)
    if err != nil {
        panic(err)
    }

    //glfw.Enable(glfw.MOUSE_CURSOR)
    //glfw.Enable(glfw.KEY_REPEAT)
    window.MakeContextCurrent()

    gl.Init()

    // Initialize AntTweakBar
    res := atb.TwInit("opengl", "")
    fmt.Printf("res=%d\n", res)

    atb.TwWindowSize(640, 480)

    // Create a tweak bar
    bar = atb.TwNewBar("TweakBar")
    fmt.Printf("bar=%T\n", bar)
    atb.TwDefine(" GLOBAL help='This example shows how to integrate AntTweakBar with GLFW and OpenGL.' ") // Message added to the help bar.

    // Add 'speed' to 'bar': it is a modifable (RW) variable of type TW_TYPE_DOUBLE. Its key shortcuts are [s] and [S].
    bar.AddVarRW("speed", &speed, " label='Rot speed' min=0 max=2 step=0.01 keyIncr=s keyDecr=S help='Rotation speed (turns/second)' ");

    // Add 'wire' to 'bar': it is a modifable variable of type TW_TYPE_BOOL32 (32 bits boolean). Its key shortcut is [w].
    bar.AddVarRW("wire", &wire, " label='Wireframe mode' key=w help='Toggle wireframe display mode.' ")

    bar.AddVarRW("tw64", &tw64, " label='tw64 test.' ")

    // Add 'time' to 'bar': it is a read-only (RO) variable of type TW_TYPE_DOUBLE, with 1 precision digit
    bar.AddVarRW("time", &time, " label='Time' precision=1 help='Time (in seconds).' ") // was RO

    // Add 'bgColor' to 'bar': it is a modifable variable of type TW_TYPE_COLOR3F (3 floats color)
    bar.AddVarRW("bgColor", &bgColor, " label='Background color' ") // TW_TYPE_COLOR3F

    // Add 'cubeColor' to 'bar': it is a modifable variable of type TW_TYPE_COLOR32 (32 bits color) with alpha
    bar.AddVarRW("cubeColor", &cubeColor, " label='Cube color' alpha help='Color and transparency of the cube.' ")

    // Set GLFW event callbacks
    // - Redirect window size changes to the callback function WindowSizeCB

    window.SetFramebufferSizeCallback(WindowSize)
    window.SetCharacterCallback(Char)
    window.SetKeyCallback(Key)
    window.SetMouseButtonCallback(MouseButton)
    window.SetCursorPositionCallback(Cursor)
    window.SetCursorEnterCallback(Enter)
    window.SetScrollCallback(Scroll)

    //atb.SetAllTWCallBacks(window)

/*
    // - Directly redirect GLFW mouse button events to AntTweakBar
    glfw.SetMouseButtonCallback((GLFWmousebuttonfun)TwEventMouseButtonGLFW);
    // - Directly redirect GLFW mouse position events to AntTweakBar
    glfw.SetMousePosCallback((GLFWmouseposfun)TwEventMousePosGLFW);
    // - Directly redirect GLFW mouse wheel events to AntTweakBar
    glfw.SetMouseWheelCallback((GLFWmousewheelfun)TwEventMouseWheelGLFW);
    // - Directly redirect GLFW key events to AntTweakBar
    glfw.SetKeyCallback((GLFWkeyfun)TwEventKeyGLFW);
    // - Directly redirect GLFW char events to AntTweakBar
    glfw.SetCharCallback((GLFWcharfun)TwEventCharGLFW);
*/

    // Initialize time
    time = glfw.GetTime();

    // Main loop (repeated while window is not closed and [ESC] is not pressed)
    for !window.ShouldClose() && !(window.GetKey(glfw.KeyEscape) == glfw.Press) {
        // Clear frame buffer using bgColor
        //gl.ClearColor(gl.GLclampf(bgColor[0]), gl.GLclampf(bgColor[1]), gl.GLclampf(bgColor[2]), gl.GLclampf(1))
        gl.ClearColor(0, 0, 0, 0)
        gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT )

        // Rotate model
        dt = glfw.GetTime() - time
        if dt < 0 {
            dt = 0
        } 
        time += dt
        turn += speed*dt
        gl.MatrixMode(gl.MODELVIEW)
        gl.LoadIdentity()
        gl.Rotated(360.0*turn, 0.4, 1, 0.2)
        gl.Translated(-0.5, -0.5, -0.5)   
    
        // Set color and draw model
        gl.Color4ubv(&cubeColor)
        //fmt.Printf("main: wire=%v, tw64=%d\n", wire, tw64)
        DrawModel(wire)

        // Draw tweak bars
        atb.TwDraw()

        // Present frame buffer
        window.SwapBuffers()
        glfw.PollEvents()
        dt++
        turn++
    }

    // Terminate AntTweakBar and GLFW
    atb.Terminate()
    glfw.Terminate()
    return
}