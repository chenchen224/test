package model

import (
	"unsafe"

	"github.com/minio/minio-go/v7"
	"github.com/spf13/viper"
	cfg "gitlab.mvalley.com/datapack/cain/pkg/config"
)

type Response struct {
	Score       float64
	PramaryName string
	PersonName  string
}

type MigrateConfig struct {
	ESConfig            cfg.ESConfiguration
	QuickSearchESConfig cfg.ESConfiguration
	PEVCMySQLConfig     cfg.MySQLConfiguration
	QuickSearchBoost    QuickSearchBoostConfiguration
}

type QuickSearchBoostConfiguration struct {
	PrimaryName  int
	KeyWord      int
	Description  int
	RankingScore float32
}

type CompanyKeyWord struct {
	CompanyName string
	PrimaryName string
}

func (c *MigrateConfig) Read(configName string, configPath string) error {
	vp := viper.New()
	vp.SetConfigName(configName)
	vp.AutomaticEnv()
	vp.AddConfigPath(configPath)

	if err := vp.ReadInConfig(); err != nil {
		return err
	}

	err := vp.Unmarshal(c)
	if err != nil {
		return err
	}

	return nil
}

type MyObjectInfo struct {
	minio.ObjectInfo
}

type SliceMock struct {
	addr uintptr
	len  int
	cap  int
}

func (o *MyObjectInfo) TransLateToBytes() []byte {
	Len := unsafe.Sizeof(*o)
	data := &SliceMock{
		addr: uintptr(unsafe.Pointer(o)),
		cap:  int(Len),
		len:  int(Len),
	}
	bytes := *(*[]byte)(unsafe.Pointer(data))
	return bytes
}
