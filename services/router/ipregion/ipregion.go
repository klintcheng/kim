package ipregion

import (
	"github.com/lionsoul2014/ip2region/binding/golang/ip2region"
)

type IpInfo struct {
	Country string
	Region  string
	City    string
	ISP     string
}

type IpRegion interface {
	Search(ip string) (*IpInfo, error)
}

type Ip2region struct {
	region *ip2region.Ip2Region
}

func NewIp2region(path string) (IpRegion, error) {
	if path == "" {
		path = "ip2region.db"
	}
	region, err := ip2region.New(path)
	if err != nil {
		return nil, err
	}

	return &Ip2region{
		region: region,
	}, nil
}

func (r *Ip2region) Search(ip string) (*IpInfo, error) {
	info, err := r.region.MemorySearch(ip)
	if err != nil {
		return nil, err
	}
	return &IpInfo{
		Country: info.Country,
		Region:  info.Region,
		City:    info.City,
		ISP:     info.ISP,
	}, nil
}
