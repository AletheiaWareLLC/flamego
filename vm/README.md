# Virtual Machine

## Processor

- 64bit RISC Architecture
- Barrel Design
- 8 Cores per Processor
- 8 Contexts per Core
- 8-Stage Pipeline
    - Fetch Instruction
    - Load Instruction
    - Decode Instruction
    - Load Data
    - Execute Operation
    - Format Data
    - Store Data
    - Retire Instruction
- 32 Registers per Context
    - r0 : r15 - Special Purpose
    - r16 : r31 - General Purpose

## Cache

- 8 x 256KB L1 Instruction (1 per Core)
- 8 x 256KB L1 Data (1 per Core)
- 1 x 8MB L2 (1 shared between 8 Cores)

## Memory

- 1GB

## IO Devices

- Storage
- Display
- Keyboard
