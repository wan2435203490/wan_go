package music

type RandomSong struct {
	TypeName string   `json:"typeName"`
	Songs    *[]*Song `json:"songs"`
}
