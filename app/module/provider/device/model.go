package device

type Model struct {
	Device struct {
		Name  string
		Type  string
		Brand string
	}
	Browser struct {
		Name    string
		Type    string
		Version string
	}
	Platform struct {
		Name    string
		Short   string
		Version string
	}
}
