import {
    OAUTH_CLIENT_CREATE_FAIL,
    OAUTH_CLIENT_CREATE_SUCCESS,
    SET_MESSAGE,
} from "./types";

import {
    prepareErrorMessage,
} from "../actions/errors";

import OAuthService from "../services/oauth.service";

export const addOAuthClient = (name) => (dispatch) => {
    return OAuthService.createClient(name).then(
        (data) => {
            dispatch({
                type: OAUTH_CLIENT_CREATE_SUCCESS,
            });

            dispatch({
                type: SET_MESSAGE,
                payload: { client: data },
            });

            return Promise.resolve();
        },
        (error) => {
            let message = prepareErrorMessage(error);

            dispatch({
                type: OAUTH_CLIENT_CREATE_FAIL,
            });

            dispatch({
                type: SET_MESSAGE,
                payload: message,
            });

            return Promise.reject();
        }
    );
};
