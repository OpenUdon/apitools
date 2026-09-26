package operationlifecycle

import (
	"slices"
	"testing"

	"github.com/OpenUdon/apitools"
)

func TestExpandRequiresScopedItemPaths(t *testing.T) {
	tests := []struct {
		name       string
		operations []apitools.OperationSummary
		want       []string
	}{
		{
			name: "reject collection-level read",
			operations: []apitools.OperationSummary{
				op("widgets", "createWidget", "POST", "/projects/{project}/widgets"),
				op("widgets", "getWidgets", "GET", "/projects/{project}/widgets"),
			},
			want: []string{"post:createWidget"},
		},
		{
			name: "reject item with missing parent scope",
			operations: []apitools.OperationSummary{
				op("widgets", "createWidget", "POST", "/projects/{project}/widgets"),
				op("widgets", "getWidget", "GET", "/projects/widgets/{widget}"),
			},
			want: []string{"post:createWidget"},
		},
		{
			name: "accept same scoped collection and trailing item parameter",
			operations: []apitools.OperationSummary{
				op("widgets", "createWidget", "POST", "/projects/{project}/widgets"),
				op("widgets", "fetchWidget", "GET", "/projects/{project}/widgets/{widget}"),
			},
			want: []string{"create:createWidget", "read:fetchWidget"},
		},
		{
			name: "accept resource named list",
			operations: []apitools.OperationSummary{
				op("widgets", "createList", "POST", "/lists"),
				op("widgets", "getList", "GET", "/lists/{list_id}"),
			},
			want: []string{"create:createList", "read:getList"},
		},
		{
			name: "accept resource named collection",
			operations: []apitools.OperationSummary{
				op("widgets", "createCollection", "POST", "/collections"),
				op("widgets", "deleteCollection", "DELETE", "/collections/{collection_id}"),
			},
			want: []string{"create:createCollection", "delete:deleteCollection"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			expanded := Expand(test.operations, test.operations[0], Options{DesiredState: true})
			if got := roleIDs(expanded); !slices.Equal(got, test.want) {
				t.Fatalf("roles = %#v, want %#v; diagnostics = %#v", got, test.want, expanded.Diagnostics)
			}
		})
	}
}

func TestExpandGoogleDiscoveryNameAndParentUseMethodResourceIdentity(t *testing.T) {
	operations := []apitools.OperationSummary{
		googleOp("compute", "compute.projects.instances.insert", "POST", "/v1/{+parent}/instances"),
		googleOp("compute", "compute.projects.instances.get", "GET", "/v1/{+name}"),
		googleOp("compute", "compute.projects.disks.get", "GET", "/v1/{+name}"),
	}
	expanded := Expand(operations, operations[0], Options{DesiredState: true})
	if got, want := roleIDs(expanded), []string{
		"create:compute.projects.instances.insert",
		"read:compute.projects.instances.get",
	}; !slices.Equal(got, want) {
		t.Fatalf("roles = %#v, want %#v; diagnostics = %#v", got, want, expanded.Diagnostics)
	}
}

func TestExpandGoogleDiscoveryParentParameterIsNotAnItemID(t *testing.T) {
	operations := []apitools.OperationSummary{
		googleOp("compute", "compute.projects.instances.insert", "POST", "/v1/{+parent}"),
		googleOp("compute", "compute.projects.instances.get", "GET", "/v1/{+parent}"),
	}
	expanded := Expand(operations, operations[0], Options{DesiredState: true})
	if got, want := roleIDs(expanded), []string{"post:compute.projects.instances.insert"}; !slices.Equal(got, want) {
		t.Fatalf("roles = %#v, want %#v; diagnostics = %#v", got, want, expanded.Diagnostics)
	}
}

