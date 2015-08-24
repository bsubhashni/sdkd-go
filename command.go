package main

const (
	NEWHANDLE            = "NEWHANDLE"
	MC_DS_MUTATE_SET     = "MC_DS_MUTATE_SET"
	MC_DS_MUTATE_REPLACE = "MC_DS_MUTATE_REPLACE"
	MC_DS_GET            = "MC_DS_GET"
	CB_VIEW_QUERY        = "CB_VIEW_QUERY"
	CB_VIEW_LOAD         = "CB_VIEW_LOAD"
	CB_N1QL_CREATE_INDEX = "CB_N1QL_CREATE_INDEX"
	CB_N1QL_QUERY        = "CB_N1QL_QUERY"
	DSTYPE_SEEDED        = "SEEDED"
	CANCEL               = "CANCEL"
	CLOSEHANDLE          = "CLOSEHANDLE"
)

type RequestCommand struct {
	Command string      `json:"Command"`
	ReqID   int         `json:"ReqID"`
	CmdData CommandData `json:"CommandData"`
	Handle  int         `json:"Handle"`
}

type CommandData struct {
	DSType              string              `json:"DSType,omitempty"`
	DS                  DS                  `json:"DS,omitempty"`
	Options             Options             `json:"Options,omitempty"`
	Bucket              string              `json:"Bucket,omitempty"`
	Port                int                 `json:"Port,omitempty"`
	Hostname            string              `json:"Hostname,omitempty"`
	VSchema             ViewSchema          `json:"Schema,omitempty"`
	ViewQueryParameters ViewQueryParameters `json:"ViewQueryParameters"`
	ViewName            string              `json:"ViewName"`
	DesignName          string              `json:"DesignName"`
}

type Options struct {
	Username           string `json:"Username"`
	ClusterCertificate string `json:"ClusterCertificate"`
	SSL                bool   `json:"SSL"`
	Password           string `json:"Password"`
	ReplicateTo        int    `json:"ReplicateTo"`
	TimeRes            int64  `json:"TimeRes"`
	PersistTo          int    `json:"PersistTo"`
	ReplicaRead        bool   `json:"ReplicaRead"`
	IterWait           uint64 `json:"IterWait"`
	DelayMax           int    `json:"DelayMax"`
	DelayMin           int    `json:"DelayMin"`
	Full               bool   `json:"Full"`
	ViewQueryCount     int    `json:"ViewQueryCount"`
	ViewQueryDelay     int    `json:"ViewQueryDelay"`
}

type DS struct {
	KSize      int    `json:"KSize"`
	KSeed      string `json:"KSeed"`
	VSize      int    `json:"VSize"`
	VSeed      string `json:"VSeed"`
	Repeat     string `json:"Repeat"`
	Count      int    `json:"Count"`
	Continuous bool   `json:"Continuous"`
}

type ViewSchema struct {
	InflateLevel   int    `json:"InflateLevel"`
	InflateContent string `json:"InflateContent"`
	KIdent         string `json:"KIdent,omitEmpty"`
	KSequence      int    `json:"KVSequence,omitEmpty"`
}

type ViewQueryParameters struct {
	Limit       uint `json:"limit"`
	Stale       bool `json:"stale"`
	Continue    bool `json:"continue"`
	IncludeDocs bool `json:"include_docs"`
	Skip        uint `json:"skip"`
	UpdateAfter bool `json:"update_after"`
}

type ResponseCommand struct {
	Command string      `json:"Command"`
	ReqID   int         `json:"ReqID"`
	ResData interface{} `json:"ResponseData"`
	Handle  int         `json:"Handle"`
	Status  int         `json:"Status"`
}

type EmptyObject struct {
}

type InfoResponse struct {
	CAPS      Caps       `json:"CAPS"`
	CONFIG    Config     `json:"CONFIG"`
	HEADERS   Headers    `json:"HEADERS"`
	TIME      uint64     `json:"TIME"`
	RUNTIME   SDKRuntime `json:"RUNTIME"`
	SDK       string     `json:"SDK"`
	Changeset string     `json:"CHANGESET"`
}

type Caps struct {
	DS_SHARED  bool `json:"DS_SHARED"`
	CANCEL     bool `json:"CANCEL"`
	CONTINUOUS bool `json:"CONTINUOUS"`
	PREAMBLE   bool `json:"PREAMBLE"`
	VIEWS      bool `json:"VIEWS"`
}

type Config struct {
	conncache string `json:"CONNCACHE"`
	ioplugin  string `json:IO_PLUGIN"`
	reconnect int    `json:"RECONNECT"`
}

type Headers struct {
	changeset string `json:"CHANGESET"`
	SDK       string `json:"SDK"`
}

type SDKRuntime struct {
	SDK string `json:"SDK"`
}

type ResultResponse struct {
	Summary map[string]int `json:"Summary"`
	Timings *Timings       `json:"Timings,omitempty"`
}

type Timings struct {
	Base    int64       `json:"Base"`
	Step    int64       `json:"Step"`
	Windows interface{} `json:"Windows"`
}

type Window struct {
	Count  int64          `json:"Count"`
	Avg    int64          `json:"Avg"`
	Min    int64          `json:"Min"`
	Max    int64          `json:"Max"`
	Errors map[string]int `json:"Errors"`
}
