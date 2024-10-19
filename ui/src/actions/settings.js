import {
    ORGANIZATION_UPDATE_SUCCESS,
    ORGANIZATION_UPDATE_FAIL,
    USER_UPDATE_SUCCESS,
    USER_UPDATE_FAIL,
    PASSWORD_UPDATE_SUCCESS,
    PASSWORD_UPDATE_FAIL,
    SET_MESSAGE,
} from "./types";

import {
    prepareErrorMessage
} from "./errors";

import OrganizationService from "../services/organization.service";
import UserService from "../services/user.service";

export const updateOrganization = (organizationUuid, organizationName) => (dispatch) => {
    return OrganizationService.updateOrganization(organizationUuid, organizationName).then(
        (response) => {
            dispatch({
                type: ORGANIZATION_UPDATE_SUCCESS,
            });

            dispatch({
                type: SET_MESSAGE,
                payload: "Organization successfully updated",
            });

            return Promise.resolve();
        },
        (error) => {
            let message = prepareErrorMessage(error);

            dispatch({
                type: ORGANIZATION_UPDATE_FAIL,
            });

            dispatch({
                type: SET_MESSAGE,
                payload: message,
            });

            return Promise.reject();
        }
    );
};

export const updateMe = (uuid, firstName, lastName, email, phone) => (dispatch) => {
    return UserService.updateMe(uuid, firstName, lastName, email, phone).then(
        (response) => {
            dispatch({
                type: USER_UPDATE_SUCCESS,
            });

            dispatch({
                type: SET_MESSAGE,
                payload: "Data successfully updated",
            });

            return Promise.resolve();
        },
        (error) => {
            let message = prepareErrorMessage(error);

            dispatch({
                type: USER_UPDATE_FAIL,
            });

            dispatch({
                type: SET_MESSAGE,
                payload: message,
            });

            return Promise.reject();
        }
    );
};

export const updatePassword = (uuid, password, confirmPassword) => (dispatch) => {
    return UserService.updatePassword(uuid, password, confirmPassword).then(
        (response) => {
            dispatch({
                type: PASSWORD_UPDATE_SUCCESS,
            });

            dispatch({
                type: SET_MESSAGE,
                payload: "Password successfully updated. Please log in again",
            });

            return Promise.resolve();
        },
        (error) => {
            let message = prepareErrorMessage(error);

            dispatch({
                type: PASSWORD_UPDATE_FAIL,
            });

            dispatch({
                type: SET_MESSAGE,
                payload: message,
            });

            return Promise.reject();
        }
    );
};
