package serverfunc

import (
	"log"
	conf "et-config/src/statusconfig"
)

func ThirdPackageError(err error) (retcode int32, msg string) {
	log.Println(err)
	retcode = int32(conf.StatusThirdPackageErr)
	msg = err.Error()
	return
}