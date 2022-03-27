package epub

import (
	"encoding/xml"
	"fmt"
	"path/filepath"
	"time"
)

const (
	SchemeMARCRelators  = "marc:relators"
	SchemeONIXCodeList5 = "onix:codelist5"
	SchemeXSDString     = "xsd:string"
)

const (
	// Content uses SchemeMARCRelators,
	// use PropertyRole* constants,
	// see https://www.loc.gov/marc/relators/relaterm.html
	PropertyRole = "role"

	PropertyTitleType         = "title-type"
	PropertyDisplaySequence   = "display-seq"
	PropertyMetadataAuthority = "meta-auth"

	// Content uses SchemeONIXCodeList5 or SchemeXSDString,
	// use PropertyIdentifierType* constants,
	// see https://onix-codelists.io/codelist/5
	PropertyIdentifierType = "identifier-type"
	// Content is a timestamp in UTC, format 2011-01-01T12:00:00Z (formal specification CCYY-MM-DDThh:mm:ssZ)
	PropertyModified = "dcterms:modified"
)

const (
	PropertyRoleAuthor       = "aut"
	PropertyRoleBookProducer = "bkp"
	PropertyRoleTranslator   = "trl"
	PropertyRoleArtist       = "art"
	// ... many more in original list
)

// XSD String
const (
	PropertyIdentifierTypeUUID = "uuid"
)

// ONIX Codelist 5
const (
	PropertyIdentifierTypeProprietary        = "01"
	PropertyIdentifierTypeISBN10             = "02"
	PropertyIdentifierTypeGTIN13             = "03"
	PropertyIdentifierTypeUPC                = "04"
	PropertyIdentifierTypeISMN10             = "05"
	PropertyIdentifierTypeDOI                = "06"
	PropertyIdentifierTypeLCCN               = "13"
	PropertyIdentifierTypeGTIN14             = "14"
	PropertyIdentifierTypeISBN13             = "15"
	PropertyIdentifierTypeLegalDepositNumber = "17"
	PropertyIdentifierTypeURN                = "22"
	PropertyIdentifierTypeOCLC               = "23"
	PropertyIdentifierTypeCoPublisherISB13   = "24"
	PropertyIdentifierTypeISMN13             = "25"
	PropertyIdentifierTypeISBNA              = "26"
	PropertyIdentifierTypeJPecode            = "27"
	PropertyIdentifierTypeOLCC               = "28"
	PropertyIdentifierTypeJPMagazineID       = "29"
	PropertyIdentifierTypeUPC125             = "30"
	PropertyIdentifierTypeBNFControlNumber   = "31"
	PropertyIdentifierTypeARK                = "35"
)

const (
	pkgCreatorID    = "creator"
	pkgPublisherID  = "publisher"
	pkgIdentifierID = "pub-id"

	pkgFileTemplate = `<?xml version="1.0" encoding="UTF-8"?>
<package version="3.0" unique-identifier="pub-id" xmlns="http://www.idpf.org/2007/opf">
  <metadata xmlns:dc="http://purl.org/dc/elements/1.1/">
    <dc:identifier id="pub-id"></dc:identifier>
    <dc:title></dc:title>
    <dc:language></dc:language>
    <dc:description></dc:description>
  </metadata>
  <manifest>
  </manifest>
  <spine toc="ncx">
  </spine>
</package>
`

	xmlnsDc = "http://purl.org/dc/elements/1.1/"
)

// pkg implements the package document file (package.opf), which contains
// metadata about the EPUB (title, author, etc) as well as a list of files the
// EPUB contains.
//
// Sample: https://github.com/bmaupin/epub-samples/blob/master/minimal-v3plus2/EPUB/package.opf
// Spec: http://www.idpf.org/epub/301/spec/epub-publications.html
type Pkg struct {
	xml *PkgRoot
}

// This holds the actual XML for the package file
type PkgRoot struct {
	XMLName          xml.Name    `xml:"http://www.idpf.org/2007/opf package"`
	UniqueIdentifier string      `xml:"unique-identifier,attr"`
	Version          string      `xml:"version,attr"`
	Metadata         PkgMetadata `xml:"metadata"`
	ManifestItems    []PkgItem   `xml:"manifest>item"`
	Spine            PkgSpine    `xml:"spine"`
}

// <dc:creator>, e.g. the author
type PkgCreator struct {
	XMLName xml.Name `xml:"dc:creator"`
	ID      string   `xml:"id,attr"`
	Data    string   `xml:",chardata"`
}

// <dc:contributor>, e.g. the generating program
type PkgContributor struct {
	XMLName xml.Name `xml:"dc:contributor"`
	ID      string   `xml:"id,attr"`
	Data    string   `xml:",chardata"`
}

