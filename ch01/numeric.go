package main

import "fmt"

func main() {
    // Signed integers
    var a int8 = -128
    var b int16 = -32768
    var c int32 = -2147483648
    var d int64 = -9223372036854775808

    // Unsigned integers
    var ua uint8 = 255
    var ub uint16 = 65535
    var uc uint32 = 4294967295
    var ud uint64 = 18446744073709551615

    // Machine-dependent integers
    var e int = -42       // int (32-bit or 64-bit depending on architecture)
    var f uint = 42       // uint

    // Floating-point numbers
    var g float32 = 3.14159
    var h float64 = 2.718281828459045

    // Complex numbers
    var i complex64 = complex(1.5, -2.5)
    var j complex128 = complex(3.0, 4.0)

    // Aliases
    var k byte = 0xFF     // byte is alias for uint8
    var l rune = 'λ'      // rune is alias for int32 (Unicode code point)

    // Print all values
    fmt.Println("Signed Integers:")
    fmt.Println("int8:", a, "int16:", b, "int32:", c, "int64:", d)

    fmt.Println("Unsigned Integers:")
    fmt.Println("uint8:", ua, "uint16:", ub, "uint32:", uc, "uint64:", ud)

    fmt.Println("Machine-Dependent Integers:")
    fmt.Println("int:", e, "uint:", f)

    fmt.Println("Floating-Point Numbers:")
    fmt.Println("float32:", g, "float64:", h)

    fmt.Println("Complex Numbers:")
    fmt.Println("complex64:", i, "complex128:", j)

    fmt.Println("Aliases:")
    fmt.Printf("byte: %d (0x%X)\n", k, k)
    fmt.Printf("rune: %c (%U)\n", l, l)
}

