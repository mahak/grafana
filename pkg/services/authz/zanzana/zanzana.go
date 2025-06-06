package zanzana

import (
	"fmt"
	"strings"

	openfgav1 "github.com/openfga/api/proto/openfga/v1"

	authlib "github.com/grafana/authlib/types"

	"github.com/grafana/grafana/pkg/services/authz/zanzana/common"
)

const (
	TypeUser           = common.TypeUser
	TypeServiceAccount = common.TypeServiceAccount
	TypeRenderService  = common.TypeRenderService
	TypeAnonymous      = common.TypeAnonymous
	TypeTeam           = common.TypeTeam
	TypeRole           = common.TypeRole
	TypeFolder         = common.TypeFolder
	TypeResource       = common.TypeResource
	TypeNamespace      = common.TypeGroupResouce
)

const (
	RelationTeamMember = common.RelationTeamMember
	RelationTeamAdmin  = common.RelationTeamAdmin
	RelationParent     = common.RelationParent
	RelationAssignee   = common.RelationAssignee

	RelationSetView  = common.RelationSetView
	RelationSetEdit  = common.RelationSetEdit
	RelationSetAdmin = common.RelationSetAdmin

	RelationGet    = common.RelationGet
	RelationUpdate = common.RelationUpdate
	RelationCreate = common.RelationCreate
	RelationDelete = common.RelationDelete

	RelationSubresourceSetView  = common.RelationSubresourceSetView
	RelationSubresourceSetEdit  = common.RelationSubresourceSetEdit
	RelationSubresourceSetAdmin = common.RelationSubresourceSetAdmin

	RelationSubresourceRead   = common.RelationSubresourceGet
	RelationSubresourceWrite  = common.RelationSubresourceUpdate
	RelationSubresourceCreate = common.RelationSubresourceCreate
	RelationSubresourceDelete = common.RelationSubresourceDelete
)

var (
	RelationsFolder      = common.RelationsTyped
	RelationsResouce     = common.RelationsResource
	RelationsSubresource = common.RelationsSubresource
)

const (
	KindDashboards string = "dashboards"
	KindFolders    string = "folders"
)

var (
	ToAuthzExtTupleKey                  = common.ToAuthzExtTupleKey
	ToAuthzExtTupleKeys                 = common.ToAuthzExtTupleKeys
	ToAuthzExtTupleKeyWithoutCondition  = common.ToAuthzExtTupleKeyWithoutCondition
	ToAuthzExtTupleKeysWithoutCondition = common.ToAuthzExtTupleKeysWithoutCondition

	ToOpenFGATuple                    = common.ToOpenFGATuple
	ToOpenFGATuples                   = common.ToOpenFGATuples
	ToOpenFGATupleKey                 = common.ToOpenFGATupleKey
	ToOpenFGATupleKeyWithoutCondition = common.ToOpenFGATupleKeyWithoutCondition
)

// NewTupleEntry constructs new openfga entry type:name[#relation].
// Relation allows to specify group of users (subjects) related to type:name
// (for example, team:devs#member refers to users which are members of team devs)
func NewTupleEntry(objectType, name, relation string) string {
	obj := fmt.Sprintf("%s:%s", objectType, name)
	if relation != "" {
		obj = fmt.Sprintf("%s#%s", obj, relation)
	}
	return obj
}

func TranslateToResourceTuple(subject string, action, kind, name string) (*openfgav1.TupleKey, bool) {
	translation, ok := resourceTranslations[kind]

	if !ok {
		return nil, false
	}

	m, ok := translation.mapping[action]
	if !ok {
		return nil, false
	}

	if name == "*" {
		return common.NewGroupResourceTuple(subject, m.relation, translation.group, translation.resource, m.subresource), true
	}

	if translation.typ == TypeResource {
		return common.NewResourceTuple(subject, m.relation, translation.group, translation.resource, m.subresource, name), true
	}

	if translation.typ == TypeFolder {
		if m.group != "" && m.resource != "" {
			return common.NewFolderResourceTuple(subject, m.relation, m.group, m.resource, m.subresource, name), true
		}

		return common.NewFolderTuple(subject, m.relation, name), true
	}

	return common.NewTypedTuple(translation.typ, subject, m.relation, name), true
}

func IsFolderResourceTuple(t *openfgav1.TupleKey) bool {
	return strings.HasPrefix(t.Object, TypeFolder) && strings.HasPrefix(t.Relation, "resource_")
}

func MergeFolderResourceTuples(a, b *openfgav1.TupleKey) {
	va := a.Condition.Context.Fields["subresources"]
	vb := b.Condition.Context.Fields["subresources"]
	va.GetListValue().Values = append(va.GetListValue().Values, vb.GetListValue().Values...)
}

func TranslateToCheckRequest(namespace, action, kind, folder, name string) (*authlib.CheckRequest, bool) {
	translation, ok := resourceTranslations[kind]

	if !ok {
		return nil, false
	}

	m, ok := translation.mapping[action]
	if !ok {
		return nil, false
	}

	verb, ok := common.RelationToVerbMapping[m.relation]
	if !ok {
		return nil, false
	}

	req := &authlib.CheckRequest{
		Namespace: namespace,
		Verb:      verb,
		Group:     translation.group,
		Resource:  translation.resource,
		Name:      name,
		Folder:    folder,
	}

	return req, true
}

func TranslateToListRequest(namespace, action, kind string) (*authlib.ListRequest, bool) {
	translation, ok := resourceTranslations[kind]

	if !ok {
		return nil, false
	}

	// FIXME: support different verbs
	req := &authlib.ListRequest{
		Namespace: namespace,
		Group:     translation.group,
		Resource:  translation.resource,
	}

	return req, true
}

func TranslateToGroupResource(kind string) string {
	translation, ok := resourceTranslations[kind]
	if !ok {
		return ""
	}
	return common.FormatGroupResource(translation.group, translation.resource, "")
}

func TranslateBasicRole(name string) string {
	return basicRolesTranslations[name]
}
