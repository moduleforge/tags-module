package service

// This file is intentionally minimal. Service tests wire mock queriers via the
// injectable newCoreQuerier / newTagQuerier factories on TagService and a
// fakeDB that provides a fakeTx satisfying pgx.Tx. The fakeTx and fakeDB
// types are defined in mock_test.go.
//
// Transaction-level behaviour (commit/rollback) is exercised by the tests in
// tag_test.go that check observer calls and error propagation.
