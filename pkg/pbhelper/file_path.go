package pbhelper

import "fmt"

type FilePathOpts struct {
	Thumb *FileThumb
}

func (f *FilePathOpts) build(path string) string {
	path += "?"

	if f.Thumb != nil {
		path += f.Thumb.build()
	}

	return path
}

type FileThumb struct {
	Width  int
	Height int
}

func (f *FileThumb) build() string {
	if f == nil {
		return ""
	}

	return fmt.Sprintf("thumb=%vx%v", f.Width, f.Height)
}

func BuildFilePath(collection, id, filename string, opts ...FilePathOpts) string {
	if filename == "" {
		return ""
	}

	path := "/api/files/" + collection + "/" + id + "/" + filename

	if len(opts) > 0 {
		opt := opts[0]
		path = opt.build(path)
	}

	return path
}
