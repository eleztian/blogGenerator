package generator

import (
	"bytes"
	"fmt"
	"html/template"
	"path/filepath"
	"strings"
)

// ListingData holds the data for the listing page
type ListingData struct {
	Title      string
	Date       string
	Short      string
	Link       string
	TimeToRead string
	Tags       []*Tag
}

// ListingGenerator Object
type ListingGenerator struct {
	Config *ListingConfig
}

// ListingConfig holds the configuration for the listing page
type ListingConfig struct {
	NPG                    int
	Posts                  []*Post
	Template               *template.Template
	Destination, PageTitle string
	IsIndex                bool
	Writer                 *IndexWriter
}

// Generate starts the listing generation
func (g *ListingGenerator) Generate() error {
	fmt.Println("\tGenerating List ", g.Config.PageTitle, "...")
	defer fmt.Println("\tFinished List ", g.Config.PageTitle, "...")
	shortTemplatePath := filepath.Join("static", "short.html")
	archiveLinkTemplatePath := filepath.Join("static", "archiveLink.html")
	npg := g.Config.NPG
	posts := g.Config.Posts
	t := g.Config.Template
	destination := g.Config.Destination
	pageTitle := g.Config.PageTitle
	short, err := getTemplate(shortTemplatePath)
	if err != nil {
		return err
	}
	var postBlocks []string
	for _, post := range posts {
		meta := post.Meta
		link := fmt.Sprintf("/%s/", post.Name)
		ld := ListingData{
			Title:      meta.Title,
			Date:       meta.Date,
			Short:      meta.Short,
			Link:       link,
			Tags:       createTags(meta.Tags),
			TimeToRead: calculateTimeToRead(string(post.HTML)),
		}
		block := bytes.Buffer{}
		if err := short.Execute(&block, ld); err != nil {
			return fmt.Errorf("error executing template %s: %v", shortTemplatePath, err)
		}
		postBlocks = append(postBlocks, block.String())
	}

	if g.Config.IsIndex {
		htmlBlocks := template.HTML(strings.Join(postBlocks, "\n"))
		archiveLink, err := getTemplate(archiveLinkTemplatePath)
		if err != nil {
			return err
		}
		lastBlock := bytes.Buffer{}
		if err := archiveLink.Execute(&lastBlock, nil); err != nil {
			return fmt.Errorf("error executing template %s: %v", archiveLinkTemplatePath, err)
		}
		htmlBlocks = template.HTML(fmt.Sprintf("%s%s", htmlBlocks, template.HTML(lastBlock.String())))
		if err := g.Config.Writer.WriteIndexHTML(destination, pageTitle, pageTitle, htmlBlocks, t); err != nil {
			return err
		}
		return nil
	}
	nPage := len(postBlocks) / npg
	if npg*nPage < len(postBlocks) {
		nPage += 1
	}
	fmt.Println("Page num: ", nPage)
	archiveLink, err := getTemplate(archiveLinkTemplatePath)
	for i := 0; i < nPage; i++ {
		s := i * npg
		e := (i + 1) * npg
		if i == nPage-1 {
			e = len(postBlocks)
		}
		data := fmt.Sprintf("page/%d/", i+2)
		htmlBlocks := template.HTML(strings.Join(postBlocks[s:e], "\n"))
		if i != 0 {
			fmt.Println(data)
			destination = filepath.Join(g.Config.Destination, fmt.Sprintf("page/%d/", i+1))
		}
		if i+1 < nPage {
			lastBlock := bytes.Buffer{}
			if err := archiveLink.Execute(&lastBlock, data); err != nil {
				return fmt.Errorf("error executing template %s: %v", archiveLinkTemplatePath, err)
			}
			htmlBlocks = template.HTML(fmt.Sprintf("%s%s", htmlBlocks, template.HTML(lastBlock.String())))
		}

		if err := g.Config.Writer.WriteIndexHTML(destination, pageTitle, pageTitle, htmlBlocks, t); err != nil {
			return err
		}
	}

	return nil
}

func calculateTimeToRead(input string) string {
	// an average human reads about 200 wpm
	var secondsPerWord = 60.0 / 200.0
	// multiply with the amount of words
	words := secondsPerWord * float64(len(strings.Split(input, " ")))
	// add 12 seconds for each image
	images := 12.0 * strings.Count(input, "<img")
	result := (words + float64(images)) / 60.0
	if result < 1.0 {
		result = 1.0
	}
	return fmt.Sprintf("%.0fm", result)
}
