package undra

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/tusklang/goat"
)

func staticsend(res http.ResponseWriter, req *http.Request) {

	//otherwise, just render the file with no handling (static)
	htmfile := path.Join("public", req.URL.Path)

	//render static (un-templated) html
	f, _ := ioutil.ReadFile(htmfile)
	res.Header().Set("Content-Type", "text/html")
	fmt.Fprint(res, string(f))
	res.Header().Set("Content-Type", "text/plain")
	///////////////////////////////////
}

func getfmt(fpath string) string {

	file, e := os.Open(path.Join("./public", fpath))
	if e != nil {
		return ".oat"
	}
	read, e := ioutil.ReadAll(file)
	if e != nil {
		return ".oat"
	}

	if strings.HasPrefix(string(read), "<!--fmt:tusk-->") {
		return ".tusk"
	} else if strings.HasPrefix(string(read), "<!--fmt:klr-->") {
		return ".klr"
	}

	return ".oat"
}

func handle(res http.ResponseWriter, req *http.Request) {

	if req.URL.Path == "/" {
		req.URL.Path = "/index.html"
	}

	//if the request is not for html
	if filepath.Ext(req.URL.Path) != ".html" {
		fmt.Println(req.URL.Path)
		http.ServeFile(res, req, path.Join("public", req.URL.Path))
		return
	}

	//remove the extension, and replace it with .oat (or .tusk or .klr)
	oatname := strings.TrimSuffix(req.URL.Path, filepath.Ext(req.URL.Path)) + getfmt(req.URL.Path)

	//prepend the server path
	oatf := path.Join("server", oatname)

	if _, f := os.Stat(oatf); !os.IsNotExist(f) {

		var tmp = params
		tmp.Name = oatname

		//load the oat (or tusk or kayl) file using goat
		lib, e := goat.LoadLibrary(oatf, tmp)

		if e != nil {
			fmt.Println(e)
			os.Exit(1)
		}

		var kareq = createRequest(*req)
		var kares = createResponse(res, req)

		_, err := goat.CallOatFunc(lib, "handle", kareq, kares)

		if err != nil {
			fmt.Println(e)
			os.Exit(1)
		}
	} else {
		staticsend(res, req)
	}
}