func TestExpandClassifiesPostUpdatesAndActions(t *testing.T) {
	updateOperations := []apitools.OperationSummary{
		op("pets", "updatePetWithForm", "POST", "/pets/{petId}"),
		op("pets", "getPet", "GET", "/pets/{petId}"),
	}
	updateExpansion := Expand(updateOperations, updateOperations[0], Options{DesiredState: true})
	if got, want := roleIDs(updateExpansion), []string{"update:updatePetWithForm", "read:getPet"}; !slices.Equal(got, want) {
		t.Fatalf("POST update roles = %#v, want %#v; diagnostics = %#v", got, want, updateExpansion.Diagnostics)
	}

	actionOperations := []apitools.OperationSummary{
		op("tickets", "closeTicket", "POST", "/tickets/{ticketId}"),
		op("tickets", "getTicket", "GET", "/tickets/{ticketId}"),
	}
	actionExpansion := Expand(actionOperations, actionOperations[0], Options{DesiredState: true})
	if got, want := roleIDs(actionExpansion), []string{"post:closeTicket", "read:getTicket"}; !slices.Equal(got, want) {
		t.Fatalf("POST action roles = %#v, want %#v; diagnostics = %#v", got, want, actionExpansion.Diagnostics)
	}

	unknownActionOperations := []apitools.OperationSummary{
		op("invoices", "renewInvoice", "POST", "/invoices/{invoiceId}"),
		op("invoices", "getInvoice", "GET", "/invoices/{invoiceId}"),
	}
	unknownActionExpansion := Expand(unknownActionOperations, unknownActionOperations[0], Options{DesiredState: true})
	if got, want := roleIDs(unknownActionExpansion), []string{"post:renewInvoice", "read:getInvoice"}; !slices.Equal(got, want) {
		t.Fatalf("unlisted POST action roles = %#v, want %#v; diagnostics = %#v", got, want, unknownActionExpansion.Diagnostics)
	}

	createOperations := []apitools.OperationSummary{
		op("projects", "createChild", "POST", "/projects/{projectId}/children"),
		op("projects", "getChild", "GET", "/projects/{projectId}/children/{childId}"),
	}
	if got, want := roleIDs(Expand(createOperations, createOperations[0], Options{DesiredState: true})), []string{"create:createChild", "read:getChild"}; !slices.Equal(got, want) {
		t.Fatalf("explicit POST create roles = %#v, want %#v", got, want)
	}
}

func TestFamilyTokensUseOnlyMeaningfulOperationIdentity(t *testing.T) {
	operation := op("billing", "getInvoice", "GET", "/billing/invoices/{invoice_id}")
	operation.Summary = "Tenant account invoice"
	operation.Tags = []string{"tenant", "account"}
	tokens := familyTokens(operation)
	if !slices.Contains(tokens, "invoice") && !slices.Contains(tokens, "invoices") {
		t.Fatalf("family tokens = %#v, want invoice identity from operation ID or path", tokens)
	}
	if slices.Contains(tokens, "tenant") || slices.Contains(tokens, "account") {
		t.Fatalf("family tokens include free-text summary/tag terms: %#v", tokens)
	}
}

func TestSameFamilyRejectsStopWordOnlyOverlap(t *testing.T) {
	seed := op("generic", "createOperation", "POST", "/api/v1/operations/{id}")
	candidate := op("generic", "getOperation", "GET", "/api/v1/operations/{id}")
	if sameFamily(seed, candidate) {
		t.Fatalf("generic operation/API/version/ID tokens established a family: %#v / %#v", familyTokens(seed), familyTokens(candidate))
	}
}

func TestGoalWantsUpdateUsesWordBoundaries(t *testing.T) {
	for _, goal := range []string{"dispatch a thing", "the updater flow", "no patching required"} {
		if goalWantsUpdate(goal) {
			t.Errorf("goal %q unexpectedly requests update", goal)
		}
	}
	for _, goal := range []string{"update a thing", "please patch this resource", "modify-or-replace the item"} {
		if !goalWantsUpdate(goal) {
			t.Errorf("goal %q should request update", goal)
		}
	}
}