// <dc:identifier>, where the unique identifier is stored
// Ex: <dc:identifier id="pub-id">urn:uuid:fe93046f-af57-475a-a0cb-a0d4bc99ba6d</dc:identifier>
type PkgIdentifier struct {
	ID   string `xml:"id,attr"`
	Data string `xml:",chardata"`
}

// <item> elements, one per each file stored in the EPUB
// Ex: <item id="nav" href="nav.xhtml" media-type="application/xhtml+xml" properties="nav" />
//     <item id="ncx" href="toc.ncx" media-type="application/x-dtbncx+xml" />
//     <item id="section0001.xhtml" href="xhtml/section0001.xhtml" media-type="application/xhtml+xml" />
type PkgItem struct {
	ID         string `xml:"id,attr"`
	Href       string `xml:"href,attr"`
	MediaType  string `xml:"media-type,attr"`
	Properties string `xml:"properties,attr,omitempty"`
}

// <itemref> elements, which define the reading order
// Ex: <itemref idref="section0001.xhtml" />
type PkgItemref struct {
	Idref string `xml:"idref,attr"`
}

// The <meta> element, which contains modified date, role of the creator (e.g.
// author), etc
// Ex: <meta refines="#creator" property="role" scheme="marc:relators" id="role">aut</meta>
//     <meta property="dcterms:modified">2011-01-01T12:00:00Z</meta>
type PkgMeta struct {
	Refines  string `xml:"refines,attr,omitempty"`
	Property string `xml:"property,attr,omitempty"`
	Scheme   string `xml:"scheme,attr,omitempty"`
	ID       string `xml:"id,attr,omitempty"`
	Data     string `xml:",chardata"`
	Name     string `xml:"name,attr,omitempty"`
	Content  string `xml:"content,attr,omitempty"`
}

// The <metadata> element
type PkgMetadata struct {
	XmlnsDc    string          `xml:"xmlns:dc,attr"`
	Identifier []PkgIdentifier `xml:"dc:identifier"`
	// Ex: <dc:title>Your title here</dc:title>
	Title string `xml:"dc:title"`
	// Ex: <dc:language>en</dc:language>
	Language    string `xml:"dc:language"`
	Description string `xml:"dc:description,omitempty"`
	Publisher   string `xml:"dc:publisher,omitempty"`
	// e.g. a URL
	Source string `xml:"dc:source,omitempty"`
	Date   string `xml:"dc:date,omitempty"`
	// Tags
	Subject     []string `xml:"dc:subject,omitempty"`
	Creator     []PkgCreator
	Contributor []PkgContributor
	Meta        []PkgMeta `xml:"meta"`
}

// The <spine> element
type PkgSpine struct {
	Items []PkgItemref `xml:"itemref"`
	Toc   string       `xml:"toc,attr"`
	Ppd   string       `xml:"page-progression-direction,attr,omitempty"`
}

// Constructor for pkg
func NewPkg() *Pkg {
	p := &Pkg{
		xml: &PkgRoot{
			Metadata: PkgMetadata{
				XmlnsDc: xmlnsDc,
			},
		},
	}

	err := xml.Unmarshal([]byte(pkgFileTemplate), &p.xml)
	if err != nil {
		panic(fmt.Sprintf(
			"Error unmarshalling package file XML: %s\n"+
				"\tp.xml=%#v\n"+
				"\tpkgFileTemplate=%s",
			err,
			*p.xml,
			pkgFileTemplate))
	}

	return p
}

func (p *Pkg) AddToManifest(id string, href string, mediaType string, properties string) {
	href = filepath.ToSlash(href)
	i := &PkgItem{
		ID:         id,
		Href:       href,
		MediaType:  mediaType,
		Properties: properties,
	}
	p.xml.ManifestItems = append(p.xml.ManifestItems, *i)
}

func (p *Pkg) AddToSpine(id string) {
	i := &PkgItemref{
		Idref: id,
	}

	p.xml.Spine.Items = append(p.xml.Spine.Items, *i)
}

func (p *Pkg) AddAuthor(author, role string) {
	id := fmt.Sprintf("%s%d", pkgCreatorID, len(p.xml.Metadata.Creator))

	p.xml.Metadata.Creator = append(p.xml.Metadata.Creator, PkgCreator{
		Data: author,
		ID:   id,
	})
	meta := PkgMeta{
		Refines:  "#" + id,
		ID:       id,
		Property: PropertyRole,
		Data:     role,
		Scheme:   SchemeMARCRelators,
	}

	p.xml.Metadata.Meta = updateMeta(p.xml.Metadata.Meta, meta)
}

