Assembler
=========

# Registers

Registers are specified with 'r' + index (ie. r0, r1, r16, r31).

## Nicknames

In addition to being specified with 'r' + index, the special purpose registers can also be specified by their nickname;

- rZero - r0 - Zero Register
- rOne - r1 - One Register
- rCID - r2 - Core Identifier Register
- rXID - r3 - Context Identifier Register
- rIVT - r4 - Interrupt Vector Table Register
- rPID - r5 - Process Identifier Register
- rPC - r6 - Program Counter Register
- rPS - r7 - Program Start Register
- rPL - r8 - Program Limit Register
- rSP - r9 - Stack Pointer Register
- rSS - r10 - Stack Start Register
- rSL - r11 - Stack Limit Register
- rDS - r12 - Data Start Register
- rDL - r13 - Data Limit Register

# Instructions

## Bitwise

All bitwise instructions take three registers; source register 1, source register 2, and destination register. 'not' is the exception, which only takes two registers; source register, and destination register.

### Not

```
not r16 r18
```

### And

```
and r16 r17 r18
```

### Or

```
or r16 r17 r18
```

### Xor

```
xor r16 r17 r18
```

### Left Shift

```
leftshift r16 r17 r18
```

### Right Shift

```
rightshift r16 r17 r18
```

## Arithmetic

All arithmetic instructions take three registers; source register 1, source register 2, and destination register.

### Add

```
add r16 r17 r18
```

### Subtract

```
subtract r16 r17 r18
```

### Multiply

```
multiply r16 r17 r18
```

### Divide

```
divide r16 r17 r18
```

### Modulo

```
modulo r16 r17 r18
```

## Control Flow

### Conditional Jump

The four conditional jumps are based on two conditions; 'equal to zero', and 'less than zero'.

```
jez r16 #Label
jnz r16 #Label
jle r16 #Label
jlz r16 #Label
```

### Call

```
call r31
```

### Return

```
return
```

## Data Movement

### Load Constant

Load Constant puts the given value into the specified destination register.

```
loadc 93 r16                // Decimal Literal
loadc 0xF8 r16              // Hexadecimal Literal
loadc 0b10 r16              // Binary Literal
loadc FOOBAR r16            // Constant Value
loadc #Label r16            // Label Address
```

### Load

Load reads the value at the given address plus offset into the specified destination register.

```
load r16 93 r17              // Decimal Literal Offset
load r16 0xF8 r17            // Hexadecimal Literal Offset
load r16 0b10 r17            // Binary Literal Offset
load r16 FOOBAR r17          // Constant Value Offset
load r16 #Label r17          // Label Address Offset
```

### Store

Store writes the value at the given address plus offset from the specified source register.

```
store r16 93 r17             // Decimal Literal Offset
store r16 0xF8 r17           // Hexadecimal Literal Offset
store r16 0b10 r17           // Binary Literal Offset
store r16 FOOBAR r17         // Constant Value Offset
store r16 #Label r17         // Label Address Offset
```

### Clear

Clear invalidates the value in the cache at the given address plus offset.
```
clear r16 93                 // Decimal Literal Offset
clear r16 0xF8               // Hexadecimal Literal Offset
clear r16 0b10               // Binary Literal Offset
clear r16 FOOBAR             // Constant Value Offset
clear r16 #Label             // Label Address Offset
```

### Flush

Flush writes the value in the cache at the given address plus offset back to main memory.

```
flush r16 93                 // Decimal Literal Offset
flush r16 0xF8               // Hexadecimal Literal Offset
flush r16 0b10               // Binary Literal Offset
flush r16 FOOBAR             // Constant Value Offset
flush r16 #Label             // Label Address Offset
```

### Push

Push writes the set of general purpose registers to the stack.

Registers must be specified in ascending order.

```
push r16,r17,r18,r19
```

### Pop

Pop read the set of general purpose registers from the stack.

Registers must be specified in decending order.

```
pop r19,r18,r17,r16
```

## Special

### Halt

```
halt
```

### Noop

```
noop
```

### Sleep

```
sleep
```

### Signal

```
signal r16
```

### Lock

```
lock
```

### Unlock

```
unlock
```

### Interrupt

```
interrupt 93                // Decimal Literal
interrupt 0xF8              // Hexadecimal Literal
interrupt 0b10              // Binary Literal
```

### Uninterrupt

```
uninterrupt
```

## Sugar

Sugar are statements supported by the assembler which aren't supported by the underlying architecture, instead the desired operation is achieved by another instruction.

### Copy

Sometimes called Move (weirdly).
Copy duplicates the data in the source register in the destination register - eg 'copy rS rD'.
Achieved by adding zero to the source register and storing the result in the destination register - eg 'add r0 rS rD'.

### Unconditional Jump

Unconditional jumps redirect program flow to the specified offset - eg 'jump #label'.
Achieved by conditionally jumping if the zero register is zero, which should always be true unless something terrible has gone wrong - eg 'jez r0 #label'.

# Constants

Constants associate an uppercase name with a value.

```
FOOBAR 93                   // Decimal Literal
FOOBAR 0xF8                 // Hexadecimal Literal
FOOBAR 0b10                 // Binary Literal
```

# Data

Data defines a memory location associated either with a value, or a label address.

```
data 93                     // Decimal Literal
data 0xF8                   // Hexadecimal Literal
data 0b10                   // Binary Literal

data #Label                 // Label Address Literal
```

# Padding

Padding inserts the given number of data statements into the output - eg 'padding 10' is equivalent to writting 'data 0' 10 times.

# Align

Align is used to produce a variable padding such that the next statement is given the specified address - eg 'align 0x200' will ensure the next statement has address 0x200 (512).
