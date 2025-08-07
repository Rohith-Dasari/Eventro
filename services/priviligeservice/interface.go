package privilegeservice

import "context"

type PrivilegeServiceInterface interface {
	EscalatePrivilege(ctx context.Context)
}
