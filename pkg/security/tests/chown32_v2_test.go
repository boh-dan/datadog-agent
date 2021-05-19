// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

// +build functionaltests,!386

package tests

import (
	"os"
	"os/exec"
	"testing"

	"github.com/DataDog/datadog-agent/pkg/security/rules"
	"gotest.tools/assert"
)

func TestChown32(t *testing.T) {
	ruleDef := &rules.RuleDefinition{
		ID:         "test_rule",
		Expression: `chown.file.path == "{{.Root}}/test-chown" && chown.file.destination.uid in [100, 101, 102, 103, 104, 105, 106] && chown.file.destination.gid in [200, 201, 202, 203, 204, 205, 206]`,
	}

	ruleDef2 := &rules.RuleDefinition{
		ID:         "test_rule2",
		Expression: `chown.file.path == "{{.Root}}/test-symlink" && chown.file.destination.uid in [100, 101, 102, 103, 104, 105, 106] && chown.file.destination.gid in [200, 201, 202, 203, 204, 205, 206]`,
	}

	test, err := newTestModule(nil, []*rules.RuleDefinition{ruleDef, ruleDef2}, testOpts{})
	if err != nil {
		t.Fatal(err)
	}
	defer test.Close()

	syscallTester, err := loadSyscallTester(t)
	if err != nil {
		t.Fatal(err)
	}

	prevUID := 98
	prevGID := 99
	fileMode := 0o447
	expectedMode := uint32(applyUmask(fileMode))
	testFile, _, err := test.CreateWithOptions("test-chown", 98, 99, fileMode)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("chown", func(t *testing.T) {
		sideTester := exec.Command(syscallTester, "chown", testFile, "100", "200")
		if output, err := sideTester.CombinedOutput(); err != nil {
			t.Error(string(output))
			t.Error(err)
		}

		defer func() {
			prevUID = 100
			prevGID = 200
		}()

		event, _, err := test.GetEvent()
		if err != nil {
			t.Error(err)
		} else {
			assert.Equal(t, event.GetType(), "chown", "wrong event type")
			assert.Equal(t, event.Chown.UID, uint32(100), "wrong user")
			assert.Equal(t, event.Chown.GID, uint32(200), "wrong user")
			assert.Equal(t, event.Chown.File.Inode, getInode(t, testFile), "wrong inode")
			assertRights(t, event.Chown.File.Mode, uint16(expectedMode), "wrong initial mode")
			assert.Equal(t, event.Chown.File.UID, uint32(prevUID), "wrong initial user")
			assert.Equal(t, event.Chown.File.GID, uint32(prevGID), "wrong initial group")

			assertNearTime(t, event.Chown.File.MTime)
			assertNearTime(t, event.Chown.File.CTime)

			if testEnvironment == DockerEnvironment {
				testContainerPath(t, event, "chown.file.container_path")
			}

			if !validateChownSchema(t, event) {
				t.Fatal(event.String())
			}
		}
	})

	t.Run("fchown", func(t *testing.T) {
		sideTester := exec.Command(syscallTester, "fchown", testFile, "101", "201")
		if output, err := sideTester.CombinedOutput(); err != nil {
			t.Error(string(output))
			t.Error(err)
		}

		defer func() {
			prevUID = 101
			prevGID = 201
		}()

		event, _, err := test.GetEvent()
		if err != nil {
			t.Error(err)
		} else {
			assert.Equal(t, event.GetType(), "chown", "wrong event type")
			assert.Equal(t, event.Chown.UID, uint32(101), "wrong user")
			assert.Equal(t, event.Chown.GID, uint32(201), "wrong user")
			assert.Equal(t, event.Chown.File.Inode, getInode(t, testFile), "wrong inode")
			assertRights(t, event.Chown.File.Mode, uint16(expectedMode), "wrong initial mode")
			assert.Equal(t, event.Chown.File.UID, uint32(prevUID), "wrong initial user")
			assert.Equal(t, event.Chown.File.GID, uint32(prevGID), "wrong initial group")

			assertNearTime(t, event.Chown.File.MTime)
			assertNearTime(t, event.Chown.File.CTime)

			if testEnvironment == DockerEnvironment {
				testContainerPath(t, event, "chown.file.container_path")
			}

			if !validateChownSchema(t, event) {
				t.Fatal(event.String())
			}
		}
	})

	t.Run("fchownat", func(t *testing.T) {
		sideTester := exec.Command(syscallTester, "fchownat", testFile, "102", "202")
		if output, err := sideTester.CombinedOutput(); err != nil {
			t.Error(string(output))
			t.Error(err)
		}

		defer func() {
			prevUID = 102
			prevGID = 202
		}()

		event, _, err := test.GetEvent()
		if err != nil {
			t.Error(err)
		} else {
			assert.Equal(t, event.GetType(), "chown", "wrong event type")
			assert.Equal(t, event.Chown.UID, uint32(102), "wrong user")
			assert.Equal(t, event.Chown.GID, uint32(202), "wrong user")
			assert.Equal(t, event.Chown.File.Inode, getInode(t, testFile), "wrong inode")
			assertRights(t, event.Chown.File.Mode, uint16(expectedMode), "wrong initial mode")
			assert.Equal(t, event.Chown.File.UID, uint32(prevUID), "wrong initial user")
			assert.Equal(t, event.Chown.File.GID, uint32(prevGID), "wrong initial group")

			assertNearTime(t, event.Chown.File.MTime)
			assertNearTime(t, event.Chown.File.CTime)

			if testEnvironment == DockerEnvironment {
				testContainerPath(t, event, "chown.file.container_path")
			}

			if !validateChownSchema(t, event) {
				t.Fatal(event.String())
			}
		}
	})

	t.Run("lchown", func(t *testing.T) {
		testSymlink, _, err := test.Path("test-symlink")
		if err != nil {
			t.Fatal(err)
		}

		if err := os.Symlink(testFile, testSymlink); err != nil {
			t.Fatal(err)
		}
		defer os.Remove(testSymlink)

		sideTester := exec.Command(syscallTester, "lchown", testSymlink, "103", "203")
		if output, err := sideTester.CombinedOutput(); err != nil {
			t.Error(string(output))
			t.Error(err)
		}

		event, _, err := test.GetEvent()
		if err != nil {
			t.Error(err)
		} else {
			assert.Equal(t, event.GetType(), "chown", "wrong event type")
			assert.Equal(t, event.Chown.UID, uint32(103), "wrong user")
			assert.Equal(t, event.Chown.GID, uint32(203), "wrong user")
			assert.Equal(t, event.Chown.File.Inode, getInode(t, testSymlink), "wrong inode")
			assertRights(t, event.Chown.File.Mode, uint16(0o777), "wrong initial mode")
			assert.Equal(t, event.Chown.File.UID, uint32(0), "wrong initial user")
			assert.Equal(t, event.Chown.File.GID, uint32(0), "wrong initial group")

			assertNearTime(t, event.Chown.File.MTime)
			assertNearTime(t, event.Chown.File.CTime)

			if testEnvironment == DockerEnvironment {
				testContainerPath(t, event, "chown.file.container_path")
			}

			if !validateChownSchema(t, event) {
				t.Fatal(event.String())
			}
		}
	})

	t.Run("lchown32", func(t *testing.T) {
		testSymlink, _, err := test.Path("test-symlink")
		if err != nil {
			t.Fatal(err)
		}

		if err := os.Symlink(testFile, testSymlink); err != nil {
			t.Fatal(err)
		}
		defer os.Remove(testSymlink)

		sideTester := exec.Command(syscallTester, "lchown32", testSymlink, "104", "204")
		if output, err := sideTester.CombinedOutput(); err != nil {
			t.Error(string(output))
			t.Error(err)
		}

		event, _, err := test.GetEvent()
		if err != nil {
			t.Error(err)
		} else {
			assert.Equal(t, event.GetType(), "chown", "wrong event type")
			assert.Equal(t, event.Chown.UID, uint32(104), "wrong user")
			assert.Equal(t, event.Chown.GID, uint32(204), "wrong user")
			assert.Equal(t, event.Chown.File.Inode, getInode(t, testSymlink), "wrong inode")
			assertRights(t, event.Chown.File.Mode, uint16(0o777), "wrong initial mode")
			assert.Equal(t, event.Chown.File.UID, uint32(0), "wrong initial user")
			assert.Equal(t, event.Chown.File.GID, uint32(0), "wrong initial group")

			assertNearTime(t, event.Chown.File.MTime)
			assertNearTime(t, event.Chown.File.CTime)

			if testEnvironment == DockerEnvironment {
				testContainerPath(t, event, "chown.file.container_path")
			}

			if !validateChownSchema(t, event) {
				t.Fatal(event.String())
			}
		}
	})

	t.Run("fchown32", func(t *testing.T) {
		sideTester := exec.Command(syscallTester, "fchown32", testFile, "105", "205")
		if output, err := sideTester.CombinedOutput(); err != nil {
			t.Error(string(output))
			t.Error(err)
		}

		defer func() {
			prevUID = 105
			prevGID = 205
		}()

		event, _, err := test.GetEvent()
		if err != nil {
			t.Error(err)
		} else {
			assert.Equal(t, event.GetType(), "chown", "wrong event type")
			assert.Equal(t, event.Chown.UID, uint32(105), "wrong user")
			assert.Equal(t, event.Chown.GID, uint32(205), "wrong user")
			assert.Equal(t, event.Chown.File.Inode, getInode(t, testFile), "wrong inode")
			assertRights(t, event.Chown.File.Mode, uint16(expectedMode), "wrong initial mode")
			assert.Equal(t, event.Chown.File.UID, uint32(prevUID), "wrong initial user")
			assert.Equal(t, event.Chown.File.GID, uint32(prevGID), "wrong initial group")

			assertNearTime(t, event.Chown.File.MTime)
			assertNearTime(t, event.Chown.File.CTime)

			if testEnvironment == DockerEnvironment {
				testContainerPath(t, event, "chown.file.container_path")
			}

			if !validateChownSchema(t, event) {
				t.Fatal(event.String())
			}
		}
	})

	t.Run("chown32", func(t *testing.T) {
		sideTester := exec.Command(syscallTester, "chown32", testFile, "106", "206")
		if output, err := sideTester.CombinedOutput(); err != nil {
			t.Error(string(output))
			t.Error(err)
		}

		defer func() {
			prevUID = 106
			prevGID = 206
		}()

		event, _, err := test.GetEvent()
		if err != nil {
			t.Error(err)
		} else {
			assert.Equal(t, event.GetType(), "chown", "wrong event type")
			assert.Equal(t, event.Chown.UID, uint32(106), "wrong user")
			assert.Equal(t, event.Chown.GID, uint32(206), "wrong user")
			assert.Equal(t, event.Chown.File.Inode, getInode(t, testFile), "wrong inode")
			assertRights(t, event.Chown.File.Mode, uint16(expectedMode), "wrong initial mode")
			assert.Equal(t, event.Chown.File.UID, uint32(prevUID), "wrong initial user")
			assert.Equal(t, event.Chown.File.GID, uint32(prevGID), "wrong initial group")

			assertNearTime(t, event.Chown.File.MTime)
			assertNearTime(t, event.Chown.File.CTime)

			if testEnvironment == DockerEnvironment {
				testContainerPath(t, event, "chown.file.container_path")
			}

			if !validateChownSchema(t, event) {
				t.Fatal(event.String())
			}
		}
	})
}
