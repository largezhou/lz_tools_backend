package service

import (
	"github.com/go-redis/redis/v8"
	"github.com/largezhou/lz_tools_backend/app/dto/code_dto"
	"strconv"
)

// CodeListSortable 场所码按照 redis 坐标排序对象
type CodeListSortable struct {
	CodeList []*code_dto.CodeListDto
	LocMap   map[string]redis.GeoLocation
}

func NewCodeListSortable(codeList []*code_dto.CodeListDto, locList []redis.GeoLocation) *CodeListSortable {
	c := &CodeListSortable{
		CodeList: codeList,
	}
	locMap := make(map[string]redis.GeoLocation)
	for _, location := range locList {
		locMap[location.Name] = location
	}

	c.LocMap = locMap

	return c
}

func (c CodeListSortable) Len() int {
	return len(c.CodeList)
}

func (c CodeListSortable) Less(i, j int) bool {
	iLoc, iOk := c.LocMap[strconv.Itoa(int(c.CodeList[i].Id))]
	jLoc, jOk := c.LocMap[strconv.Itoa(int(c.CodeList[j].Id))]
	if iOk {
		c.CodeList[i].Dist = iLoc.Dist
	}
	if jOk {
		c.CodeList[j].Dist = jLoc.Dist
	}

	// 都有位置信息，则距离小的在前面
	if iOk && jOk {
		return iLoc.Dist < jLoc.Dist
	}

	// 否则，有位置的在前面
	if iOk && !jOk {
		return true
	}

	if !iOk && jOk {
		return false
	}

	// 都没有位置信息，则相同
	return false
}

func (c CodeListSortable) Swap(i, j int) {
	c.CodeList[i], c.CodeList[j] = c.CodeList[j], c.CodeList[i]
}
