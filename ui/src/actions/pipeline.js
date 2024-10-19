export const PERMISSION_LOGGED_IN = 'PERMISSION_LOGGED_IN';
export const PERMISSION_LOGGED_OUT = 'PERMISSION_LOGGED_OUT';

export const validatePermissions = (isLoggedIn, user, requiredPermission) => {
    if (
        requiredPermission === PERMISSION_LOGGED_IN
        && !isLoggedIn
    ) {
        return false;
    } else if (
        requiredPermission === PERMISSION_LOGGED_OUT
        && isLoggedIn
    ) {
        return false;
    } 

    return true;
};
