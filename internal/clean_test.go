/*
Copyright © 2025 PATRICK HERMANN PATRICK.HERMANN@SVA.DE
*/

package internal

import (
	"testing"
)

func TestCleanString(t *testing.T) {
	input := `
playbooks:
- sthings.baseos.prepare_env
- sthings.baseos.setup
- sthings.baseos.dev
-
ansibleVarsFile:
- '''manage_filesystem+-true'''
- '''update_packages+-true'''
- hello
-
`
	expected := `playbooks:
- sthings.baseos.prepare_env
- sthings.baseos.setup
- sthings.baseos.dev
ansibleVarsFile:
- hello`

	output := CleanString(input)

	if output != expected {
		t.Errorf("cleanString failed:\nExpected:\n%s\nGot:\n%s", expected, output)
	}
}
