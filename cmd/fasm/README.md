Flame Assembler
===============

# Install

```
go install aletheiaware.com/flamego/cmd/fasm
```

# Usage

Assemble the given file(s) only.

```
fasm file.asm
```

Assemble the given file(s) and log verbose information to standard out.

```
fasm -v file.asm
```

Assemble the given file(s) and write the binary to the given file.

```
fasm -o file.bin file.asm
```

Assemble the given file(s) and write the addresses to the given file.

```
fasm -a file.address file.asm
```
