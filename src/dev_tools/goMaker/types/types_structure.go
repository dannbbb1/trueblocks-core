package types

import (
	"encoding/json"
	"path/filepath"
	"sort"
	"strings"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
)

type Structure struct {
	Class       string    `json:"class,omitempty" toml:"class" csv:"class"`
	DocGroup    string    `json:"doc_group,omitempty" toml:"doc_group" csv:"doc_group"`
	DocRoute    string    `json:"doc_route,omitempty" toml:"doc_route"`
	DocDescr    string    `json:"doc_descr,omitempty" toml:"doc_descr" csv:"doc_descr"`
	DocNotes    string    `json:"doc_notes,omitempty" toml:"doc_notes" csv:"doc_notes"`
	UiRoute     string    `json:"ui_route,omitempty" toml:"ui_route"`
	ProducedBy  string    `json:"produced_by,omitempty" toml:"produced_by"`
	ContainedBy string    `json:"contained_by,omitempty" toml:"contained_by"`
	CacheAs     string    `json:"cache_as,omitempty" toml:"cache_as"`
	CacheBy     string    `json:"cache_by,omitempty" toml:"cache_by"`
	CacheType   string    `json:"cache_type,omitempty" toml:"cache_type"`
	DisableGo   bool      `json:"disable_go,omitempty" toml:"disable_go"`
	DisableDocs bool      `json:"disable_docs,omitempty" toml:"disable_docs"`
	Attributes  string    `json:"attributes,omitempty" toml:"attributes"`
	Members     []Member  `json:"members,omitempty" toml:"members"`
	Route       string    `json:"-" toml:"-"`
	Producers   []string  `json:"-" toml:"-"`
	cbPtr       *CodeBase `json:"-" toml:"-"`
}

func (s *Structure) executeTemplate(name, tmplCode string) string {
	return executeTemplate(s, "structure", name, tmplCode)
}

func (s *Structure) String() string {
	bytes, _ := json.MarshalIndent(s, "", "  ")
	return string(bytes)
}

func (s Structure) Validate() bool {
	return true
}

func (s *Structure) Name() string {
	return Lower(s.Class)
}

func (s *Structure) Num() int {
	parts := strings.Split(s.DocGroup, "-")
	if len(parts) > 1 {
		return int(base.MustParseInt64(parts[0]))
	}
	logger.Fatal("unknown group: " + s.DocGroup)
	return 0
}

func (s *Structure) IsCachable() bool {
	return s.CacheType != ""
}

func (s *Structure) IsFilenameCache() bool {
	return s.CacheBy == "filename"
}

func (s *Structure) IsMarshalOnly() bool {
	return s.CacheType == "marshal_only"
}

func (s *Structure) ClassOrClassGroup() string {
	if s.IsCacheAsGroup() {
		return s.Class + "Group"
	}
	return s.Class
}

func (s *Structure) IsCacheAsGroup() bool {
	return s.CacheAs == "group"
}

func (s *Structure) HasNotes() bool {
	notePath := filepath.Join(GetTemplatePath(), "model-intros/", CamelCase(s.Class)+".notes.md")
	return file.FileExists(notePath)
}

func (s *Structure) HasTimestamp() bool {
	for _, m := range s.Members {
		if m.Name == "timestamp" {
			return true
		}
	}
	return false
}

func (s *Structure) HasSorts() bool {
	for _, m := range s.Members {
		if strings.Contains(m.Attributes, "sorts") {
			return true
		}
	}
	return strings.Contains(s.Attributes, "sorts=")
}

func (s *Structure) NeedsAddress() bool {
	return strings.Contains(s.CacheBy, "address")
}

func (s *Structure) GroupName() string {
	parts := strings.Split(s.DocGroup, "-")
	if len(parts) > 1 {
		return LowerNoSpaces(parts[1])
	}
	logger.Fatal("unknown group: " + s.DocGroup)
	return ""
}

func (s *Structure) ModelIntro() string {
	tmplName := "modelIntro" + s.Class
	tmpl := strings.Trim(getTemplateContents(filepath.Join("model-intros", CamelCase(s.Class))), ws)
	return s.executeTemplate(tmplName, tmpl)
}

func (s *Structure) ModelProducers() string {
	ret := []string{}
	for _, route := range s.Producers { // producers are stored as the name of the route that produces it
		tmplName := "modelProducers"
		tmpl := `- [chifra {{.Route}}](/chifra/{{.Group}}/#chifra-{{.Route}})`
		c := Command{
			Route: route,
			Group: s.cbPtr.RouteToGroup(route),
		}
		ret = append(ret, c.executeTemplate(tmplName, tmpl))
	}
	return strings.Join(ret, "\n")
}

func (s *Structure) ModelMembers() string {
	sort.Slice(s.Members, func(i, j int) bool {
		return s.Members[i].DocOrder < s.Members[j].DocOrder
	})
	header := []string{"Field", "Description", "Type"}
	rows := [][]string{}
	for _, m := range s.Members {
		if m.DocOrder > 0 {
			rows = append(rows, []string{m.Name, m.MarkdownDescription(), m.MarkdownType()})
		}
	}
	return MarkdownTable(header, rows)
}

func (s *Structure) ModelNotes() string {
	tmplName := "Notes" + s.Class
	tmpl := strings.Trim(getTemplateContents(filepath.Join("model-intros", CamelCase(s.Class)+".notes")), ws)
	return strings.Trim(s.executeTemplate(tmplName, tmpl), ws)
}

