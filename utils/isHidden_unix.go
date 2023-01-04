//go:build !windows

package utils

// IsHidden returns true if the file/directory is hidden
//
//	@param path - the path to the file/directory
//	@param forceHidden - true if the file/directory is forced to be hidden
//	@return bool - true if the file/directory is hidden
func IsHidden(path string, forceHidden bool) bool {
	// check if the file/directory is forced to be hidden
	if forceHidden {
		return true
	}

	// check if the file/directory is hidden by the OS (e.g. .git, .vscode)
	return int(path[0]) == DotCharacter
}
