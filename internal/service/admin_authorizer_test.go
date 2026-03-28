package service

import "testing"

func TestAdminAuthorizerAssistantCannotWriteSettings(t *testing.T) {
	t.Parallel()

	authz, err := NewInMemoryAdminAuthorizer()
	if err != nil {
		t.Fatalf("NewInMemoryAdminAuthorizer() error = %v", err)
	}
	if err := authz.AssignRole("2", AdminRoleAssistant); err != nil {
		t.Fatalf("AssignRole() error = %v", err)
	}

	ok, err := authz.Enforce("2", AdminDomain, "settings", "write")
	if err != nil {
		t.Fatalf("Enforce() error = %v", err)
	}
	if ok {
		t.Fatal("assistant should not have settings:write")
	}
}

func TestAdminAuthorizerAssistantCanWriteProblems(t *testing.T) {
	t.Parallel()

	authz, err := NewInMemoryAdminAuthorizer()
	if err != nil {
		t.Fatalf("NewInMemoryAdminAuthorizer() error = %v", err)
	}
	if err := authz.AssignRole("2", AdminRoleAssistant); err != nil {
		t.Fatalf("AssignRole() error = %v", err)
	}

	ok, err := authz.Enforce("2", AdminDomain, "problems", "write")
	if err != nil {
		t.Fatalf("Enforce() error = %v", err)
	}
	if !ok {
		t.Fatal("assistant should have problems:write")
	}
}
