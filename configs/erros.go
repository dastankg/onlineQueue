package configs

func ErrMissingEnvVar(varName string) error {
	return &MissingEnvVarError{VarName: varName}
}

type MissingEnvVarError struct {
	VarName string
}

func (e *MissingEnvVarError) Error() string {
	return "Missing required environment variable: " + e.VarName
}
