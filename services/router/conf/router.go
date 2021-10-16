package conf

import (
	"encoding/json"
	"io/ioutil"
)

type IDC struct {
	ID     string
	Weight int
}

type Region struct {
	ID    string
	Idcs  []IDC
	Slots []byte
}

type Country string

type Mapping struct {
	Region    string
	Locations []string
}

type Router struct {
	Mapping map[Country]string
	Regions map[string]*Region
}

func LoadMapping(path string) (map[Country]string, error) {
	bts, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var mps []Mapping
	err = json.Unmarshal(bts, &mps)
	if err != nil {
		return nil, err
	}
	mp := make(map[Country]string)
	for _, v := range mps {
		region := v.Region
		for _, loc := range v.Locations {
			mp[Country(loc)] = region
		}
	}
	return mp, nil
}

func LoadRegions(path string) (map[string]*Region, error) {
	bts, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var regions []*Region
	err = json.Unmarshal(bts, &regions)
	if err != nil {
		return nil, err
	}
	res := make(map[string]*Region)
	for _, region := range regions {
		res[region.ID] = region
		for i, idc := range region.Idcs {
			// 1.通过权重生成分片中的slots
			shard := make([]byte, idc.Weight)
			// 2. 给当前slots设置值，指向索引i
			for j := 0; j < idc.Weight; j++ {
				shard[j] = byte(i)
			}
			// 2. 追加到Slots中
			region.Slots = append(region.Slots, shard...)
		}
	}
	return res, nil
}