func (s *Structure) CacheLoc() string {
	if s.Class == "LightBlock" {
		return "Block"
	}
	return s.Class
}

func (s *Structure) CacheIdStr() string {
	switch s.CacheBy {
	case "address,block":
		return "\"%s-%09d\", s.Address.Hex()[2:], s.BlockNumber"
	case "address,block,fourbyte":
		return "\"%s-%s-%09d\", s.Address.Hex()[2:], s.Encoding[2:], s.BlockNumber"
	case "address,tx":
		return "\"%s-%09d-%05d\", s.Address.Hex()[2:], s.BlockNumber, s.TransactionIndex"
	case "block":
		return "\"%09d\", s.BlockNumber"
	case "tx":
		return "\"%09d-%05d\", s.BlockNumber, s.TransactionIndex"
	case "filename":
		return "\"%0s\", s.Filename"
	default:
		logger.Fatal("Unknown cache by format:", s.CacheBy)
		return ""
	}
}

func (s *Structure) HasTsTypesTypes() bool {
	return len(s.TsTypesTypes()) > 0
}

func (s *Structure) TsTypesTypes() string {
	theMap := map[string]bool{}
	for _, m := range s.Members {
		t := m.BaseType()
		common := t == "string" || t == "any" || t == "boolean" || t == "number"
		isParam := s.Class == "Parameter" && t == "Parameter"
		if !common && !isParam {
			theMap[t] = true
		}
	}
	ret := []string{}
	for key := range theMap {
		ret = append(ret, key)
	}
	sort.Slice(ret, func(i, j int) bool {
		return strings.ToLower(ret[i]) < strings.ToLower(ret[j])
	})
	return strings.Join(ret, ", ")
}

func (s *Structure) TsTypeMembers() string {
	ret := []string{}
	for _, m := range s.Members {
		ret = append(ret, m.TsType())
	}
	return strings.Join(ret, "\n")
}

func (s *Structure) parseType(which string) (string, string) {
	which += "="
	attrs := s.Attributes
	if !strings.Contains(attrs, which) {
		return "", ""
	}
	attrs = strings.ReplaceAll(attrs, which, "`")
	parts := strings.Split(attrs, "`")
	parts = strings.Split(parts[1], "|")
	parts = strings.Split(parts[0], ".")

	first, second := parts[0], ""
	if len(parts) > 1 {
		second = parts[1]
	}

	return first, second
}

func (s *Structure) ItemName() string {
	_, typ := s.parseType("itemType")
	return typ
}

func (s *Structure) ItemType() string {
	pkg, typ := s.parseType("itemType")
	return FirstLower(pkg) + "." + FirstUpper(typ)
}

func (s *Structure) InputType() string {
	if s.Class == "Manifest" {
		return "coreTypes.Manifest"
	}
	return s.ItemType()
}

func (s *Structure) EmbedName() string {
	_, typ := s.parseType("embedType")
	return Lower(typ)
}

func (s *Structure) EmbedType() string {
	pkg, typ := s.parseType("embedType")
	if pkg == "types" {
		return FirstUpper(typ)
	}
	return FirstLower(pkg) + "." + FirstUpper(typ)
}

func (s *Structure) OtherName() string {
	_, typ := s.parseType("otherType")
	return Lower(typ)
}

func (s *Structure) OtherType() string {
	pkg, typ := s.parseType("otherType")
	return FirstLower(pkg) + "." + FirstUpper(typ)
}

func (s *Structure) HasItems() bool {
	return strings.Contains(s.Attributes, "itemType=")
}

func (s *Structure) HasEmbed() bool {
	return strings.Contains(s.Attributes, "embedType=")
}

func (s *Structure) HasOther() bool {
	return strings.Contains(s.Attributes, "otherType=")
}

func (s *Structure) NeedsChain() bool {
	return !strings.Contains(s.Attributes, "noChain")
}

func (s *Structure) NeedsFetch() bool {
	return !strings.Contains(s.Attributes, "noFetch")
}

func (s *Structure) IsEditable() bool {
	return strings.Contains(s.Attributes, "editable")
}

func (s *Structure) SortsInstance() string {
	if !s.HasSorts() {
		return ""
	}

	ret := "sdk.SortSpec {\n"
	spec, _ := s.parseType("sorts")
	fields := strings.Split(spec, ",")
	// logger.Info(fields)
	// logger.Info()
	flds := []string{}
	orders := []string{}
	for _, field := range fields {
		// logger.Info(field)
		// logger.Info()
		parts := strings.Split(field, "+")
		if len(parts) > 1 {
			flds = append(flds, "\""+parts[0]+"\"")
			orders = append(orders, "sdk."+parts[1])
		}
	}
	ret += "\tFields: []string{" + strings.Join(flds, ",") + "},\n"
	ret += "\tOrder: []sdk.SortOrder{" + strings.Join(orders, ",") + "},\n"
	ret += "},"
	return ret
}

func (s *Structure) UiRouteNum() uint64 {
	parts := strings.Split(s.UiRoute, "-")
	return base.MustParseUint64(parts[0])
}

func (s *Structure) UiRouteName() string {
	parts := strings.Split(s.UiRoute, "-")
	return parts[1]
}

func (s *Structure) UiRouteStr() string {
	ret := s.UiRouteName()
	if ret == "project" {
		return ""
	}
	return ret
}

func (s *Structure) IsHistory() bool {
	return s.UiRouteName() == "history"
}

func (s *Structure) IsProject() bool {
	return s.UiRouteName() == "project"
}

func (s *Structure) IsWizard() bool {
	return s.UiRouteName() == "wizard"
}
