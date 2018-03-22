package create

import (
	"fmt"
	"os"
	"encoding/json"
	"io/ioutil"
	"github.com/g4stly/gofast/common"
)

/* leaf and branch types are used to represent
 * the target directory layout */
type leaf struct {
	name string
}

type branch struct {
	name		string
	files		[]leaf
	directories	[]branch
}

func (self *branch) print(prefix string) {
	prefix = fmt.Sprintf("%v%v->", prefix, self.name)
	for _, f := range self.files {
		common.Log("%v%v", prefix, f.name)
	}
	for _, d := range self.directories {
		d.print(prefix)
	}
}

func (self *branch) create(dirname string) {
	dirname = fmt.Sprintf("%v%v/", dirname, self.name)
	common.Log("creating directory `%v`", dirname)
	os.Mkdir(dirname, 0755)
	for _, leaf := range self.files {
		filename := fmt.Sprintf("%v%v", dirname, leaf.name)
		common.Log("creating file `%v`", filename)
		err := ioutil.WriteFile(filename, []byte("test"), 0755)
		if err != nil {
			common.Fatal("create: directoryTreeCreate(): %v", err)
		}
	}
	for _, dir := range self.directories {
		dir.create(dirname)
	}
}

/* creating the command type makes it easy to 
 * pass around the methods we actually want
 * to make available to other parts of our program,
 * the command type is simply the thing we're passing */
type command struct {
	directoryTree	branch
	targetName	string
	targetType	string
}

// this gets called from the outside world
func (self *command) Exec(args []string) int {
	switch (len(args)) {
	case 0:
		return self.Help()
		break
	case 1:
		self.targetType = "generic"
		self.targetName = args[0]
		return self.createNewProject()
		break
	}
	//default case for now
	return self.Help()
}

func (self *command) createNewProject() int {
	common.Log("creating new %v project with name %v", self.targetType, self.targetName)
	var directoryTree interface{}

	// read layout.json TODO: read json from a database or something
	filename := fmt.Sprintf("templates/%v.json", self.targetType)
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		common.Fatal("create: ReadFile(): %v", err)
	}

	// parse the json, put it in directoryTree
	err = json.Unmarshal(file, &directoryTree)
	if err != nil {
		common.Fatal("create: Unmarshal(): %v", err)
	}

	// parse the json (named tree) into our tree-ish thing (named directoryTree)
	tree := directoryTree.(map[string]interface{})
	self.directoryTree = jsonToBranch(self.targetName, tree)
	//self.directoryTree.print("")
	// now create those directories and touch those files
	self.directoryTree.create("")

	return 0
}

func jsonToBranch(name string, tree map[string]interface{}) branch {
	result := branch{name: name}
	for k, v := range tree {
		if k == "leaf" {
			file := leaf{name: v.(string)}
			result.files = append(result.files, file)
			continue
		}
		result.directories = append(result.directories, jsonToBranch(k, v.(map[string]interface{})))
	}
	return result
}

func (self *command) Help() int {
	common.Out("create <template> <name>")
	return 1
}

func New() common.Command {
	cmd := command{}
	return &cmd
}
