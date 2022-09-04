package v1

import (
	"fmt"
	"io"

	"github.com/lunixbochs/struc"
)

type PriceType struct {
	value int
}

func (c *PriceType) Pack(p []byte, opt *struc.Options) (int, error) {
	return 0, nil
}
func (c *PriceType) readOneByte(r io.Reader) byte {
	data := make([]byte, 1)
	r.Read(data)
	return data[0]
}
func (c *PriceType) Unpack(r io.Reader, length int, opt *struc.Options) error {
	posbyte := 6
	bdata := c.readOneByte(r)
	intdata := int(bdata & 0x3f)
	sign := (bdata & 0x40) != 0

	if (bdata & 0x80) != 0 {
		for {
			bdata = c.readOneByte(r)
			intdata += (int(bdata&0x7f) << posbyte)
			posbyte += 7
			if (bdata & 0x80) == 0 {
				break
			}
		}
	}
	if sign {
		intdata = -intdata
	}

	c.value = intdata
	return nil
}
func (c *PriceType) Size(opt *struc.Options) int {
	return -1
}
func (c *PriceType) String() string {
	return fmt.Sprintf("%v", c.value)
}
func (c *PriceType) getValue() int {
	return c.value
}
