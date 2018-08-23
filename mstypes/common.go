package mstypes

// LPWSTR implements https://msdn.microsoft.com/en-us/library/cc230355.aspx
type LPWSTR struct {
	String string `ndr:"pointer,conformant,varying"`
}
