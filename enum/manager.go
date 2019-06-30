package enum

const LEVEL_NORMAL int = 0
const LEVEL_SUCCESS int = 1
const LEVEL_INFO int = 2
const LEVEL_WARN int = 3
const LEVEL_ERROR int = 4

type EnumItem struct {
	Value interface{} `json:"value"`
	Label string      `json:"label"`
	Level int         `json:"level"`
}
type EnumSet struct {
	KeyName string
	m       []EnumItem
}

type EnumManager struct {
	Sets []*EnumSet
}

func NewEnumManager() *EnumManager {

	return &EnumManager{
		Sets: []*EnumSet{},
	}
}

func (e *EnumManager) AddEnumSet(key string, m []EnumItem) *EnumSet {

	instance := EnumSet{
		KeyName: key,
		m:       m,
	}
	e.Sets = append(e.Sets, &instance)
	return &instance
}

func (e *EnumSet) ContainsLabel(value string) bool {
	for _, val := range e.m {
		if value == val.Label {
			return true
		}
	}
	return false
}
func (e *EnumSet) GetData() []EnumItem {
	return e.m
}
func (e *EnumSet) ContainsValue(value interface{}) bool {
	for _, item := range e.m {
		if item.Value == value {
			return true
		}
	}
	return false
}

func (e *EnumSet) GetLabel(key interface{}) string {
	for _, item := range e.m {
		if key == item.Value {
			return item.Label
		}
	}
	return "[ERROR]:LABEL_NOT_FOUND"

}
