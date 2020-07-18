package replays

import "fmt"

type ResourceGuid uint64

type resourceFlag uint64

const (
	Index    resourceFlag = 0x00000000FFFFFFFF
	Locale                = 0x0000001F00000000
	Reserved              = 0x0000006000000000
	Region                = 0x00000F8000000000
	Platform              = 0x0000F00000000000
	Type                  = 0x0FFF000000000000
	Engine                = 0xF000000000000000
	Key                   = 0x0000FFFFFFFFFFFF
)

func (guid ResourceGuid) Attribute(flag resourceFlag) uint64 {
	return uint64(guid) & uint64(flag)
}

func (guid ResourceGuid) Key() uint64 {
	return guid.Attribute(Key)
}

func (guid ResourceGuid) Type() uint16 {
	f := guid.Attribute(Type) >> 48
	f = flipBits(f)

	return uint16(f) + 1
}

func (guid ResourceGuid) String() string {
	return fmt.Sprintf("%012X.%03X", guid.Key(), guid.Type())
}

func flipBits(num uint64) uint64 {
	num = ((num >> 1) & 0x55555555) | ((num & 0x55555555) << 1)
	num = ((num >> 2) & 0x33333333) | ((num & 0x33333333) << 2)
	num = ((num >> 4) & 0x0F0F0F0F) | ((num & 0x0F0F0F0F) << 4)
	num = ((num >> 8) & 0x00FF00FF) | ((num & 0x00FF00FF) << 8)
	num = (num >> 16) | (num << 16)
	num >>= 20
	return num
}
