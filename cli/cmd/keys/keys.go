package keys

const (
	Init_Flag = "init-flag"
)

// type Policy interface {
// 	GetValue() string
// 	GetShortHand() string
// }
// type IOrganization struct {
// 	Name, ShortHand string
// }

// func (i *IOrganization) GetValue() string {
// 	return i.Name
// }
// func (i *IOrganization) GetShortHand() string {
// 	return i.Name
// }

// func Get[T Policy](key string) T {
// 	switch key {
// 	case OrganizationKey:
// 		p := &IOrganization{
// 			Name:      "org",
// 			ShortHand: "o",
// 		}
// 		return any(&p).(T)
// 	}
// 	var t T
// 	return t
// }
