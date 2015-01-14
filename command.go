package main

const (
	NEWHANDLE        = "NEWHANDLE"
	MC_DS_MUTATE_SET = "MC_DS_MUTATE_SET"
	MC_DS_MUTATE_GET = "MC_DS_MUTATE_GET"
	DSTYPE_SEEDED    = "SEEDED"
)

type RequestCommand struct {
	Command string      `json:"Command"`
	ReqID   int         `json:"ReqID"`
	CmdData CommandData `json:"CommandData"`
	Handle  int         `json:"Handle"`
}

type CommandData struct {
	DSType   string  `json:"DSType,omitempty"`
	DS       DS      `json:"DS,omitempty"`
	Options  Options `json:"Options,omitempty"`
	Bucket   string  `json:"Bucket,omitempty"`
	Port     int     `json:"Port,omitempty"`
	Hostname string  `json:"Hostname,omitempty"`
}

type Options struct {
	OtherNodes         []string `json:"OtherNodes"`
	Username           string   `json:"Username"`
	ClusterCertificate string   `json:"ClusterCertificate"`
	SSL                bool     `json:"SSL"`
	Password           string   `json:"Password"`
	DelayMin           int      `json:"DelayMin"`
	ReplicateTo        int      `json:"ReplicateTo"`
	TimeRes            int64    `json:"TimeRes"`
	PersistTo          int      `json:"PersistTo"`
	ReplicaRead        bool     `json:"ReplicaRead"`
	IterWait           int      `json:"IterWait"`
	DelayMax           int      `json:"DelayMax"`
	Full               bool     `json:"Full"`
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

type Schema struct {
	InflateLevel   int    `json:"InflateLevel"`
	InflateContent string `json:"InflateContent"`
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
	CAPS    Caps    `json:"CAPS"`
	CONFIG  Config  `json:"CONFIG"`
	HEADERS Headers `json:"HEADERS"`
	TIME    int     `json:"TIME"`
	RUNTIME Runtime `json:"RUNTIME"`
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

type Runtime struct {
	SDK string `json:"SDK"`
}

type ResultResponse struct {
	Summary map[string]int `json:"Summary"`
	Timings Timings        `json:"Timings"`
}

type Timings struct {
	Base    int64    `json:"Base"`
	Step    int64    `json:"Step"`
	Windows []Window `json:"Windows"`
}

type Window struct {
	Count  int64          `json:"Count"`
	Avg    int64          `json:"Avg"`
	Min    int64          `json:"Min"`
	Max    int64          `json:"Max"`
	Errors map[string]int `json:"Errors"`
}
