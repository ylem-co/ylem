import {
    REGISTER_SUCCESS,
    REGISTER_FAIL,
    LOGIN_SUCCESS,
    LOGIN_FAIL,
    LOGOUT,
    SET_MESSAGE,
} from "./types";

import {
    prepareErrorMessage,
} from "../actions/errors";

import AuthService from "../services/auth.service";

export const register = (firstName, lastName, email, password, confirmPassword, phone, organizationName = null, invitationKey = null) => (dispatch) => {
    return AuthService.register(firstName, lastName, email, password, confirmPassword, phone, organizationName, invitationKey).then(
        (response) => {
            dispatch({
                type: REGISTER_SUCCESS,
            });

            dispatch({
                type: SET_MESSAGE,
                payload: "Your profile successfully created",
            });

            return Promise.resolve();
        },
        (error) => {
            let message = prepareErrorMessage(error);

            dispatch({
                type: REGISTER_FAIL,
            });

            dispatch({
                type: SET_MESSAGE,
                payload: message,
            });

            return Promise.reject();
        }
    );
};

export const login = (email, password) => (dispatch) => {
    return AuthService.login(email, password).then(
        (data) => {
            dispatch({
                type: LOGIN_SUCCESS,
                payload: { user: data },
            });

            return Promise.resolve();
        },
        (error) => {
            const message =
                (error.response &&
                    error.response.data &&
                    error.response.data.error) ||
                (error.response &&
                    error.response.data &&
                    error.response.data.message) ||
                error.message ||
                error.toString();

            dispatch({
                type: LOGIN_FAIL,
            });

            dispatch({
                type: SET_MESSAGE,
                payload: message,
            });

            return Promise.reject();
        }
    );
};

export const signInWithGoogle = (queryString) => (dispatch) => {
    return AuthService.signInWithGoogle(queryString).then(
        (data) => {
            dispatch({
                type: LOGIN_SUCCESS,
                payload: { user: data },
            });

            return Promise.resolve();
        },
        (error) => {
            const message =
                (error.response &&
                    error.response.data &&
                    error.response.data.error) ||
                error.message ||
                error.toString();

            dispatch({
                type: LOGIN_FAIL,
            });

            dispatch({
                type: SET_MESSAGE,
                payload: message,
            });

            return Promise.reject();
        }
    );
};

export const logout = () => async(dispatch) => {
    await AuthService.logout();

    dispatch({
        type: LOGOUT,
    });
};
