package research

import "fmt"

type AddressSpace string

const (
	ROMCPU  AddressSpace = "ROMCPU"
	ROMFile AddressSpace = "ROMFILE"
	WRAM    AddressSpace = "WRAM"
	VRAM    AddressSpace = "VRAM"
	CGRAM   AddressSpace = "CGRAM"
	OAM     AddressSpace = "OAM"
	ARAM    AddressSpace = "ARAM"
	DSP     AddressSpace = "DSP"
)

type Address struct {
	Space  AddressSpace
	Offset uint32
}

func (a Address) String() string {
	return fmt.Sprintf("%s:$%X", a.Space, a.Offset)
}
