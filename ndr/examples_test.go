package ndr

//
//import "github.com/jcmturner/rpc/mstypes"
//
//const (
//	PAC_Kerb_Validation_Info_MS = "01100800cccccccca00400000000000000000200d186660f656ac601ffffffffffffff7fffffffffffffff7f17d439fe784ac6011794a328424bc601175424977a81c60108000800040002002400240008000200120012000c0002000000000010000200000000001400020000000000180002005410000097792c00010200001a0000001c000200200000000000000000000000000000000000000016001800200002000a000c002400020028000200000000000000000010000000000000000000000000000000000000000000000000000000000000000d0000002c0002000000000000000000000000000400000000000000040000006c007a00680075001200000000000000120000004c0069007100690061006e00670028004c006100720072007900290020005a00680075000900000000000000090000006e0074006400730032002e0062006100740000000000000000000000000000000000000000000000000000000000000000000000000000001a00000061c433000700000009c32d00070000005eb4320007000000010200000700000097b92c00070000002bf1320007000000ce30330007000000a72e2e00070000002af132000700000098b92c000700000062c4330007000000940133000700000076c4330007000000aefe2d000700000032d22c00070000001608320007000000425b2e00070000005fb4320007000000ca9c35000700000085442d0007000000c2f0320007000000e9ea310007000000ed8e2e0007000000b6eb310007000000ab2e2e0007000000720e2e00070000000c000000000000000b0000004e0054004400450056002d00440043002d003000350000000600000000000000050000004e0054004400450056000000040000000104000000000005150000005951b81766725d2564633b0b0d0000003000020007000000340002000700002038000200070000203c000200070000204000020007000020440002000700002048000200070000204c000200070000205000020007000020540002000700002058000200070000205c00020007000020600002000700002005000000010500000000000515000000b9301b2eb7414c6c8c3b351501020000050000000105000000000005150000005951b81766725d2564633b0b74542f00050000000105000000000005150000005951b81766725d2564633b0be8383200050000000105000000000005150000005951b81766725d2564633b0bcd383200050000000105000000000005150000005951b81766725d2564633b0b5db43200050000000105000000000005150000005951b81766725d2564633b0b41163500050000000105000000000005150000005951b81766725d2564633b0be8ea3100050000000105000000000005150000005951b81766725d2564633b0bc1193200050000000105000000000005150000005951b81766725d2564633b0b29f13200050000000105000000000005150000005951b81766725d2564633b0b0f5f2e00050000000105000000000005150000005951b81766725d2564633b0b2f5b2e00050000000105000000000005150000005951b81766725d2564633b0bef8f3100050000000105000000000005150000005951b81766725d2564633b0b075f2e0000000000"
//)
//
//type KerbValidationInfo struct {
//	LogOnTime               mstypes.FileTime
//	LogOffTime              mstypes.FileTime
//	KickOffTime             mstypes.FileTime
//	PasswordLastSet         mstypes.FileTime
//	PasswordCanChange       mstypes.FileTime
//	PasswordMustChange      mstypes.FileTime
//	EffectiveName           mstypes.RPCUnicodeString
//	FullName                mstypes.RPCUnicodeString
//	LogonScript             mstypes.RPCUnicodeString
//	ProfilePath             mstypes.RPCUnicodeString
//	HomeDirectory           mstypes.RPCUnicodeString
//	HomeDirectoryDrive      mstypes.RPCUnicodeString
//	LogonCount              uint16
//	BadPasswordCount        uint16
//	UserID                  uint32
//	PrimaryGroupID          uint32
//	GroupCount              uint32
//	pGroupIDs               uint32
//	GroupIDs                []mstypes.GroupMembership
//	UserFlags               uint32
//	UserSessionKey          mstypes.UserSessionKey
//	LogonServer             mstypes.RPCUnicodeString
//	LogonDomainName         mstypes.RPCUnicodeString
//	pLogonDomainID          uint32
//	LogonDomainID           mstypes.RPCSID
//	Reserved1               []uint32 // Has 2 elements
//	UserAccountControl      uint32
//	SubAuthStatus           uint32
//	LastSuccessfulILogon    mstypes.FileTime
//	LastFailedILogon        mstypes.FileTime
//	FailedILogonCount       uint32
//	Reserved3               uint32
//	SIDCount                uint32
//	pExtraSIDs              uint32
//	ExtraSIDs               []mstypes.KerbSidAndAttributes
//	pResourceGroupDomainSID uint32
//	ResourceGroupDomainSID  mstypes.RPCSID
//	ResourceGroupCount      uint32
//	pResourceGroupIDs       uint32
//	ResourceGroupIDs        []mstypes.GroupMembership
//}