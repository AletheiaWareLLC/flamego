# Instruction Set Architecture

32-bit instructions

LoadConstant:           1CCCCCCC CCCCCCCC CCCCCCCC CCCDDDDD
Jump:                   01CCBOOO OOOOOOOO OOOOOOOO OOORRRRR
Load/Store:             001TTOOO OOOOOOOO OOOOOOAA AAARRRRR
Bitwise & Arithmetic:   0001TTTT -------- -2222211 111DDDDD
Reserved:               00001--- -------- -------- --------
Push/Pop:               000001T- -------- MMMMMMMM MMMMMMMM
Call:                   00000010 -------- -------- ---AAAAA
Return:                 00000011 -------- -------- --------
Special:                00000001 TTTT---- -------- --------

## Bitwise

Assembly: operation source1 source2 destination
Opcode: 00010TTT -------- -2222211 111DDDDD

T: type;
 - 000 - Not
 - 001 - And
 - 010 - Or
 - 011 - Xor
 - 100 - Left Shift
 - 101 - Right Shift
 - 110 - Reserved
 - 111 - Reserved

1: first source register

2: second source register

D: destination register

### Not

Performs a bitwise not.

```
register[destination] = !register[source1]
```

### And

Performs a bitwise and.

```
register[destination] = register[source1] & register[source2]
```

### Or

Performs a bitwise or.

```
register[destination] = register[source1] | register[source2]
```

### Xor

Performs a bitwise xor.

```
register[destination] = register[source1] ^ register[source2]
```

### Left Shift

Performs a left logical shift.

```
register[destination] = register[source1] << register[source2]
```

### Right Shift

Performs a right logical shift.

```
register[destination] = register[source1] >> register[source2]
```

## Arithmetic

Assembly: operation source1 source2 destination
Opcode: 00011TTT -------- -2222211 111DDDDD

T: type;
 - 000 - Add
 - 001 - Subtract
 - 010 - Multiply
 - 011 - Divide
 - 100 - Modulo
 - 101 - Reserved
 - 110 - Reserved
 - 111 - Reserved

1: first source register

2: second source register

D: destination register

### Add

Performs an addition.

```
register[destination] = register[source1] + register[source2]
```

### Subtract

Performs a subtraction.

```
register[destination] = register[source1] - register[source2]
```

### Multiply

Performs a multiplication.

```
register[destination] = register[source1] * register[source2]
```

### Divide

Performs a division.

```
register[destination] = register[source1] / register[source2]
```

Triggers InterruptArithmeticError if contents of source2 is 0.

### Modulo

Performs a modulo.

```
register[destination] = register[source1] % register[source2]
```

Triggers InterruptArithmeticError if contents of source2 is 0.

## Control Flow

### Jump

Assembly: operation conditional addressoffset
Opcode: 01CCBOOO OOOOOOOO OOOOOOOO OOORRRRR

C: conditioncode:
 - 00: ez: jump if regster[conditional] is equal to zero (all bits zero)
 - 01: nz: jump if regster[conditional] is not equal to zero (not all bits zero)
 - 10: le: jump if regster[conditional] is less than or equal to zero (most significant bit set or all bits zero)
 - 11: lz: jump if regster[conditional] is less than zero (most significant bit set)

R: conditional register

B: backwards
 - set if addressoffset should be subtracted from programcounter

O: addressoffset
 - 22bit

```
if condition {
    if backwards {
        programcounter = programcounter - addressoffset
    } else {
        programcounter = programcounter + addressoffset
    }
}
```

### Call

Assembly: call addressregister
Opcode: 00000010 -------- -------- ---AAAAA

A: address register

```
stack[stackpointer] = programcounter + instructionsize
stackpointer += datasize
programcounter = register[address]
```

Pushes the address of the next instruction (sum of RProgramCounter and InstructionSize) onto the stack.

Sets RProgramCounter to the contents of address register.

### Return

Assembly: return
Opcode: 00000011 -------- -------- --------

```
stackpointer -= datasize
programcounter = stack[stackpointer]
```

Pops the return address off the stack.

Sets RProgramCounter to the return address.

## Data Movement

