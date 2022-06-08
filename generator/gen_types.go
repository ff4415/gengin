package generator

import (
	"fmt"
	"io"
	"os"
	"strings"

	. "github.com/MasterJoyHunan/go-zero-gin-plugin/prepare"
	"github.com/MasterJoyHunan/go-zero-gin-plugin/tpl"

	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

var groupTypes []TypeBelongGroup
var typeGroup = make(map[string]string)

type TypeBelongGroup struct {
	GroupName string
	TypeStr   string
	TypeMap   []spec.Type
}

func GenTypes() error {
	types, err := buildGroupTypes()
	if err != nil {
		return err
	}

	typeFilename := typesPacket + ".go"

	for _, t := range types {
		typeGroupInfo := parseGroupName(t.GroupName, typesDir, typesPacket)
		filename := pathx.JoinPackages(PluginInfo.Dir, typeGroupInfo.DirPath, typeFilename)
		os.Remove(filename)

		err = genFile(fileGenConfig{
			dir:             PluginInfo.Dir,
			subDir:          typeGroupInfo.DirPath,
			filename:        typeFilename,
			templateName:    "typesTemplate",
			builtinTemplate: tpl.TypesTemplate,
			data: map[string]interface{}{
				"pkgName": typeGroupInfo.PkgName,
				"types":   t.TypeStr,
			},
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// buildTypes gen types to string
func buildTypes(types []spec.Type) (string, error) {
	var builder strings.Builder
	first := true
	for _, tp := range types {
		if first {
			first = false
		} else {
			builder.WriteString("\n\n")
		}
		if err := writeType(&builder, tp); err != nil {
			return "", err
		}
	}

	return builder.String(), nil
}

// buildGroupTypes gen types to string
func buildGroupTypes() ([]TypeBelongGroup, error) {
	// 用于保存 type 被哪几个 group 用到
	container := make(map[string]map[string]int, 0)
	for _, group := range PluginInfo.Api.Service.Groups {
		for _, route := range group.Routes {
			joinContainer(container, route.RequestTypeName(), group.GetAnnotation(groupProperty))
			joinContainer(container, route.ResponseTypeName(), group.GetAnnotation(groupProperty))
		}
	}

	// 用于保存 group 下包含几个 type， 如果一个 type 被多个 group 用到了，则放入公共的 type 文件中
	// 如果一个 type 没有设定 group ，则也会放到 公共的 type 文件中
	groupType := make(map[string]map[string]int)
	for typeName, groups := range container {
		_, ok := groups[commonGroup]
		groupName := ""
		if ok || len(groups) >= 2 {
			groupName = commonGroup
		} else {
			for group := range groups {
				groupName = group
			}
		}
		_, ok = groupType[groupName]
		if !ok {
			groupType[groupName] = make(map[string]int)
		}
		groupType[groupName][typeName] = 1
	}

	// 包装返回
	for group, typeNames := range groupType {
		var temp []spec.Type
		for _, t := range PluginInfo.Api.Types {
			_, ok := typeNames[t.Name()]
			if ok {
				temp = append(temp, t)
			}
		}
		typeStr, err := buildTypes(temp)
		if err != nil {
			return nil, err
		}
		groupTypes = append(groupTypes, TypeBelongGroup{
			GroupName: group,
			TypeStr:   typeStr,
			TypeMap:   temp,
		})
	}

	// TODO
	for _, g := range groupTypes {
		for _, s := range g.TypeMap {
			if g.GroupName == commonGroup {
				typeGroup[s.Name()] = ""
				continue
			}
			typeGroup[s.Name()] = g.GroupName
		}
	}

	return groupTypes, nil
}

func joinContainer(container map[string]map[string]int, typeName string, group string) {
	if typeName == "" {
		return
	}
	if group == "" {
		group = commonGroup
	}
	_, ok := container[typeName]
	if !ok {
		container[typeName] = make(map[string]int, 0)
	}
	container[typeName][group] = 1
}

func writeType(writer io.Writer, tp spec.Type) error {
	structType, ok := tp.(spec.DefineStruct)
	if !ok {
		return fmt.Errorf("unspport struct type: %s", tp.Name())
	}

	fmt.Fprintf(writer, "type %s struct {\n", title.String(tp.Name()))
	for _, member := range structType.Members {
		if member.IsInline {
			if _, err := fmt.Fprintf(writer, "%s\n", title.String(member.Type.Name())); err != nil {
				return err
			}

			continue
		}
		tag := member.Tag
		before, _, found := strings.Cut(tag, ":")
		if found && strings.HasSuffix(before, "path") {
			tag = strings.Replace(tag, "path", "uri", 1)
		}
		if err := writeProperty(writer, member.Name, member.Tag, member.GetComment(), member.Type); err != nil {
			return err
		}
	}
	fmt.Fprintf(writer, "}")
	return nil
}

func writeProperty(writer io.Writer, name, tag, comment string, tp spec.Type) error {
	var err error
	if len(comment) > 0 {
		comment = strings.TrimPrefix(comment, "//")
		comment = "//" + comment
		_, err = fmt.Fprintf(writer, "%s %s %s %s\n", title.String(name), tp.Name(), tag, comment)
	} else {
		_, err = fmt.Fprintf(writer, "%s %s %s\n", title.String(name), tp.Name(), tag)
	}
	return err
}

//
//func getTypesImportAlias(pkg GroupInfo) string {
//	if pkg.DirPath == typesPacket {
//		return ""
//	}
//	return pkg.PkgName + typePkgAlias + " "
//}
//
//func getTypesUseAlias(pkg GroupInfo) string {
//	if pkg.DirPath == typesPacket {
//		return typesPacket + "."
//	}
//	return pkg.PkgName + typePkgAlias + "."
//}