func TestDispatchGoalDoesNotExpandUpdateRole(t *testing.T) {
	operations := []apitools.OperationSummary{
		op("things", "createThing", "POST", "/things"),
		op("things", "patchThing", "PATCH", "/things/{id}"),
	}
	if got, want := roleIDs(Expand(operations, operations[0], Options{Goal: "dispatch a thing", DesiredState: true})), []string{"post:createThing"}; !slices.Equal(got, want) {
		t.Fatalf("dispatch roles = %#v, want %#v", got, want)
	}
	if got, want := roleIDs(Expand(operations, operations[0], Options{Goal: "please patch the thing", DesiredState: true})), []string{"create:createThing", "update:patchThing"}; !slices.Equal(got, want) {
		t.Fatalf("patch roles = %#v, want %#v", got, want)
	}
}

func TestExpandResolvesSeedByOperationID(t *testing.T) {
	operations := []apitools.OperationSummary{
		op("widgets", "createWidget", "POST", "/widgets"),
		op("widgets", "getWidget", "GET", "/widgets/{id}"),
	}
	seed := apitools.OperationSummary{OperationID: "createWidget"}
	expanded := Expand(operations, seed, Options{DesiredState: true})
	if got, want := roleIDs(expanded), []string{"create:createWidget", "read:getWidget"}; !slices.Equal(got, want) {
		t.Fatalf("roles = %#v, want %#v; diagnostics = %#v", got, want, expanded.Diagnostics)
	}
}

func TestOperationIDOnlyAmbiguousSeedDoesNotChooseDocument(t *testing.T) {
	first := op("api-a", "createWidget", "POST", "/widgets")
	first.DocumentPath = "/tmp/api-a/openapi.yaml"
	second := op("api-b", "createWidget", "POST", "/widgets")
	second.DocumentPath = "/tmp/api-b/openapi.yaml"
	expanded := Expand([]apitools.OperationSummary{first, second}, apitools.OperationSummary{OperationID: "createWidget"}, Options{DesiredState: true})
	if len(expanded.Roles) != 0 || len(expanded.Diagnostics) != 1 || expanded.Diagnostics[0].Code != "operation_lifecycle.seed_ambiguous" {
		t.Fatalf("expansion = %#v, want an ambiguity diagnostic with no selected role", expanded)
	}
}

func TestSameRelativePathDifferentAbsoluteDocumentsAreNotSiblings(t *testing.T) {
	seed := op("api", "createWidget", "POST", "/widgets")
	seed.DocumentPath = "/tmp/project-a/openapi/api.yaml"
	seed.DocumentRelativePath = "openapi/api.yaml"
	candidate := op("api", "getWidget", "GET", "/widgets/{id}")
	candidate.DocumentPath = "/tmp/project-b/openapi/api.yaml"
	candidate.DocumentRelativePath = "openapi/api.yaml"
	expanded := Expand([]apitools.OperationSummary{seed, candidate}, seed, Options{DesiredState: true})
	if got, want := roleIDs(expanded), []string{"post:createWidget"}; !slices.Equal(got, want) {
		t.Fatalf("roles = %#v, want %#v", got, want)
	}
}

func TestAbsoluteDocumentPathAndURLOutrankRelativePath(t *testing.T) {
	tests := []struct {
		name string
		set  func(*apitools.OperationSummary, *apitools.OperationSummary)
	}{
		{
			name: "absolute file path",
			set: func(seed, candidate *apitools.OperationSummary) {
				seed.DocumentPath, candidate.DocumentPath = "/tmp/project/openapi/api.yaml", "/tmp/project/openapi/api.yaml"
				seed.DocumentRelativePath, candidate.DocumentRelativePath = "seed.yaml", "candidate.yaml"
			},
		},
		{
			name: "absolute URL",
			set: func(seed, candidate *apitools.OperationSummary) {
				seed.DocumentURL, candidate.DocumentURL = "https://example.test/openapi.json", "https://example.test/openapi.json"
				seed.DocumentRelativePath, candidate.DocumentRelativePath = "seed.yaml", "candidate.yaml"
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			seed := op("api", "createWidget", "POST", "/widgets")
			candidate := op("api", "getWidget", "GET", "/widgets/{id}")
			test.set(&seed, &candidate)
			expanded := Expand([]apitools.OperationSummary{seed, candidate}, seed, Options{DesiredState: true})
			if got, want := roleIDs(expanded), []string{"create:createWidget", "read:getWidget"}; !slices.Equal(got, want) {
				t.Fatalf("roles = %#v, want %#v; diagnostics = %#v", got, want, expanded.Diagnostics)
			}
		})
	}
}

