package ole

/*
Date：2023.03.02
Author：scl
Description：解析MiniFat结构
*/
import (
	"encoding/binary"
)

//获取minifat信息
func (ole *OleInfo) getMiniFATSectors() (err error) {
	var section = make([]byte, 0)

	position := ole.calculateOffset(ole.header.MiniFatSect[:])

	//循环读MiniFat，一般都是一个
	for i := uint32(0); i < binary.LittleEndian.Uint32(ole.header.NumberMiniFATSectors[:]); i++ {
		sector := NewSector(&ole.header)
		err := ole.getData(position, &sector.Data)

		if err != nil {
			return err
		}

		//循环获取数据
		for _, value := range sector.getMiniFatFATSectorLocations() {
			section = append(section, value)
			//4个字节一个数据
			if len(section) == 4 {
				ole.miniFatPositions = append(ole.miniFatPositions, binary.LittleEndian.Uint32(section))
				section = make([]byte, 0)
			}
		}
		position = position + sector.SectorSize
	}

	return
}