func (p *Pkg) AddContributor(author, role string) {
	id := fmt.Sprintf("%s%d", pkgCreatorID, len(p.xml.Metadata.Creator))

	p.xml.Metadata.Creator = append(p.xml.Metadata.Creator, PkgCreator{
		Data: author,
		ID:   id,
	})
	meta := PkgMeta{
		Refines:  "#" + id,
		ID:       id,
		Property: PropertyRole,
		Data:     role,
		Scheme:   SchemeMARCRelators,
	}

	p.xml.Metadata.Meta = updateMeta(p.xml.Metadata.Meta, meta)
}

// Add an EPUB 2 cover meta element for backward compatibility (http://idpf.org/forum/topic-715)
func (p *Pkg) SetCover(coverRef string) {
	meta := PkgMeta{
		Name:    "cover",
		Content: coverRef,
	}
	p.xml.Metadata.Meta = updateMeta(p.xml.Metadata.Meta, meta)
}

func (p *Pkg) AddCustomMeta(name, content string) {
	meta := PkgMeta{
		Name:    name,
		Content: content,
	}
	p.xml.Metadata.Meta = updateMeta(p.xml.Metadata.Meta, meta)
}

// AddIdentifier adds an identifier of the EPUB, such as a UUID, DOI,
// ISBN or ISSN. If no identifier is set, a UUID will be automatically
// generated.
func (p *Pkg) AddIdentifier(identifier, typeSchema, typeContent string) {
	id := fmt.Sprintf("%s%d", pkgCreatorID, len(p.xml.Metadata.Creator))

	p.xml.Metadata.Identifier = append(p.xml.Metadata.Identifier, PkgIdentifier{
		ID:   id,
		Data: identifier,
	})
	meta := PkgMeta{
		Refines:  "#" + id,
		ID:       id,
		Property: PropertyIdentifierType,
		Data:     typeContent,
		Scheme:   typeSchema,
	}
	p.xml.Metadata.Meta = updateMeta(p.xml.Metadata.Meta, meta)
}

func (p *Pkg) SetLang(lang string) {
	p.xml.Metadata.Language = lang
}

func (p *Pkg) SetDescription(desc string) {
	p.xml.Metadata.Description = desc
}

func (p *Pkg) SetPublisher(publisher string) {
	p.xml.Metadata.Publisher = publisher
}

func (p *Pkg) SetSource(source string) {
	p.xml.Metadata.Source = source
}

func (p *Pkg) SetDate(dt time.Time) {
	p.xml.Metadata.Date = dt.Format(time.RFC3339)
}

func (p *Pkg) SetSubject(subject []string) {
	p.xml.Metadata.Subject = subject
}

func (p *Pkg) AddSubject(subject string) {
	p.xml.Metadata.Subject = append(p.xml.Metadata.Subject, subject)
}

func (p *Pkg) SetPpd(direction string) {
	p.xml.Spine.Ppd = direction
}

func (p *Pkg) SetModified(timestamp string) {
	meta := PkgMeta{
		Data:     timestamp,
		Property: PropertyModified,
	}

	p.xml.Metadata.Meta = updateMeta(p.xml.Metadata.Meta, meta)
}

func (p *Pkg) SetTitle(title string) {
	p.xml.Metadata.Title = title
}

// Update the <meta> element
func updateMeta(a []PkgMeta, m PkgMeta) []PkgMeta {
	indexToReplace := -1

	if len(a) > 0 {
		// If we've already added the modified meta element to the meta array
		for i, meta := range a {
			if meta == m {
				indexToReplace = i
				break
			}
		}
	}

	// If the array is empty or the meta element isn't in it
	if indexToReplace == -1 {
		// Add the meta element to the array of meta elements
		a = append(a, m)

		// If the meta element is found
	} else {
		// Replace it
		a[indexToReplace] = m
	}

	return a
}

// Write the package file to the temporary directory
func (p *Pkg) write(tempDir string) {
	now := time.Now().UTC().Format("2006-01-02T15:04:05Z")
	p.SetModified(now)

	pkgFilePath := filepath.Join(tempDir, contentFolderName, pkgFilename)

	output, err := xml.MarshalIndent(p.xml, "", "  ")
	if err != nil {
		panic(fmt.Sprintf(
			"Error marshalling XML for package file: %s\n"+
				"\tXML=%#v",
			err,
			p.xml))
	}
	// Add the xml header to the output
	pkgFileContent := append([]byte(xml.Header), output...)
	// It's generally nice to have files end with a newline
	pkgFileContent = append(pkgFileContent, "\n"...)

	if err := filesystem.WriteFile(pkgFilePath, []byte(pkgFileContent), filePermissions); err != nil {
		panic(fmt.Sprintf("Error writing package file: %s", err))
	}
}
