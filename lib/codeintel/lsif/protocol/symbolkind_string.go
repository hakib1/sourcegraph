// Code generated by "stringer -output=bazel-out/darwin_arm64-fastbuild/bin/lib/codeintel/lsif/protocol/symbolkind_string.go -type=SymbolKind lib/codeintel/lsif/protocol/symbol.go"; DO NOT EDIT.

package protocol

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[File-1]
	_ = x[Module-2]
	_ = x[Namespace-3]
	_ = x[Package-4]
	_ = x[Class-5]
	_ = x[Method-6]
	_ = x[Property-7]
	_ = x[Field-8]
	_ = x[Constructor-9]
	_ = x[Enum-10]
	_ = x[Interface-11]
	_ = x[Function-12]
	_ = x[Variable-13]
	_ = x[Constant-14]
	_ = x[String-15]
	_ = x[Number-16]
	_ = x[Boolean-17]
	_ = x[Array-18]
	_ = x[Object-19]
	_ = x[Key-20]
	_ = x[Null-21]
	_ = x[EnumMember-22]
	_ = x[Struct-23]
	_ = x[Event-24]
	_ = x[Operator-25]
	_ = x[TypeParameter-26]
}

const _SymbolKind_name = "FileModuleNamespacePackageClassMethodPropertyFieldConstructorEnumInterfaceFunctionVariableConstantStringNumberBooleanArrayObjectKeyNullEnumMemberStructEventOperatorTypeParameter"

var _SymbolKind_index = [...]uint8{0, 4, 10, 19, 26, 31, 37, 45, 50, 61, 65, 74, 82, 90, 98, 104, 110, 117, 122, 128, 131, 135, 145, 151, 156, 164, 177}

func (i SymbolKind) String() string {
	i -= 1
	if i < 0 || i >= SymbolKind(len(_SymbolKind_index)-1) {
		return "SymbolKind(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _SymbolKind_name[_SymbolKind_index[i]:_SymbolKind_index[i+1]]
}