// TestUnderSpecifiedSeedResolvesAgainstFullyDescribedOperation proves that a
// seed supplied with only a relative path (no absolute document path) still
// resolves against a fully described operation from the same document,
// instead of being declared a source mismatch merely because the seed lacks
// the higher-priority absolute-path identifier the pool operation carries.
func TestUnderSpecifiedSeedResolvesAgainstFullyDescribedOperation(t *testing.T) {
	full := apitools.OperationSummary{
		ID: "createWidget", OperationID: "createWidget", Method: "POST", Path: "/widgets",
		DocumentPath: "/abs/openapi/w.yaml", DocumentRelativePath: "openapi/w.yaml",
	}
	sibling := apitools.OperationSummary{
		ID: "getWidget", OperationID: "getWidget", Method: "GET", Path: "/widgets/{id}",
		DocumentPath: "/abs/openapi/w.yaml", DocumentRelativePath: "openapi/w.yaml",
	}
	seed := apitools.OperationSummary{OperationID: "createWidget", DocumentRelativePath: "openapi/w.yaml"}
	expanded := Expand([]apitools.OperationSummary{full, sibling}, seed, Options{DesiredState: true})
	if got, want := roleIDs(expanded), []string{"create:createWidget", "read:getWidget"}; !slices.Equal(got, want) {
		t.Fatalf("roles = %#v, want %#v; diagnostics = %#v", got, want, expanded.Diagnostics)
	}
}

// TestFullySpecifiedSeedResolvesDuplicateOperationIDByMethodAndPath proves
// that a seed built from a real operation (carrying its own method and path)
// resolves to its exact match among several operations that share its
// operation ID in the same document, instead of being declared ambiguous.
func TestFullySpecifiedSeedResolvesDuplicateOperationIDByMethodAndPath(t *testing.T) {
	create := apitools.OperationSummary{
		ID: "createWidget", OperationID: "createWidget", Method: "POST", Path: "/widgets",
		DocumentPath: "/abs/openapi/w.yaml",
	}
	copyAction := apitools.OperationSummary{
		ID: "createWidget", OperationID: "createWidget", Method: "POST", Path: "/widgets/{id}/copy",
		DocumentPath: "/abs/openapi/w.yaml",
	}
	getWidget := apitools.OperationSummary{
		ID: "getWidget", OperationID: "getWidget", Method: "GET", Path: "/widgets/{id}",
		DocumentPath: "/abs/openapi/w.yaml",
	}
	operations := []apitools.OperationSummary{create, copyAction, getWidget}
	expanded := Expand(operations, create, Options{DesiredState: true})
	if got, want := roleIDs(expanded), []string{"create:createWidget", "read:getWidget"}; !slices.Equal(got, want) {
		t.Fatalf("roles = %#v, want %#v; diagnostics = %#v", got, want, expanded.Diagnostics)
	}

	// A bare operation-ID-only seed (no method/path to disambiguate) for the
	// same duplicated ID must still be reported as ambiguous.
	bareSeed := apitools.OperationSummary{OperationID: "createWidget", DocumentPath: "/abs/openapi/w.yaml"}
	ambiguous := Expand(operations, bareSeed, Options{DesiredState: true})
	if len(ambiguous.Roles) != 0 || len(ambiguous.Diagnostics) != 1 || ambiguous.Diagnostics[0].Code != "operation_lifecycle.seed_ambiguous" {
		t.Fatalf("expansion = %#v, want an ambiguity diagnostic with no selected role", ambiguous)
	}
}

