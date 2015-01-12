package main

const (
	MC_DS_MUTATE_SET = "MC_DS_MUTATE_SET"
	MC_DS_MUTATE_GET = "MC_DS_MUTATE_GET"
	DSTYPE_SEEDED    = "SEEDED"
)

type RequestCommand struct {
	Command string      `json:"Command"`
	ReqID   int         `json:"ReqID"`
	cmdData interface{} `json:"CommandData"`
	Handle  int         `json:"Handle"`
}

type CommandData struct {
	DSType  string  `json:"DSType"`
	ds      DS      `json:"DS"`
	Options options `json:"Options"`

}

type Options struct {
    OtherNodes []string `json:"OtherNodes"`
    Username string `json:"Username"`
    ClusterCertificate `json:"Clu"
}

type DS struct {
	KSize  string `json:"KSize"`
	KSeed  string `json:"KSeed"`
	VSize  string `json:"VSize"`
	VSeed  string `json:"VSeed"`
	Repeat string `json:"Repeat"`
	Count  int    `json:"Count"`
}

type ResponseCommand struct {
	Command string      `json:"Command"`
	ReqId   int         `json:"ReqID"`
	ResData interface{} `json:"ResponseData"`
	Handle  int         `json:"HANDLE"`
	Status  int         `json:"STATUS"`
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
