// Go MySQL Driver - A MySQL-Driver for Go's database/sql package
//
// Copyright 2013 The Go-MySQL-Driver Authors. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at http://mozilla.org/MPL/2.0/.

package mysql

import (
	"bytes"
	"errors"
	"log"
	"testing"
)

func TestErrorsSetLogger(t *testing.T) {
	previous := getLogger()
	defer func() {
		SetLogger(previous)
	}()

	// set up logger
	const expected = "prefix: test\n"
	buffer := bytes.NewBuffer(make([]byte, 0, 64))
	logger := log.New(buffer, "prefix: ", 0)

	// print
	SetLogger(logger)
	getLogger().Print("test")

	// check result
	if actual := buffer.String(); actual != expected {
		t.Errorf("expected %q, got %q", expected, actual)
	}
}

func TestErrorsStrictIgnoreNotes(t *testing.T) {
	runTests(t, dsn+"&sql_notes=false", func(dbt *DBTest) {
		dbt.mustExec("DROP TABLE IF EXISTS does_not_exist")
	})
}

func TestMySQLErrIs(t *testing.T) {
	infraErr := &MySQLError{Number: 1234, Message: "the server is on fire"}
	otherInfraErr := &MySQLError{Number: 1234, Message: "the datacenter is flooded"}
	if !errors.Is(infraErr, otherInfraErr) {
		t.Errorf("expected errors to be the same: %+v %+v", infraErr, otherInfraErr)
	}

	differentCodeErr := &MySQLError{Number: 5678, Message: "the server is on fire"}
	if errors.Is(infraErr, differentCodeErr) {
		t.Fatalf("expected errors to be different: %+v %+v", infraErr, differentCodeErr)
	}

	nonMysqlErr := errors.New("not a mysql error")
	if errors.Is(infraErr, nonMysqlErr) {
		t.Fatalf("expected errors to be different: %+v %+v", infraErr, nonMysqlErr)
	}
}
