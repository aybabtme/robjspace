package rubyobj

type flagSchema struct {
	WbProtected bool `json:"wb_protected,omitempty"`
	Old         bool `json:"old,omitempty"`
	Marked      bool `json:"marked,omitempty"`
}

func (f *flagSchema) clear() {
	f.WbProtected = false
	f.Old = false
	f.Marked = false
}

type objectSchema struct {
	Address    string     `json:"address,omitempty"`
	Class      string     `json:"class,omitempty"`
	NodeType   string     `json:"node_type,omitempty"`
	References []string   `json:"references,omitempty"`
	Type       string     `json:"type,omitempty"`
	Value      string     `json:"value,omitempty"`
	Line       uint64     `json:"line,omitempty"`
	Method     string     `json:"method,omitempty"`
	File       string     `json:"file,omitempty"`
	Fd         int        `json:"fd,omitempty"`
	Bytesize   uint64     `json:"bytesize,omitempty"`
	Capacity   uint64     `json:"capacity,omitempty"`
	Length     uint64     `json:"length,omitempty"`
	Size       uint64     `json:"size,omitempty"`
	Encoding   string     `json:"encoding,omitempty"`
	Default    string     `json:"default,omitempty"`
	Name       string     `json:"name,omitempty"`
	Struct     string     `json:"struct,omitempty"`
	Ivars      uint64     `json:"ivars,omitempty"`
	Generation uint64     `json:"generation,omitempty"`
	Memsize    uint64     `json:"memsize,omitempty"`
	Frozen     bool       `json:"frozen,omitempty"`
	Embedded   bool       `json:"embedded,omitempty"`
	Broken     bool       `json:"broken,omitempty"`
	Fstring    bool       `json:"fstring,omitempty"`
	Shared     bool       `json:"shared,omitempty"`
	Flags      flagSchema `json:"flags,omitempty"`
}

func (o *objectSchema) clear() {
	o.Address = ""
	o.Class = ""
	o.NodeType = ""
	o.References = nil
	o.Type = ""
	o.Value = ""
	o.Line = 0
	o.Method = ""
	o.File = ""
	o.Fd = 0
	o.Bytesize = 0
	o.Capacity = 0
	o.Length = 0
	o.Size = 0
	o.Encoding = ""
	o.Default = ""
	o.Name = ""
	o.Struct = ""
	o.Ivars = 0
	o.Generation = 0
	o.Memsize = 0
	o.Frozen = false
	o.Embedded = false
	o.Broken = false
	o.Fstring = false
	o.Shared = false
	o.Flags.clear()
}