func TestExpandCollectionItemLifecycle(t *testing.T) {
	operations := []apitools.OperationSummary{
		op("k8s", "createCoreV1NamespacedConfigMap", "POST", "/api/v1/namespaces/{namespace}/configmaps"),
		op("k8s", "readCoreV1NamespacedConfigMap", "GET", "/api/v1/namespaces/{namespace}/configmaps/{name}"),
		op("k8s", "replaceCoreV1NamespacedConfigMap", "PUT", "/api/v1/namespaces/{namespace}/configmaps/{name}"),
		op("k8s", "deleteCoreV1NamespacedConfigMap", "DELETE", "/api/v1/namespaces/{namespace}/configmaps/{name}"),
	}
	expanded := Expand(operations, operations[0], Options{Goal: "create, read, update, and delete configmaps", DesiredState: true})
	if got := roleIDs(expanded); !slices.Equal(got, []string{"create:createCoreV1NamespacedConfigMap", "read:readCoreV1NamespacedConfigMap", "update:replaceCoreV1NamespacedConfigMap", "delete:deleteCoreV1NamespacedConfigMap"}) {
		t.Fatalf("roles = %#v", got)
	}
}

func TestDiscoveryUploadNormalizationRequiresExplicitProvenance(t *testing.T) {
	seed := op("storage", "storage.objects.insert", "POST", "/upload/storage/v1/b/{bucket}/o")
	read := op("storage", "storage.objects.get", "GET", "/storage/v1/b/{bucket}/o/{object}")
	if got := roleIDs(Expand([]apitools.OperationSummary{seed, read}, seed, Options{DesiredState: true})); !slices.Equal(got, []string{"post:storage.objects.insert"}) {
		t.Fatalf("unprovenanced roles = %#v", got)
	}
	seed.Extensions = map[string]string{"x-uws-source-kind": apitools.APISourceKindGoogleDiscovery}
	read.Extensions = map[string]string{"x-uws-source-kind": apitools.APISourceKindGoogleDiscovery}
	if got := roleIDs(Expand([]apitools.OperationSummary{seed, read}, seed, Options{DesiredState: true})); !slices.Equal(got, []string{"create:storage.objects.insert", "read:storage.objects.get"}) {
		t.Fatalf("Discovery roles = %#v", got)
	}
}

func TestExpandGoogleDotOperationIDs(t *testing.T) {
	operations := []apitools.OperationSummary{
		op("storage", "storage.buckets.insert", "POST", "/b"),
		op("storage", "storage.buckets.get", "GET", "/b/{bucket}"),
		op("storage", "storage.buckets.patch", "PATCH", "/b/{bucket}"),
		op("storage", "storage.buckets.update", "PUT", "/b/{bucket}"),
		op("storage", "storage.buckets.delete", "DELETE", "/b/{bucket}"),
	}
	expanded := Expand(operations, operations[0], Options{Goal: "create, read, update, and delete buckets", DesiredState: true})
	if got := roleIDs(expanded); !slices.Equal(got, []string{"create:storage.buckets.insert", "read:storage.buckets.get", "update:storage.buckets.patch", "delete:storage.buckets.delete"}) {
		t.Fatalf("roles = %#v", got)
	}
}

func TestNormalizePathScopesUploadHandling(t *testing.T) {
	if got := normalizePath("/upload/storage/v1/b/{bucket}/o", false); got != "/upload/storage/v1/b/{bucket}/o" {
		t.Fatalf("unprovenanced upload path = %q", got)
	}
	if got := normalizePath("/upload/storage/v1/b/{bucket}/o", true); got != "/storage/v1/b/{bucket}/o" {
		t.Fatalf("Discovery upload path = %q", got)
	}
	if got := normalizePath("/uploading/storage/v1/b/{bucket}/o", true); got != "/uploading/storage/v1/b/{bucket}/o" {
		t.Fatalf("uploading path = %q", got)
	}
}

