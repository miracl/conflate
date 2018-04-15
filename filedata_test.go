package conflate

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func testFiledataNew(t *testing.T, data []byte, path string) (filedata, error) {
	inc, err := newIncludeFromPath(emptyURL, path)
	assert.Nil(t, err)
	return newFiledata(data, inc)
}

func testFiledataNewAssert(t *testing.T, data []byte, path string) filedata {
	fd, err := testFiledataNew(t, data, path)
	assert.Nil(t, err)
	return fd
}

func TestFiledata_WrapErrorNil(t *testing.T) {
	fd, err := testFiledataNew(t, testMarshalJSON, "myurl")
	assert.Nil(t, err)
	err = fd.wrapError(nil)
	assert.Nil(t, err)
}

func TestFiledata_WrapError(t *testing.T) {
	fd := filedata{}
	url, err := toURL(emptyURL, "myurl")
	assert.Nil(t, err)
	fd.include.URL = url
	err = errors.New("My Error")
	err = fd.wrapError(err)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "My Error")
	assert.Regexp(t, "Error processing.*myurl : My Error", err.Error())
}

func TestFiledata_WrapErrorBlank(t *testing.T) {
	fd := filedata{}
	err1 := errors.New("My Error")
	err2 := fd.wrapError(err1)
	assert.NotNil(t, err2)
	assert.Equal(t, err1, err2)
}

func TestFiledata_JSONAsAny(t *testing.T) {
	fd, err := testFiledataNew(t, testMarshalJSON, "file")
	assert.Nil(t, err)
	assert.Equal(t, testMarshalData, fd.obj)
}

func TestFiledata_JSONAsUnknown(t *testing.T) {
	fd, err := testFiledataNew(t, testMarshalJSON, "file.unknown")
	assert.Nil(t, err)
	assert.Equal(t, testMarshalData, fd.obj)
}

func TestFiledata_JSONAsJSON(t *testing.T) {
	fd, err := testFiledataNew(t, testMarshalJSON, "file.json")
	assert.Nil(t, err)
	assert.Equal(t, testMarshalData, fd.obj)
}

func TestFiledata_JSONAsJSN(t *testing.T) {
	fd, err := testFiledataNew(t, testMarshalJSON, "file.jsn")
	assert.Nil(t, err)
	assert.Equal(t, testMarshalData, fd.obj)
}

func TestFiledata_JSONAsTOML(t *testing.T) {
	fd, err := testFiledataNew(t, testMarshalJSON, "file.toml")
	assert.NotNil(t, err)
	assert.Nil(t, fd.obj)
}

func TestFiledata_YAMLAsAny(t *testing.T) {
	fd, err := testFiledataNew(t, testMarshalYAML, "file")
	assert.Nil(t, err)
	assert.Equal(t, testMarshalData, fd.obj)
}

func TestFiledata_YAMLAsUnknown(t *testing.T) {
	fd, err := testFiledataNew(t, testMarshalYAML, "file.unknown")
	assert.Nil(t, err)
	assert.Equal(t, testMarshalData, fd.obj)
}

func TestFiledata_YAMLAsYAML(t *testing.T) {
	fd, err := testFiledataNew(t, testMarshalYAML, "file.yaml")
	assert.Nil(t, err)
	assert.Equal(t, testMarshalData, fd.obj)
}

func TestFiledata_YAMLAsYML(t *testing.T) {
	fd, err := testFiledataNew(t, testMarshalYAML, "file.yml")
	assert.Nil(t, err)
	assert.Equal(t, testMarshalData, fd.obj)
}

func TestFiledata_YAMLAsTOML(t *testing.T) {
	fd, err := testFiledataNew(t, testMarshalYAML, "file.toml")
	assert.NotNil(t, err)
	assert.Nil(t, fd.obj)
}

func TestFiledata_TOMLAsAny(t *testing.T) {
	fd, err := testFiledataNew(t, testMarshalTOML, "file")
	assert.Nil(t, err)
	assert.Equal(t, testMarshalData, fd.obj)
}

func TestFiledata_TOMLAsUnknown(t *testing.T) {
	fd, err := testFiledataNew(t, testMarshalTOML, "file.unknown")
	assert.Nil(t, err)
	assert.Equal(t, testMarshalData, fd.obj)
}

func TestFiledata_TOMLAsTOML(t *testing.T) {
	fd, err := testFiledataNew(t, testMarshalTOML, "file.toml")
	assert.Nil(t, err)
	assert.Equal(t, testMarshalData, fd.obj)
}

func TestFiledata_TOMLAsTML(t *testing.T) {
	fd, err := testFiledataNew(t, testMarshalTOML, "file.tml")
	assert.Nil(t, err)
	assert.Equal(t, testMarshalData, fd.obj)
}

