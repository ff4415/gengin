package generator

import (
	"path"

	"github.com/MasterJoyHunan/gengin/prepare"
	"github.com/MasterJoyHunan/gengin/tpl"
	"github.com/zeromicro/go-zero/tools/goctl/util/format"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func GenSvcContext() error {
	for _, group := range prepare.ApiSpec.Service.Groups {
		subDir := group.GetAnnotation(groupProperty)
		subDir, err := format.FileNamingFormat(dirStyle, subDir)
		if err != nil {
			return err
		}

		handlePkg := path.Join("handler", subDir)
		handleBase := path.Base(handlePkg)

		fileName := `svc_` + handleBase + `_context.go`
		err = GenFile(fileName,
			tpl.SvcContextTemplate,
			WithSubDir("svc"),
			WithData(map[string]any{
				"svcName": cases.Title(language.English, cases.NoLower).String(handleBase),
			}),
		)
		if err != nil {
			return err
		}
	}
	return nil
}