func TestExpandHyphenAndCreateOrUpdateFamilies(t *testing.T) {
	cloudflare := []apitools.OperationSummary{
		op("cloudflare", "r2-create-bucket", "POST", "/accounts/{account_id}/r2/buckets"),
		op("cloudflare", "r2-get-bucket", "GET", "/accounts/{account_id}/r2/buckets/{bucket_name}"),
		op("cloudflare", "r2-patch-bucket", "PATCH", "/accounts/{account_id}/r2/buckets/{bucket_name}"),
		op("cloudflare", "r2-delete-bucket", "DELETE", "/accounts/{account_id}/r2/buckets/{bucket_name}"),
	}
	if got := roleIDs(Expand(cloudflare, cloudflare[0], Options{Goal: "create, read, update, and delete buckets", DesiredState: true})); !slices.Equal(got, []string{"create:r2-create-bucket", "read:r2-get-bucket", "update:r2-patch-bucket", "delete:r2-delete-bucket"}) {
		t.Fatalf("Cloudflare roles = %#v", got)
	}
	azure := []apitools.OperationSummary{
		op("azure", "Databases_CreateOrUpdate", "PUT", "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Sql/servers/{serverName}/databases/{databaseName}"),
		op("azure", "Databases_Get", "GET", "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Sql/servers/{serverName}/databases/{databaseName}"),
		op("azure", "Databases_Delete", "DELETE", "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Sql/servers/{serverName}/databases/{databaseName}"),
	}
	if got := roleIDs(Expand(azure, azure[0], Options{DesiredState: true})); !slices.Equal(got, []string{"create:Databases_CreateOrUpdate", "read:Databases_Get", "delete:Databases_Delete"}) {
		t.Fatalf("Azure roles = %#v", got)
	}
}

func TestExpandRejectsAmbiguousSibling(t *testing.T) {
	operations := []apitools.OperationSummary{op("widgets", "createWidget", "POST", "/widgets"), op("widgets", "getWidget", "GET", "/widgets/{id}"), op("widgets", "readWidget", "GET", "/widgets/{id}")}
	expanded := Expand(operations, operations[0], Options{DesiredState: true})
	if len(expanded.Diagnostics) != 1 || expanded.Diagnostics[0].Code != "operation_lifecycle.ambiguous_read" {
		t.Fatalf("diagnostics = %#v", expanded.Diagnostics)
	}
}

func TestExpandPreservesAPIFirstSingleOperationRoles(t *testing.T) {
	for _, test := range []struct{ id, method, role string }{
		{id: "createWidget", method: "POST", role: "post"},
		{id: "putWidget", method: "PUT", role: "put"},
		{id: "deleteWidget", method: "DELETE", role: "delete"},
	} {
		operation := op("api", test.id, test.method, "/widgets/{id}")
		if got := roleIDs(Expand([]apitools.OperationSummary{operation}, operation, Options{})); !slices.Equal(got, []string{test.role + ":" + test.id}) {
			t.Fatalf("%s roles = %#v", test.id, got)
		}
	}
}

func op(source, id, method, path string) apitools.OperationSummary {
	return apitools.OperationSummary{ID: id, OperationID: id, DocumentName: source, Method: method, Path: path}
}

func googleOp(source, id, method, path string) apitools.OperationSummary {
	operation := op(source, id, method, path)
	operation.Extensions = map[string]string{"x-uws-source-kind": apitools.APISourceKindGoogleDiscovery}
	return operation
}

func roleIDs(expanded Expansion) []string {
	out := make([]string, 0, len(expanded.Roles))
	for _, role := range expanded.Roles {
		out = append(out, role.Role+":"+role.Operation.OperationID)
	}
	return out
}
