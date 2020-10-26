package gorm_v2

import (
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
	"goskeleton/app/global/variable"
	"time"
)

func getSqlserverDriver() (*gorm.DB, error) {
	writeDb := getDsn("sqlserver", "Write")
	gormDb, err := gorm.Open(sqlserver.Open(writeDb), &gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 redefineLog(), //本项目骨架接管 gorm v2 自带日志
	})
	if err != nil {
		//gorm 数据库驱动初始化失败
		return nil, err
	}

	// 如果开启了读写分离，配置读数据库（resource、read、replicas）
	if variable.ConfigGormv2Yml.GetInt("Gormv2.SqlServer.IsOpenReadDb") == 1 {
		readDb := getDsn("SqlServer", "Read")
		err := gormDb.Use(dbresolver.Register(dbresolver.Config{
			//Sources:  []gorm.Dialector{sqlserver.Open(writeDb)}, //  写 操作库， 执行类
			Replicas: []gorm.Dialector{sqlserver.Open(readDb)}, //  读 操作库，查询类
			Policy:   dbresolver.RandomPolicy{},                // sources/replicas 负载均衡策略适用于
		}, "").SetConnMaxIdleTime(time.Minute).
			SetConnMaxLifetime(variable.ConfigGormv2Yml.GetDuration("Gormv2.SqlServer.SetConnMaxLifetime") * time.Second).
			SetMaxIdleConns(variable.ConfigGormv2Yml.GetInt("Gormv2.SqlServer.SetMaxIdleConns")).
			SetMaxOpenConns(variable.ConfigGormv2Yml.GetInt("Gormv2.SqlServer.SetMaxOpenConns")))
		if err != nil {
			return nil, err
		}
	}
	return gormDb, nil
}
