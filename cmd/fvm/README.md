Flame Virtual Machine
=====================

# Install

```
go install aletheiaware.com/flamego/cmd/fvm
```

# Usage

Invoke the virtual machine only.

```
fvm
```

Invoke the virtual machine with the given bootloader in memory.

```
fvm -b bootloader.bin
```

Invoke the virtual machine with the given file in storage.

```
fvm -s storage.bin
```

Invoke the virtual machine with the given bootloader in memory and kernel in storage.

```
fvm -b bootloader.bin -s kernel.bin
```
