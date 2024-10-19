export const ROLE_ORGANIZATION_ADMIN = 'ROLE_ORGANIZATION_ADMIN';
export const ROLE_TEAM_MEMBER = 'ROLE_TEAM_MEMBER';
export const ROLE_ALLOWED_TO_SWITCH = 'ROLE_ALLOWED_TO_SWITCH';

export const USER_FRIENDLY_ROLES = [
	{system: ROLE_ORGANIZATION_ADMIN, user_friendly: "Administrator"},
	{system: ROLE_TEAM_MEMBER, user_friendly: "Team member"},
];

export const showUserFriendlyRoles = (roles) => {
	var rolesToReturn = [];

	for (const e of USER_FRIENDLY_ROLES) {
  		if (roles.includes(e.system)) {
  			rolesToReturn.push(e.user_friendly);
  		}
	}

	return rolesToReturn.join(', ');
}
