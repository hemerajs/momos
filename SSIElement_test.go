package momos

import (
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/nats-io/nuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateElement(t *testing.T) {

	html := `
	<ssi
	name="basket"
	timeout="2000"
	template="true"
	src="http://localhost:8081">
		Default content!
		
		<ssi-timeout>
		<span>Please try it again! {{.DateLocal}}</span>
		</ssi-timeout>
		
		<ssi-error>
		<span>Please call the support!</span>
		</ssi-error>
	</ssi>
	`

	r := strings.NewReader(html)
	doc, _ := goquery.NewDocumentFromReader(r)
	element := doc.Find("ssi")
	se := SSIElement{Element: element}

	se.SetTimeout(element.AttrOr("timeout", "2000"))
	se.SetSrc(element.AttrOr("src", ""))
	se.SetName(element.AttrOr("name", nuid.Next()))
	se.SetTemplate(element.AttrOr("template", "false"))
	se.SetFilterIncludes(element.AttrOr("no-scripts", "true"))

	assert.Equal(t, se.timeout.String(), "2s", "should be 2s")
	assert.Equal(t, se.name, "basket", "should be the name of the name attribute")
	assert.Equal(t, se.hasTemplate, true, "should be the timeout of the timeout attribute")
	assert.Equal(t, se.src, "http://localhost:8081", "should be the url from url attribute")
}

func TestReplaceWithDefaultHTML(t *testing.T) {

	html := `
	<ssi
	name="basket"
	timeout="2000"
	template="true"
	src="http://localhost:8081">
		Default content!
		
		<ssi-timeout>
		<span>Please try it again! {{.DateLocal}}</span>
		</ssi-timeout>
		
		<ssi-error>
		<span>Please call the support!</span>
		</ssi-error>
	</ssi>
	`

	r := strings.NewReader(html)
	doc, _ := goquery.NewDocumentFromReader(r)
	element := doc.Find("ssi")
	se := SSIElement{Element: element}

	se.SetTimeout(element.AttrOr("timeout", "2000"))
	se.SetSrc(element.AttrOr("src", ""))
	se.SetName(element.AttrOr("name", nuid.Next()))
	se.SetTemplate(element.AttrOr("template", "false"))
	se.SetFilterIncludes(element.AttrOr("no-scripts", "true"))

	assert.Equal(t, se.timeout.String(), "2s", "should be 2s")
	assert.Equal(t, se.name, "basket", "should be the name of the name attribute")
	assert.Equal(t, se.hasTemplate, true, "should be the timeout of the timeout attribute")
	assert.Equal(t, se.src, "http://localhost:8081", "should be the url from url attribute")

	se.replaceWithDefaultHTML()

	html, _ = doc.Html()
	assert.Equal(t, html, "<html><head></head><body>\n\t\tDefault content!\n\t\t\n\t\t\n\t\t\n\t\t\n\t\n\t</body></html>", "should contain the default html")
}
func TestReplaceWithErrorHTML(t *testing.T) {

	html := `
		<ssi
		name="basket"
		timeout="2000"
		template="true"
		src="http://localhost:8081">
			Default content!
			
			<ssi-timeout>
			<span>Please try it again! {{.DateLocal}}</span>
			</ssi-timeout>
			
			<ssi-error>
			<span>Please call the support!</span>
			</ssi-error>
		</ssi>
		`

	r := strings.NewReader(html)
	doc, _ := goquery.NewDocumentFromReader(r)
	element := doc.Find("ssi")
	se := SSIElement{Element: element}

	se.SetTimeout(element.AttrOr("timeout", "2000"))
	se.SetSrc(element.AttrOr("src", ""))
	se.SetName(element.AttrOr("name", nuid.Next()))
	se.SetTemplate(element.AttrOr("template", "false"))
	se.SetFilterIncludes(element.AttrOr("no-scripts", "true"))

	assert.Equal(t, se.timeout.String(), "2s", "should be 2s")
	assert.Equal(t, se.name, "basket", "should be the name of the name attribute")
	assert.Equal(t, se.hasTemplate, true, "should be the timeout of the timeout attribute")
	assert.Equal(t, se.src, "http://localhost:8081", "should be the url from url attribute")

	se.GetErrorTag()

	assert.Equal(t, se.HasErrorTag, true, "should be true")

	se.replaceWithErrorHTML()

	html, _ = doc.Html()
	assert.Equal(t, html, "<html><head></head><body>\n\t\t\t<span>Please call the support!</span>\n\t\t\t\n\t\t</body></html>", "should contain the default html")
}

func TestReplaceWithTimeoutHTML(t *testing.T) {

	html := `
	<ssi
	name="basket"
	timeout="2000"
	template="true"
	src="http://localhost:8081">
		Default content!
		
		<ssi-timeout>
		<span>Please try it again! {{.DateLocal}}</span>
		</ssi-timeout>
		
		<ssi-error>
		<span>Please call the support!</span>
		</ssi-error>
	</ssi>
	`

	r := strings.NewReader(html)
	doc, _ := goquery.NewDocumentFromReader(r)
	element := doc.Find("ssi")
	se := SSIElement{Element: element}

	se.SetTimeout(element.AttrOr("timeout", "2000"))
	se.SetSrc(element.AttrOr("src", ""))
	se.SetName(element.AttrOr("name", nuid.Next()))
	se.SetTemplate(element.AttrOr("template", "false"))
	se.SetFilterIncludes(element.AttrOr("no-scripts", "true"))

	assert.Equal(t, se.timeout.String(), "2s", "should be 2s")
	assert.Equal(t, se.name, "basket", "should be the name of the name attribute")
	assert.Equal(t, se.hasTemplate, true, "should be the timeout of the timeout attribute")
	assert.Equal(t, se.src, "http://localhost:8081", "should be the url from url attribute")

	se.GetTimeoutTag()

	assert.Equal(t, se.HasTimeoutTag, true, "should be true")

	se.replaceWithTimeoutHTML()

	html, _ = doc.Html()
	assert.Equal(t, html, "<html><head></head><body>\n\t\t<span>Please try it again! </span>\n\t\t\n\t</body></html>", "should contain the default html")
}
