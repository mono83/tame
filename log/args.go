package log

// Method is an xray.Arg interface implementation
type Method string

// Name is xray.Arg interface implementation
func (Method) Name() string { return "method" }

// Value is xray.Arg interface implementation
func (m Method) Value() string { return string(m) }

// Scalar is xray.Arg interface implementation
func (m Method) Scalar() interface{} { return string(m) }

// URL is an xray.Arg interface implementation
type URL string

// Name is xray.Arg interface implementation
func (URL) Name() string { return "url" }

// Value is xray.Arg interface implementation
func (u URL) Value() string { return string(u) }

// Scalar is xray.Arg interface implementation
func (u URL) Scalar() interface{} { return string(u) }