func TestFiledata_TOMLAsJSON(t *testing.T) {
	fd, err := testFiledataNew(t, testMarshalTOML, "file.json")
	assert.NotNil(t, err)
	assert.Nil(t, fd.obj)
}

func TestFiledata_NoIncludes(t *testing.T) {
	fd, err := testLoader.wrapFiledata([]byte(`{"x": 1}`))
	assert.Nil(t, err)
	assert.Nil(t, fd.obj[Includes])
	assert.Equal(t, map[string]interface{}{"x": 1.0}, fd.obj)
}

func TestFiledata_BlankIncludes(t *testing.T) {
	fd, err := testLoader.wrapFiledata([]byte(`{"includes":[], "x": 1}`))
	assert.Nil(t, err)
	assert.Nil(t, fd.obj[Includes])
	assert.Equal(t, map[string]interface{}{"x": 1.0}, fd.obj)
}

func TestFiledata_NullIncludes(t *testing.T) {
	fd, err := testLoader.wrapFiledata([]byte(`{"includes":null, "x": 1}`))
	assert.Nil(t, err)
	assert.Nil(t, fd.obj[Includes])
	assert.Equal(t, map[string]interface{}{"x": 1.0}, fd.obj)
}

func TestFiledata_Includes(t *testing.T) {
	fd, err := testLoader.wrapFiledata([]byte(`{"includes":["test1", {"path": "test2"}], "x": 1}`))
	assert.Nil(t, err)
	assert.Len(t, fd.includes, 2)
	assert.Regexp(t, "^file://.*/test1", fd.includes[0].URL.String())
	assert.Regexp(t, "test2", fd.includes[1].Path)
	assert.Nil(t, fd.obj[Includes])
	assert.Equal(t, map[string]interface{}{"x": 1.0}, fd.obj)
}

func TestFiledata_IncludesError(t *testing.T) {
	_, err := testLoader.wrapFiledata([]byte(`{"includes": "not array"}`))
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Could not extract includes")
}

func TestFiledata_Expand(t *testing.T) {
	w := os.Getenv("W")
	x := os.Getenv("X")
	y := os.Getenv("Y")
	z := os.Getenv("Z")
	os.Setenv("W", "$W")
	os.Setenv("X", `"x"`)
	os.Setenv("Y", `y`)
	os.Setenv("Z", `$Y`)
	defer func() {
		os.Setenv("W", w)
		os.Setenv("X", x)
		os.Setenv("Y", y)
		os.Setenv("Z", z)
	}()
	b := recursiveExpand([]byte(`{"W":"$W","X":$X,"Y":"$Y","Z":"$Z"}`))
	assert.Equal(t, string(`{"W":"$W","X":"x","Y":"y","Z":"y"}`), string(b))
}

func TestFiledatas_Unmarshal(t *testing.T) {
	fds := filedatas{
		testFiledataNewAssert(t, testMarshalJSON, "file.json"),
		testFiledataNewAssert(t, testMarshalYAML, "file.yaml"),
		testFiledataNewAssert(t, testMarshalTOML, "file.toml"),
	}
	assert.Equal(t, []interface{}{testMarshalData, testMarshalData, testMarshalData}, fds.objs())
}

func TestFiledatas_DifferentIncludes(t *testing.T) {
	old := Includes
	Includes = "using"
	defer func() { Includes = old }()
	fd, err := testLoader.wrapFiledata([]byte(`{"using":["test1", "test2"], "x": 1}`))
	assert.Nil(t, err)
	assert.Len(t, fd.includes, 2)
	assert.Regexp(t, "^file://.*/test1", fd.includes[0].URL.String())
	assert.Regexp(t, "^file://.*/test2", fd.includes[1].URL.String())
	assert.Nil(t, fd.obj[Includes])
	assert.Equal(t, map[string]interface{}{"x": 1.0}, fd.obj)
}

func TestFiledatas_NoIncludes(t *testing.T) {
	old := Includes
	Includes = "using"
	defer func() { Includes = old }()
	fd, err := testLoader.wrapFiledata([]byte(`{"includes":["test1", "test2"]}`))
	assert.Nil(t, err)
	assert.Empty(t, fd.includes)
	assert.Len(t, fd.obj["includes"], 2)
}

func TestFiledatas_IgnoreIncludes(t *testing.T) {
	old := Includes
	Includes = ""
	defer func() { Includes = old }()
	fd, err := testLoader.wrapFiledata([]byte(`{"":["test1", "test2"]}`))
	assert.Nil(t, err)
	assert.Empty(t, fd.includes)
	assert.Equal(t, map[string]interface{}{"": []interface{}{"test1", "test2"}}, fd.obj)
}
