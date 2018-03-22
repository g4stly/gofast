package create

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"github.com/g4stly/gofast/common"
	"io"
	"os"
	"text/template"
)

/* creating the command type makes it easy to
 * pass around the methods we actually want
 * to make available to other parts of our program,
 * the command type is simply the thing we're passing */

type command struct {
	directoryTree     branch
	templates         map[string]*template.Template
	targetName        string
	targetProjectType string
}

// this type shouldn't even exist

type stringBuf struct {
	asString string
	asBytes  []byte
}

func (self *stringBuf) Write(p []byte) (n int, err error) {
	//common.Log("string buffer just got %v", string(p))
	self.asString = string(p)
	self.asBytes = p
	return len(p), nil
}

/* branch type is used to represent
 * the target directory layout */

type branch struct {
	name        string
	files       []string // just the file names
	directories []branch
}

func (self *branch) print(prefix string) {
	prefix = fmt.Sprintf("%v%v->", prefix, self.name)
	for _, f := range self.files {
		common.Log("%v%v", prefix, f)
	}
	for _, d := range self.directories {
		d.print(prefix)
	}
}

func (self *branch) create(dirname string, templates map[string]*template.Template) {
	dirname = fmt.Sprintf("%v%v/", dirname, self.name)
	common.Log("creating directory `%v`", dirname)
	os.Mkdir(dirname, 0755)
	for _, filename := range self.files {
		// get absolute-ish filename
		absoluteFilename := fmt.Sprintf("%v%v", dirname, filename)
		common.Log("creating file `%v`", absoluteFilename)

		// open file for writing to
		file, err := os.OpenFile(absoluteFilename, os.O_WRONLY|os.O_CREATE, 0755)
		if err != nil {
			common.Fatal("create: OpenFile(): %v", err)
		}
		defer file.Close()

		// pick our template using our non-absolute filename
		// execute the template into the file we opened
		err = templates[filename].Execute(file, self)
		if err != nil {
			common.Fatal("create: template.Execute(): %v", err)
		}
	}
	// recursion is fucking sexy btw
	for _, dir := range self.directories {
		dir.create(dirname, templates)
	}
}

// this gets called from the outside world
func (self *command) Exec(args []string) int {
	switch len(args) {
	case 0:
		return self.Help()
		break
	case 1:
		self.targetProjectType = "generic"
		self.targetName = args[0]
		return self.createNewProject()
		break
	}
	//default case for now
	return self.Help()
}

func (self *command) createNewProject() int {
	common.Log("creating new %v project with name %v", self.targetProjectType, self.targetName)

	// parse the template package
	zipReader := getZipReader(self.targetProjectType)
	defer zipReader.Close()

	self.directoryTree = branch{name: "NO_LAYOUT"}
	for _, zippedFile := range zipReader.File {
		// read fileContents from the zipped file
		fileContents := readZippedFile(zippedFile)

		// if it's the layout.json, parse that
		if zippedFile.Name == "layout.json" {
			common.Log("parsing layout.json")
			self.directoryTree = jsonToBranch(readJson(fileContents))
			continue
		}

		// otherwise, create a new template
		common.Log("creating template with name `%v`", zippedFile.Name)

		tempTemplate := template.New(zippedFile.Name)
		self.templates[zippedFile.Name] = template.Must(tempTemplate.Parse(fileContests.asString))
	}

	// if we didn't find a layout.json, banic!
	if self.directoryTree.name == "NO_LAYOUT" {
		common.Fatal("create: failed to find `layout.json`")
	}

	// create the junk
	self.directoryTree.create("", self.templates)

	return 0
}

func getZipReader(templateName string) *zip.ReadCloser {
	reader, err := zip.OpenReader(fmt.Sprintf("templates/%v.zip", templateName))
	if err != nil {
		common.Fatal("create: getZipReader(): %v", err)
	}
	return reader
}

func readZippedFile(zippedFile *zip.ReadCloser) *stringBuf {
	file, err := zippedFile.Open()
	if err != nil {
		common.Fatal("create: readZippedFile(): %v", err)
	}
	defer file.Close()

	buffer := &stringBuf{}
	size := int64(zippedFile.UncompressedSize64)

	bytes, err := io.CopyN(buffer, file, size)
	common.Log("read %v bytes from layout.json", bytes)
	if err != nil {
		common.Fatal("create: file.Read(): %v", err)
	}

	return buffer
}

func readJson(jsonData []byte, size int64) map[string]interface{} {
	var jsonObject interface{}
	err = json.Unmarshal(jsonData, &jsonObject)
	if err != nil {
		common.Fatal("create: json.Unmarshal(): %v", err)
	}
	return jsonObject.(map[string]interface{})
}

func jsonToBranch(name string, tree map[string]interface{}) branch {
	result := branch{name: name}
	for k, v := range tree {
		if k == "leaf" {
			result.files = append(result.files, v.(string))
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
	cmd := command{templates: make(map[string]*template.Template)}
	return &cmd
}
