package zabbix

import (
	"testing"
)

func TestProxyGroupLifecycle(t *testing.T) {
	api := testGetAPI(t)

	// Create
	group := ProxyGroup{Name: "test_proxygroup_create"}
	groupsToCreate := ProxyGroups{group}
	err := api.ProxyGroupsCreate(groupsToCreate)
	if err != nil {
		t.Fatalf("ProxyGroupsCreate failed: %v", err)
	}
	if groupsToCreate[0].ProxyGroupID == "" {
		t.Fatal("ProxyGroupsCreate did not return a proxy group ID")
	}
	createdGroupID := groupsToCreate[0].ProxyGroupID
	t.Logf("Created proxy group with ID: %s", createdGroupID)

	// Get by ID
	retrievedGroupByID, err := api.ProxyGroupGetByID(createdGroupID)
	if err != nil {
		t.Fatalf("ProxyGroupGetByID failed: %v", err)
	}
	if retrievedGroupByID == nil {
		t.Fatal("ProxyGroupGetByID returned nil")
	}
	if retrievedGroupByID.Name != group.Name {
		t.Errorf("ProxyGroupGetByID: Expected name %s, got %s", group.Name, retrievedGroupByID.Name)
	}

	// Get by Name
	retrievedGroupByName, err := api.ProxyGroupGetByName(group.Name)
	if err != nil {
		t.Fatalf("ProxyGroupGetByName failed: %v", err)
	}
	if retrievedGroupByName == nil {
		t.Fatal("ProxyGroupGetByName returned nil")
	}
	if retrievedGroupByName.ProxyGroupID != createdGroupID {
		t.Errorf("ProxyGroupGetByName: Expected ID %s, got %s", createdGroupID, retrievedGroupByName.ProxyGroupID)
	}

	// Update
	updatedGroup := ProxyGroup{
		ProxyGroupID: createdGroupID,
		Name:         "test_proxygroup_update",
	}
	groupsToUpdate := ProxyGroups{updatedGroup}
	err = api.ProxyGroupsUpdate(groupsToUpdate)
	if err != nil {
		t.Fatalf("ProxyGroupsUpdate failed: %v", err)
	}

	// Verify Update
	verifiedGroup, err := api.ProxyGroupGetByID(createdGroupID)
	if err != nil {
		t.Fatalf("ProxyGroupGetByID after update failed: %v", err)
	}
	if verifiedGroup.Name != updatedGroup.Name {
		t.Errorf("ProxyGroupsUpdate: Expected name %s, got %s", updatedGroup.Name, verifiedGroup.Name)
	}
	t.Logf("Updated proxy group name to: %s", verifiedGroup.Name)

	// Delete
	err = api.ProxyGroupsDelete([]string{createdGroupID})
	if err != nil {
		t.Fatalf("ProxyGroupsDelete failed: %v", err)
	}

	// Verify Delete
	_, err = api.ProxyGroupGetByID(createdGroupID)
	if err == nil {
		t.Fatal("ProxyGroupGetByID succeeded after delete, group was not deleted")
	}
	t.Logf("Successfully deleted proxy group ID: %s", createdGroupID)
}