### Load Constant

Assembly: loadc destination constant
Opcode: 1CCCCCCC CCCCCCCC CCCCCCCC CCCDDDDD

D: destination register

C: constant
 - 26bit

```
register[destination] = constant
```

### Load

Assembly: load address offset destination
Opcode: 00100OOO OOOOOOOO OOOOOOAA AAARRRRR

A: address register

O: offset
 - 17bit

D: destination register

```
register[destination] = memory[register[address] + offset]
```

Retryable if L1 Data Cache is unavailable or unsuccessful (cache miss).

### Store

Assembly: store address offset source
Opcode: 00101OOO OOOOOOOO OOOOOOAA AAASSSSS

A: address register

O: offset
 - 17bit

S: source register

```
memory[register[address] + offset] = register[source]
```

Retryable if L1 Data Cache is unavailable or unsuccessful (cache miss).

### Clear

Assembly: clear address offset
Opcode: 00110OOO OOOOOOOO OOOOOOAA AAA00000

A: address register

O: offset
 - 17bit

Invalidates the data stored in the cache(s) at the address, forcing the next load to come from main memory.

Retryable if L1 Instruction, L1 Data, or L2 Cache is unavailable or unsuccessful (cache miss).

### Flush

Assembly: flush address offset
Opcode: 00111OOO OOOOOOOO OOOOOOAA AAA00000

A: address register

O: offset
 - 17bit

Writes the data stored in the cache(s) at the address to main memory.

Retryable if L1 Data, or L2 Cache is unavailable or unsuccessful (cache miss).

### Push

Assembly: push register
Opcode: 0000010- -------- MMMMMMMM MMMMMMMM

M: register mask

Pushes the given registers onto stack.

Retryable while each register in mask is pushed.

Retryable if L1 Data Cache is unavailable or unsuccessful (cache miss).

Triggers InterruptStackOverflowError if RStackPointer > RStackLimit.

### Pop

Assembly: pop register
Opcode: 0000011- -------- MMMMMMMM MMMMMMMM

M: register mask

Pops the given registers from stack.

Retryable while each register in mask is popped.

Retryable if L1 Data Cache is unavailable or unsuccessful (cache miss).

Triggers InterruptStackUnderflowError if RStackPointer < RStackStart.

## Special

### Halt

Assembly: halt
Opcode: 00000001 0000---- -------- --------

Halts the processor.

Only callable during an interrupt - triggers InterruptUnsupportedOperationError otherwise.

### Noop

Assembly: noop
Opcode: 00000001 0001---- -------- --------

Does nothing

### Sleep

Assembly: sleep
Opcode: 00000001 0010---- -------- --------

Puts the processor to sleep, to be awoken by the next signal.

Only callable during an interrupt - triggers InterruptUnsupportedOperationError otherwise.

Takes context out of interrupted state as only an uninterrupted context can be interrupted.

### Signal

Assembly: signal device
Opcode: 00000001 0011---- -------- ---DDDDD

D: device register

Sends a signal the given device.

Device addressing;
 - 0-7 core
 - 8-65535 io device

Only callable during an interrupt - triggers InterruptUnsupportedOperationError otherwise.

### Lock

Assembly: lock
Opcode: 00000001 0100---- -------- --------

Acquires the hardware lock.

Only callable during an interrupt - triggers InterruptUnsupportedOperationError otherwise.

Retryable if lock is not acquired.

### Unlock

Assembly: unlock
Opcode: 00000001 0101---- -------- --------

Releases the hardware lock.

Only callable during an interrupt - triggers InterruptUnsupportedOperationError otherwise.

Retryable if lock is not released.

### Interrupt

Assembly: interrupt addressregister
Opcode: 00000001 0110---- IIIIIIII IIIIIIII

I: interrupt identifier
 - 16bit

```
programcounter = interruptvectortable + interruptidentifier
```

### Uninterrupt

Assembly: uninterrupt addressregister
Opcode: 00000001 0111---- -------- 000AAAAA

A: address register

```
programcounter = register[address]
```

Only callable during an interrupt - triggers InterruptUnsupportedOperationError otherwise.
