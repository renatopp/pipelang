package main

import (
	"bytes"
	"html/template"
	"os"
	"pagegen/vv"
	"pagegen/vv/processor"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type TourToc struct {
	Title   string
	Chapter bool
	Link    string
}

var TourPath = regexp.MustCompile(`chapter\d+_(.*)\.section\d+_([^\.]+)`)

func main() {
	tourFile, err := os.ReadFile("../../.page/tour/layout.html")
	if err != nil {
		panic(err)
	}
	tourTemplate, err := template.New("tour").Parse(string(tourFile))
	if err != nil {
		panic(err)
	}
	tourFiles, err := vv.ListFiles("../../.docs/tour")
	if err != nil {
		panic(err)
	}
	tourToc := convertTourToc(tourFiles)

	data := map[string]any{
		"Title":   "PIPE Lang",
		"TourToc": tourToc,
	}

	err = vv.Run(&vv.Ctx{
		RootDir: "../../",
		DistDir: "dist",
		SourceDirs: []string{
			".page/static",
			".page/parts",
			".page/pages",
			".docs/tour",
		},
		Processors: []vv.ProcessorFn{
			// Store tour code
			func(config *vv.Ctx, file *vv.File) error {
				if file.Extension != ".pp" {
					return nil
				}

				dir := strings.ReplaceAll(file.Dir, string(filepath.Separator), ".")
				name := "code_" + dir

				data[name] = string(file.Body)
				file.Ignored = true
				file.Hidden = true
				return nil
			},

			// Store html parts
			processor.NewHtmlPart(data, func(file *vv.File) bool {
				return strings.Contains(file.SourcePath, ".page"+string(filepath.Separator)+"parts")
			}),

			// Process html pages
			processor.NewHtml(data, func(file *vv.File) bool {
				return file.Extension == ".html"
			}),

			// Process tour files
			func(config *vv.Ctx, file *vv.File) error {
				if file.Extension != ".md" || !strings.Contains(file.SourcePath, "tour") {
					return nil
				}

				dir := strings.ReplaceAll(file.Dir, string(filepath.Separator), ".")
				data["tour_content"] = template.HTML(mdToHTML(file.Body))
				data["tour_code"] = data["code_"+dir].(string)
				data["tour_prev"] = getTourPrevPath(file.SourcePath, tourFiles)
				data["tour_cur"] = getTourCurPath(dir)
				data["tour_next"] = getTourNextPath(file.SourcePath, tourFiles)
				var buf bytes.Buffer
				tourTemplate.Execute(&buf, data)

				file.Dir = getTourCurPath(dir)
				file.Name = "index"
				file.Extension = ".html"
				file.Body = buf.Bytes()

				return nil
			},

			// Print all
			func(ctx *vv.Ctx, file *vv.File) error {
				println(file.SourcePath, "=>", file.Dir, file.Name, file.Extension)
				return nil
			},

			// Ignore all
			// func(ctx *vv.Ctx, file *vv.File) error {
			// 	file.Ignored = true
			// 	return nil
			// },
		},
	})
	if err != nil {
		panic(err)
	}
}

func mdToHTML(md []byte) []byte {
	// create markdown parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	// create HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return markdown.Render(doc, renderer)
}

func getTourCurPath(path string) string {
	dir := strings.ReplaceAll(path, string(filepath.Separator), ".")
	matches := TourPath.FindStringSubmatch(dir)
	if len(matches) != 3 {
		return ""
	}
	chapter := matches[1]
	section := matches[2]

	return strings.Join([]string{"/tour", chapter, section}, "/")
}

func getTourPrevPath(path string, paths []string) string {
	idx := slices.Index(paths, path)
	idx -= 2
	if idx < 0 || !strings.Contains(paths[idx], "tour") {
		return ""
	}

	return getTourCurPath(paths[idx])
}

func getTourNextPath(path string, paths []string) string {
	idx := slices.Index(paths, path)
	idx += 2
	if idx >= len(paths) || !strings.Contains(paths[idx], "tour") {
		return ""
	}

	return getTourCurPath(paths[idx])
}

var titleCaser = cases.Title(language.English)

func getTourSectionName(path string) string {
	dir := strings.ReplaceAll(path, string(filepath.Separator), ".")
	matches := TourPath.FindStringSubmatch(dir)
	if len(matches) != 3 {
		return ""
	}
	slug := matches[2]
	slug = strings.ReplaceAll(slug, "-", " ")
	return titleCaser.String(slug)
}

func getTourChapterName(path string) string {
	dir := strings.ReplaceAll(path, string(filepath.Separator), ".")
	matches := TourPath.FindStringSubmatch(dir)
	if len(matches) != 3 {
		return ""
	}
	slug := matches[1]
	slug = strings.ReplaceAll(slug, "-", " ")
	return titleCaser.String(slug)
}

func convertTourToc(paths []string) []TourToc {
	curChapter := getTourChapterName(paths[0])
	toc := []TourToc{
		{Title: curChapter, Chapter: true},
	}

	for _, path := range paths {
		if !strings.Contains(path, ".md") {
			continue
		}

		dir := strings.ReplaceAll(path, string(filepath.Separator), ".")

		chapter := getTourChapterName(path)
		if chapter != curChapter {
			curChapter = chapter
			toc = append(toc, TourToc{Title: curChapter, Chapter: true})
		}

		toc = append(toc, TourToc{
			Title: getTourSectionName(path),
			Link:  getTourCurPath(dir),
		})
	}
	return toc
}
